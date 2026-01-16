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

// Data //

// Always check the current config endpoint after resources have been applied.
// This ensures the output reports the actual current configuration.
data "toadlester_type" "current" {
  depends_on = [toadlester.exp-limit, toadlester.exp-mod, toadlester.exp-size, toadlester.exp-tail,
  toadlester.float-limit, toadlester.float-mod, toadlester.float-size, toadlester.float-tail,
  toadlester.int-limit, toadlester.int-mod, toadlester.int-size, toadlester.int-tail]
}

// Resources //

// Integer Loop //
// Integer range Limit
resource "toadlester" "int-limit" {
  name  = "INT_LIMIT"
  value = "10000"
  algo  = "up"
}

// Integer Limit Modifier
resource "toadlester" "int-mod" {
  name  = "INT_MOD"
  value = "2"
  algo  = "up"
}

// Integer loop Size
resource "toadlester" "int-size" {
  name  = "INT_SIZE"
  value = "100"
  algo  = "up"
}

// Integer Tail is not used
resource "toadlester" "int-tail" {
  name  = "INT_TAIL"
  value = "1"
  algo  = "up"
}

// Float Loop //
// Float range Limit
resource "toadlester" "float-limit" {
  name  = "FLOAT_LIMIT"
  value = "100"
  algo  = "up"
}

// Float Limit Modifier
resource "toadlester" "float-mod" {
  name  = "FLOAT_MOD"
  value = "1.123"
  algo  = "up"
}

// Float loop Size
resource "toadlester" "float-size" {
  name  = "FLOAT_SIZE"
  value = "10"
  algo  = "up"
}

// Float decimal Tail
resource "toadlester" "float-tail" {
  name  = "FLOAT_TAIL"
  value = "5"
  algo  = "up"
}

// Exponent Loop //
// Exponent range Limit
resource "toadlester" "exp-limit" {
  name  = "EXP_LIMIT"
  value = "250"
  algo  = "up"
}

// Exponent Limit Modifier
resource "toadlester" "exp-mod" {
  name  = "EXP_MOD"
  value = "250.43"
  algo  = "up"
}

// Exponent loop Size
resource "toadlester" "exp-size" {
  name  = "EXP_SIZE"
  value = "50"
  algo  = "up"
}

// Exponent decimal Tail
resource "toadlester" "exp-tail" {
  name  = "EXP_TAIL"
  value = "3"
  algo  = "up"
}

// Outputs //

output "toadlester_endpoint" {
  value = data.toadlester_type.current.endpoint
}

output "current_type_settings" {
  value = data.toadlester_type.current.config
}

output "example_resource_output" {
  value = "${toadlester.int-size.name}: ${data.toadlester_type.current.config["INT_SIZE"]}"
}
