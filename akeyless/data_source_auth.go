package akeyless

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAuth() *schema.Resource {
	return &schema.Resource{
		Description: "Authenticate to Akeyless and returns a token to be used by the provider",
		Read:        dataSourceAuthRead,
		Schema: map[string]*schema.Schema{
			"api_key_login":  apiKeyLoginSchema,
			"aws_iam_login":  awsIamLoginSchema,
			"gcp_login":      gcpLoginSchema,
			"azure_ad_login": azureAdLoginSchema,
			"jwt_login":      jwtLoginSchema,
			"email_login":    emailLoginSchema,
			"uid_login":      uidLoginSchema,
			"cert_login":     certLoginSchema,
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The token",
			},
			"complete_auth_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Complete authentication link",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Token expiration time",
			},
			"creds": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System access credentials",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ID",
						},
						"auth_creds": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Temporary credentials for accessing Auth",
						},
						"expiry": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Credentials expiration date",
						},
						"kfm_creds": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Temporary credentials for accessing the KFMs instances",
						},
						"need_mfa_app_first_config": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If the user didn't complete to configure the MFA app",
						},
						"required_mfa": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Required MFA",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Credentials tmp token",
						},
						"uam_creds": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Temporary credentials for accessing the UAM service",
						},
					},
				},
			},
		},
	}
}

func dataSourceAuthRead(d *schema.ResourceData, m interface{}) error {

	provider := m.(*providerMeta)
	client := *provider.client

	ctx := context.Background()

	authBody, err := getAuthInfo(d)
	if err != nil {
		return err
	}

	authOut, _, err := client.Auth(ctx).Body(*authBody).Execute()
	if err != nil {
		return err
	}

	token := authOut.GetToken()
	err = d.Set("token", token)
	if err != nil {
		return err
	}

	if authOut.CompleteAuthLink != nil {
		err = d.Set("complete_auth_link", authOut.GetCompleteAuthLink())
		if err != nil {
			return err
		}
	}

	if authOut.Expiration != nil {
		err = d.Set("expiration", authOut.GetExpiration())
		if err != nil {
			return err
		}
	}

	if authOut.Creds != nil {
		creds := authOut.GetCreds()
		credsMap := make([]map[string]interface{}, 0, 1)
		credItem := make(map[string]interface{})

		if creds.AccessId != nil {
			credItem["access_id"] = creds.GetAccessId()
		}
		if creds.AuthCreds != nil {
			credItem["auth_creds"] = creds.GetAuthCreds()
		}
		if creds.Expiry != nil {
			credItem["expiry"] = int(creds.GetExpiry())
		}
		if creds.KfmCreds != nil {
			credItem["kfm_creds"] = creds.GetKfmCreds()
		}
		if creds.NeedMfaAppFirstConfig != nil {
			credItem["need_mfa_app_first_config"] = creds.GetNeedMfaAppFirstConfig()
		}
		if creds.RequiredMfa != nil {
			credItem["required_mfa"] = creds.GetRequiredMfa()
		}
		if creds.Token != nil {
			credItem["token"] = creds.GetToken()
		}
		if creds.UamCreds != nil {
			credItem["uam_creds"] = creds.GetUamCreds()
		}

		if len(credItem) > 0 {
			credsMap = append(credsMap, credItem)
			err = d.Set("creds", credsMap)
			if err != nil {
				return err
			}
		}
	}

	provider.token = &token

	d.SetId("dummy_id")
	return nil
}
