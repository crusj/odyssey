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
func TestPrintRouteTable(t *testing.T) {
	defaultRouter.PreMiddleware([]Middleware{middleware1, middleware2}...).AfterMiddleware([]Middleware{middleware3}...)
	defaultRouter.Register([]*Route{
		&Route{
			Method:       "GET",
			Path:         "/test_register",
			HandleFunc:   handleFuncTest,
			PreMiddles:   nil,
			AfterMiddles: nil,
		},
		&Route{
			Method:       "GET",
			Path:         "/test_register",
			HandleFunc:   handleFuncTest,
			PreMiddles:   []Middleware{middleware2},
			AfterMiddles: []Middleware{middleware2},
		},
	}...)
	defaultRouter.PrintRouteTable()
}
func handleFuncTest(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("handleFunc test"))
}
func middleware1(handleFunc HandleFunc) HandleFunc {
	return func(ctx *fasthttp.RequestCtx) {
		handleFunc(ctx)
	}
}
func middleware2(handleFunc HandleFunc) HandleFunc {
	return func(ctx *fasthttp.RequestCtx) {
		handleFunc(ctx)
	}
}
func middleware3(handleFunc HandleFunc) HandleFunc {
	return func(ctx *fasthttp.RequestCtx) {
		handleFunc(ctx)
	}
}
