package main

import (
	"github.com/MaksymBilenko/terraform-provider-aws-synthetics/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: synthetics.Provider})
}
