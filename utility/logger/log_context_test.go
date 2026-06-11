package logger

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

func TestLogContext(t *testing.T) {
	t.Log(4)
	var (
		ctx  = gctx.New()
		now  = gtime.Now()
		lctx = NewLogContext(WithCost(now))
		err  error
	)
	defer lctx.Debug(ctx, "logCtx")
	defer lctx.SetLogErr(&err)
	lctx.SetLogValue("a", 1)
	err = gerror.NewCode(gcode.CodeInternalError)
}

type A struct {
	B *B `json:"b"`
}

type B struct {
	V1 string `json:"v1"`
	V2 int    `json:"v2"`
	C  *C     `json:"c"`
}

type C struct {
	V1 string `json:"v1"`
}

func TestLogGjson(t *testing.T) {
	var (
		ctx  = gctx.New()
		data = &A{
			B: &B{
				V1: "1",
				V2: 2,
				C: &C{
					V1: "2",
				},
			},
		}
		list = []*A{
			data,
		}
		m = g.Map{
			"data": list,
		}
	)
	Debug(ctx, data)
	Debug(ctx, list)
	Debug(ctx, m)
}
