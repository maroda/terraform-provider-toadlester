# Toadlester Provider

[![Tests](https://github.com/maroda/terraform-provider-toadlester/actions/workflows/test.yml/badge.svg)](https://github.com/maroda/terraform-provider-toadlester/actions/workflows/test.yml)
[![Release](https://github.com/maroda/terraform-provider-toadlester/actions/workflows/release.yml/badge.svg)](https://github.com/maroda/terraform-provider-toadlester/actions/workflows/release.yml)

Terraform Provider for the [ToadLester](https://github.com/maroda/toadlester) metrics creation utility.

See [main.tf](main.tf) for the full example configuration.

Quickstart to get `toadlester` running; `base_url=localhost:8899` is the provider configuration.
```shell
$ docker run -d --rm -it --network host --name toadlester ghcr.io/maroda/toadlester:latest
$ curl localhost:8899/metrics
Metric_float_up: 1.6
Metric_float_down: 2.8
Metric_int_down: 24
Metric_int_up: 6
Metric_exp_up: 6.7e+00
Metric_exp_down: 5.2e+00
```

## Which to use
The official one is: `source = "registry.terraform.io/maroda/toadlester"`

For local MacOS development, place this in `~/.terraformrc`:
```shell
provider_installation {
  # This points to the directory that holds the binary.
  dev_overrides {
    "hashicorp.com/edu/toadlester" = "<FULL_PATH_TO>/terraform-provider-toadlester"
  }

  # For all other providers, install them directly as normal.
  direct {}
}
```

And in terraform use: `source = "hashicorp.com/edu/toadlester"`

Running `go build` will output a binary called `terraform-provider-toadlester` by default, which will be used according to the binary location given in `~/.terraformrc`.
