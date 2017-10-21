package adserver

type Request struct {
	AppId uint32
	SlotId uint32
	AdNum uint32
	Ip string
	DeviceId string
	Os uint32 // 0:android, 1:ios
	OsVersion string
	UnitId uint32
	CreativeId uint32
	SearchId string
	ClickUrl string
}
