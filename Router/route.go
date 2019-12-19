package Router

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"reflect"
	"runtime"
	"strings"
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
	//是否重复注册
	repeat bool
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
	ColorRed     = "\x1b[91m"
	ColorDefault = "\x1b[39m"

	routeTable *RouteTable
	//用来判断是否存在重复注册情况
	shorRouteTable map[string]string

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
	defer router.afterMiddleware.lock.Unlock()
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
		var repeat bool
		if method, ok := shorRouteTable[route.Path]; ok {
			if method == route.Method {
				repeat = true
			}
		}
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
		if repeat == false {
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
		}
		t := &table{
			method:       route.Method,
			path:         route.Path,
			name:         "",
			handleFunc:   GetFunctionName(route.HandleFunc),
			preMiddles:   preName,
			afterMiddles: afterName,
			repeat:       repeat,
		}
		router.routeTable.tables = append(router.routeTable.tables, t)
		shorRouteTable[route.Path] = route.Method
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
	shorRouteTable = make(map[string]string)
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

//终端打印已注册的路由表
func (router *Router) PrintRouteTable() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{
				Align: simpletable.AlignLeft,
				Text:  "Method",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "Path",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "Name",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "HandleFunc",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "Pre-middleware",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "After-middleware",
			},
		},
	}
	for _, v := range router.routeTable.tables {
		var method, path string = v.method, v.path
		if v.repeat == true {
			method = red(method)
			path = red(path)
		}
		row := []*simpletable.Cell{
			{
				Align: simpletable.AlignLeft,
				Text:  method,
			},
			{
				Align: simpletable.AlignLeft,
				Text:  path,
			},
			{
				Align: simpletable.AlignLeft,
				Text:  v.name,
			},
			{
				Align: simpletable.AlignLeft,
				Text:  v.handleFunc,
			},
			{
				Align: simpletable.AlignLeft,
				Text:  strings.Join(v.preMiddles, "\r\n"),
			},
			{
				Align: simpletable.AlignLeft,
				Text:  strings.Join(v.afterMiddles, "\r\n"),
			},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	table.SetStyle(simpletable.StyleMarkdown)
	fmt.Println(table)
}

//获取函数名
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}
