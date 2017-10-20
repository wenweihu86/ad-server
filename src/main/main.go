package main

import (
	"net/http"
	"adserver"
	"adhandler"
)

func main() {
	adserver.LoadLocationDict(
		"./data/GeoLiteCity-Blocks.csv",
		"./data/GeoLiteCity-Location.csv")
	adserver.ReadAdDict("./data/ad_info.txt")
	http.HandleFunc("/ad/search", adhandler.SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.ImpressionHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.ListenAndServe(":8001", nil)
}
