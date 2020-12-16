package Router

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/crusj/odyssey/Utils"
	"github.com/valyala/fasthttp"
	"sync"
)

var (
	ColorRed     = "\x1b[91m"
	ColorDefault = "\x1b[39m"

	routeTable *RouteTable

	// 用来判断是否存在重复注册情况
	RegisteredRoutes map[string]string

	DefaultRouter *Router
)

func init() {
	RegisteredRoutes = make(map[string]string)
	DefaultRouter = NewRouter()
}

// 路由表类型
type RouteTable struct {
	lock   sync.Mutex
	tables []*table
}

type Router struct {
	// 公共前置中间件
	preMiddleware *middles
	// 公共后置中间件
	afterMiddleware *middles
	// 路由表
	routeTable *RouteTable
	FastRouter *fasthttprouter.Router
}

func NewRouter() *Router {
	return &Router{
		preMiddleware:   &middles{},
		afterMiddleware: &middles{},
		routeTable:      new(RouteTable),
		FastRouter:      fasthttprouter.New(),
	}
}

// 处理器
type HandleFunc = func(ctx *fasthttp.RequestCtx)

// 中间件前置
func (router *Router) PreMiddleware(middleware ...Middleware) *Router {
	router.preMiddleware.lock.Lock()
	defer router.preMiddleware.lock.Unlock()
	router.preMiddleware.middles = middleware

	return router
}

// 中间件后置
func (router *Router) AfterMiddleware(middleware ...Middleware) *Router {
	router.afterMiddleware.lock.Lock()
	defer router.afterMiddleware.lock.Unlock()
	router.afterMiddleware.middles = middleware

	return router
}

// 注册
func (router *Router) Register(routes ...*Route) *Router {
	for _, route := range routes {
		preMiddleNameSet, afterMiddleNameSet := route.ChainFunc()
		isRepeat := router.checkPathIsRepeat(route)
		if isRepeat == false {
			route.RegisterToFastHttp(router.FastRouter)
		}
		// 记录路由表
		router.appendTable(route, isRepeat, preMiddleNameSet, afterMiddleNameSet)
	}

	return router
}

// checkPathIsRepeat 检查注册路由路径是否重复
func (router *Router) checkPathIsRepeat(route *Route) bool {
	method, exists := RegisteredRoutes[route.Path]

	return exists && method == route.Method
}

func (router *Router) appendTable(route *Route, isRepeat bool, routePreNameSet, routeAfterNameSet []string) {
	t := &table{
		method:       route.Method,
		path:         route.Path,
		name:         "",
		handleFunc:   Utils.GetFunctionName(route.HandleFunc),
		preMiddles:   routePreNameSet,
		afterMiddles: routeAfterNameSet,
		repeat:       isRepeat,
	}
	router.routeTable.tables = append(router.routeTable.tables, t)
	RegisteredRoutes[route.Path] = route.Method
}
