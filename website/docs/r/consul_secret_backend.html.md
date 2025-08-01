---
layout: "vault"
page_title: "Vault: vault_consul_secret_backend resource"
sidebar_current: "docs-vault-resource-consul-secret-backend"
description: |-
  Creates a Consul secret backend for Vault.
---

# vault\_consul\_secret\_backend

Creates a Consul Secret Backend for Vault. Consul secret backends can then issue Consul tokens, once a role has been
added to the backend.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

## Example Usage

#### Creating a standard backend resource:
```hcl
resource "vault_consul_secret_backend" "test" {
  path        = "consul"
  description = "Manages the Consul backend"
  address     = "127.0.0.1:8500"
  token       = "4240861b-ce3d-8530-115a-521ff070dd29"
}
```

#### Creating a backend resource to bootstrap a new Consul instance:
```hcl
resource "vault_consul_secret_backend" "test" {
  path        = "consul"
  description = "Bootstrap the Consul backend"
  address     = "127.0.0.1:8500"
  bootstrap   = true
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
   *Available only for Vault Enterprise*.

* `token` - (Optional) The Consul management token this backend should use to issue new tokens. This field is required
when `bootstrap` is false.

~> **Important** Because Vault does not support reading the configured token back from the API, Terraform cannot detect
and correct drift on `token`. Changing the value, however, _will_ overwrite the previously stored values.

* `bootstrap` - (Optional) Denotes that the resource is used to bootstrap the Consul ACL system.

~> **Important** When `bootstrap` is true, Vault will attempt to bootstrap the Consul server. The token returned from
this operation will only ever be known to Vault. If the resource is ever destroyed, the bootstrap token will be lost
and a [Consul reset may be required.](https://learn.hashicorp.com/tutorials/consul/access-control-troubleshoot#reset-the-acl-system)

* `path` - (Optional) The unique location this backend should be mounted at. Must not begin or end with a `/`. Defaults
to `consul`.

* `disable_remount` - (Optional) If set, opts out of mount migration on path updates.
  See here for more info on [Mount Migration](https://www.vaultproject.io/docs/concepts/mount-migration)

* `description` - (Optional) A human-friendly description for this backend.

* `address` - (Required) Specifies the address of the Consul instance, provided as "host:port" like "127.0.0.1:8500".

* `scheme` - (Optional) Specifies the URL scheme to use. Defaults to `http`.

* `ca_cert` - (Optional) CA certificate to use when verifying Consul server certificate, must be x509 PEM encoded.

* `client_cert` - (Optional) Client certificate used for Consul's TLS communication, must be x509 PEM encoded and if
this is set you need to also set client_key.

* `client_key` - (Optional) Client key used for Consul's TLS communication, must be x509 PEM encoded and if this is set
you need to also set client_cert.

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

## Import

Consul secret backends can be imported using the `path`, e.g.

```
$ terraform import vault_consul_secret_backend.example consul
```
