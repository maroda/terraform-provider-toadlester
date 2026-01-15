terraform {
  required_providers {
    toadlester = {
      // source = "registry.terraform.io/maroda/toadlester"
      // This is a dev override configured in ~/.terraformrc - `init` is not used, go straight to `apply`
      source = "hashicorp.com/edu/toadlester"
    }
  }
}

provider "toadlester" {
  base_url = "http://grogu:8899"
}

data "toadlester_type" "current" {
  name = "INT_SIZE"
  algo = "up"
}

resource "toadlester" "configure" {
  name  = "INT_SIZE"
  value = "100"
  algo  = "up"
}

output "endpoint" {
  value = data.toadlester_type.current.endpoint
}
output "name" {
  value = data.toadlester_type.current.name
}

