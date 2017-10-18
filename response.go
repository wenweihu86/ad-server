package main

type AdInfo struct {
	UnitId uint
	CreativeId uint
	IconImageUrl string
	MainImageUrl string
	Title string
	Description string
	AppPackageName string
	ClickUrl string
}

type Response struct {
	ResCode int
	AdList []AdInfo
}
