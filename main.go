package main

import (
	"context"
	"flag"
	"log"

	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: fusionauth.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/gpsinsight/fusionauth", opts) //nolint:staticcheck
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
