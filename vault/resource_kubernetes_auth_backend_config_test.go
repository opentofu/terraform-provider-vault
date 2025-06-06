// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

const (
	kubernetesJWT               = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWF1dGgtdmF1bHQtb3BlcmF0b3IiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoia3ViZXJuZXRlcy1hdXRoLXZhdWx0LW9wZXJhdG9yLXRva2VuLWZycmc3Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Imt1YmVybmV0ZXMtYXV0aC12YXVsdC1vcGVyYXRvciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjMwYzRiZjdkLTMwZmYtMTFlOC04ODdkLTA4MDAyNzZhYmI4OCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWF1dGgtdmF1bHQtb3BlcmF0b3I6a3ViZXJuZXRlcy1hdXRoLXZhdWx0LW9wZXJhdG9yIn0.V6lWrH6rgNfghn5Qc9IdPwxrAV0E8cdVNvGh3KmVCZpZVwOnL4eCQ3R6V37pO7ssTs-0aYYWc2NYcGnLiXvUPah89uK2wkE_Eod3NgWDqlutcM-LJuIK6xubuH0y2Bpb7ZddZmtc5MOa8e88iwiZmQ_zKhifwESdwFWaA5Nn1QNzwIPu2kOZU0Wz9sVN4i_NETUGqaEQYVU6DF--gErCLeloUDERW-QyrCRZ-ymTFt7UWRiPi3zAZ0-BG8j4TsjNYLiifGiMiaD6Ss-pd0brVMzQylpVlnZ7Of6ywIv-BWVa278ki3cd1RMqQj8tzHNg2tlbBSLMn92Gxh16jBW90w"
	kubernetesAnotherJWT        = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tNzhsNXAiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjUzZjZhY2Y5LTMwZmEtMTFlOC04ODdkLTA4MDAyNzZhYmI4OCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.pM5ugxaTX22vBsb4IOUz-pgwKM_ZfhgFhS1PKfpcSJYs4h-C4OujjKpF8j-Lw5oxaHOIxbROUxurlh-9eqvYqKREQOVxZvhoqxiflWCuAtu4RxHI-x4COSqV0H7pc_JNnDbgEqBhbFW1UiKfoye3QiqiqwYBaxvdpyH3uarv5yi26FT2lNvy6rHWMaMg3VZLzbZPOPY-v1C0RUbiyCz100A2UvaU5QbdHfwFab18vqgB_FN4aFXP9TKrcDUPyFyAhoC6h4Tb_ounuQ1u8UWtLL_KwDK7KEAgwg-FfayzHtw41PneS9nNtNm7bZJsLtzsvSzuMJpwehnN1GtUCMvDzA"
	envVarTFAccK8sSkipInCluster = "TF_ACC_K8S_SKIP_IN_CLUSTER"
)

const kubernetesCAcert = `-----BEGIN CERTIFICATE-----
MIIFXTCCA0WgAwIBAgIJAOLBeuu/P2O6MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTgwNDExMTEwNDE5WhcNMTkwNDExMTEwNDE5WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAx8AtbNS5rSO8ZlHE+s7yzJI2qBNnOZsfyWApfw4rie0u3lXm2aHYMK4B
gg1vBVeUe/LbrnRDMMkZ5cKU24CnUNa+2A6ITQPnyHXEJxzR/O6fPgxW88itCha1
ZWOlTe8dhcNWh95Id+4m5H9RztdCpI7hwvo8HuCdlWSP/HsceI6nAOyw2dPBKZ1B
Z9FaqA6r+ET/kGs/iU9a7ZJtNlPBmAHm54hj+eNVmroP5JALpf8zcoPGIVhvv9Mu
GeveSg7CPMVpnjVnjQfI5ZnJmmLI6wqgpvuvhycMbpvbOvsJ6JIfudwIQv/Suk10
R+K4booiFs0yPkrZ+NjbmdwRMqV7R3cPqTOzzHtG+HC+AJB4t4ecAULWOzqH/KY7
ICpNQZm3o2sT0ZeL5p1v/up3cncASgJBGMapLB5rrBdxykMifeioZHt+U55O0bq3
X0tRmaVSLMwWVa8UQQNK6pEKFsDYCC+knYUmT7Fyt6xOstE0Zwrpbda33BULVLlM
jGtY4ZqMEmxRd6iwuX+XziPJ+kSjyBi8SBTeP3bVh4OXkLnddNL45fwIrb7VB8EJ
C5RYUNH+nCJS3NW16ZFM/TSvn+r5K/iRFIMy+YUHloL4hiVKkA4bF17zEXfqynpA
Xw7wxwEOGnNnE1rz4XhQp5LqZgRZlRWEdZ9naFNIj9WsUcKHqRECAwEAAaNQME4w
HQYDVR0OBBYEFFSgoZGl4H8Y3/pb2aDAWopCW/nvMB8GA1UdIwQYMBaAFFSgoZGl
4H8Y3/pb2aDAWopCW/nvMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIB
AFwauMwjtDSLxt4LvzTE/7Zz9vZ1iavGisRFQmAH3pK/RAB4830b3Y5C98abZRbt
7qS+J1xia7MSPQr2ex7EHK6eE1U3apuOgcckFrR4DiieZEjRxOlqHRUXnrtHHNsi
v3fx6IeoB5/685lwAWyS65R5lJsggnLiG4gwl3t0uN50/pjegN+iAzn3krnLC85c
u1dQAhY9XiiBy1jcX+zDQBVi+YRp5gk4KiBipwE0gB0aUJIlODdSqe5Tl0DqD7/6
W95ABX9ksuaREoFDVWMEgsPQlXj5cD7nQsX6Ghnrb60m9s60/Wnftjw+UHcm/QK3
28QlfcBvA4cfjAW+8WCuhKMjyUPEdD8tHG9C/PN/i+c5MOkmCGMIGbay6AN0nQM6
TlrjK0wHk1eBIs+u5GqL1Cg+GKddeqFP9slMCyuSKyDJzn4uZB8j+NwrdM9Ma3UY
dwWpxrJOJY5k2f/GFMmOgk4qite+PZT5nCT8YLiO5nN6nCxQyEtOwRIbDrZjJSeH
g+ra+LRtnP6DGO9r/2EO5XLUgLp8hPZbq3+xd8TWMpv5TnuqPGJ5gM4rQzH3+H75
2K27ycqQbqM7ceOoihM6hb2VhEeoq8nYxWKZ4OGONpoaHA0tdYFGyVRYBBMUK/DB
1B4wS06RrGt6oud4fLHZpuvspPQTlLRjHvuXDi/cbIIE
-----END CERTIFICATE-----`

const kubernetesPEMfile = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAx8AtbNS5rSO8ZlHE+s7y
zJI2qBNnOZsfyWApfw4rie0u3lXm2aHYMK4Bgg1vBVeUe/LbrnRDMMkZ5cKU24Cn
UNa+2A6ITQPnyHXEJxzR/O6fPgxW88itCha1ZWOlTe8dhcNWh95Id+4m5H9RztdC
pI7hwvo8HuCdlWSP/HsceI6nAOyw2dPBKZ1BZ9FaqA6r+ET/kGs/iU9a7ZJtNlPB
mAHm54hj+eNVmroP5JALpf8zcoPGIVhvv9MuGeveSg7CPMVpnjVnjQfI5ZnJmmLI
6wqgpvuvhycMbpvbOvsJ6JIfudwIQv/Suk10R+K4booiFs0yPkrZ+NjbmdwRMqV7
R3cPqTOzzHtG+HC+AJB4t4ecAULWOzqH/KY7ICpNQZm3o2sT0ZeL5p1v/up3cncA
SgJBGMapLB5rrBdxykMifeioZHt+U55O0bq3X0tRmaVSLMwWVa8UQQNK6pEKFsDY
CC+knYUmT7Fyt6xOstE0Zwrpbda33BULVLlMjGtY4ZqMEmxRd6iwuX+XziPJ+kSj
yBi8SBTeP3bVh4OXkLnddNL45fwIrb7VB8EJC5RYUNH+nCJS3NW16ZFM/TSvn+r5
K/iRFIMy+YUHloL4hiVKkA4bF17zEXfqynpAXw7wxwEOGnNnE1rz4XhQp5LqZgRZ
lRWEdZ9naFNIj9WsUcKHqRECAwEAAQ==
-----END PUBLIC KEY-----`

func TestAccKubernetesAuthBackendConfig_import(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	jwt := kubernetesJWT
	issuer := "kubernetes/serviceaccount"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, jwt, issuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", jwt),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, issuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				ResourceName:      "vault_kubernetes_auth_backend_config.config",
				ImportState:       true,
				ImportStateVerify: true,
				// NOTE: The API can't serve these fields, so ignore them.
				ImportStateVerifyIgnore: []string{"backend", "token_reviewer_jwt"},
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, jwt, kubernetesCAcert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", jwt),
				),
			},
			{
				ResourceName:      "vault_kubernetes_auth_backend_config.config",
				ImportState:       true,
				ImportStateVerify: true,
				// NOTE: The API can't serve these fields, so ignore them.
				ImportStateVerifyIgnore: []string{"backend", "token_reviewer_jwt"},
			},
		},
	})
}

func TestAccKubernetesAuthBackendConfig_basic(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	jwt := kubernetesJWT

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, jwt, kubernetesCAcert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", jwt),
				),
			},
		},
	})
}

func testAccCheckKubernetesAuthBackendConfigDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vault_kubernetes_auth_backend_config" {
			continue
		}

		client, e := provider.GetClient(rs.Primary, testProvider.Meta())
		if e != nil {
			return e
		}

		secret, err := client.Logical().Read(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error checking for Kubernetes auth backend config %q: %s", rs.Primary.ID, err)
		}
		if secret != nil {
			return fmt.Errorf("Kubernetes auth backend config %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccKubernetesAuthBackendConfig_update(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	oldJWT := kubernetesJWT
	newJWT := kubernetesAnotherJWT

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, oldJWT, kubernetesCAcert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, newJWT, kubernetesCAcert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", newJWT),
				),
			},
		},
	})
}

func TestAccKubernetesAuthBackendConfig_full(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	jwt := kubernetesJWT
	issuer := "api"
	testResource := "vault_kubernetes_auth_backend_config.config"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, jwt, issuer,
					true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResource, "backend", backend),
					resource.TestCheckResourceAttr(testResource, consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr(testResource, consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr(testResource, "token_reviewer_jwt", jwt),
					resource.TestCheckResourceAttr(testResource, "pem_keys.#", "1"),
					resource.TestCheckResourceAttr(testResource, "pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr(testResource, consts.FieldIssuer, "api"),
					resource.TestCheckResourceAttr(testResource, consts.FieldDisableISSValidation, strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(testResource, consts.FieldDisableLocalCAJWT, strconv.FormatBool(true)),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					meta := testProvider.Meta().(*provider.ProviderMeta)
					return !meta.IsAPISupported(provider.VaultVersion116), nil
				},
				Config: testAccKubernetesAuthBackendConfig_useAnnotations(backend, jwt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResource, fieldUseAnnotationsAsAliasMetadata, strconv.FormatBool(true)),
				),
			},
		},
	})
}

func TestAccKubernetesAuthBackendConfig_fullUpdate(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	oldJWT := kubernetesJWT
	newJWT := kubernetesAnotherJWT
	oldIssuer := "kubernetes/serviceaccount"
	newIssuer := "api"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, oldJWT, oldIssuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, newJWT, newIssuer,
					true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", newJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, newIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(true)),
				),
			},
			{
				// ensure we can set disable_iss_validation to false
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, newJWT, newIssuer,
					false, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", newJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, newIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(true)),
				),
			},
		},
	})
}

func TestAccKubernetesAuthBackendConfig_localCA(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	jwt := kubernetesJWT

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutil.TestAccPreCheck(t)
			SkipIfAPIVersionGTE(t, testProvider.Meta(), vaultVersion193)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, jwt, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, ""),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", jwt),
				),
			},
			{
				PreConfig: func() {
					client := testProvider.Meta().(*provider.ProviderMeta).MustGetClient()

					path := kubernetesAuthBackendConfigPath(backend)
					if _, err := client.Logical().Write(path, map[string]interface{}{
						consts.FieldKubernetesCACert: kubernetesCAcert,
						consts.FieldKubernetesHost:   "http://example.com:443",
						"token_reviewer_jwt":         jwt,
					}); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccKubernetesAuthBackendConfigConfig_basic(backend, jwt, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", jwt),
				),
			},
		},
	})
}

func testAccKubernetesAuthBackendConfigConfig_basic(backend, jwt, ca string) string {
	config := fmt.Sprintf(`
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
  path = "%s"
}

resource "vault_kubernetes_auth_backend_config" "config" {
  backend = vault_auth_backend.kubernetes.path
  kubernetes_host = "http://example.com:443"
  token_reviewer_jwt = %q
`, backend, jwt)

	if ca != "" {
		config += fmt.Sprintf(`
  kubernetes_ca_cert = %q
`, ca)
	}

	return config + "}"
}

func testAccKubernetesAuthBackendConfig_useAnnotations(backend, jwt string) string {
	return fmt.Sprintf(`
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
  path = "%s"
}

resource "vault_kubernetes_auth_backend_config" "config" {
  backend = vault_auth_backend.kubernetes.path
  kubernetes_host = "http://example.com:443"
  token_reviewer_jwt = %q
  use_annotations_as_alias_metadata = true
}
`, backend, jwt)
}

func testAccKubernetesAuthBackendConfigConfig_full(backend, caCert, jwt, issuer string,
	disableIssValidation, disableLocalCaJwt, omitCA bool,
) string {
	var caConfig string
	if !omitCA {
		caConfig = fmt.Sprintf(`kubernetes_ca_cert = %q`, caCert)
	}

	config := fmt.Sprintf(`
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
  path = "%s"
}

resource "vault_kubernetes_auth_backend_config" "config" {
  backend = vault_auth_backend.kubernetes.path
  kubernetes_host = "http://example.com:443"
  %s
  token_reviewer_jwt = %q
  pem_keys = [%q]
  issuer = %q
  disable_iss_validation = %t
  disable_local_ca_jwt = %t
}`, backend, caConfig, jwt, kubernetesPEMfile, issuer, disableIssValidation, disableLocalCaJwt)

	return config
}

func testAccKubernetesAuthBackendConfigConfig_fullUnsetCA(backend, jwt, issuer string, disableIssValidation, disableLocalCaJwt bool) string {
	config := fmt.Sprintf(`
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
  path = "%s"
}

resource "vault_kubernetes_auth_backend_config" "config" {
  backend = vault_auth_backend.kubernetes.path
  kubernetes_host = "http://example.com:443"
  token_reviewer_jwt = %q
  pem_keys = [%q]
  issuer = %q
  disable_iss_validation = %t
  disable_local_ca_jwt = %t
}`, backend, jwt, kubernetesPEMfile, issuer, disableIssValidation, disableLocalCaJwt)

	return config
}

func TestAccKubernetesAuthBackendConfig_fullInK8sCluster(t *testing.T) {
	backend := acctest.RandomWithPrefix("kubernetes")
	oldJWT := kubernetesJWT
	// newJWT := kubernetesAnotherJWT
	oldIssuer := "kubernetes/serviceaccount"
	// newIssuer := "api"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutil.TestAccPreCheck(t)
			testutil.SkipTestEnvSet(t, envVarTFAccK8sSkipInCluster)
			SkipIfAPIVersionLT(t, testProvider.Meta(), vaultVersion193)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		CheckDestroy:             testAccCheckKubernetesAuthBackendConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, "", oldJWT, oldIssuer,
					false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, ""),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, oldJWT, oldIssuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, "", oldJWT, oldIssuer,
					false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, ""),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, "", oldJWT, oldIssuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, ""),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				PreConfig: func() {
					client := testProvider.Meta().(*provider.ProviderMeta).MustGetClient()

					path := kubernetesAuthBackendConfigPath(backend)
					// mimic pre vault-1.9.3 behaviour, whereby the in-cluster K8s CA was stored in
					// Vault config.
					if _, err := client.Logical().Write(path, map[string]interface{}{
						consts.FieldKubernetesCACert: kubernetesCAcert,
						consts.FieldKubernetesHost:   "http://example.com:443",
					}); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, "", oldJWT, oldIssuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, ""),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
			{
				Config: testAccKubernetesAuthBackendConfigConfig_full(backend, kubernetesCAcert, oldJWT, oldIssuer,
					false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"backend", backend),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesHost, "http://example.com:443"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldKubernetesCACert, kubernetesCAcert),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"token_reviewer_jwt", oldJWT),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.#", "1"),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						"pem_keys.0", kubernetesPEMfile),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldIssuer, oldIssuer),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableISSValidation, strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("vault_kubernetes_auth_backend_config.config",
						consts.FieldDisableLocalCAJWT, strconv.FormatBool(false)),
				),
			},
		},
	})
}
