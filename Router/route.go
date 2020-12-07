package Router

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/crusj/odyssey/Utils"
)

type Route struct {
	Method       string
	Path         string
	HandleFunc   HandleFunc
	PreMiddles   []Middleware
	AfterMiddles []Middleware
}

// 串联中间件与Handle
func (r *Route) ChainFunc() (preMiddleNameSet, afterMiddleNameSet []string) {
	// 串联后置中间件
	if count := len(r.AfterMiddles); count > 0 {
		for i := 0; i <= count-1; i++ {
			r.HandleFunc = r.AfterMiddles[i](r.HandleFunc)
			afterMiddleNameSet = append(afterMiddleNameSet, Utils.GetFunctionName(r.AfterMiddles[i]))
		}
	}
	// 串联前置中间件
	if count := len(r.PreMiddles); count > 0 {
		for i := count - 1; i >= 0; i-- {
			r.HandleFunc = r.PreMiddles[i](r.HandleFunc)
			preMiddleNameSet = append(preMiddleNameSet, Utils.GetFunctionName(r.PreMiddles[i]))
		}
	}

	return
}

// RegisterToFastHttp
func (r *Route) RegisterToFastHttp(fastRouter *fasthttprouter.Router) {
	switch r.Method {
	case "GET":
		fastRouter.GET(r.Path, r.HandleFunc)
	case "POST":
		fastRouter.POST(r.Path, r.HandleFunc)
	case "PUT":
		fastRouter.PUT(r.Path, r.HandleFunc)
	case "DELETE":
		fastRouter.DELETE(r.Path, r.HandleFunc)
	default:
		panic(fmt.Sprintf("Bad Request Method: %s", r.Method))
	}
}
