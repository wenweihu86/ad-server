package main

import (
	"net/http"
	"strconv"
	"log"
	"encoding/json"
	"time"
	"math/rand"

	"adserver"
	"adhandler"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

	unitNum := len(adserver.AdUnits)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randIndex := random.Intn(unitNum)
	adUnit := adserver.AdUnits[randIndex]
	adCreative := adserver.AdCreativeMap[adUnit.CreativeId]

	adInfo := adserver.AdInfo{
		UnitId: adUnit.UnitId,
		CreativeId: adCreative.CreativeId,
		IconImageUrl: adCreative.IconImageUrl,
		MainImageUrl: "",
		Title: adCreative.Title,
		Description: "",
		AppPackageName: "",
		ClickUrl: adCreative.ClickUrl,
	}
	adList := make([]adserver.AdInfo, 0, 1)
	adList = append(adList, adInfo)
	res := adserver.Response{
		ResCode: 0,
		AdList: adList,
	}
	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)

	log.Printf("appId=%d slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s " +
		"unitId=%d creativeId=%d IconImageUrl=%s ClickUrl=%s\n",
		req.AppId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion,
		adInfo.UnitId, adInfo.CreativeId, adInfo.IconImageUrl, adInfo.ClickUrl)
}

func main() {
	adserver.ReadAdDict()
	http.HandleFunc("/ad/search", SearchHandler)
	http.HandleFunc("/ad/impression",adhandler.DisplayHandler)
	http.HandleFunc("/ad/click",adhandler.ClickHandler)
	http.ListenAndServe(":8001", nil)
}
