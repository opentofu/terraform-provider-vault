// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-vault/testutil"
)

// TestAccDataSourceAzureAccessCredentials_basic tests the creation of dynamic
// service principals using azure_roles
func TestAccDataSourceAzureAccessCredentialsAzureRoles_basic(t *testing.T) {
	// This test takes a while because it's testing a loop that
	// retries real credentials until they're eventually consistent.
	if testing.Short() {
		t.SkipNow()
	}
	mountPath := acctest.RandomWithPrefix("tf-test-azure")
	conf := testutil.GetTestAzureConf(t)
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureAccessCredentialsConfigBasicAzureRoles(mountPath, conf, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_secret"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "lease_id"),
				),
			},
		},
	})
}

// TestAccDataSourceAzureAccessCredentialsAzureGroups_basic tests the creation of dynamic
// service principals using azure_groups
// Requires AZURE_GROUP_NAME to be set to a group that the service principal will be assigned to
func TestAccDataSourceAzureAccessCredentialsAzureGroups_basic(t *testing.T) {
	// This test takes a while because it's testing a loop that
	// retries real credentials until they're eventually consistent.
	if testing.Short() {
		t.SkipNow()
	}

	groupName := os.Getenv("AZURE_GROUP_NAME")
	if groupName == "" {
		t.Skip("AZURE_GROUP_NAME must be set to run this test")
	}

	mountPath := acctest.RandomWithPrefix("tf-test-azure")
	conf := testutil.GetTestAzureConfForGroups(t)
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureAccessCredentialsConfigBasicAzureGroups(mountPath, conf, groupName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_secret"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "lease_id"),
				),
			},
		},
	})
}

// TestAccDataSourceAzureAccessCredentials_ExistingSP tests the credential
// generation for existing service principals
func TestAccDataSourceAzureAccessCredentials_ExistingSP(t *testing.T) {
	// This test takes a while because it's testing a loop that
	// retries real credentials until they're eventually consistent.
	if testing.Short() {
		t.SkipNow()
	}
	mountPath := acctest.RandomWithPrefix("tf-test-azure")
	conf := testutil.GetTestAzureConfExistingSP(t)
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureAccessCredentialsConfig_existingSP(mountPath, conf, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "client_secret"),
					resource.TestCheckResourceAttrSet("data.vault_azure_access_credentials.test", "lease_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureAccessCredentialsConfig_existingSP(mountPath string, conf *testutil.AzureTestConf, maxSecs int) string {
	template := `
resource "vault_azure_secret_backend" "test" {
	path = "{{mountPath}}"
	subscription_id = "{{subscriptionID}}"
	tenant_id = "{{tenantID}}"
	client_id = "{{clientID}}"
	client_secret = "{{clientSecret}}"
}

resource "vault_azure_secret_backend_role" "test" {
	backend = vault_azure_secret_backend.test.path
	role = "my-role"
	application_object_id = "{{appObjectID}}"
	ttl = 300
	max_ttl = 600
}

data "vault_azure_access_credentials" "test" {
    backend = vault_azure_secret_backend.test.path
    role = vault_azure_secret_backend_role.test.role
    validate_creds = true
	num_seconds_between_tests = 1
	max_cred_validation_seconds = {{maxCredValidationSeconds}}
}`

	parsed := strings.Replace(template, "{{mountPath}}", mountPath, -1)
	parsed = strings.Replace(parsed, "{{subscriptionID}}", conf.SubscriptionID, -1)
	parsed = strings.Replace(parsed, "{{tenantID}}", conf.TenantID, -1)
	parsed = strings.Replace(parsed, "{{clientID}}", conf.ClientID, -1)
	parsed = strings.Replace(parsed, "{{clientSecret}}", conf.ClientSecret, -1)
	parsed = strings.Replace(parsed, "{{appObjectID}}", conf.AppObjectID, -1)
	parsed = strings.Replace(parsed, "{{maxCredValidationSeconds}}", strconv.Itoa(maxSecs), -1)
	return parsed
}

func testAccDataSourceAzureAccessCredentialsConfigBasicAzureRoles(mountPath string, conf *testutil.AzureTestConf, maxSecs int) string {
	template := `
resource "vault_azure_secret_backend" "test" {
	path = "{{mountPath}}"
	subscription_id = "{{subscriptionID}}"
	tenant_id = "{{tenantID}}"
	client_id = "{{clientID}}"
	client_secret = "{{clientSecret}}"
}

resource "vault_azure_secret_backend_role" "test" {
	backend = vault_azure_secret_backend.test.path
	role = "my-role"
	azure_roles {
		role_name = "Reader"
		scope = "{{scope}}"
	}
	ttl = 300
	max_ttl = 600
}

data "vault_azure_access_credentials" "test" {
    backend = vault_azure_secret_backend.test.path
    role = vault_azure_secret_backend_role.test.role
    validate_creds = true
	num_seconds_between_tests = 1
	max_cred_validation_seconds = {{maxCredValidationSeconds}}
}`

	parsed := strings.Replace(template, "{{mountPath}}", mountPath, -1)
	parsed = strings.Replace(parsed, "{{subscriptionID}}", conf.SubscriptionID, -1)
	parsed = strings.Replace(parsed, "{{tenantID}}", conf.TenantID, -1)
	parsed = strings.Replace(parsed, "{{clientID}}", conf.ClientID, -1)
	parsed = strings.Replace(parsed, "{{clientSecret}}", conf.ClientSecret, -1)
	parsed = strings.Replace(parsed, "{{scope}}", conf.Scope, -1)
	parsed = strings.Replace(parsed, "{{maxCredValidationSeconds}}", strconv.Itoa(maxSecs), -1)
	return parsed
}

func testAccDataSourceAzureAccessCredentialsConfigBasicAzureGroups(mountPath string, conf *testutil.AzureTestConf, groupName string, maxSecs int) string {
	template := `
resource "vault_azure_secret_backend" "test" {
	path = "{{mountPath}}"
	subscription_id = "{{subscriptionID}}"
	tenant_id = "{{tenantID}}"
	client_id = "{{clientID}}"
	client_secret = "{{clientSecret}}"
}

resource "vault_azure_secret_backend_role" "test" {
	backend = vault_azure_secret_backend.test.path
	role = "my-role"
	azure_groups {
		group_name = "{{groupName}}"
	}
	ttl = 300
	max_ttl = 600
}

data "vault_azure_access_credentials" "test" {
    backend = vault_azure_secret_backend.test.path
    role = vault_azure_secret_backend_role.test.role
    validate_creds = true
	num_sequential_successes = 1
	num_seconds_between_tests = 30
	max_cred_validation_seconds = {{maxCredValidationSeconds}}
}`

	parsed := strings.Replace(template, "{{mountPath}}", mountPath, -1)
	parsed = strings.Replace(parsed, "{{subscriptionID}}", conf.SubscriptionID, -1)
	parsed = strings.Replace(parsed, "{{tenantID}}", conf.TenantID, -1)
	parsed = strings.Replace(parsed, "{{clientID}}", conf.ClientID, -1)
	parsed = strings.Replace(parsed, "{{clientSecret}}", conf.ClientSecret, -1)
	parsed = strings.Replace(parsed, "{{groupName}}", groupName, -1)
	parsed = strings.Replace(parsed, "{{maxCredValidationSeconds}}", strconv.Itoa(maxSecs), -1)
	return parsed
}

func Test_getAzureCloudConfigFromName(t *testing.T) {
	t.Parallel()

	mixedCap := func(s string) string {
		var r string
		s = strings.ToUpper(s)
		for i := 0; i < len(s); i++ {
			l := fmt.Sprintf("%c", s[i])
			if i%2 == 0 {
				r += strings.ToLower(l)
			} else {
				r += l
			}
		}
		return r
	}

	type test struct {
		name      string
		cloudName string
		want      cloud.Configuration
		wantErr   bool
	}
	tests := []test{
		{
			name:      "invalid",
			cloudName: "unknown",
			wantErr:   true,
		},
		{
			name:      "empty",
			cloudName: "",
			want:      cloud.AzurePublic,
			wantErr:   false,
		},
	}
	for k, v := range azureCloudConfigMap {
		tests = append(tests, test{
			name:      "mixed-" + k,
			cloudName: mixedCap(k),
			want:      v,
			wantErr:   false,
		})
		tests = append(tests, test{
			name:      "default-" + k,
			cloudName: k,
			want:      v,
			wantErr:   false,
		})
		tests = append(tests, test{
			name:      "lower-" + strings.ToLower(k),
			cloudName: strings.ToLower(k),
			want:      v,
			wantErr:   false,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAzureCloudConfigFromName(tt.cloudName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAzureCloudConfigFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAzureCloudConfigFromName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
