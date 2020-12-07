package odyssey

import (
	"github.com/crusj/odyssey/Router"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	bodyMsg = []byte("handleFunc test")
)

func TestRun(t *testing.T) {
	Router.DefaultRouter.Register([]*Router.Route{
		&Router.Route{
			Method:       "GET",
			Path:         "/test_run",
			HandleFunc:   handleFuncTest,
			PreMiddles:   nil,
			AfterMiddles: nil,
		},
	}...)
	server := NewServer(Router.DefaultRouter, "8080")
	server.Run()
	time.Sleep(time.Second)
}

func handleFuncTest(ctx *fasthttp.RequestCtx) {
	ctx.Write(bodyMsg)
}

func TestGet(t *testing.T) {
	testFunc(t, "get", nil, nil, "")
}
func TestPost(t *testing.T) {
	testFunc(t, "post", nil, nil, "")
}
func TestPut(t *testing.T) {
	testFunc(t, "put", nil, nil, "")
}
func TestDelete(t *testing.T) {
	testFunc(t, "delete", nil, nil, "")

}

func testFunc(t *testing.T, method string, preMiddleWare []Router.Middleware, afterMiddleWare []Router.Middleware, expect string) {
	Router.DefaultRouter.Register([]*Router.Route{
		&Router.Route{
			Method:       strings.ToUpper(method),
			Path:         "/test_run",
			HandleFunc:   handleFuncTest,
			PreMiddles:   preMiddleWare,
			AfterMiddles: afterMiddleWare,
		},
	}...)
	server := NewServer(Router.DefaultRouter, "8080")
	server.Run()
	time.Sleep(time.Second)

	req, _ := http.NewRequest(strings.ToUpper(method), "http://localhost:8080/test_run", nil)
	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	if expect != "" {
		bodyMsg = []byte(expect)
	}
	assert.Equal(t, bodyMsg, body)
}

// 测试中间件
func TestMiddleware(t *testing.T) {
	testFunc(t, "get",
		[]Router.Middleware{
			middlewareMsgOne,
		},
		[]Router.Middleware{
			middlewareMsgTwo,
		},
		"MsgOne"+string(bodyMsg)+"MsgTwo",
	)
}

func middlewareMsgOne(handleFunc Router.HandleFunc) Router.HandleFunc {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Write([]byte("MsgOne"))
		handleFunc(ctx)
	}
}
func middlewareMsgTwo(handleFunc Router.HandleFunc) Router.HandleFunc {
	return func(ctx *fasthttp.RequestCtx) {
		handleFunc(ctx)
		ctx.Write([]byte("MsgTwo"))
	}
}
