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

func resourceProducerMssql() *schema.Resource {
	return &schema.Resource{
		Description:        "Microsoft SQL Server producer resource",
		DeprecationMessage: "Deprecated: Please use new resource: akeyless_dynamic_secret_mssql",
		Create:             resourceProducerMssqlCreate,
		Read:               resourceProducerMssqlRead,
		Update:             resourceProducerMssqlUpdate,
		Delete:             resourceProducerMssqlDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProducerMssqlImport,
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
			"custom_username_template": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Customize how temporary usernames are generated using go template",
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
			"mssql_allowed_db_names": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "CSV of allowed DB names for runtime selection when getting the secret value. Empty => use target DB only; \"*\" => any DB allowed; One or more names => user must choose from this list",
			},
			"password_length": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The length of the password to be generated",
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
			"mssql_dbname": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MSSQL Server DB Name",
			},
			"mssql_username": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MS SQL Server user",
			},
			"mssql_password": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MS SQL Server password",
			},
			"mssql_host": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MS SQL Server host name",
				Default:     "127.0.0.1",
			},
			"mssql_port": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MS SQL Server port",
				Default:     "1433",
			},
			"mssql_create_statements": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MSSQL Server Creation Statements",
				Default:     "CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';",
			},
			"mssql_revocation_statements": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "MSSQL Server Revocation Statements",
				Default:     "DROP LOGIN [{{name}}];",
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
			"secure_access_bastion_issuer": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Path to the SSH Certificate Issuer for your Akeyless Bastion",
			},
			"secure_access_host": {
				Type:        schema.TypeSet,
				Required:    false,
				Optional:    true,
				Description: "Target DB servers for connections., For multiple values repeat this flag.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"secure_access_db_schema": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The db schema",
			},
			"secure_access_web": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Enable Web Secure Remote Access ",
				Default:     "false",
			},
			"secure_access_db_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Enable Web Secure Remote Access ",
				Computed:    true,
			},
		},
	}
}

func resourceProducerMssqlCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetName := d.Get("target_name").(string)
	customUsernameTemplate := d.Get("custom_username_template").(string)
	deleteProtection := d.Get("delete_protection").(string)
	itemCustomFields := d.Get("item_custom_fields").(map[string]interface{})
	mssqlAllowedDbNames := d.Get("mssql_allowed_db_names").(string)
	mssqlDbname := d.Get("mssql_dbname").(string)
	mssqlUsername := d.Get("mssql_username").(string)
	mssqlPassword := d.Get("mssql_password").(string)
	mssqlHost := d.Get("mssql_host").(string)
	mssqlPort := d.Get("mssql_port").(string)
	mssqlCreateStatements := d.Get("mssql_create_statements").(string)
	mssqlRevocationStatements := d.Get("mssql_revocation_statements").(string)
	passwordLength := d.Get("password_length").(string)
	producerEncryptionKeyName := d.Get("producer_encryption_key_name").(string)
	userTtl := d.Get("user_ttl").(string)
	tagsSet := d.Get("tags").(*schema.Set)
	tags := common.ExpandStringList(tagsSet.List())
	secureAccessEnable := d.Get("secure_access_enable").(string)
	secureAccessBastionIssuer := d.Get("secure_access_bastion_issuer").(string)
	secureAccessCertificateIssuer := d.Get("secure_access_certificate_issuer").(string)
	secureAccessDelay := d.Get("secure_access_delay").(int)
	secureAccessHostSet := d.Get("secure_access_host").(*schema.Set)
	secureAccessHost := common.ExpandStringList(secureAccessHostSet.List())
	secureAccessDbSchema := d.Get("secure_access_db_schema").(string)
	secureAccessWeb := d.Get("secure_access_web").(bool)

	body := akeyless_api.GatewayCreateProducerMSSQL{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.TargetName, targetName)
	common.GetAkeylessPtr(&body.CustomUsernameTemplate, customUsernameTemplate)
	common.GetAkeylessPtr(&body.DeleteProtection, deleteProtection)
	if len(itemCustomFields) > 0 {
		fields := make(map[string]string)
		for k, v := range itemCustomFields {
			fields[k] = v.(string)
		}
		common.GetAkeylessPtr(&body.ItemCustomFields, fields)
	}
	common.GetAkeylessPtr(&body.MssqlAllowedDbNames, mssqlAllowedDbNames)
	common.GetAkeylessPtr(&body.MssqlDbname, mssqlDbname)
	common.GetAkeylessPtr(&body.MssqlUsername, mssqlUsername)
	common.GetAkeylessPtr(&body.MssqlPassword, mssqlPassword)
	common.GetAkeylessPtr(&body.MssqlHost, mssqlHost)
	common.GetAkeylessPtr(&body.MssqlPort, mssqlPort)
	common.GetAkeylessPtr(&body.MssqlCreateStatements, mssqlCreateStatements)
	common.GetAkeylessPtr(&body.MssqlRevocationStatements, mssqlRevocationStatements)
	common.GetAkeylessPtr(&body.PasswordLength, passwordLength)
	common.GetAkeylessPtr(&body.ProducerEncryptionKeyName, producerEncryptionKeyName)
	common.GetAkeylessPtr(&body.UserTtl, userTtl)
	common.GetAkeylessPtr(&body.Tags, tags)
	common.GetAkeylessPtr(&body.SecureAccessEnable, secureAccessEnable)
	common.GetAkeylessPtr(&body.SecureAccessBastionIssuer, secureAccessBastionIssuer)
	common.GetAkeylessPtr(&body.SecureAccessCertificateIssuer, secureAccessCertificateIssuer)
	if secureAccessDelay > 0 {
		delay := int64(secureAccessDelay)
		common.GetAkeylessPtr(&body.SecureAccessDelay, delay)
	}
	common.GetAkeylessPtr(&body.SecureAccessHost, secureAccessHost)
	common.GetAkeylessPtr(&body.SecureAccessDbSchema, secureAccessDbSchema)
	common.GetAkeylessPtr(&body.SecureAccessWeb, secureAccessWeb)

	_, _, err := client.GatewayCreateProducerMSSQL(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create Secret: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create Secret: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceProducerMssqlRead(d *schema.ResourceData, m interface{}) error {
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
	if rOut.MssqlRevocationStatements != nil {
		err = d.Set("mssql_revocation_statements", *rOut.MssqlRevocationStatements)
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
	if rOut.DbName != nil {
		err = d.Set("mssql_dbname", *rOut.DbName)
		if err != nil {
			return err
		}
	}
	if rOut.DbUserName != nil {
		err = d.Set("mssql_username", *rOut.DbUserName)
		if err != nil {
			return err
		}
	}
	if rOut.DbPwd != nil {
		err = d.Set("mssql_password", *rOut.DbPwd)
		if err != nil {
			return err
		}
	}
	if rOut.DbHostName != nil {
		err = d.Set("mssql_host", *rOut.DbHostName)
		if err != nil {
			return err
		}
	}
	if rOut.DbPort != nil {
		err = d.Set("mssql_port", *rOut.DbPort)
		if err != nil {
			return err
		}
	}
	if rOut.MssqlCreationStatements != nil {
		err = d.Set("mssql_create_statements", *rOut.MssqlCreationStatements)
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

func resourceProducerMssqlUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	targetName := d.Get("target_name").(string)
	customUsernameTemplate := d.Get("custom_username_template").(string)
	deleteProtection := d.Get("delete_protection").(string)
	itemCustomFields := d.Get("item_custom_fields").(map[string]interface{})
	mssqlAllowedDbNames := d.Get("mssql_allowed_db_names").(string)
	mssqlDbname := d.Get("mssql_dbname").(string)
	mssqlUsername := d.Get("mssql_username").(string)
	mssqlPassword := d.Get("mssql_password").(string)
	mssqlHost := d.Get("mssql_host").(string)
	mssqlPort := d.Get("mssql_port").(string)
	mssqlCreateStatements := d.Get("mssql_create_statements").(string)
	mssqlRevocationStatements := d.Get("mssql_revocation_statements").(string)
	passwordLength := d.Get("password_length").(string)
	producerEncryptionKeyName := d.Get("producer_encryption_key_name").(string)
	userTtl := d.Get("user_ttl").(string)
	tagsSet := d.Get("tags").(*schema.Set)
	tags := common.ExpandStringList(tagsSet.List())
	secureAccessEnable := d.Get("secure_access_enable").(string)
	secureAccessBastionIssuer := d.Get("secure_access_bastion_issuer").(string)
	secureAccessCertificateIssuer := d.Get("secure_access_certificate_issuer").(string)
	secureAccessDelay := d.Get("secure_access_delay").(int)
	secureAccessHostSet := d.Get("secure_access_host").(*schema.Set)
	secureAccessHost := common.ExpandStringList(secureAccessHostSet.List())
	secureAccessDbSchema := d.Get("secure_access_db_schema").(string)
	secureAccessWeb := d.Get("secure_access_web").(bool)

	body := akeyless_api.GatewayUpdateProducerMSSQL{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.TargetName, targetName)
	common.GetAkeylessPtr(&body.CustomUsernameTemplate, customUsernameTemplate)
	common.GetAkeylessPtr(&body.DeleteProtection, deleteProtection)
	if len(itemCustomFields) > 0 {
		fields := make(map[string]string)
		for k, v := range itemCustomFields {
			fields[k] = v.(string)
		}
		common.GetAkeylessPtr(&body.ItemCustomFields, fields)
	}
	common.GetAkeylessPtr(&body.MssqlAllowedDbNames, mssqlAllowedDbNames)
	common.GetAkeylessPtr(&body.MssqlDbname, mssqlDbname)
	common.GetAkeylessPtr(&body.MssqlUsername, mssqlUsername)
	common.GetAkeylessPtr(&body.MssqlPassword, mssqlPassword)
	common.GetAkeylessPtr(&body.MssqlHost, mssqlHost)
	common.GetAkeylessPtr(&body.MssqlPort, mssqlPort)
	common.GetAkeylessPtr(&body.MssqlCreateStatements, mssqlCreateStatements)
	common.GetAkeylessPtr(&body.MssqlRevocationStatements, mssqlRevocationStatements)
	common.GetAkeylessPtr(&body.PasswordLength, passwordLength)
	common.GetAkeylessPtr(&body.ProducerEncryptionKeyName, producerEncryptionKeyName)
	common.GetAkeylessPtr(&body.UserTtl, userTtl)
	common.GetAkeylessPtr(&body.Tags, tags)
	common.GetAkeylessPtr(&body.SecureAccessEnable, secureAccessEnable)
	common.GetAkeylessPtr(&body.SecureAccessBastionIssuer, secureAccessBastionIssuer)
	common.GetAkeylessPtr(&body.SecureAccessCertificateIssuer, secureAccessCertificateIssuer)
	if secureAccessDelay > 0 {
		delay := int64(secureAccessDelay)
		common.GetAkeylessPtr(&body.SecureAccessDelay, delay)
	}
	common.GetAkeylessPtr(&body.SecureAccessHost, secureAccessHost)
	common.GetAkeylessPtr(&body.SecureAccessDbSchema, secureAccessDbSchema)
	common.GetAkeylessPtr(&body.SecureAccessWeb, secureAccessWeb)

	_, _, err := client.GatewayUpdateProducerMSSQL(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update : %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update : %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceProducerMssqlDelete(d *schema.ResourceData, m interface{}) error {
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

func resourceProducerMssqlImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceProducerMssqlRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
