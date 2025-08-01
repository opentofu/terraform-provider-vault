---
layout: "vault"
page_title: "Vault: vault_ssh_secret_backend_ca resource"
sidebar_current: "docs-vault-resource-ssh-secret-backend-ca"
description: |-
  Managing CA information in an SSH secret backend in Vault
---

# vault\_ssh\_secret\_backend\_ca

Provides a resource to manage CA information in an SSH secret backend
[SSH secret backend within Vault](https://www.vaultproject.io/docs/secrets/ssh/index.html).

## Example Usage

```hcl
resource "vault_mount" "example" {
    type = "ssh"
}

resource "vault_ssh_secret_backend_ca" "foo" {
    backend = vault_mount.example.path
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
   *Available only for Vault Enterprise*.

* `backend` - (Optional) The path where the SSH secret backend is mounted. Defaults to 'ssh'

* `generate_signing_key` - (Optional) Whether Vault should generate the signing key pair internally. Defaults to true

* `key_type` - (Optional) Specifies the desired key type for the generated SSH CA key when `generate_signing_key` is set to `true`.

* `key_bits` - (Optional) Specifies the desired key bits for the generated SSH CA key when `generate_signing_key` is set to `true`.

* `public_key` - (Optional) The public key part the SSH CA key pair; required if generate_signing_key is false.

* `private_key` - (Optional) The private key part the SSH CA key pair; required if generate_signing_key is false.

* `managed_key_name` - (Optional) The name of the managed key to use. When using a managed key, this field or managed_key_id is required.

* `managed_key_id` - (Optional) The ID of the managed key to use. When using a managed key, this field or managed_key_name is required.

~> **Important** Because Vault does not support reading the private_key back from the API, Terraform cannot detect
and correct drift on `private_key`. Changing the values, however, _will_ overwrite the previously stored values.


## Attributes Reference

No additional attributes are exposed by this resource.

## Import

SSH secret backend CAs can be imported using the `path`, e.g.

```
$ terraform import vault_ssh_secret_backend_ca.foo ssh
```
