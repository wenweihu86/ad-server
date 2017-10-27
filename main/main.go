package main

import (
	"github.com/wenweihu86/ad-server/adserver"
	"github.com/wenweihu86/ad-server/adhandler"
	"os"
	"github.com/valyala/fasthttp"
	"strconv"
)

func main() {
	adserver.LoadGlobalConf("./conf", "ad_server")
	adserver.InitLog(adserver.GlobalConfObject)
	//加载位置字典
	err := adserver.LocationDict.Load()
	if err != nil {
		os.Exit(-1)
	}
	adserver.LocationDict.StartReloadTimer()
	
	// 初始化并加载广告信息
	adserver.AdDictObject = adserver.NewAdDict(adserver.GlobalConfObject.AdFileName)
	err = adserver.AdDictObject.Load()
	if err != nil {
		os.Exit(-1)
	}
	adserver.AdDictObject.StartReloadTimer()

    requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
			case "/ad/search":
				adhandler.SearchHandler(ctx)
			case "/ad/impression":
				adhandler.ImpressionHandler(ctx)
			case "/ad/click":
				adhandler.ClickHandler(ctx)
			case "/ad/conversion":
				adhandler.ConversionHandler(ctx)
			default:
				ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	    }
    }
    listenPort := ":" + strconv.Itoa(adserver.GlobalConfObject.AdServerPort)
    fasthttp.ListenAndServe(listenPort, requestHandler)
}
