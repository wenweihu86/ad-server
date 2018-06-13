package core

import "github.com/ibbd-dev/go-async-log"

var AdServerLog *asyncLog.LogFile // 除了检索、展现、点击监控外，所以日志都用这个log打印
var SearchLog *asyncLog.LogFile // 每个检索请求打印一行日志
var ImpressionLog *asyncLog.LogFile // 每个展现请求打印一行日志
var ClickLog *asyncLog.LogFile // 每个点击请求打印一行日志
var ConversionLog *asyncLog.LogFile // 每个转化请求打印一行日志

/*
*  ibbd-dev/go-async-log 日志框架说明
*  1. 自动切割周期：默认按小时
*  2. 默认全部写入文件
*  3. 批量写入周期：默认每秒写入一次
*  4. 是否需要Flags：默认需要
*/
// TODO: 让日志支持format打印
func InitLog(globalConfObject *GlobalConf) {
	AdServerLog = asyncLog.NewLevelLog(globalConfObject.AdServerLogFileName,
		asyncLog.Priority(globalConfObject.LogLevel))
	SearchLog = asyncLog.NewLevelLog(globalConfObject.SearchLogFileName,
		asyncLog.LevelInfo)
	ImpressionLog = asyncLog.NewLevelLog(globalConfObject.ImpressionLogFileName,
		asyncLog.LevelInfo)
	ClickLog = asyncLog.NewLevelLog(globalConfObject.ClickLogFileName,
		asyncLog.LevelInfo)
	ConversionLog = asyncLog.NewLevelLog(globalConfObject.ConversionLogFileName,
		asyncLog.LevelInfo)
}
