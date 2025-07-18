---
layout: "vault"
page_title: "Vault: vault_pki_secret_backend_config_scep data source"
sidebar_current: "docs-vault-datasource-pki-secret-backend-config-scep"
description: |-
 Reads the PKI SCEP configuration from Vault Enterprise. 
---

# vault\_pki\_secret\_backend\_config\_scep

Reads the PKI SCEP configuration from Vault Enterprise.

~> **Important** All data retrieved from Vault will be
written in cleartext to state file generated by Terraform, will appear in
the console output when Terraform runs, and may be included in plan files
if secrets are interpolated into any resource attributes.
Protect these artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

## Example Usage

```hcl
resource "vault_mount" "pki" {
  path        = "pki"
  type        = "pki"
  description = "PKI secret engine mount"
}

data "vault_pki_secret_backend_config_scep" "scep_config" {
  backend     = vault_mount.pki.path
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace of the target resource.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
  *Available only for Vault Enterprise*.

* `backend` - (Required) The path to the PKI secret backend to
  read the SCEP configuration from, with no leading or trailing `/`s.
 
## Attributes Reference

* `allowed_digest_algorithms` - List of allowed digest algorithms for SCEP requests.

* `allowed_encryption_algorithms` - List of allowed encryption algorithms for SCEP requests.

* `authenticators` - Lists the mount accessors SCEP should delegate authentication requests towards (see [below for nested schema](#nestedatt--authenticators)).
 
* `default_path_policy` - Specifies the policy to be used for non-role-qualified SCEP requests; valid values are 'sign-verbatim', or "role:<role_name>" to specify a role to use as this policy.

* `enabled` - Specifies whether SCEP is enabled.

* `external_validation` - Lists the 3rd party validation of SCEP requests (see [below for nested schema](#nestedatt--externalvalidation)).

* `last_updated` - A read-only timestamp representing the last time the configuration was updated.

* `restrict_ca_chain_to_issuer` - If true, only return the issuer CA, otherwise the entire CA certificate chain will be returned if available from the PKI mount.


<a id="nestedatt--authenticators"></a>
### Nested Schema for `authenticators`

* `cert` - The accessor and cert_role properties for cert auth backends.
 
* `scep` - The accessor property for scep auth backends.

<a id="nestedatt--externalvalidation"></a>
### Nested Schema for `external_validation`

* `intune` - The tenant_id, client_id, client_secret and environment properties for Microsoft Intune validation of SCEP requests.

