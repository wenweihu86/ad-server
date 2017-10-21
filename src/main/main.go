package main

import (
	"net/http"
	"adserver"
	"adhandler"
)

func main() {
	adserver.LoadGlobalConf("./conf", "ad_server")
	adserver.LoadLocationDict(
		adserver.GlobalConfObject.GeoBlockFileName,
		adserver.GlobalConfObject.GeoLocationFileName)
	adserver.ReadAdDict(adserver.GlobalConfObject.AdFileName)
	http.HandleFunc("/ad/search", adhandler.SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.ImpressionHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.ListenAndServe(":8001", nil)
}
