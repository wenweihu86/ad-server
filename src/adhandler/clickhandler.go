package adhandler
import (
	"net/http"
	"strconv"
	
	"encoding/json"
	"time"
	//"math/rand"
	"ad-server/src/adserver"
	"fmt"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"ad-server/src/utils"
	"os"
)
var adlog = logrus.New()
var logpath="./data/logfile/"
func ClickHandler(w http.ResponseWriter, r *http.Request) {
	//获得编码后的查询字符串
	queryStringEncoded:=r.URL.RawQuery
	//解码
    queryStringDecodedBytes,err:=base64.StdEncoding.DecodeString(queryStringEncoded)
    if err!=nil{//异常处理
    	fmt.Println("error:",err)
        res := adserver.Response{
			ResCode: 4,
			AdList: nil,
		}
		resBytes, _ := json.Marshal(res)
		w.Write(resBytes)
    }
    r.URL.RawQuery=string(queryStringDecodedBytes)

	r.ParseForm()
	req := new(adserver.Request)
	// app_id
	if len(r.Form["app_id"]) > 0 {
		appId, _ := strconv.ParseUint(r.Form["app_id"][0], 10, 32)
		req.AppId = uint(appId)
	}
	// slot_id
	if len(r.Form["slot_id"]) > 0 {
		slotId, _ := strconv.ParseUint(r.Form["slot_id"][0], 10, 32)
		req.SlotId = uint(slotId)
	}
	// ad_num
	if len(r.Form["ad_num"]) > 0 {
		adNum, _ := strconv.ParseUint(r.Form["ad_num"][0], 10, 32)
		req.AdNum = uint(adNum)
	}
	// ip
	if len(r.Form["ip"]) > 0 {
		req.Ip = r.Form["ip"][0]
	}
	// device_id
	if len(r.Form["device_id"]) > 0 {
		req.DeviceId = r.Form["device_id"][0]
	}
	// os
	if len(r.Form["os"]) > 0 {
		os, _ := strconv.ParseUint(r.Form["os"][0], 10, 32)
		req.Os = uint(os)
	}
	// os_version
	if len(r.Form["os_version"]) > 0 {

		req.OsVersion = r.Form["os_version"][0]
	}
	//unit_id
	if len(r.Form["unit_id"]) > 0 {
		req.UnitId = r.Form["unit_id"][0]
	}
	//creative_id
	if len(r.Form["creative_id"]) > 0 {
		req.CreativeId = r.Form["creative_id"][0]
	}
	//search_id
	if len(r.Form["search_id"]) > 0 {
		req.SearchId = r.Form["search_id"][0]
	}
	//click_url
	if len(r.Form["click_url"]) > 0 {
		req.ClickUrl = r.Form["click_url"][0]
	}
    //log设置输出
    adlog.Out = os.Stdout
    dateStr:=strconv.Itoa(time.Now().Year())+strconv.Itoa(int(time.Now().Month()))+strconv.Itoa(time.Now().Day())
    logFileName:=dateStr+"log.log"
    logFile:=logpath+logFileName
    fmt.Println(logFile)
    //判断日志文件是否存在
    if !utils.CheckFileIsExist(logFile){
    	_,err:= os.Create(logFile) 
		if err!=nil{
		   fmt.Println(err)
		   res := adserver.Response{
				ResCode: 4,
				AdList: nil,
			}
			resBytes, _ := json.Marshal(res)
			w.Write(resBytes)
		}
    }
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
	    adlog.Out = file
	} else {
	    adlog.Info("Failed to log to file, using default stderr")
	    fmt.Println(err)
	}
    
    adlog.WithFields(logrus.Fields{
	    "appId": req.AppId,
	    "slotId":  req.SlotId,
	    "adNum":req.AdNum,
	    "iP":req.Ip,
        "deviceId":req.DeviceId,
        "oS":req.Os,
        "osVersion":req.OsVersion,
        "unitId":req.UnitId,
        "creativeId":req.CreativeId,
        "searchId":req.SearchId,
        "clickUrl":req.ClickUrl,
    }).Info("test")
}