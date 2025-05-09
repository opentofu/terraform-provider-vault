// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
)

const DefaultMaxHTTPRetriesCCC = 10

// Description is essentially a DataSource or Resource with some additional metadata
// that helps with maintaining the Terraform Vault Provider.
type Description struct {
	// PathInventory is used for taking an inventory of the supported endpoints in the
	// Terraform Vault Provider and comparing them to the endpoints noted as available in
	// Vault's OpenAPI description. A list of Vault's endpoints can be obtained by,
	// from Vault's home directory, running "$ ./scripts/gen_openapi.sh", and then by
	// drilling into the paths with "$ cat openapi.json | jq ".paths" | jq 'keys[]'".
	// Here's a short example of how paths and their path variables should be represented:
	//		"/transit/keys/{name}/config"
	//		"/transit/random"
	//		"/transit/random/{urlbytes}"
	//		"/transit/sign/{name}/{urlalgorithm}"
	PathInventory []string

	// EnterpriseOnly defaults to false, but should be marked true if a resource is enterprise only.
	EnterpriseOnly bool

	Resource *schema.Resource
}

type ResourceRegistry map[string]*Description

type ResourcesMap map[string]*schema.Resource

func NewProvider(
	dataRegistry ResourceRegistry,
	resourceRegistry ResourceRegistry,
	extraResourcesMaps ...ResourcesMap,
) *schema.Provider {
	dataSourcesMap, err := parse(dataRegistry)
	if err != nil {
		panic(err)
	}

	coreResourcesMap, err := parse(resourceRegistry)
	if err != nil {
		panic(err)
	}

	for _, m := range extraResourcesMaps {
		MustAddSchemaResource(m, coreResourcesMap, nil)
	}

	r := &schema.Provider{
		// This schema must match exactly the fwprovider (Terraform Plugin Framework) schema.
		// Notably the attributes can have no Default values.
		Schema: map[string]*schema.Schema{
			// Not `Required` but must be set via config or env. Otherwise we
			// return an error.
			consts.FieldAddress: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the root of the target Vault server.",
			},
			"add_address_to_env": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If true, adds the value of the `address` argument to the Terraform process environment.",
			},
			// Not `Required` but must be set via config, env, or token helper.
			// Otherwise we return an error.
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Token to use to authenticate to Vault.",
			},
			"token_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Token name to use for creating the Vault child token.",
			},
			"skip_child_token": {
				Type:     schema.TypeBool,
				Optional: true,
				// Setting to true will cause max_lease_ttl_seconds and token_name to be ignored (not used).
				// Note that this is strongly discouraged due to the potential of exposing sensitive secret data.
				Description: "Set this to true to prevent the creation of ephemeral child token used by this provider.",
			},
			consts.FieldCACertFile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to a CA certificate file to validate the server's certificate.",
			},
			consts.FieldCACertDir: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to directory containing CA certificate files to validate the server's certificate.",
			},
			consts.FieldSkipTLSVerify: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this to true only if the target Vault server is an insecure development instance.",
			},
			consts.FieldTLSServerName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name to use as the SNI host when connecting via TLS.",
			},
			"max_lease_ttl_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum TTL for secret leases requested by this provider.",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of retries when a 5xx error code is encountered.",
			},
			"max_retries_ccc": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of retries for Client Controlled Consistency related operations",
			},
			consts.FieldNamespace: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The namespace to use. Available only for Vault Enterprise.",
			},
			"headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The headers to send with each Vault request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The header name",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The header value",
						},
					},
				},
			},
			consts.FieldSkipGetVaultVersion: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip the dynamic fetching of the Vault server version.",
			},
			consts.FieldVaultVersionOverride: {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Override the Vault server version, " +
					"which is normally determined dynamically from the target Vault server",
				ValidateDiagFunc: ValidateDiagSemVer,
			},
			consts.FieldSetNamespaceFromToken: {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "In the case where the Vault token is for a specific namespace " +
					"and the provider namespace is not configured, use the token namespace " +
					"as the root namespace for all resources.",
			},
			consts.FieldClientAuth: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Client authentication credentials.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						consts.FieldCertFile: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to a file containing the client certificate.",
						},
						consts.FieldKeyFile: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to a file containing the private key that the certificate was issued for.",
						},
					},
				},
			},
		},
		ConfigureFunc:  NewProviderMeta,
		DataSourcesMap: dataSourcesMap,
		ResourcesMap:   coreResourcesMap,
	}

	MustAddAuthLoginSchema(r.Schema)

	// Set the provider Meta (instance data) here.
	// It will be overwritten by the result of the call to ConfigureFunc,
	// but can be used pre-configuration by other (non-primary) provider servers.
	r.SetMeta(&ProviderMeta{})

	return r
}

func parse(descs map[string]*Description) (map[string]*schema.Resource, error) {
	var errs error
	resourceMap := make(map[string]*schema.Resource)
	for k, desc := range descs {
		resourceMap[k] = desc.Resource
		if len(desc.PathInventory) == 0 {
			errs = multierror.Append(errs, fmt.Errorf("%q needs its paths inventoried", k))
		}
	}
	return resourceMap, errs
}

// ReadWrapper provides common read operations to the wrapped schema.ReadFunc.
func ReadWrapper(f schema.ReadFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, i interface{}) error {
		if err := importNamespace(d); err != nil {
			return err
		}

		return f(d, i)
	}
}

// ReadContextWrapper provides common read operations to the wrapped schema.ReadContextFunc.
func ReadContextWrapper(f schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
		if err := importNamespace(d); err != nil {
			return diag.FromErr(err)
		}
		return f(ctx, d, i)
	}
}

// MountCreateContextWrapper performs a minimum version requirement check prior to the
// wrapped schema.CreateContextFunc.
func MountCreateContextWrapper(f schema.CreateContextFunc, minVersion *version.Version) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		currentVersion := meta.(*ProviderMeta).GetVaultVersion()

		if !IsAPISupported(meta, minVersion) {
			return diag.Errorf("feature not enabled on current Vault version. min version required=%s; "+
				"current vault version=%s", minVersion, currentVersion)
		}

		return f(ctx, d, meta)
	}
}

// UpdateContextWrapper performs a minimum version requirement check prior to the
// wrapped schema.UpdateContextFunc.
func UpdateContextWrapper(f schema.UpdateContextFunc, minVersion *version.Version) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		currentVersion := meta.(*ProviderMeta).GetVaultVersion()

		if !IsAPISupported(meta, minVersion) {
			return diag.Errorf("feature not enabled on current Vault version. min version required=%s; "+
				"current vault version=%s", minVersion, currentVersion)
		}

		return f(ctx, d, meta)
	}
}

func importNamespace(d *schema.ResourceData) error {
	if ns := os.Getenv(consts.EnvVarVaultNamespaceImport); ns != "" {
		s := d.State()
		var attemptNamespaceImport bool
		if s.Empty() {
			// state does not yet exist or is empty
			// import is acceptable
			attemptNamespaceImport = true
		} else {
			// only import if namespace
			// is not already set in state
			s.Lock()
			defer s.Unlock()
			_, ok := s.Attributes[consts.FieldNamespace]
			attemptNamespaceImport = !ok
		}
		if attemptNamespaceImport {
			log.Printf(`[INFO] Environment variable %s set, `+
				`attempting TF state import "%s=%s"`,
				consts.EnvVarVaultNamespaceImport, consts.FieldNamespace, ns)
			if err := d.Set(consts.FieldNamespace, ns); err != nil {
				return fmt.Errorf("failed to import %q, err=%w",
					consts.EnvVarVaultNamespaceImport, err)
			}
		}
	}

	return nil
}
