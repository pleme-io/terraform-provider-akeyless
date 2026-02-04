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

func resourceGatewayMigrationGcp() *schema.Resource {
	return &schema.Resource{
		Description: "GCP Migration resource",
		Create:      resourceGatewayMigrationGcpCreate,
		Read:        resourceGatewayMigrationGcpRead,
		Update:      resourceGatewayMigrationGcpUpdate,
		Delete:      resourceGatewayMigrationGcpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationGcpImport,
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
			"gcp_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Base64-encoded GCP Service Account private key text with sufficient permissions to Secrets Manager, Minimum required permission is Secret Manager Secret Accessor, e.g. 'roles/secretmanager.secretAccessor' (relevant only for GCP migration)",
			},
			"protection_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the key that protects the classic key value (if empty, the account default key will be used)",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Migration ID",
			},
		},
	}
}

func resourceGatewayMigrationGcpCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	gcpKey := d.Get("gcp_key").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayCreateMigration("", name, "", "", targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.GcpKey, gcpKey)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration GCP: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration GCP: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationGcpRead(d *schema.ResourceData, m interface{}) error {
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
		if rOut.Body.GcpSecretsMigrations != nil && len(rOut.Body.GcpSecretsMigrations) > 0 {
			for _, migration := range rOut.Body.GcpSecretsMigrations {
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

func resourceGatewayMigrationGcpUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	gcpKey := d.Get("gcp_key").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration("", "", "", targetLocation)
	body.Token = &token
	body.Name = &name
	common.GetAkeylessPtr(&body.GcpKey, gcpKey)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationGcpRead(d, m)
		if err != nil {
			return err
		}
	}
	id = d.Get("id").(string)
	body.Id = &id

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update Gateway Migration GCP: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update Gateway Migration GCP: %v", err)
	}

	return nil
}

func resourceGatewayMigrationGcpDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationGcpRead(d, m)
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

func resourceGatewayMigrationGcpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationGcpRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
