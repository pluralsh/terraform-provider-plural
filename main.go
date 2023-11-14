package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-plural/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// version is managed by GoReleaser, see: https://goreleaser.com/cookbooks/using-main.version/
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/plural",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
