package main

import (
	"context"
	"os"

	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/examples"
)

func main() {
	process := examples.NewDummyProcess()
	process.Run(context.Background(), os.Args)
}
