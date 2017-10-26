package adserver

type Request struct {
	SearchId string
	SlotId uint32
	ReqAdNum uint32
	Ip string
	DeviceId string
	Os uint32 // 0:android, 1:ios
	OsVersion string
	UnitId uint32
	CreativeId uint32
	ClickUrl string
}
