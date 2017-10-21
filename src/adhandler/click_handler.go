package adhandler

import (
	"net/http"
	"strconv"
	"adserver"
    "encoding/base64"
	"fmt"
)

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	//获得编码后的查询字符串
	queryStringEncoded := r.URL.RawQuery
	//解码
    queryStringDecodedBytes , err := base64.StdEncoding.DecodeString(queryStringEncoded)
    if err != nil {
        return
    }
    r.URL.RawQuery = string(queryStringDecodedBytes)

	r.ParseForm()
	req := new(adserver.Request)
	// app_id
	if len(r.Form["app_id"]) > 0 {
		appId, _ := strconv.ParseUint(r.Form["app_id"][0], 10, 32)
		req.AppId = uint32(appId)
	}
	// slot_id
	if len(r.Form["slot_id"]) > 0 {
		slotId, _ := strconv.ParseUint(r.Form["slot_id"][0], 10, 32)
		req.SlotId = uint32(slotId)
	}
	// ad_num
	if len(r.Form["ad_num"]) > 0 {
		adNum, _ := strconv.ParseUint(r.Form["ad_num"][0], 10, 32)
		req.AdNum = uint32(adNum)
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
		req.Os = uint32(os)
	}
	// os_version
	if len(r.Form["os_version"]) > 0 {

		req.OsVersion = r.Form["os_version"][0]
	}
	//unit_id
	if len(r.Form["unit_id"]) > 0 {
		unit, _ := strconv.ParseUint(r.Form["unit_id"][0], 10, 32)
		req.UnitId = uint32(unit)
	}
	//creative_id
	if len(r.Form["creative_id"]) > 0 {
		creativeId, _ := strconv.ParseUint(r.Form["creative_id"][0], 10, 32)
		req.CreativeId = uint32(creativeId)
	}
	//search_id
	if len(r.Form["search_id"]) > 0 {
		req.SearchId = r.Form["search_id"][0]
	}
	//click_url
	if len(r.Form["click_url"]) > 0 {
		req.ClickUrl = r.Form["click_url"][0]
	}
	adserver.ClickLog.Info(fmt.Sprintf(
		"searchId=%s slotId=%d ip=%s os=%d unit_id=%d creativeId=%d",
		req.SearchId, req.SlotId, req.Ip, req.Os, req.UnitId, req.CreativeId))
	// TODO: 302跳转到click url
}