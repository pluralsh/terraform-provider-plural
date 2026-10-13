package main

import (
	"flag"
	"log"

	"terraform-provider-plural/internal/provider"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

var (
	// version is managed by GoReleaser, see: https://goreleaser.com/cookbooks/using-main.version/
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	var opts []tf6server.ServeOpt
	if debug {
		opts = append(opts, tf6server.WithManagedDebug())
	}

	err := tf6server.Serve("registry.terraform.io/pluralsh/plural", provider.NewServer(version), opts...)

	if err != nil {
		log.Fatal(err.Error())
	}
}
