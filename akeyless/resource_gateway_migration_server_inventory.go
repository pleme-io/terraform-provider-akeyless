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

func resourceGatewayMigrationServerInventory() *schema.Resource {
	return &schema.Resource{
		Description: "Server Inventory Migration resource",
		Create:      resourceGatewayMigrationServerInventoryCreate,
		Read:        resourceGatewayMigrationServerInventoryRead,
		Update:      resourceGatewayMigrationServerInventoryUpdate,
		Delete:      resourceGatewayMigrationServerInventoryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationServerInventoryImport,
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
			"si_target_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH, Windows or Linked Target Name",
			},
			"si_users_path_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path location template for migrating users as Rotated Secrets",
			},
			"si_auto_rotate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable/Disable automatic/recurrent rotation for migrated secrets",
			},
			"si_rotation_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The hour of the scheduled rotation in UTC",
			},
			"si_rotation_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of days to wait between every automatic rotation [1-365]",
			},
			"si_sra_enable_rdp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable/Disable RDP Secure Remote Access for the migrated local users rotated secrets",
			},
			"si_user_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma-separated list of groups to migrate users from",
			},
			"si_users_ignore": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma-separated list of Local Users which should not be migrated",
			},
			"protection_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a key that used to encrypt the secret value (if empty, the account default protectionKey key will be used)",
			},
		},
	}
}

func resourceGatewayMigrationServerInventoryCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hosts := d.Get("hosts").(string)
	siTargetName := d.Get("si_target_name").(string)
	siUsersPathTemplate := d.Get("si_users_path_template").(string)
	siAutoRotate := d.Get("si_auto_rotate").(string)
	siRotationHour := d.Get("si_rotation_hour").(int)
	siRotationInterval := d.Get("si_rotation_interval").(int)
	siSraEnableRdp := d.Get("si_sra_enable_rdp").(string)
	siUserGroups := d.Get("si_user_groups").(string)
	siUsersIgnore := d.Get("si_users_ignore").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayCreateMigration(hosts, name, siTargetName, siUsersPathTemplate, targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.SiAutoRotate, siAutoRotate)
	if siRotationHour != 0 {
		body.SiRotationHour = common.PtrInt32(int32(siRotationHour))
	}
	if siRotationInterval != 0 {
		body.SiRotationInterval = common.PtrInt32(int32(siRotationInterval))
	}
	common.GetAkeylessPtr(&body.SiSraEnableRdp, siSraEnableRdp)
	common.GetAkeylessPtr(&body.SiUserGroups, siUserGroups)
	common.GetAkeylessPtr(&body.SiUsersIgnore, siUsersIgnore)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration Server Inventory: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration Server Inventory: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationServerInventoryRead(d *schema.ResourceData, m interface{}) error {
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

	_, res, err := client.GatewayGetMigration(ctx).Body(body).Execute()
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

	d.SetId(path)

	return nil
}

func resourceGatewayMigrationServerInventoryUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	hosts := d.Get("hosts").(string)
	siTargetName := d.Get("si_target_name").(string)
	siUsersPathTemplate := d.Get("si_users_path_template").(string)
	siAutoRotate := d.Get("si_auto_rotate").(string)
	siRotationHour := d.Get("si_rotation_hour").(int)
	siRotationInterval := d.Get("si_rotation_interval").(int)
	siSraEnableRdp := d.Get("si_sra_enable_rdp").(string)
	siUserGroups := d.Get("si_user_groups").(string)
	siUsersIgnore := d.Get("si_users_ignore").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration(hosts, siTargetName, siUsersPathTemplate, targetLocation)
	body.Token = &token
	body.Name = &name
	common.GetAkeylessPtr(&body.SiAutoRotate, siAutoRotate)
	if siRotationHour != 0 {
		body.SiRotationHour = common.PtrInt32(int32(siRotationHour))
	}
	if siRotationInterval != 0 {
		body.SiRotationInterval = common.PtrInt32(int32(siRotationInterval))
	}
	common.GetAkeylessPtr(&body.SiSraEnableRdp, siSraEnableRdp)
	common.GetAkeylessPtr(&body.SiUserGroups, siUserGroups)
	common.GetAkeylessPtr(&body.SiUsersIgnore, siUsersIgnore)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update Gateway Migration Server Inventory: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update Gateway Migration Server Inventory: %v", err)
	}

	return nil
}

func resourceGatewayMigrationServerInventoryDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	name := d.Get("name").(string)

	deleteItem := akeyless_api.GatewayDeleteMigration{
		Token: &token,
		Id:    name,
	}

	ctx := context.Background()
	_, _, err := client.GatewayDeleteMigration(ctx).Body(deleteItem).Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceGatewayMigrationServerInventoryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationServerInventoryRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
