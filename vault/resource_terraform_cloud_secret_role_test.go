// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

func TestTerraformCloudSecretRole(t *testing.T) {
	backend := acctest.RandomWithPrefix("tf-test-backend")
	name := acctest.RandomWithPrefix("tf-test-name")
	organization := "hashicorp-vault-testing"
	vals := testutil.SkipTestEnvUnset(t, "TEST_TF_TOKEN", "TEST_TF_TEAM_ID", "TEST_TF_USER_ID")
	token, teamID, userID := vals[0], vals[1], vals[2]
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		PreCheck: func() {
			testutil.TestAccPreCheck(t)
		},
		CheckDestroy: testAccTerraformCloudSecretRoleCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTerraformCloudSecretRole_initialConfig(backend, token, name, organization, teamID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "name", name),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "organization", organization),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "ttl", "0"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "max_ttl", "0"),

					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "name", name+"_team_id"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "organization", organization),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "team_id", teamID),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "ttl", "0"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "max_ttl", "0"),

					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "name", name+"_user_id"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "user_id", userID),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "ttl", "0"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "max_ttl", "0"),
				),
			},
			{
				Config: testTerraformCloudSecretRole_updateConfig(backend, token, name, organization, teamID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "name", name),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "organization", organization),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "ttl", "120"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_org", "max_ttl", "240"),

					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "name", name+"_team_id"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "organization", organization),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "team_id", teamID),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "ttl", "120"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_team", "max_ttl", "0"),

					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "name", name+"_user_id"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "user_id", userID),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "ttl", "120"),
					resource.TestCheckResourceAttr("vault_terraform_cloud_secret_role.test_user", "max_ttl", "0"),
				),
			},
		},
	})
}

func testAccTerraformCloudSecretRoleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vault_terraform_cloud_secret_role" {
			continue
		}

		client, e := provider.GetClient(rs.Primary, testProvider.Meta())
		if e != nil {
			return e
		}

		secret, err := client.Logical().Read(rs.Primary.ID)
		if err != nil {
			return err
		}
		if secret != nil {
			return fmt.Errorf("role %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testTerraformCloudSecretRole_initialConfig(backend, token, name, organization, teamId, userId string) string {
	return fmt.Sprintf(`
resource "vault_terraform_cloud_secret_backend" "test" {
  backend = "%s"
  description = "test description"
  token = "%s"
}

resource "vault_terraform_cloud_secret_role" "test_org" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%s"
  organization = "%s"
}

resource "vault_terraform_cloud_secret_role" "test_team" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%[3]s_team_id"
  organization = "%[4]s"
  team_id = "%[5]s"
}

resource "vault_terraform_cloud_secret_role" "test_user" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%[3]s_user_id"
  user_id = "%[6]s"
}
`, backend, token, name, organization, teamId, userId)
}

func testTerraformCloudSecretRole_updateConfig(backend, token, name, organization, teamId, userId string) string {
	return fmt.Sprintf(`
resource "vault_terraform_cloud_secret_backend" "test" {
  backend = "%s"
  description = "test description"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds = 86400
  address = "https://app.terraform.io"
  token = "%s"
}

resource "vault_terraform_cloud_secret_role" "test_org" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%s"
  organization = "%s"

  ttl = 120
  max_ttl = 240
}

resource "vault_terraform_cloud_secret_role" "test_team" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%[3]s_team_id"
  organization = "%[4]s"
  team_id = "%[5]s"

  ttl = 120
}

resource "vault_terraform_cloud_secret_role" "test_user" {
  backend = vault_terraform_cloud_secret_backend.test.backend
  name = "%[3]s_user_id"
  user_id = "%[6]s"

  ttl = 120
}
`, backend, token, name, organization, teamId, userId)
}

func TestTerraformCloudSecretBackendRoleNameFromPath(t *testing.T) {
	{
		name, err := terraformCloudSecretRoleNameFromPath("foo/role/bar")
		if err != nil {
			t.Fatalf("error getting name: %v", err)
		}
		if name != "bar" {
			t.Fatalf("expected name 'bar', but got %s", name)
		}
	}

	{
		name, err := terraformCloudSecretRoleNameFromPath("no match")
		if err == nil {
			t.Fatal("Expected error getting name but got nil")
		}
		if name != "" {
			t.Fatalf("expected empty name, but got %s", name)
		}
	}
}

func TestTerraformCloudSecretBackendRoleBackendFromPath(t *testing.T) {
	{
		backend, err := terraformCloudSecretRoleBackendFromPath("foo/role/bar")
		if err != nil {
			t.Fatalf("error getting backend: %v", err)
		}
		if backend != "foo" {
			t.Fatalf("expected backend 'foo', but got %s", backend)
		}
	}

	{
		backend, err := terraformCloudSecretRoleBackendFromPath("no match")
		if err == nil {
			t.Fatal("Expected error getting backend but got nil")
		}
		if backend != "" {
			t.Fatalf("expected empty backend, but got %s", backend)
		}
	}
}
