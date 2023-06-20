package examples

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

type DummyController struct {
	timer    *exec.FixedRateTimer
	logger   *log.CompositeLogger
	message  string
	counter1 int
	counter2 int64
}

func NewDummyController() *DummyController {
	c := &DummyController{
		logger:   log.NewCompositeLogger(),
		message:  "Hello World!",
		counter1: 0,
	}

	c.timer = exec.NewFixedRateTimerFromTask(c, 1000, 1000, 5)

	return c
}

func (c *DummyController) Message() string {
	return c.message
}

func (c *DummyController) SetMessage(value string) {
	c.message = value
}

func (c *DummyController) Counter() int {
	return c.counter1
}

func (c *DummyController) SetCounter(value int) {
	c.counter1 = value
}

func (c *DummyController) Configure(ctx context.Context, config *config.ConfigParams) {
	c.message = config.GetAsStringWithDefault("message", c.message)
}

func (c *DummyController) SetReferences(ctx context.Context, references refer.IReferences) {
	c.logger.SetReferences(ctx, references)
}

func (c *DummyController) IsOpen() bool {
	return c.timer.IsStarted()
}

func (c *DummyController) Open(ctx context.Context) error {
	c.timer.Start(ctx)
	c.logger.Trace(ctx, "Dummy controller opened")
	return nil
}

func (c *DummyController) Close(ctx context.Context) error {
	c.timer.Stop(ctx)
	c.logger.Trace(ctx, "Dummy controller closed")
	return nil
}

func (c *DummyController) Notify(ctx context.Context, args *exec.Parameters) {
	go func(c *DummyController) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					msg := convert.StringConverter.ToString(r)
					err = errors.New(msg)
				}
				// Send shutdown signal with err to container
				// and close all components
				cctx.SendShutdownSignalWithErr(ctx, err)
			}
		}()
		atomic.AddInt64(&c.counter2, 1)
	}(c)
	c.logger.Info(ctx, "%d - %s", c.counter1, c.message)
	c.counter1++
}
