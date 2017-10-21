package adserver

import "github.com/ibbd-dev/go-async-log"

var AdServerLog *asyncLog.LogFile
var SearchLog *asyncLog.LogFile
var ImpressionLog *asyncLog.LogFile
var ClickLog *asyncLog.LogFile

/*
*  ibbd-dev/go-async-log 日志框架说明
*  1. 自动切割周期：默认按小时
*  2. 默认全部写入文件
*  3. 批量写入周期：默认每秒写入一次
*  4. 是否需要Flags：默认需要
*/
func InitLog(globalConfObject *GlobalConf) {
	AdServerLog = asyncLog.NewLevelLog(globalConfObject.AdServerLogFileName,
		asyncLog.Priority(globalConfObject.LogLevel))
	SearchLog = asyncLog.NewLevelLog(globalConfObject.SearchLogFileName,
		asyncLog.LevelInfo)
	ImpressionLog = asyncLog.NewLevelLog(globalConfObject.ImpressionLogFileName,
		asyncLog.LevelInfo)
	ClickLog = asyncLog.NewLevelLog(globalConfObject.ClickLogFileName,
		asyncLog.LevelInfo)
}
