package test

import (
	"context"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

type Shutdown struct {
	started    bool
	runCode    int
	mode       string
	minTimeout int
	maxTimeout int
}

func NewShutdown() *Shutdown {
	return &Shutdown{
		started:    false,
		runCode:    0,
		mode:       "exception",
		minTimeout: 300000,
		maxTimeout: 900000,
	}
}

func (c *Shutdown) Configure(ctx context.Context, config config.ConfigParams) {
	c.mode = config.GetAsStringWithDefault("mode", c.mode)
	c.minTimeout = config.GetAsIntegerWithDefault("min_timeout", c.minTimeout)
	c.maxTimeout = config.GetAsIntegerWithDefault("max_timeout", c.maxTimeout)
}

func (c *Shutdown) IsOpen() bool {
	return c.started
}

func (c *Shutdown) Open(ctx context.Context) error {
	if c.started {
		return nil
	}

	delay := int(float32(c.maxTimeout-c.minTimeout)*rand.Float32() + float32(c.minTimeout))
	c.runCode++
	go c.doShutdown(ctx, delay, c.runCode)
	c.started = true

	return nil
}

func (c *Shutdown) Close(ctx context.Context) error {
	// Todo: Properly interrupt the go proc
	c.started = false
	return nil
}

func (c *Shutdown) Shutdown(ctx context.Context) {
	if c.mode == "null" || c.mode == "nullpointer" {
		var obj io.Writer
		obj.Write([]byte{})
	} else if c.mode == "zero" || c.mode == "dividebyzero" {
		_ = 0 / 100
	} else if c.mode == "exit" || c.mode == "processexit" {
		os.Exit(1)
	} else {
		panic("Crash test exception")
	}
}

func (c *Shutdown) doShutdown(ctx context.Context, delay int, runCode int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	if c.started && c.runCode == runCode {
		c.Shutdown(ctx)
	}
}
