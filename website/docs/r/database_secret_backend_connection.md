---
layout: "vault"
page_title: "Vault: vault_database_secret_backend_connection resource"
sidebar_current: "docs-vault-resource-database-secret-backend-connection"
description: |-
  Configures a database secret backend connection for Vault.
---

# vault\_database\_secret\_backend\_connection

Creates a Database Secret Backend connection in Vault. Database secret backend
connections can be used to generate dynamic credentials for the database.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

## Example Usage

```hcl
resource "vault_mount" "db" {
  path = "postgres"
  type = "database"
}

resource "vault_database_secret_backend_connection" "postgres" {
  backend           = vault_mount.db.path
  name              = "postgres"
  allowed_roles     = ["dev", "prod"]
  rotation_schedule = "0 * * * SAT"
  rotation_window   = 3600

  postgresql {
    connection_url = "postgres://username:password@host:port/database"
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](../index.html#namespace).
   *Available only for Vault Enterprise*.

* `name` - (Required) A unique name to give the database connection.

* `backend` - (Required) The unique name of the Vault mount to configure.

* `plugin_name` - (Optional) Specifies the name of the plugin to use.

* `verify_connection` - (Optional) Whether the connection should be verified on
  initial configuration or not.

* `allowed_roles` - (Optional) A list of roles that are allowed to use this
  connection.

* `root_rotation_statements` - (Optional) A list of database statements to be executed to rotate the root user's credentials.

* `data` - (Optional) A map of sensitive data to pass to the endpoint. Useful for templated connection strings.

* `rotation_period` - (Optional) The amount of time in seconds Vault should wait before rotating the root credential.
  A zero value tells Vault not to rotate the root credential. The minimum rotation period is 10 seconds. Requires Vault Enterprise 1.19+.

* `rotation_schedule` - (Optional) The schedule, in [cron-style time format](https://en.wikipedia.org/wiki/Cron),
  defining the schedule on which Vault should rotate the root token. Requires Vault Enterprise 1.19+.

* `rotation_window` - (Optional) The maximum amount of time in seconds allowed to complete
  a rotation when a scheduled token rotation occurs. The default rotation window is
  unbound and the minimum allowable window is `3600`. Requires Vault Enterprise 1.19+.

* `disable_automated_rotation` - (Optional) Cancels all upcoming rotations of the root credential until unset. Requires Vault Enterprise 1.19+.

* `cassandra` - (Optional) A nested block containing configuration options for Cassandra connections.

* `couchbase` - (Optional) A nested block containing configuration options for Couchbase connections.

* `mongodb` - (Optional) A nested block containing configuration options for MongoDB connections.

* `mongodbatlas` - (Optional) A nested block containing configuration options for MongoDB Atlas connections.

* `hana` - (Optional) A nested block containing configuration options for SAP HanaDB connections.

* `mssql` - (Optional) A nested block containing configuration options for MSSQL connections.

* `mysql` - (Optional) A nested block containing configuration options for MySQL connections.

* `mysql_rds` - (Optional) A nested block containing configuration options for RDS MySQL connections.

* `mysql_aurora` - (Optional) A nested block containing configuration options for Aurora MySQL connections.

* `mysql_legacy` - (Optional) A nested block containing configuration options for legacy MySQL connections.

* `postgresql` - (Optional) A nested block containing configuration options for PostgreSQL connections.

* `oracle` - (Optional) A nested block containing configuration options for Oracle connections.

* `elasticsearch` - (Optional) A nested block containing configuration options for Elasticsearch connections.

* `snowflake` - (Optional) A nested block containing configuration options for Snowflake connections.

* `influxdb` - (Optional) A nested block containing configuration options for InfluxDB connections.

* `redis` - (Optional) A nested block containing configuration options for Redis connections.

* `redis_elasticache` - (Optional) A nested block containing configuration options for Redis ElastiCache connections.

Exactly one of the nested blocks of configuration options must be supplied.

### Cassandra Configuration Options

* `hosts` - (Required) The hosts to connect to.

* `username` - (Required) The username to authenticate with.

* `password` - (Required) The password to authenticate with.

* `port` - (Optional) The default port to connect to if no port is specified as
  part of the host.

* `tls` - (Optional) Whether to use TLS when connecting to Cassandra.

* `insecure_tls` - (Optional) Whether to skip verification of the server
  certificate when using TLS.

* `pem_bundle` - (Optional) Concatenated PEM blocks configuring the certificate
  chain.

* `pem_json` - (Optional) A JSON structure configuring the certificate chain.

* `protocol_version` - (Optional) The CQL protocol version to use.

* `connect_timeout` - (Optional) The number of seconds to use as a connection
  timeout.

* `skip_verification` - (Optional) Skip permissions checks when a connection to Cassandra is first created.
  These checks ensure that Vault is able to create roles, but can be resource intensive in clusters with many roles.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Couchbase Configuration Options

* `hosts` - (Required) A set of Couchbase URIs to connect to. Must use `couchbases://` scheme if `tls` is `true`.

* `username` - (Required) Specifies the username for Vault to use.

* `password` - (Required) Specifies the password corresponding to the given username.

* `tls` - (Optional) Whether to use TLS when connecting to Couchbase.

* `insecure_tls` - (Optional) Whether to skip verification of the server
  certificate when using TLS.

* `base64_pem` - (Optional) Required if `tls` is `true`. Specifies the certificate authority of the Couchbase server, as a PEM certificate that has been base64 encoded.

* `bucket_name` - (Optional) Required for Couchbase versions prior to 6.5.0. This is only used to verify vault's connection to the server.

* `username_template` - (Optional) Template describing how dynamic usernames are generated.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### InfluxDB Configuration Options

* `host` - (Required) The host to connect to.

* `username` - (Required) The username to authenticate with.

* `password` - (Required) The password to authenticate with.

* `port` - (Optional) The default port to connect to if no port is specified as
  part of the host.

* `tls` - (Optional) Whether to use TLS when connecting to Cassandra.

* `insecure_tls` - (Optional) Whether to skip verification of the server
  certificate when using TLS.

* `pem_bundle` - (Optional) Concatenated PEM blocks configuring the certificate
  chain.

* `pem_json` - (Optional) A JSON structure configuring the certificate chain.

* `username_template` - (Optional) Template describing how dynamic usernames are generated.

* `connect_timeout` - (Optional) The number of seconds to use as a connection
  timeout.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Redis Configuration Options

* `host` - (Required) The host to connect to.

* `username` - (Required) The username to authenticate with.

* `password` - (Required) The password to authenticate with.

* `port` - (Required) The default port to connect to if no port is specified as
  part of the host.

* `tls` - (Optional) Whether to use TLS when connecting to Redis.

* `insecure_tls` - (Optional) Whether to skip verification of the server
  certificate when using TLS.

* `ca_cert` - (Optional) The contents of a PEM-encoded CA cert file to use to verify the Redis server's identity.

### Redis ElastiCache Configuration Options

* `url` - (Required) The url to connect to including the port; e.g. master.my-cluster.xxxxxx.use1.cache.amazonaws.com:6379.

* `username` - (Optional) The AWS access key id to authenticate with. If omitted Vault tries to infer from the credential provider chain instead.

* `password` - (Optional) The AWS secret access key to authenticate with. If omitted Vault tries to infer from the credential provider chain instead.

* `region` - (Optional) The region where the ElastiCache cluster is hosted. If omitted Vault tries to infer from the environment instead.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### MongoDB Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/mongodb.html#sample-payload)
  for an example.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation.
See the [Vault
  docs](https://www.vaultproject.io/docs/concepts/username-templating)

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### MongoDB Atlas Configuration Options

* `public_key` - (Required) The Public Programmatic API Key used to authenticate with the MongoDB Atlas API.

* `private_key` - (Required) The Private Programmatic API Key used to connect with MongoDB Atlas API.

* `project_id` - (Required) The Project ID the Database User should be created within.

### SAP HanaDB Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/hanadb.html#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `disable_escaping` - (Optional) Disable special character escaping in username and password.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### MSSQL Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/mssql.html#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation.
See the [Vault
  docs](https://www.vaultproject.io/docs/concepts/username-templating)

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `disable_escaping` - (Optional) Disable special character escaping in username and password.

* `contained_db` - (Optional bool: false) For Vault v1.9+. Set to true when the target is a
  Contained Database, e.g. AzureSQL.
  See the [Vault
  docs](https://www.vaultproject.io/api/secret/databases/mssql#contained_db)

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### MySQL Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/mysql-maria.html#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `auth_type` - (Optional) Enable IAM authentication to a Google Cloud instance when set to `gcp_iam`

* `service_account_json` - (Optional) JSON encoding of an IAM access key. Requires `auth_type` to be `gcp_iam`.

* `tls_certificate_key` - (Optional) x509 certificate for connecting to the database. This must be a PEM encoded version of the private key and the certificate combined.

* `tls_ca` - (Optional) x509 CA file for validating the certificate presented by the MySQL server. Must be PEM encoded.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation.
See the [Vault
  docs](https://www.vaultproject.io/docs/concepts/username-templating)

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### PostgreSQL Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/postgresql.html#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `self_managed` - (Optional)  If set, allows onboarding static roles with a rootless
  connection configuration. Mutually exclusive with `username` and `password`.
  If set, will force `verify_connection` to be false. Requires Vault 1.18+ Enterprise.

* `tls_ca` - (Optional) The x509 CA file for validating the certificate
  presented by the PostgreSQL server. Must be PEM encoded.

* `tls_certificate` - (Optional) The x509 client certificate for connecting to
  the database. Must be PEM encoded.

* `password_authentication` - (Optional) When set to `scram-sha-256`, passwords will be
  hashed by Vault before being sent to PostgreSQL. See the [Vault docs](https://www.vaultproject.io/api-docs/secret/databases/postgresql.html#sample-payload)
  for an example. Requires Vault 1.14+.

* `private_key` - (Optional) The secret key used for the x509 client
  certificate. Must be PEM encoded.

* `auth_type` - (Optional) Enable IAM authentication to a Google Cloud instance when set to `gcp_iam`

* `service_account_json` - (Optional) JSON encoding of an IAM access key. Requires `auth_type` to be `gcp_iam`.

* `disable_escaping` - (Optional) Disable special character escaping in username and password.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation.
See the [Vault
  docs](https://www.vaultproject.io/docs/concepts/username-templating)

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Oracle Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/oracle.html#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation.
See the [Vault
  docs](https://www.vaultproject.io/docs/concepts/username-templating)

* `split_statements` - (Optional) Enable spliting statements after semi-colons.

* `disconnect_sessions` - (Optional) Enable the built-in session disconnect mechanism.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Elasticsearch Configuration Options

* `url` - (Required) The URL for Elasticsearch's API. https requires certificate
  by trusted CA if used.

* `username` - (Required) The username to be used in the connection.

* `password` - (Required) The password to be used in the connection.

* `ca_cert` - (Optional) The path to a PEM-encoded CA cert file to use to verify the Elasticsearch server's identity.

* `ca_path` - (Optional) The path to a directory of PEM-encoded CA cert files to use to verify the Elasticsearch server's identity.

* `client_cert` - (Optional) The path to the certificate for the Elasticsearch client to present for communication.

* `client_key` - (Optional) The path to the key for the Elasticsearch client to use for communication.

* `tls_server_name` - (Optional) This, if set, is used to set the SNI host when connecting via TLS.

* `insecure` - (Optional) Whether to disable certificate verification.

* `username_template` - (Optional) For Vault v1.7+. The template to use for username generation. See [Vault docs](https://www.vaultproject.io/docs/concepts/username-templating) for more details.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Snowflake Configuration Options

* `connection_url` - (Required) A URL containing connection information. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/snowflake#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  use.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  maintain.

* `max_connection_lifetime` - (Optional) The maximum number of seconds to keep
  a connection alive for.

* `username` - (Optional) The username to be used in the connection (the account admin level).

* `password` - **Deprecated** (Optional) The password to be used in the connection. Please migrate to key-pair authentication by [November 2025](https://www.snowflake.com/en/blog/blocking-single-factor-password-authentification/).

* `private_key_wo_version` - (Optional)  The version of the `private_key_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

* `username_template` - (Optional) - [Template](https://www.vaultproject.io/docs/concepts/username-templating) describing how dynamic usernames are generated.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

### Redshift Configuration Options

* `connection_url` - (Required) Specifies the Redshift DSN. See
  the [Vault
  docs](https://www.vaultproject.io/api-docs/secret/databases/redshift#sample-payload)
  for an example.

* `max_open_connections` - (Optional) The maximum number of open connections to
  the database.

* `max_idle_connections` - (Optional) The maximum number of idle connections to
  the database.

* `max_connection_lifetime` - (Optional) The maximum amount of time a connection may be reused.

* `username` - (Optional) The root credential username used in the connection URL.

* `password` - (Optional) The root credential password used in the connection URL.

* `disable_escaping` - (Optional) Disable special character escaping in username and password.

* `username_template` - (Optional) - [Template](https://www.vaultproject.io/docs/concepts/username-templating) describing how dynamic usernames are generated.

* `password_wo_version` - (Optional)  The version of the `password_wo`. For more info see [updating write-only attributes](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/guides/using_write_only_attributes.html#updating-write-only-attributes).

## Ephemeral Attributes Reference

The following write-only attributes are supported for all DBs that support username/password:

* `password_wo` - (Optional) The password for the user. Can be updated.
  **Note**: This property is write-only and will not be read from the API.

The following write-only attribute is supported only for Snowflake DB:

* `private_key_wo` - (Optional) The private key associated with the Snowflake user.
  **Note**: This property is write-only and will not be read from the API.

## Attributes Reference

No additional attributes are exported by this resource.

## Import

Database secret backend connections can be imported using the `backend`, `/config/`, and the `name` e.g.

```
$ terraform import vault_database_secret_backend_connection.example postgres/config/postgres
```
