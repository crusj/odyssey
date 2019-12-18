package Router

import (
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestRouter_Register(t *testing.T) {
	defaultRouter.Register([]*Route{
		&Route{
			Method:       "GET",
			Path:         "/test_register",
			HandleFunc:   handleFuncTest,
			PreMiddles:   nil,
			AfterMiddles: nil,
		},
	}...)
	time.Sleep(time.Second * 1)
	if len(defaultRouter.routeTable.tables) != 1 {
		t.Fatal("路由注册失败")
	}
}
func handleFuncTest(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("handleFunc test"))
}
