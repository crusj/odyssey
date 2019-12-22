package Service

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func Upload(ctx fasthttp.RequestCtx) {
	file, err := ctx.FormFile("file")
	if err != nil {
		panic(fmt.Sprintf("上传文件失败,%v", err))
	}
	
}
