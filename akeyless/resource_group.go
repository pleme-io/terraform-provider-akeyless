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

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Group resource",
		Create:      resourceGroupCreate,
		Read:        resourceGroupRead,
		Update:      resourceGroupUpdate,
		Delete:      resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
				ForceNew:    true,
			},
			"group_alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A short group alias",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the object",
			},
			"user_assignment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A json string defining the access permission assignment for this client",
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	groupAlias := d.Get("group_alias").(string)
	description := d.Get("description").(string)
	userAssignment := d.Get("user_assignment").(string)

	body := akeyless_api.CreateGroup{
		Name:       name,
		GroupAlias: groupAlias,
		Token:      &token,
	}
	common.GetAkeylessPtr(&body.Description, description)
	common.GetAkeylessPtr(&body.UserAssignment, userAssignment)

	_, _, err := client.CreateGroup(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("failed to create group: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("failed to create group: %w", err)
	}

	d.SetId(name)

	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()

	path := d.Id()

	body := akeyless_api.GetGroup{
		Name:  path,
		Token: &token,
	}

	rOut, res, err := client.GetGroup(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			if res.StatusCode == http.StatusNotFound {
				// The resource was deleted outside of the current Terraform workspace, so invalidate this resource
				d.SetId("")
				return nil
			}
			return fmt.Errorf("failed to get group: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("failed to get group: %w", err)
	}

	if rOut.GroupAlias != nil {
		err = d.Set("group_alias", *rOut.GroupAlias)
		if err != nil {
			return err
		}
	}

	if rOut.Description != nil {
		err = d.Set("description", *rOut.Description)
		if err != nil {
			return err
		}
	}

	// Note: UserAssignments is an array in the API response but user_assignment
	// in the resource is a JSON string. This would require marshaling the array to JSON.
	// For now, we skip reading it back to avoid complexity.

	d.SetId(path)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	var apiErr akeyless_api.GenericOpenAPIError
	ctx := context.Background()
	name := d.Get("name").(string)
	groupAlias := d.Get("group_alias").(string)
	description := d.Get("description").(string)
	userAssignment := d.Get("user_assignment").(string)

	body := akeyless_api.UpdateGroup{
		Name:       name,
		GroupAlias: groupAlias,
		Token:      &token,
	}
	common.GetAkeylessPtr(&body.Description, description)
	common.GetAkeylessPtr(&body.UserAssignment, userAssignment)

	_, _, err := client.UpdateGroup(ctx).Body(body).Execute()
	if err != nil {
		if errors.As(err, &apiErr) {
			return fmt.Errorf("failed to update group: %v", string(apiErr.Body()))
		}
		return fmt.Errorf("failed to update group: %w", err)
	}

	d.SetId(name)

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*providerMeta)
	client := *provider.client
	token := *provider.token

	path := d.Id()

	deleteItem := akeyless_api.DeleteGroup{
		Token: &token,
		Name:  path,
	}

	ctx := context.Background()
	_, _, err := client.DeleteGroup(ctx).Body(deleteItem).Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	err := resourceGroupRead(d, m)
	if err != nil {
		return nil, err
	}

	err = d.Set("name", id)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
