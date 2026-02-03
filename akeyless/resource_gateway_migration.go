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

func resourceGatewayMigration() *schema.Resource {
	return &schema.Resource{
		Description: "Gateway Migration resource",
		Create:      resourceGatewayMigrationCreate,
		Read:        resourceGatewayMigrationRead,
		Update:      resourceGatewayMigrationUpdate,
		Delete:      resourceGatewayMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Migration name",
				ForceNew:    true,
			},
			"hosts": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A comma separated list of IPs, CIDR ranges, or DNS names to scan",
			},
			"si_target_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH, Windows or Linked Target Name (Relevant only for Server Inventory migration)",
			},
			"si_users_path_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path location template for migrating users as Rotated Secrets (Relevant only for Server Inventory migration)",
			},
			"target_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target location in Akeyless for imported secrets",
			},
			"aws_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Access Key (relevant only for AWS migration)",
			},
			"aws_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Access Key ID (relevant only for AWS migration)",
			},
			"aws_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS region (relevant only for AWS migration)",
			},
			"azure_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Key Vault Access client ID (relevant only for Azure Key Vault migration)",
			},
			"azure_kv_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Key Vault name (relevant only for Azure Key Vault migration)",
			},
			"azure_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Key Vault secret (relevant only for Azure Key Vault migration)",
			},
			"azure_tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Key Vault Access tenant ID (relevant only for Azure Key Vault migration)",
			},
			"gcp_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Base64-encoded GCP Service Account private key text (relevant only for GCP migration)",
			},
			"json": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set output format to JSON",
				Default:     false,
			},
			"k8s_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s namespace (relevant only for K8s migration)",
			},
			"k8s_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "K8s Client password (relevant only for K8s migration with Password Authentication method)",
			},
			"k8s_skip_system": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "K8s skip system secrets (relevant only for K8s migration)",
			},
			"k8s_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "K8s Bearer Token (relevant only for K8s migration with Token Authentication method)",
			},
			"k8s_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s API server URL (relevant only for K8s migration)",
			},
			"k8s_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s Client username (relevant only for K8s migration with Password Authentication method)",
			},
			"protection_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a key that used to encrypt the secret value (if empty, the account default protectionKey key will be used)",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Migration type (hashi/aws/gcp/k8s/azure_kv/active_directory/server_inventory/certificate)",
			},
		},
	}
}

func resourceGatewayMigrationCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	hosts := d.Get("hosts").(string)
	siTargetName := d.Get("si_target_name").(string)
	siUsersPathTemplate := d.Get("si_users_path_template").(string)
	targetLocation := d.Get("target_location").(string)
	awsKey := d.Get("aws_key").(string)
	awsKeyId := d.Get("aws_key_id").(string)
	awsRegion := d.Get("aws_region").(string)
	azureClientId := d.Get("azure_client_id").(string)
	azureKvName := d.Get("azure_kv_name").(string)
	azureSecret := d.Get("azure_secret").(string)
	azureTenantId := d.Get("azure_tenant_id").(string)
	gcpKey := d.Get("gcp_key").(string)
	jsonOutput := d.Get("json").(bool)
	k8sNamespace := d.Get("k8s_namespace").(string)
	k8sPassword := d.Get("k8s_password").(string)
	k8sSkipSystem := d.Get("k8s_skip_system").(bool)
	k8sToken := d.Get("k8s_token").(string)
	k8sUrl := d.Get("k8s_url").(string)
	k8sUsername := d.Get("k8s_username").(string)
	protectionKey := d.Get("protection_key").(string)
	migrationType := d.Get("type").(string)

	body := akeyless_api.NewGatewayCreateMigration(hosts, name, siTargetName, siUsersPathTemplate, targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.AwsKey, awsKey)
	common.GetAkeylessPtr(&body.AwsKeyId, awsKeyId)
	common.GetAkeylessPtr(&body.AwsRegion, awsRegion)
	common.GetAkeylessPtr(&body.AzureClientId, azureClientId)
	common.GetAkeylessPtr(&body.AzureKvName, azureKvName)
	common.GetAkeylessPtr(&body.AzureSecret, azureSecret)
	common.GetAkeylessPtr(&body.AzureTenantId, azureTenantId)
	common.GetAkeylessPtr(&body.GcpKey, gcpKey)
	common.GetAkeylessPtr(&body.Json, jsonOutput)
	common.GetAkeylessPtr(&body.K8sNamespace, k8sNamespace)
	common.GetAkeylessPtr(&body.K8sPassword, k8sPassword)
	common.GetAkeylessPtr(&body.K8sSkipSystem, k8sSkipSystem)
	common.GetAkeylessPtr(&body.K8sToken, k8sToken)
	common.GetAkeylessPtr(&body.K8sUrl, k8sUrl)
	common.GetAkeylessPtr(&body.K8sUsername, k8sUsername)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)
	common.GetAkeylessPtr(&body.Type, migrationType)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationRead(d *schema.ResourceData, m interface{}) error {
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

	// The migration exists, keep the ID
	d.SetId(path)

	return nil
}

func resourceGatewayMigrationUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	hosts := d.Get("hosts").(string)
	siTargetName := d.Get("si_target_name").(string)
	siUsersPathTemplate := d.Get("si_users_path_template").(string)
	targetLocation := d.Get("target_location").(string)
	awsKey := d.Get("aws_key").(string)
	awsKeyId := d.Get("aws_key_id").(string)
	awsRegion := d.Get("aws_region").(string)
	azureClientId := d.Get("azure_client_id").(string)
	azureKvName := d.Get("azure_kv_name").(string)
	azureSecret := d.Get("azure_secret").(string)
	azureTenantId := d.Get("azure_tenant_id").(string)
	gcpKey := d.Get("gcp_key").(string)
	jsonOutput := d.Get("json").(bool)
	k8sNamespace := d.Get("k8s_namespace").(string)
	k8sPassword := d.Get("k8s_password").(string)
	k8sSkipSystem := d.Get("k8s_skip_system").(bool)
	k8sToken := d.Get("k8s_token").(string)
	k8sUrl := d.Get("k8s_url").(string)
	k8sUsername := d.Get("k8s_username").(string)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration(hosts, siTargetName, siUsersPathTemplate, targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.AwsKey, awsKey)
	common.GetAkeylessPtr(&body.AwsKeyId, awsKeyId)
	common.GetAkeylessPtr(&body.AwsRegion, awsRegion)
	common.GetAkeylessPtr(&body.AzureClientId, azureClientId)
	common.GetAkeylessPtr(&body.AzureKvName, azureKvName)
	common.GetAkeylessPtr(&body.AzureSecret, azureSecret)
	common.GetAkeylessPtr(&body.AzureTenantId, azureTenantId)
	common.GetAkeylessPtr(&body.GcpKey, gcpKey)
	common.GetAkeylessPtr(&body.Json, jsonOutput)
	common.GetAkeylessPtr(&body.K8sNamespace, k8sNamespace)
	common.GetAkeylessPtr(&body.K8sPassword, k8sPassword)
	common.GetAkeylessPtr(&body.K8sSkipSystem, k8sSkipSystem)
	common.GetAkeylessPtr(&body.K8sToken, k8sToken)
	common.GetAkeylessPtr(&body.K8sUrl, k8sUrl)
	common.GetAkeylessPtr(&body.K8sUsername, k8sUsername)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update : %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update : %v", err)
	}

	return nil
}

func resourceGatewayMigrationDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	id := d.Id()

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

func resourceGatewayMigrationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
