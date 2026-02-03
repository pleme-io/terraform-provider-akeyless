// generated fule
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

func resourceProducerRdp() *schema.Resource {
	return &schema.Resource{
		Description:        "RDP Producer resource",
		DeprecationMessage: "Deprecated: Please use new resource: akeyless_dynamic_secret_rdp",
		Create:             resourceProducerRdpCreate,
		Read:               resourceProducerRdpRead,
		Update:             resourceProducerRdpUpdate,
		Delete:             resourceProducerRdpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProducerRdpImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Producer name",
				ForceNew:    true,
			},
			"target_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Name of existing target to use in producer creation",
			},
			"rdp_user_groups": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RDP UserGroup name(s). Multiple values should be separated by comma",
			},
			"rdp_host_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RDP Host name",
			},
			"rdp_admin_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RDP Admin name",
			},
			"rdp_admin_pwd": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RDP Admin Password",
			},
			"rdp_host_port": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RDP Host port",
				Default:     "22",
			},
			"fixed_user_only": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Enable fixed user only",
				Default:     "false",
			},
			"fixed_user_claim_keyname": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "For externally provided users, denotes the key-name of IdP claim to extract the username from (relevant only for fixed-user-only=true)",
			},
			"producer_encryption_key_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Encrypt producer with following key",
			},
			"user_ttl": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "User TTL",
				Default:     "60m",
			},
			"password_length": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The length of the password to be generated",
			},
			"custom_username_template": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Customize how temporary usernames are generated using go template",
			},
			"allow_user_extend_session": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "Allow user to extend session",
			},
			"warn_user_before_expiration": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "Warn user before expiration",
			},
			"delete_protection": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Protection from accidental deletion of this object [true/false]",
			},
			"item_custom_fields": {
				Type:        schema.TypeMap,
				Required:    false,
				Optional:    true,
				Description: "Additional custom fields to associate with the item",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeSet,
				Required:    false,
				Optional:    true,
				Description: "List of the tags attached to this secret. To specify multiple tags use argument multiple times: -t Tag1 -t Tag2",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"secure_access_enable": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Enable/Disable secure remote access, [true/false]",
			},
			"secure_access_rdp_domain": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Required when the Dynamic Secret is used for a domain user",
			},
			"secure_access_rdp_user": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Override the RDP Domain username",
			},
			"secure_access_host": {
				Type:        schema.TypeSet,
				Required:    false,
				Optional:    true,
				Description: "Target servers for connections., For multiple values repeat this flag.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"secure_access_allow_external_user": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Allow providing external user for a domain users",
				Default:     "false",
			},
			"secure_access_bastion_issuer": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Deprecated. use secure-access-certificate-issuer",
			},
			"secure_access_certificate_issuer": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Path to the SSH Certificate Issuer for your Akeyless Secure Access",
			},
			"secure_access_delay": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The delay duration, in seconds, to wait after generating just-in-time credentials. Accepted range: 0-120 seconds",
			},
			"secure_access_rd_gateway_server": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "RD Gateway server",
			},
			"secure_access_web": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Enable Web Secure Remote Access ",
				Computed:    true,
			},
		},
	}
}

func resourceProducerRdpCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetName := d.Get("target_name").(string)
	rdpUserGroups := d.Get("rdp_user_groups").(string)
	rdpHostName := d.Get("rdp_host_name").(string)
	rdpAdminName := d.Get("rdp_admin_name").(string)
	rdpAdminPwd := d.Get("rdp_admin_pwd").(string)
	rdpHostPort := d.Get("rdp_host_port").(string)
	fixedUserOnly := d.Get("fixed_user_only").(string)
	fixedUserClaimKeyname := d.Get("fixed_user_claim_keyname").(string)
	producerEncryptionKeyName := d.Get("producer_encryption_key_name").(string)
	userTtl := d.Get("user_ttl").(string)
	passwordLength := d.Get("password_length").(string)
	customUsernameTemplate := d.Get("custom_username_template").(string)
	allowUserExtendSession := d.Get("allow_user_extend_session").(int)
	warnUserBeforeExpiration := d.Get("warn_user_before_expiration").(int)
	deleteProtection := d.Get("delete_protection").(string)
	itemCustomFields := d.Get("item_custom_fields").(map[string]interface{})
	tagsSet := d.Get("tags").(*schema.Set)
	tags := common.ExpandStringList(tagsSet.List())
	secureAccessEnable := d.Get("secure_access_enable").(string)
	secureAccessRdpDomain := d.Get("secure_access_rdp_domain").(string)
	secureAccessRdpUser := d.Get("secure_access_rdp_user").(string)
	secureAccessHostSet := d.Get("secure_access_host").(*schema.Set)
	secureAccessHost := common.ExpandStringList(secureAccessHostSet.List())
	secureAccessAllowExternalUser := d.Get("secure_access_allow_external_user").(bool)
	secureAccessBastionIssuer := d.Get("secure_access_bastion_issuer").(string)
	secureAccessCertificateIssuer := d.Get("secure_access_certificate_issuer").(string)
	secureAccessDelay := d.Get("secure_access_delay").(int)
	secureAccessRdGatewayServer := d.Get("secure_access_rd_gateway_server").(string)

	body := akeyless_api.GatewayCreateProducerRdp{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.TargetName, targetName)
	common.GetAkeylessPtr(&body.RdpUserGroups, rdpUserGroups)
	common.GetAkeylessPtr(&body.RdpHostName, rdpHostName)
	common.GetAkeylessPtr(&body.RdpAdminName, rdpAdminName)
	common.GetAkeylessPtr(&body.RdpAdminPwd, rdpAdminPwd)
	common.GetAkeylessPtr(&body.RdpHostPort, rdpHostPort)
	common.GetAkeylessPtr(&body.FixedUserOnly, fixedUserOnly)
	common.GetAkeylessPtr(&body.FixedUserClaimKeyname, fixedUserClaimKeyname)
	common.GetAkeylessPtr(&body.ProducerEncryptionKeyName, producerEncryptionKeyName)
	common.GetAkeylessPtr(&body.UserTtl, userTtl)
	common.GetAkeylessPtr(&body.PasswordLength, passwordLength)
	common.GetAkeylessPtr(&body.CustomUsernameTemplate, customUsernameTemplate)
	if allowUserExtendSession != 0 {
		body.AllowUserExtendSession = &[]int64{int64(allowUserExtendSession)}[0]
	}
	if warnUserBeforeExpiration != 0 {
		body.WarnUserBeforeExpiration = &[]int64{int64(warnUserBeforeExpiration)}[0]
	}
	common.GetAkeylessPtr(&body.DeleteProtection, deleteProtection)
	if len(itemCustomFields) > 0 {
		fields := make(map[string]string)
		for k, v := range itemCustomFields {
			fields[k] = v.(string)
		}
		body.ItemCustomFields = &fields
	}
	common.GetAkeylessPtr(&body.Tags, tags)
	common.GetAkeylessPtr(&body.SecureAccessEnable, secureAccessEnable)
	common.GetAkeylessPtr(&body.SecureAccessRdpDomain, secureAccessRdpDomain)
	common.GetAkeylessPtr(&body.SecureAccessRdpUser, secureAccessRdpUser)
	common.GetAkeylessPtr(&body.SecureAccessHost, secureAccessHost)
	common.GetAkeylessPtr(&body.SecureAccessAllowExternalUser, secureAccessAllowExternalUser)
	common.GetAkeylessPtr(&body.SecureAccessBastionIssuer, secureAccessBastionIssuer)
	common.GetAkeylessPtr(&body.SecureAccessCertificateIssuer, secureAccessCertificateIssuer)
	if secureAccessDelay != 0 {
		body.SecureAccessDelay = &[]int64{int64(secureAccessDelay)}[0]
	}
	common.GetAkeylessPtr(&body.SecureAccessRdGatewayServer, secureAccessRdGatewayServer)

	_, _, err := client.GatewayCreateProducerRdp(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Secret: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Secret: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceProducerRdpRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()

	path := d.Id()
	body := akeyless_api.GatewayGetProducer{
		Name:  path,
		Token: &token,
	}

	rOut, res, err := client.GatewayGetProducer(ctx).Body(body).Execute()
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
	if rOut.FixedUserOnly != nil {
		err = d.Set("fixed_user_only", *rOut.FixedUserOnly)
		if err != nil {
			return err
		}
	}
	if rOut.UserTtl != nil {
		err = d.Set("user_ttl", *rOut.UserTtl)
		if err != nil {
			return err
		}
	}
	if rOut.Tags != nil {
		err = d.Set("tags", rOut.Tags)
		if err != nil {
			return err
		}
	}

	if rOut.ItemTargetsAssoc != nil {
		targetName := common.GetTargetName(rOut.ItemTargetsAssoc)
		err = d.Set("target_name", targetName)
		if err != nil {
			return err
		}
	}
	if rOut.Groups != nil {
		err = d.Set("rdp_user_groups", *rOut.Groups)
		if err != nil {
			return err
		}
	}
	if rOut.HostName != nil {
		err = d.Set("rdp_host_name", *rOut.HostName)
		if err != nil {
			return err
		}
	}
	if rOut.AdminName != nil {
		err = d.Set("rdp_admin_name", *rOut.AdminName)
		if err != nil {
			return err
		}
	}
	if rOut.AdminPwd != nil {
		err = d.Set("rdp_admin_pwd", *rOut.AdminPwd)
		if err != nil {
			return err
		}
	}
	if rOut.HostPort != nil {
		err = d.Set("rdp_host_port", *rOut.HostPort)
		if err != nil {
			return err
		}
	}
	if rOut.DynamicSecretKey != nil {
		err = d.Set("producer_encryption_key_name", *rOut.DynamicSecretKey)
		if err != nil {
			return err
		}
	}

	common.GetSra(d, rOut.SecureRemoteAccessDetails, "DYNAMIC_SECERT")

	d.SetId(path)

	return nil
}

func resourceProducerRdpUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetName := d.Get("target_name").(string)
	rdpUserGroups := d.Get("rdp_user_groups").(string)
	rdpHostName := d.Get("rdp_host_name").(string)
	rdpAdminName := d.Get("rdp_admin_name").(string)
	rdpAdminPwd := d.Get("rdp_admin_pwd").(string)
	rdpHostPort := d.Get("rdp_host_port").(string)
	fixedUserOnly := d.Get("fixed_user_only").(string)
	fixedUserClaimKeyname := d.Get("fixed_user_claim_keyname").(string)
	producerEncryptionKeyName := d.Get("producer_encryption_key_name").(string)
	userTtl := d.Get("user_ttl").(string)
	passwordLength := d.Get("password_length").(string)
	customUsernameTemplate := d.Get("custom_username_template").(string)
	allowUserExtendSession := d.Get("allow_user_extend_session").(int)
	warnUserBeforeExpiration := d.Get("warn_user_before_expiration").(int)
	deleteProtection := d.Get("delete_protection").(string)
	itemCustomFields := d.Get("item_custom_fields").(map[string]interface{})
	tagsSet := d.Get("tags").(*schema.Set)
	tags := common.ExpandStringList(tagsSet.List())
	secureAccessEnable := d.Get("secure_access_enable").(string)
	secureAccessRdpDomain := d.Get("secure_access_rdp_domain").(string)
	secureAccessRdpUser := d.Get("secure_access_rdp_user").(string)
	secureAccessHostSet := d.Get("secure_access_host").(*schema.Set)
	secureAccessHost := common.ExpandStringList(secureAccessHostSet.List())
	secureAccessAllowExternalUser := d.Get("secure_access_allow_external_user").(bool)
	secureAccessBastionIssuer := d.Get("secure_access_bastion_issuer").(string)
	secureAccessCertificateIssuer := d.Get("secure_access_certificate_issuer").(string)
	secureAccessDelay := d.Get("secure_access_delay").(int)
	secureAccessRdGatewayServer := d.Get("secure_access_rd_gateway_server").(string)

	body := akeyless_api.GatewayUpdateProducerRdp{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.TargetName, targetName)
	common.GetAkeylessPtr(&body.RdpUserGroups, rdpUserGroups)
	common.GetAkeylessPtr(&body.RdpHostName, rdpHostName)
	common.GetAkeylessPtr(&body.RdpAdminName, rdpAdminName)
	common.GetAkeylessPtr(&body.RdpAdminPwd, rdpAdminPwd)
	common.GetAkeylessPtr(&body.RdpHostPort, rdpHostPort)
	common.GetAkeylessPtr(&body.FixedUserOnly, fixedUserOnly)
	common.GetAkeylessPtr(&body.FixedUserClaimKeyname, fixedUserClaimKeyname)
	common.GetAkeylessPtr(&body.ProducerEncryptionKeyName, producerEncryptionKeyName)
	common.GetAkeylessPtr(&body.UserTtl, userTtl)
	common.GetAkeylessPtr(&body.PasswordLength, passwordLength)
	common.GetAkeylessPtr(&body.CustomUsernameTemplate, customUsernameTemplate)
	if allowUserExtendSession != 0 {
		body.AllowUserExtendSession = &[]int64{int64(allowUserExtendSession)}[0]
	}
	if warnUserBeforeExpiration != 0 {
		body.WarnUserBeforeExpiration = &[]int64{int64(warnUserBeforeExpiration)}[0]
	}
	common.GetAkeylessPtr(&body.DeleteProtection, deleteProtection)
	if len(itemCustomFields) > 0 {
		fields := make(map[string]string)
		for k, v := range itemCustomFields {
			fields[k] = v.(string)
		}
		body.ItemCustomFields = &fields
	}
	common.GetAkeylessPtr(&body.Tags, tags)
	common.GetAkeylessPtr(&body.SecureAccessEnable, secureAccessEnable)
	common.GetAkeylessPtr(&body.SecureAccessRdpDomain, secureAccessRdpDomain)
	common.GetAkeylessPtr(&body.SecureAccessRdpUser, secureAccessRdpUser)
	common.GetAkeylessPtr(&body.SecureAccessHost, secureAccessHost)
	common.GetAkeylessPtr(&body.SecureAccessAllowExternalUser, secureAccessAllowExternalUser)
	common.GetAkeylessPtr(&body.SecureAccessBastionIssuer, secureAccessBastionIssuer)
	common.GetAkeylessPtr(&body.SecureAccessCertificateIssuer, secureAccessCertificateIssuer)
	if secureAccessDelay != 0 {
		body.SecureAccessDelay = &[]int64{int64(secureAccessDelay)}[0]
	}
	common.GetAkeylessPtr(&body.SecureAccessRdGatewayServer, secureAccessRdGatewayServer)

	_, _, err := client.GatewayUpdateProducerRdp(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update : %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update : %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceProducerRdpDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	path := d.Id()

	deleteItem := akeyless_api.GatewayDeleteProducer{
		Token: &token,
		Name:  path,
	}

	ctx := context.Background()
	_, _, err := client.GatewayDeleteProducer(ctx).Body(deleteItem).Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceProducerRdpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceProducerRdpRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
