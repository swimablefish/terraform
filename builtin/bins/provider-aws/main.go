package main

import (
	"github.com/swimablefish/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider,
	})
}
