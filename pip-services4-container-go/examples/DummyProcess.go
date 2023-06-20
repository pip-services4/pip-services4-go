package examples

import "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"

func NewDummyProcess() *container.ProcessContainer {
	c := container.NewProcessContainer("dummy", "Sample dummy process")
	c.SetConfigPath("./examples/dummy.yaml")
	c.AddFactory(NewDummyFactory())
	return c
}
