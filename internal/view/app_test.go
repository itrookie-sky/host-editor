package view

import (
	"context"
	"host-editor/utility/logger"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestCancel(t *testing.T) {

	parent := gctx.New()
	ctx, cancel := context.WithCancel(parent)
	sun := gctx.WithSpan(ctx, "sun")
	go func() {
		<-parent.Done()
		logger.Debug(parent, "parent")
	}()
	go func() {
		<-sun.Done()
		logger.Debug(sun, "sun")
	}()
	go func() {
		<-ctx.Done()
		logger.Debug(ctx, "ctx")
	}()
	time.Sleep(time.Millisecond * 50)
	cancel()
	time.Sleep(time.Millisecond * 100)
}
