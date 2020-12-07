package Router

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"strings"
)

type table struct {
	// Method
	method string
	// 路径
	path string
	// 路由命名
	name string
	// handleFunc
	handleFunc string
	// 后置中间件
	preMiddles []string
	// 后置中间件
	afterMiddles []string
	// 是否重复注册
	repeat bool
}

// 终端打印已注册的路由表
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
			method = beRedStr(method)
			path = beRedStr(path)
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

func beRedStr(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}
