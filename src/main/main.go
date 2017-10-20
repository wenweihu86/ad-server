package main

import (
	"net/http"
	"adserver"
	"adhandler"
)

func main() {
	adserver.ReadAdDict()
	http.HandleFunc("/ad/search", adhandler.SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.ImpressionHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.ListenAndServe(":8001", nil)
}
