package main

import (
	"github.com/andrewstucki/terraform-provider-circleci/circleci"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circleci.Provider,
	})
}
