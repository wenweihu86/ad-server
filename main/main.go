package main

import (
	"net/http"
	"github.com/wenweihu86/ad-server/adserver"
	"github.com/wenweihu86/ad-server/adhandler"
	"strconv"
	"os"
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

	http.HandleFunc("/ad/search", adhandler.SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.ImpressionHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.HandleFunc("/ad/conversion",adhandler.ConversionHandler)
	listenPort := ":" + strconv.Itoa(adserver.GlobalConfObject.AdServerPort)
	http.ListenAndServe(listenPort, nil)
}
