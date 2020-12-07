package odyssey

import (
	"fmt"
	"github.com/crusj/odyssey/Router"
	"github.com/valyala/fasthttp"
)

type server struct {
	port   string
	router *Router.Router
}

func NewServer(router *Router.Router, port string) *server {
	return &server{
		port:   port,
		router: router,
	}
}

// 运行http server
func (s *server) Run() {
	go func() {
		err := fasthttp.ListenAndServe(":"+s.port, s.router.FastRouter.Handler)
		if err != nil {
			panic(fmt.Sprintf("fasthttpserver run error! %v", err))
		}
	}()
	s.router.PrintRouteTable()
}
