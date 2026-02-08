package akeyless

import (
	"context"
	"fmt"
	"testing"

	akeyless_api "github.com/akeylesslabs/akeyless-go/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestGroupResource(t *testing.T) {
	t.Parallel()

	groupName := "test_group"
	groupPath := testPath(groupName)

	config := fmt.Sprintf(`
		resource "akeyless_group" "%v" {
			name 				= "%v"
			group_alias 		= "testalias"
			description 		= "Test group description"
		}
	`, groupName, groupPath)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_group" "%v" {
			name 				= "%v"
			group_alias 		= "testaliasupd"
			description 		= "Updated test group description"
		}
	`, groupName, groupPath)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkGroupExistsRemotely(groupPath),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkGroupExistsRemotely(groupPath),
				),
			},
		},
	})
}

func checkGroupExistsRemotely(path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(*providerMeta).client
		token := *testAccProvider.Meta().(*providerMeta).token

		body := akeyless_api.GetGroup{
			Name:  path,
			Token: &token,
		}

		_, _, err := client.GetGroup(context.Background()).Body(body).Execute()
		if err != nil {
			return err
		}

		return nil
	}
}
