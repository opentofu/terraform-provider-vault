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

func TestAccTransformRole(t *testing.T) {
	path := acctest.RandomWithPrefix("transform")
	role := acctest.RandomWithPrefix("test-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestEntPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testTransformRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTransformRole_basicConfig(path, role, "ccn-fpe"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_transform_role.test", "path", path),
					resource.TestCheckResourceAttr("vault_transform_role.test", "name", role),
					resource.TestCheckResourceAttr("vault_transform_role.test", "transformations.0", "ccn-fpe"),
					resource.TestCheckResourceAttr("vault_transform_role.test", "transformations.#", "1"),
				),
			},
			{
				Config: testTransformRole_basicConfig(path, role, "ccn-fpe+updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_transform_role.test", "path", path),
					resource.TestCheckResourceAttr("vault_transform_role.test", "name", role),
					resource.TestCheckResourceAttr("vault_transform_role.test", "transformations.0", "ccn-fpe+updated"),
					resource.TestCheckResourceAttr("vault_transform_role.test", "transformations.#", "1"),
				),
			},
			{
				ResourceName:      "vault_transform_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testTransformRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vault_transform_role" {
			continue
		}
		client, e := provider.GetClient(rs.Primary, testProvider.Meta())
		if e != nil {
			return e
		}
		secret, err := client.Logical().Read(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error checking for role %q: %s", rs.Primary.ID, err)
		}
		if secret != nil {
			return fmt.Errorf("role %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testTransformRole_basicConfig(path, role, tranformations string) string {
	return fmt.Sprintf(`
resource "vault_mount" "mount_transform" {
  path = "%s"
  type = "transform"
}
resource "vault_transform_role" "test" {
  path = vault_mount.mount_transform.path
  name = "%s"
  transformations = ["%s"]
}
`, path, role, tranformations)
}
