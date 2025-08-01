---
layout: "vault"
page_title: "Vault: vault_azure_secret_backend resource"
sidebar_current: "docs-vault-resource-azure-secret-backend"
description: |-
  Creates an azure secret backend for Vault.
---

# vault\_azure\_secret\_backend

Creates an Azure Secret Backend for Vault.

The Azure secrets engine dynamically generates Azure service principals and role assignments. Vault roles can be mapped to one or more Azure roles, providing a simple, flexible way to manage the permissions granted to generated service principals.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

~> It is highly recommended that one transition to the Microsoft Graph API.
See [use_microsoft_graph_api ](https://www.vaultproject.io/api-docs/secret/azure#use_microsoft_graph_api)
for more information. The example below demonstrates how to do this. 

## Example Usage:

You can setup the Azure secrets engine with Workload Identity Federation (WIF) for a secret-less configuration:
```hcl
resource "vault_azure_secret_backend" "azure" {
  subscription_id         = "11111111-2222-3333-4444-111111111111"
  tenant_id               = "11111111-2222-3333-4444-222222222222"
  client_id               = "11111111-2222-3333-4444-333333333333"
  identity_token_audience = "<TOKEN_AUDIENCE>"
  identity_token_ttl      = "<TOKEN_TTL>"
  rotation_schedule       = "0 * * * SAT"
  rotation_window         = 3600
}
```

```hcl
resource "vault_azure_secret_backend" "azure" {
  subscription_id         = "11111111-2222-3333-4444-111111111111"
  tenant_id               = "11111111-2222-3333-4444-222222222222"
  client_id               = "11111111-2222-3333-4444-333333333333"
  client_secret           = "12345678901234567890"
  environment             = "AzurePublicCloud"
  rotation_schedule       = "0 * * * SAT"
  rotation_window         = 3600
}
```

## Argument Reference

The following arguments are supported:

- `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
   *Available only for Vault Enterprise*.

- `subscription_id` (`string: <required>`) - The subscription id for the Azure Active Directory.

- `tenant_id` (`string: <required>`) - The tenant id for the Azure Active Directory.

- `client_id` (`string:""`) - The OAuth2 client id to connect to Azure.

- `client_secret` (`string:""`) - The OAuth2 client secret to connect to Azure.

- `environment` (`string:""`) - The Azure environment.

- `path` (`string: <optional>`) - The unique path this backend should be mounted at. Defaults to `azure`.

- `identity_token_audience` - (Optional) The audience claim value. Requires Vault 1.17+.
  *Available only for Vault Enterprise*

- `identity_token_ttl` - (Optional) The TTL of generated identity tokens in seconds. Requires Vault 1.17+.
  *Available only for Vault Enterprise*

- `rotation_period` - (Optional) The amount of time in seconds Vault should wait before rotating the root credential.
  A zero value tells Vault not to rotate the root credential. The minimum rotation period is 10 seconds. Requires Vault Enterprise 1.19+.
  *Available only for Vault Enterprise*

- `rotation_schedule` - (Optional) The schedule, in [cron-style time format](https://en.wikipedia.org/wiki/Cron),
  defining the schedule on which Vault should rotate the root token. Requires Vault Enterprise 1.19+.
  *Available only for Vault Enterprise*

- `rotation_window` - (Optional) The maximum amount of time in seconds allowed to complete
  a rotation when a scheduled token rotation occurs. The default rotation window is
  unbound and the minimum allowable window is `3600`. Requires Vault Enterprise 1.19+. *Available only for Vault Enterprise*

- `disable_automated_rotation` - (Optional) Cancels all upcoming rotations of the root credential until unset. Requires Vault Enterprise 1.19+.
  *Available only for Vault Enterprise*

- `disable_remount` - (Optional) If set, opts out of mount migration on path updates.
  See here for more info on [Mount Migration](https://www.vaultproject.io/docs/concepts/mount-migration)


### Common Mount Arguments
These arguments are common across all resources that mount a secret engine.

* `description` - (Optional) Human-friendly description of the mount

* `default_lease_ttl_seconds` - (Optional) Default lease duration for tokens and secrets in seconds

* `max_lease_ttl_seconds` - (Optional) Maximum possible lease duration for tokens and secrets in seconds

* `audit_non_hmac_response_keys` - (Optional) Specifies the list of keys that will not be HMAC'd by audit devices in the response data object.

* `audit_non_hmac_request_keys` - (Optional) Specifies the list of keys that will not be HMAC'd by audit devices in the request data object.

* `local` - (Optional) Boolean flag that can be explicitly set to true to enforce local mount in HA environment

* `options` - (Optional) Specifies mount type specific options that are passed to the backend

* `seal_wrap` - (Optional) Boolean flag that can be explicitly set to true to enable seal wrapping for the mount, causing values stored by the mount to be wrapped by the seal's encryption capability

* `external_entropy_access` - (Optional) Boolean flag that can be explicitly set to true to enable the secrets engine to access Vault's external entropy source

* `allowed_managed_keys` - (Optional) Set of managed key registry entry names that the mount in question is allowed to access

* `listing_visibility` - (Optional) Specifies whether to show this mount in the UI-specific
  listing endpoint. Valid values are `unauth` or `hidden`. If not set, behaves like `hidden`.

* `passthrough_request_headers` - (Optional) List of headers to allow and pass from the request to
  the plugin.

* `allowed_response_headers` - (Optional) List of headers to allow, allowing a plugin to include
  them in the response.

* `delegated_auth_accessors` - (Optional)  List of allowed authentication mount accessors the
  backend can request delegated authentication for.

* `plugin_version` - (Optional) Specifies the semantic version of the plugin to use, e.g. "v1.0.0".
  If unspecified, the server will select any matching unversioned plugin that may have been
  registered, the latest versioned plugin registered, or a built-in plugin in that order of precedence.

* `identity_token_key` - (Optional)  The key to use for signing plugin workload identity tokens. If
  not provided, this will default to Vault's OIDC default key. Requires Vault Enterprise 1.16+.


## Attributes Reference

No additional attributes are exported by this resource.
