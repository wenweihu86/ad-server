package utils
import(
  "os"
  "strconv"
  "time"
)
func CheckFileIsExist(filename string) (bool) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false;
	}
	return true;
}
func GetLogFileName(filename,logpath string) string {
	dateStr := strconv.Itoa(time.Now().Year()) + strconv.Itoa(int(time.Now().Month())) + strconv.Itoa(time.Now().Day())+strconv.Itoa(time.Now().Hour())
    var logFileName string
    if filename == "impression"{
    	logFileName = "impression." + dateStr + ".log"
    }else if filename == "click"{
    	logFileName = "click." + dateStr + ".log"
	}else{
		logFileName = "default." + dateStr + ".log"
	}
    logFile := logpath + logFileName
    return logFile
}
