package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/maroda/terraform-provider-toadlester/toadlester"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: toadlester.Provider,
	})
}
