package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/swimablefish/terraform/builtin/providers/aws"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider,
	})
}
