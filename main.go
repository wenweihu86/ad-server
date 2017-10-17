package main

import (
	"net/http"
	"strconv"
	"log"
	"encoding/json"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	req := new(Request)
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

	log.Printf("appId=%d slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s\n",
		req.AppId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion)

	adInfo := AdInfo{
		AdId: 1,
		CreativeId: 1,
		IconImageUrl: "http://www.baidu.com",
		MainImageUrl: "http://www.weibo.com",
		Title: "title",
		Description: "description",
		AppPackageName: "com.baidu.map",
		ClickUrl: "http://map.baidu.com",
	}
	adList := make([]AdInfo, 0, 1)
	adList = append(adList, adInfo)
	res := Response{
		ResCode: 0,
		AdList: adList,
	}
	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)
}

func main() {
	http.HandleFunc("/ad/search", SearchHandler)
	http.ListenAndServe(":8001", nil)
}
