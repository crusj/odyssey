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
	ctx.SetBody(bodyMsg)
}

func TestGet(t *testing.T) {
	testFunc(t, "get")
}
func TestPost(t *testing.T) {
	testFunc(t, "post")
}
func TestPut(t *testing.T) {
	testFunc(t, "put")
}
func TestDelete(t *testing.T) {
	testFunc(t, "delete")
}

func testFunc(t *testing.T, method string) {
	Router.DefaultRouter.Register([]*Router.Route{
		&Router.Route{
			Method:       strings.ToUpper(method),
			Path:         "/test_run",
			HandleFunc:   handleFuncTest,
			PreMiddles:   nil,
			AfterMiddles: nil,
		},
	}...)
	server := NewServer(Router.DefaultRouter, "8080")
	server.Run()
	time.Sleep(time.Second)

	req, _ := http.NewRequest(strings.ToUpper(method), "http://localhost:8080/test_run", nil)
	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, bodyMsg, body)
}
