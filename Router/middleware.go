package Router

import "sync"

// 中间件集合类型
type middles struct {
	lock    sync.Mutex
	middles []Middleware
}

// 中间件类型
type Middleware func(handleFunc HandleFunc) HandleFunc
