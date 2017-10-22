package main

import (
	"net/http"
	"adserver"
	"adhandler"
)

func main() {
	adserver.LoadGlobalConf("./conf", "ad_server")
	adserver.InitLog(adserver.GlobalConfObject)
	adserver.LoadLocationDict(
		adserver.GlobalConfObject.GeoBlockFileName,
		adserver.GlobalConfObject.GeoLocationFileName)

	// 初始化并加载广告信息
	adserver.AdDictObject = adserver.NewAdDict(adserver.GlobalConfObject.AdFileName)
	adserver.AdDictObject.Load()
	adserver.AdDictObject.StartReloadTimer()

	http.HandleFunc("/ad/search", adhandler.SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.ImpressionHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.HandleFunc("/ad/conversion",adhandler.ConversionHandler)
	http.ListenAndServe(":8001", nil)
}
