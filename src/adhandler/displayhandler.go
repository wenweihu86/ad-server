package adhandler
import (
	"net/http"
	"strconv"
	//"math/rand"
	"adserver"
	"encoding/base64"
	"time"
	"github.com/sirupsen/logrus"
)

//展示handler
func DisplayHandler(w http.ResponseWriter, r *http.Request) {
	ConfigLocalFilesystemLogger("./log/impression","impression.log",time.Hour*24,time.Hour)
	//获得编码后的查询字符串
	queryStringEncoded := r.URL.RawQuery
	//解码
    queryStringDecodedBytes,err := base64.StdEncoding.DecodeString(queryStringEncoded)
    if err != nil{//异常处理
    	adLog.Println(err)
        return
    }
    r.URL.RawQuery = string(queryStringDecodedBytes)
    
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
		os , _ := strconv.ParseUint(r.Form["os"][0], 10, 32)
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
    adLog.WithFields(logrus.Fields{
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
    }).Info("test displayHandler")
}