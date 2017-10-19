package adserver

type AdInfo struct {
	UnitId uint32
	CreativeId uint32
	IconImageUrl string
	MainImageUrl string
	Title string
	Description string
	AppPackageName string
	ClickUrl string
}

type Response struct {
	ResCode int32
	AdList []AdInfo
}
