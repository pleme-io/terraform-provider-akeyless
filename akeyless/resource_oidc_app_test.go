package akeyless

import (
	"context"
	"fmt"
	"testing"

	akeyless_api "github.com/akeylesslabs/akeyless-go/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOidcAppResource(t *testing.T) {
	t.Parallel()

	appName := "test_oidc_app"
	appPath := testPath(appName)

	config := fmt.Sprintf(`
		resource "akeyless_oidc_app" "%v" {
			name 				= "%v"
			description 		= "Test OIDC app"
			redirect_uris 		= ["https://localhost/callback"]
			delete_protection 	= "true"
		}
	`, appName, appPath)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_oidc_app" "%v" {
			name 				= "%v"
			description 		= "Updated OIDC app"
			redirect_uris 		= ["https://localhost/callback", "https://localhost/callback2"]
			delete_protection 	= "false"
		}
	`, appName, appPath)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkOidcAppExistsRemotely(appPath),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkOidcAppExistsRemotely(appPath),
				),
			},
		},
	})
}

func checkOidcAppExistsRemotely(path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(*providerMeta).client
		token := *testAccProvider.Meta().(*providerMeta).token

		body := akeyless_api.DescribeItem{
			Name:  path,
			Token: &token,
		}

		_, _, err := client.DescribeItem(context.Background()).Body(body).Execute()
		if err != nil {
			return err
		}

		return nil
	}
}
