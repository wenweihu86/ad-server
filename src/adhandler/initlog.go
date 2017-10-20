package adhandler
import (
    "github.com/lestrrat/go-file-rotatelogs"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
    "time"
    "github.com/pkg/errors"
    "path"
)
var adLog=logrus.New()
// config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
    baseLogPaht := path.Join(logPath, logFileName)
    writer, err := rotatelogs.New(
        baseLogPaht+".%Y%m%d%H%M",
        // rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
        rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
        rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
    )
    if err != nil {
        adLog.Errorf("config local file system logger error. %+v", errors.WithStack(err))
    }
    lfHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    })
    adLog.AddHook(lfHook)
}