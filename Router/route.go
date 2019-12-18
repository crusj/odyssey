package Router

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"reflect"
	"runtime"
	"sync"
)

type Router struct {
	preMiddleware   *middles
	afterMiddleware *middles
	//路由表
	routeTable *RouteTable
}
type table struct {
	//Method
	method string
	//路径
	path string
	//路由命名
	name string
	//handleFunc
	handleFunc string
	//后置中间件
	preMiddles []string
	//后置中间件
	afterMiddles []string
}

//路由表类型
type RouteTable struct {
	lock   sync.Mutex
	tables []*table
}
type Route struct {
	Method       string
	Path         string
	HandleFunc   HandleFunc
	PreMiddles   []Middleware
	AfterMiddles []Middleware
}

//请求处理类型
type HandleFunc = func(ctx *fasthttp.RequestCtx)

//中间件类型
type Middleware func(handleFunc HandleFunc) HandleFunc

//中间件集合类型
type middles struct {
	lock    sync.Mutex
	middles []Middleware
}

var (
	routeTable    *RouteTable
	fastRouter    *fasthttprouter.Router
	defaultRouter *Router
)

//中间件前置
func (router *Router) PreMiddleware(middleware ...Middleware) *Router {
	router.preMiddleware.lock.Lock()
	defer router.preMiddleware.lock.Unlock()
	router.preMiddleware.middles = middleware
	return router
}

//中间件后置
func (router *Router) AfterMiddleware(middleware ...Middleware) *Router {
	router.afterMiddleware.lock.Lock()
	defer router.preMiddleware.lock.Unlock()
	router.afterMiddleware.middles = middleware
	return router
}

//注册
func (router *Router) Register(routes ...*Route) *Router {
	var pre, after []Middleware
	//记录路由表
	router.routeTable.lock.Lock()
	defer router.routeTable.lock.Unlock()
	for _, route := range routes {
		var preName, afterName []string
		pre = combineMiddles(append(route.PreMiddles, router.preMiddleware.middles...)...)
		after = combineMiddles(append(route.AfterMiddles, router.afterMiddleware.middles...)...)
		handleFunc := route.HandleFunc
		preCount, afterCount := len(pre), len(after)
		//串联后置中间件
		if afterCount > 0 {
			for i := 0; i <= afterCount-1; i++ {
				handleFunc = after[i](handleFunc)
				afterName = append(afterName, GetFunctionName(after[i]))
			}
		}
		//串联前置中间件
		if preCount > 0 {
			for i := preCount - 1; i >= 0; i-- {
				handleFunc = pre[i](handleFunc)
				preName = append(preName, GetFunctionName(pre[i]))
			}
		}
		switch route.Method {
		case "GET":
			fastRouter.GET(route.Path, handleFunc)
		case "POST":
			fastRouter.POST(route.Path, handleFunc)
		case "PUT":
			fastRouter.PUT(route.Path, handleFunc)
		case "DELETE":
			fastRouter.DELETE(route.Path, handleFunc)
		default:
			panic("")
		}
		t := &table{
			method:       route.Method,
			path:         route.Path,
			name:         "",
			handleFunc:   GetFunctionName(route.HandleFunc),
			preMiddles:   preName,
			afterMiddles: afterName,
		}
		router.routeTable.tables = append(router.routeTable.tables, t)
	}
	return router
}

//运行http server
func (router *Router) Run(listenPort string) {
	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%v", listenPort), fastRouter.Handler)
		if err != nil {
			panic(fmt.Sprintf("fasthttpserver run error! %v", err))
		}
	}()
}
func init() {
	fastRouter = fasthttprouter.New()
	defaultRouter = &Router{
		preMiddleware:   &middles{},
		afterMiddleware: &middles{},
		routeTable:      new(RouteTable),
	}
}

//合并中间件，去掉重复的
func combineMiddles(middles ...Middleware) (combined []Middleware) {
	m := make(map[string]Middleware)
	for _, middle := range middles {
		address := fmt.Sprintf("%v", middle)
		if _, ok := m[address]; !ok {
			m[address] = middle
			combined = append(combined, middle)
		}
	}
	return
}

//路由表
func ListRoutes() {

}

//获取函数名
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
