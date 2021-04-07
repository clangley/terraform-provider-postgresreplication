provider "postgresreplication" {
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "postgres"
  sslmode  = "disable"
}
resource "postgresreplication_replica_identity" "test" {
  table_name = "foo"
  database = "postgres"
}

terraform {
  required_providers {

    postgresreplication = {
      source  = "clangley/postgresreplication"
      version = "0.8.0"
    }
  }
}
