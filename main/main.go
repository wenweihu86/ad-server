package main

import (
	"github.com/wenweihu86/ad-server/core"
	"os"
	"github.com/valyala/fasthttp"
	"strconv"
)

func main() {
	core.LoadGlobalConf("./conf", "ad_server")
	core.InitLog(core.GlobalConfObject)
	//加载位置字典
	err := core.LocationDict.Load()
	if err != nil {
		os.Exit(-1)
	}
	core.LocationDict.StartReloadTimer()
	
	// 初始化并加载广告信息
	core.AdDictObject = core.NewAdDict(core.GlobalConfObject.AdFileName)
	err = core.AdDictObject.Load()
	if err != nil {
		os.Exit(-1)
	}
	core.AdDictObject.StartReloadTimer()

    requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
			case "/ad/search":
				core.SearchHandler(ctx)
			case "/ad/impression":
				core.ImpressionHandler(ctx)
			case "/ad/click":
				core.ClickHandler(ctx)
			case "/ad/conversion":
				core.ConversionHandler(ctx)
			default:
				ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	    }
    }
    listenPort := ":" + strconv.Itoa(core.GlobalConfObject.AdServerPort)
    fasthttp.ListenAndServe(listenPort, requestHandler)
}
