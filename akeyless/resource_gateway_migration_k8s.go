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

func resourceGatewayMigrationK8s() *schema.Resource {
	return &schema.Resource{
		Description: "Kubernetes Migration resource",
		Create:      resourceGatewayMigrationK8sCreate,
		Read:        resourceGatewayMigrationK8sRead,
		Update:      resourceGatewayMigrationK8sUpdate,
		Delete:      resourceGatewayMigrationK8sDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGatewayMigrationK8sImport,
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
			"k8s_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s API server URL",
			},
			"k8s_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "K8s Bearer Token (relevant only for K8s migration with Token Authentication method)",
			},
			"k8s_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s Client username (relevant only for K8s migration with Password Authentication method)",
			},
			"k8s_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "K8s Client password (relevant only for K8s migration with Password Authentication method)",
			},
			"k8s_ca_certificate": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "K8s Cluster CA certificate (relevant only for K8s migration with Certificate Authentication method)",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"k8s_client_certificate": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "K8s Client certificate (relevant only for K8s migration with Certificate Authentication method)",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"k8s_client_key": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "K8s Client key (relevant only for K8s migration with Certificate Authentication method)",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"k8s_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8s namespace",
			},
			"k8s_skip_system": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "K8s skip system secrets",
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

func resourceGatewayMigrationK8sCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	k8sUrl := d.Get("k8s_url").(string)
	k8sToken := d.Get("k8s_token").(string)
	k8sUsername := d.Get("k8s_username").(string)
	k8sPassword := d.Get("k8s_password").(string)
	k8sCaCertificate := d.Get("k8s_ca_certificate").([]interface{})
	k8sClientCertificate := d.Get("k8s_client_certificate").([]interface{})
	k8sClientKey := d.Get("k8s_client_key").([]interface{})
	k8sNamespace := d.Get("k8s_namespace").(string)
	k8sSkipSystem := d.Get("k8s_skip_system").(bool)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayCreateMigration("", name, "", "", targetLocation)
	body.Token = &token
	common.GetAkeylessPtr(&body.K8sUrl, k8sUrl)
	common.GetAkeylessPtr(&body.K8sToken, k8sToken)
	common.GetAkeylessPtr(&body.K8sUsername, k8sUsername)
	common.GetAkeylessPtr(&body.K8sPassword, k8sPassword)
	if len(k8sCaCertificate) > 0 {
		k8sCaCertificateInt32 := make([]int32, len(k8sCaCertificate))
		for i, v := range k8sCaCertificate {
			k8sCaCertificateInt32[i] = int32(v.(int))
		}
		body.K8sCaCertificate = k8sCaCertificateInt32
	}
	if len(k8sClientCertificate) > 0 {
		k8sClientCertificateInt32 := make([]int32, len(k8sClientCertificate))
		for i, v := range k8sClientCertificate {
			k8sClientCertificateInt32[i] = int32(v.(int))
		}
		body.K8sClientCertificate = k8sClientCertificateInt32
	}
	if len(k8sClientKey) > 0 {
		k8sClientKeyInt32 := make([]int32, len(k8sClientKey))
		for i, v := range k8sClientKey {
			k8sClientKeyInt32[i] = int32(v.(int))
		}
		body.K8sClientKey = k8sClientKeyInt32
	}
	common.GetAkeylessPtr(&body.K8sNamespace, k8sNamespace)
	common.GetAkeylessPtr(&body.K8sSkipSystem, k8sSkipSystem)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	_, _, err := client.GatewayCreateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Gateway Migration K8s: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Gateway Migration K8s: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceGatewayMigrationK8sRead(d *schema.ResourceData, m interface{}) error {
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
		if rOut.Body.K8sMigrations != nil && len(rOut.Body.K8sMigrations) > 0 {
			for _, migration := range rOut.Body.K8sMigrations {
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

func resourceGatewayMigrationK8sUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetLocation := d.Get("target_location").(string)
	k8sUrl := d.Get("k8s_url").(string)
	k8sToken := d.Get("k8s_token").(string)
	k8sUsername := d.Get("k8s_username").(string)
	k8sPassword := d.Get("k8s_password").(string)
	k8sCaCertificate := d.Get("k8s_ca_certificate").([]interface{})
	k8sClientCertificate := d.Get("k8s_client_certificate").([]interface{})
	k8sClientKey := d.Get("k8s_client_key").([]interface{})
	k8sNamespace := d.Get("k8s_namespace").(string)
	k8sSkipSystem := d.Get("k8s_skip_system").(bool)
	protectionKey := d.Get("protection_key").(string)

	body := akeyless_api.NewGatewayUpdateMigration("", "", "", targetLocation)
	body.Token = &token
	body.Name = &name
	common.GetAkeylessPtr(&body.K8sUrl, k8sUrl)
	common.GetAkeylessPtr(&body.K8sToken, k8sToken)
	common.GetAkeylessPtr(&body.K8sUsername, k8sUsername)
	common.GetAkeylessPtr(&body.K8sPassword, k8sPassword)
	if len(k8sCaCertificate) > 0 {
		k8sCaCertificateInt32 := make([]int32, len(k8sCaCertificate))
		for i, v := range k8sCaCertificate {
			k8sCaCertificateInt32[i] = int32(v.(int))
		}
		body.K8sCaCertificate = k8sCaCertificateInt32
	}
	if len(k8sClientCertificate) > 0 {
		k8sClientCertificateInt32 := make([]int32, len(k8sClientCertificate))
		for i, v := range k8sClientCertificate {
			k8sClientCertificateInt32[i] = int32(v.(int))
		}
		body.K8sClientCertificate = k8sClientCertificateInt32
	}
	if len(k8sClientKey) > 0 {
		k8sClientKeyInt32 := make([]int32, len(k8sClientKey))
		for i, v := range k8sClientKey {
			k8sClientKeyInt32[i] = int32(v.(int))
		}
		body.K8sClientKey = k8sClientKeyInt32
	}
	common.GetAkeylessPtr(&body.K8sNamespace, k8sNamespace)
	common.GetAkeylessPtr(&body.K8sSkipSystem, k8sSkipSystem)
	common.GetAkeylessPtr(&body.ProtectionKey, protectionKey)

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationK8sRead(d, m)
		if err != nil {
			return err
		}
	}
	id = d.Get("id").(string)
	body.Id = &id

	_, _, err := client.GatewayUpdateMigration(ctx).Body(*body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update Gateway Migration K8s: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update Gateway Migration K8s: %v", err)
	}

	return nil
}

func resourceGatewayMigrationK8sDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	id := d.Get("id").(string)
	if id == "" {
		err := resourceGatewayMigrationK8sRead(d, m)
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

func resourceGatewayMigrationK8sImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGatewayMigrationK8sRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
