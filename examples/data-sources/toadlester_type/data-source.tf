# Query toadlester type settings
#
# When using along side Resource updates, check the current
# config endpoint after resources have been applied.
# This ensures the output reports the actual current configuration.
data "toadlester_type" "current" {
  depends_on = [toadlester.exp-limit, toadlester.exp-mod, toadlester.exp-size, toadlester.exp-tail, toadlester.float-limit, toadlester.float-mod, toadlester.float-size, toadlester.float-tail, toadlester.int-limit, toadlester.int-mod, toadlester.int-size, toadlester.int-tail]
}

# Example Data-Source Outputs
# Report back the endpoint used
output "toadlester_endpoint" {
  value = data.toadlester_type.current.endpoint
}

# Print the entire config
output "current_type_settings" {
  value = data.toadlester_type.current.config
}

# Use `config[<SETTING>]` to print specific items
output "example_resource_output" {
  value = "${toadlester.int-size.name}: ${data.toadlester_type.current.config["INT_SIZE"]}"
}
