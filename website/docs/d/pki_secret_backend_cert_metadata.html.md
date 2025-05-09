---
layout: "vault"
page_title: "Vault: vault_pki_secret_backend_cert_metadata data source"
sidebar_current: "docs-vault-datasource-pki-secret-backend-cert-metadata"
description: |-
  Reads certificate metadata from Vault Enterprise.
---

# vault\_pki\_secret\_backend\_cert_metadata

Reads certificate metadata from Vault Enterprise.

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

resource "vault_pki_secret_backend_root_cert" "root" {
  backend     = vault_mount.pki.path
  type        = "internal"
  common_name = "example"
  ttl         = "86400"
  issuer_name = "example"
}

resource "vault_pki_secret_backend_role" "test" {
  backend           = vault_pki_secret_backend_root_cert.test.backend
  name              = "test"
  allowed_domains   = ["test.my.domain"]
  allow_subdomains  = true
  max_ttl           = "3600"
  key_usage         = ["DigitalSignature", "KeyAgreement", "KeyEncipherment"]
  no_store_metadata = false
}

resource "vault_pki_secret_backend_cert" "test" {
  backend               = vault_pki_secret_backend_role.test.backend
  name                  = vault_pki_secret_backend_role.test.name
  common_name           = "cert.test.my.domain"
  ttl                   = "720h"
  min_seconds_remaining = 60
  cert_metadata         = "dGVzdCBtZXRhZGF0YQ=="
}

data "vault_pki_secret_backend_cert_metadata" "test" {
  path = vault_mount.test-root.path
  serial = vault_pki_secret_backend_cert.test.serial_number
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace of the target resource.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
  *Available only for Vault Enterprise*.

* `path` - (Required) The path to the PKI secret backend to
  read the cert metadata from, with no leading or trailing `/`s.

* `serial` - (Required) Specifies the serial of the certificate whose metadata to read.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `issuer_id` - ID of the issuer.

* `cert_metadata` - The metadata associated with the certificate

* `expiration` - The expiration date of the certificate in unix epoch format

* `role` - The role used to create the certificate

* `serial_number` - The serial number