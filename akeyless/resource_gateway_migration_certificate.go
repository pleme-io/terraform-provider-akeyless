package akeyless

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	akeyless_api "github.com/akeylesslabs/akeyless-go/v5"
	"github.com/akeylesslabs/terraform-provider-akeyless/akeyless/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGatewayMigrationCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "Certificate Migration resource",
		Create:      resourceGatewayMigrationCertificateCreate,
		Read:        resourceGatewayMigrationCertificateRead,
		Update:      resourceGatewayMigrationCertificateUpdate,
		Delete:      resourceGatewayMigrationCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationCertificateImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Migration name",
				ForceNew:    true,
			},
			"target_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target location in Akeyless for imported secrets",
			},
			"hosts": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A comma separated list of IPs, CIDR ranges, or DNS names to scan",
			},
			"port_ranges": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A comma separated list of port ranges",
			},
			"expiration_event_in": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "How many days before the expiration of the certificate would you like to be notified",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"protection_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a key that used to encrypt the secret value (if empty, the account default protectionKey key will be used)",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Migration ID",
			},
		},
	}
}

func resourceGatewayMigrationCertificateCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hosts := d.Get("hosts").(string)
	portRanges := d.Get("port_ranges").(string)
	expirationEventIn := d.Get("expiration_event_in").([]interface{})
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayCreateMigration(hosts, name, "", "", targetLocation)
	body.Token = &token
	body.Type = akeyless_api.PtrString("certificate")
	common.GetAkeylessPtr(&body.PortRanges, portRanges)
	if len(expirationEventIn) > 0 {
		body.ExpirationEventIn = common.ExpandStringList(expirationEventIn)
	}
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration Certificate: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration Certificate: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationCertificateRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()

	path := d.Id()

	body := akeyless_api.GatewayGetMigration{
		Name:  &path,
		Token: &token,
	}

	rOut, res, err := client.GatewayGetMigration(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			if res.StatusCode == http.StatusNotFound {
				// The resource was deleted outside of the current Terraform workspace, so invalidate this resource
				d.SetId("")
				return nil
			}
			return fmt.Errorf("can't value: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't get value: %v", err)
	}

	if rOut.Body != nil {
		if rOut.Body.CertificateMigrations != nil && len(rOut.Body.CertificateMigrations) > 0 {
			for _, migration := range rOut.Body.CertificateMigrations {
				if migration.General != nil && *migration.General.Name == path {
					id := migration.General.Id
					if id != nil {
						err = d.Set("id", *id)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	d.SetId(path)

	return nil
}

func resourceGatewayMigrationCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hosts := d.Get("hosts").(string)
	portRanges := d.Get("port_ranges").(string)
	expirationEventIn := d.Get("expiration_event_in").([]interface{})
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration(hosts, "", "", targetLocation)
	body.Token = &token
	body.Name = &name
	common.GetAkeylessPtr(&body.PortRanges, portRanges)
	if len(expirationEventIn) > 0 {
		body.ExpirationEventIn = common.ExpandStringList(expirationEventIn)
	}
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationCertificateRead(d, m)
		if err != nil {
			return err
		}
	}
	id = d.Get("id").(string)
	body.Id = &id

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update Gateway Migration Certificate: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update Gateway Migration Certificate: %v", err)
	}

	return nil
}

func resourceGatewayMigrationCertificateDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationCertificateRead(d, m)
		if err != nil {
			return err
		}
	}
	id = d.Get("id").(string)

	deleteItem := akeyless_api.GatewayDeleteMigration{
		Token: &token,
		Id:    id,
	}

	ctx := context.Background()
	_, _, err := client.GatewayDeleteMigration(ctx).Body(deleteItem).Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceGatewayMigrationCertificateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationCertificateRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
