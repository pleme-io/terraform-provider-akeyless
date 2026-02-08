package akeyless

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	akeyless_api "github.com/akeylesslabs/akeyless-go/v5"
	"github.com/akeylesslabs/terraform-provider-akeyless/akeyless/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOidcApp() *schema.Resource {
	return &schema.Resource{
		Description: "OIDC App resource",
		Create:      resourceOidcAppCreate,
		Read:        resourceOidcAppRead,
		Update:      resourceOidcAppUpdate,
		Delete:      resourceOidcAppDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOidcAppImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OIDC App name",
				ForceNew:    true,
			},
			"accessibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For personal password manager",
				Default:     "regular",
			},
			"audience": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A comma separated list of allowed audiences",
			},
			"delete_protection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protection from accidental deletion of this object [true/false]",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the object",
			},
			"item_custom_fields": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional custom fields to associate with the item",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a key that used to encrypt the OIDC application (if empty, the account default protectionKey key will be used)",
			},
			"metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Deprecated - use description",
			},
			"permission_assignment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A json string defining the access permission assignment for this app",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true if the app is public (cannot keep secrets)",
			},
			"redirect_uris": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A comma separated list of allowed redirect uris",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scopes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A comma separated list of allowed scopes",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Add tags attached to this object",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOidcAppCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	accessibility := d.Get("accessibility").(string)
	audience := d.Get("audience").(string)
	deleteProtection := d.Get("delete_protection").(string)
	description := d.Get("description").(string)
	itemCustomFields := d.Get("item_custom_fields").(map[string]interface{})
	key := d.Get("key").(string)
	metadata := d.Get("metadata").(string)
	permissionAssignment := d.Get("permission_assignment").(string)
	public := d.Get("public").(bool)
	redirectUrisSet := d.Get("redirect_uris").(*schema.Set)
	redirectUris := common.ExpandStringList(redirectUrisSet.List())
	scopesSet := d.Get("scopes").(*schema.Set)
	scopes := common.ExpandStringList(scopesSet.List())
	tagsSet := d.Get("tags").(*schema.Set)
	tags := common.ExpandStringList(tagsSet.List())

	body := akeyless_api.CreateOidcApp{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.Accessibility, accessibility)
	common.GetAkeylessPtr(&body.Audience, audience)
	common.GetAkeylessPtr(&body.DeleteProtection, deleteProtection)
	common.GetAkeylessPtr(&body.Description, description)
	if len(itemCustomFields) > 0 {
		customFieldsMap := make(map[string]string)
		for k, v := range itemCustomFields {
			customFieldsMap[k] = v.(string)
		}
		body.ItemCustomFields = &customFieldsMap
	}
	common.GetAkeylessPtr(&body.Key, key)
	common.GetAkeylessPtr(&body.Metadata, metadata)
	common.GetAkeylessPtr(&body.PermissionAssignment, permissionAssignment)
	common.GetAkeylessPtr(&body.Public, public)
	common.GetAkeylessPtr(&body.RedirectUris, redirectUris)
	common.GetAkeylessPtr(&body.Scopes, scopes)
	if len(tags) > 0 {
		body.Tags = tags
	}

	_, _, err := client.CreateOidcApp(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't create OIDC App: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't create OIDC App: %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceOidcAppRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()

	path := d.Id()

	body := akeyless_api.DescribeItem{
		Name:  path,
		Token: &token,
	}

	rOut, res, err := client.DescribeItem(ctx).Body(body).Execute()
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

	if rOut.ItemMetadata != nil {
		err = d.Set("description", *rOut.ItemMetadata)
		if err != nil {
			return err
		}
	}
	if rOut.ItemTags != nil {
		err = d.Set("tags", rOut.ItemTags)
		if err != nil {
			return err
		}
	}
	if rOut.ProtectionKeyName != nil {
		err = d.Set("key", *rOut.ProtectionKeyName)
		if err != nil {
			return err
		}
	}
	if rOut.DeleteProtection != nil {
		err = d.Set("delete_protection", strconv.FormatBool(*rOut.DeleteProtection))
		if err != nil {
			return err
		}
	}

	d.SetId(path)

	return nil
}

func resourceOidcAppUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	audience := d.Get("audience").(string)
	key := d.Get("key").(string)
	permissionAssignment := d.Get("permission_assignment").(string)
	public := d.Get("public").(bool)
	redirectUrisSet := d.Get("redirect_uris").(*schema.Set)
	redirectUris := common.ExpandStringList(redirectUrisSet.List())
	scopesSet := d.Get("scopes").(*schema.Set)
	scopes := common.ExpandStringList(scopesSet.List())

	body := akeyless_api.UpdateOidcApp{
		Name:  name,
		Token: &token,
	}
	common.GetAkeylessPtr(&body.Audience, audience)
	common.GetAkeylessPtr(&body.Key, key)
	common.GetAkeylessPtr(&body.PermissionAssignment, permissionAssignment)
	common.GetAkeylessPtr(&body.Public, public)
	if len(redirectUris) > 0 {
		redirectUrisStr := strings.Join(redirectUris, ",")
		common.GetAkeylessPtr(&body.RedirectUris, redirectUrisStr)
	}
	if len(scopes) > 0 {
		scopesStr := strings.Join(scopes, ",")
		common.GetAkeylessPtr(&body.Scopes, scopesStr)
	}

	_, _, err := client.UpdateOidcApp(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("can't update : %v", string(apiErr.Body()))
		}
		return fmt.Errorf("can't update : %v", err)
	}

	d.SetId(name)

	return nil
}

func resourceOidcAppDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	path := d.Id()

	deleteItem := akeyless_api.DeleteItem{
		Token: &token,
		Name:  path,
	}

	ctx := context.Background()
	_, _, err := client.DeleteItem(ctx).Body(deleteItem).Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceOidcAppImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceOidcAppRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
