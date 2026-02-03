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

func resourceGatewayMigrationHashi() *schema.Resource {
	return &schema.Resource{
		Description: "HashiCorp Vault Migration resource",
		Create:      resourceGatewayMigrationHashiCreate,
		Read:        resourceGatewayMigrationHashiRead,
		Update:      resourceGatewayMigrationHashiUpdate,
		Delete:      resourceGatewayMigrationHashiDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationHashiImport,
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
			"hashi_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HashiCorp Vault API URL",
			},
			"hashi_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "HashiCorp Vault access token",
			},
			"hashi_ns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HashiCorp Vault Namespaces",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"hashi_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Import secret key as json value or independent secrets",
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

func resourceGatewayMigrationHashiCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hashiUrl := d.Get("hashi_url").(string)
	hashiToken := d.Get("hashi_token").(string)
	hashiNs := d.Get("hashi_ns").([]interface{})
	hashiJson := d.Get("hashi_json").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayCreateMigration("", name, "", "", targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.HashiUrl, hashiUrl)
	common.GetAkeylessPtr(&body.HashiToken, hashiToken)
	if len(hashiNs) > 0 {
		body.HashiNs = common.ExpandStringList(hashiNs)
	}
	common.GetAkeylessPtr(&body.HashiJson, hashiJson)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration HashiCorp Vault: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration HashiCorp Vault: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationHashiRead(d *schema.ResourceData, m interface{}) error {
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
		if rOut.Body.HashiMigrations != nil && len(rOut.Body.HashiMigrations) > 0 {
			for _, migration := range rOut.Body.HashiMigrations {
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

func resourceGatewayMigrationHashiUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hashiUrl := d.Get("hashi_url").(string)
	hashiToken := d.Get("hashi_token").(string)
	hashiNs := d.Get("hashi_ns").([]interface{})
	hashiJson := d.Get("hashi_json").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration("", "", "", targetLocation)
	body.Token = &token
	body.Name = &name
	common.GetAkeylessPtr(&body.HashiUrl, hashiUrl)
	common.GetAkeylessPtr(&body.HashiToken, hashiToken)
	if len(hashiNs) > 0 {
		body.HashiNs = common.ExpandStringList(hashiNs)
	}
	common.GetAkeylessPtr(&body.HashiJson, hashiJson)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationHashiRead(d, m)
		if err != nil {
			return err
		}
	}
	id = d.Get("id").(string)
	body.Id = &id

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update Gateway Migration HashiCorp Vault: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update Gateway Migration HashiCorp Vault: %v", err)
	}

	return nil
}

func resourceGatewayMigrationHashiDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationHashiRead(d, m)
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

func resourceGatewayMigrationHashiImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationHashiRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
