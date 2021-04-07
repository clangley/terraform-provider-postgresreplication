# PostgreSQL logical replication Terraform provider

[![Build Status](https://travis-ci.com/form3tech-oss/terraform-provider-postgresreplication.svg?branch=master)](https://travis-ci.com/form3tech-oss/terraform-provider-postgresreplication)

A Terraform provider for managing PostgreSQL [logical replication slots](https://www.postgresql.org/docs/9.5/logicaldecoding-walsender.html).

## Summary

Manages the lifecycle of logical replication slots. It is useful to manage the lifecycle of replication slots outside of the
service that is consuming them. If the service is shut down you probably want the slot to stay alive so you can
resume streaming when the service is started up again, it is only if a service is decommissioned that you want
the slot to be removed. By defining the slot independently from the service lifecycle you can clean up the slot when the
application is no longer needed.

## Installation

Download the relevant binary from [releases](https://github.com/form3tech-oss/terraform-provider-postgresreplication/releases) and copy it to `$HOME/.terraform.d/plugins/`.

## Configuration

Provider configuration:

```hcl-terraform
provider postgresreplication {
  host     = "localhost"
  post     = 5432
  user     = "superuser"
  password = "superpassword"
}
```

The following provider block variables are available for configuration:

- `host` - The server host to connect to.
- `port` - The server port to connect to.
- `user` - The user to use to connect.
- `password` - The password to use to connect.

## Resources

### postgresreplication_slot

```hcl-terraform
resource "postgresreplication_slot" "test_slot" {
    slot_name 	   = "test_slot"
    output_plugin  = "wal2json"
    database       = "my_db"
}
```

The above would result in a replication created named "test_slot" for the "my_db" database.