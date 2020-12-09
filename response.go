package odyssey

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

const (
	ErrorCode       = 400
	SuccessCode     = 200
	SuccessMsg      = "ok"
	ContentTypeJson = "application/json"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Response) Error(ctx *fasthttp.RequestCtx, msg string, code ...int) {
	code = append(code, ErrorCode)
	r.Code = code[0]
	r.Msg = msg

	ctx.SetContentType(ContentTypeJson)

	ret, _ := json.Marshal(r)
	ctx.SetBody(ret)
}

func (r *Response) Success(ctx *fasthttp.RequestCtx, data interface{}) {
	r.Code = SuccessCode
	r.Msg = SuccessMsg
	r.Data = data

	ctx.SetContentType("application/json")

	ret, _ := json.Marshal(r)
	ctx.SetBody(ret)
}
