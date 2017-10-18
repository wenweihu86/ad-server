package adserver

type Request struct {
	AppId uint
	SlotId uint
	AdNum uint
	Ip string
	DeviceId string
	Os uint // 0:android, 1:ios
	OsVersion string
}
