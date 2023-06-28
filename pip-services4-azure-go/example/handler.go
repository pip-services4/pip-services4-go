package main

import (
	"context"
	"log"
	"net/http"
	"os"

	aserv "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/controllers"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
)

func main() {
	// create container
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:controller:azurefunc:default:1.0",
		// "service.descriptor", "pip-services-dummies:controller:commandable-azurefunc:default:1.0",
	)

	ctx := cctx.NewContextWithTraceId(context.Background(), "handler.main")

	funcContainer := aserv.NewDummyAzureFunction()
	funcContainer.Configure(ctx, config)
	funcContainer.Open(ctx)

	handler := funcContainer.GetHandler()

	// run server
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/api/HttpTrigger1", handler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
