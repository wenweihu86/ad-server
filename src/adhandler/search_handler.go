package adhandler

import (
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"math/rand"
	"adserver"
	"strings"
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

	// TODO: read from location dict
	country := "CN"
	city := "ALL"
	key := strings.ToLower(country) + "_" + strings.ToLower(city)
	adDict := adserver.AdDict
	unitIdList, exist := adDict.LocationUnitMap[key]

	var res = &adserver.Response{}
	if exist {
		unitNum := len(unitIdList)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		randIndex := random.Intn(unitNum)
		unitId := unitIdList[randIndex]
		unitInfo := adDict.AdUnitMap[unitId]
		adCreative := adDict.AdCreativeMap[unitInfo.CreativeId]

		adInfo := adserver.AdInfo{
			UnitId: unitInfo.UnitId,
			CreativeId: adCreative.CreativeId,
			Title: adCreative.Title,
			Description: adCreative.Description,
			AppPackageName: adCreative.AppPackageName,
			IconImageUrl: adCreative.IconImageUrl,
			MainImageUrl: adCreative.MainImageUrl,
			ClickUrl: adCreative.ClickUrl,
		}
		adList := make([]adserver.AdInfo, 0, 1)
		adList = append(adList, adInfo)
		res.ResCode = 0
		res.AdList = adList
		log.Printf("appId=%d slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s " +
			"unitId=%d creativeId=%d IconImageUrl=%s ClickUrl=%s\n",
			req.AppId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion,
			adInfo.UnitId, adInfo.CreativeId, adInfo.IconImageUrl, adInfo.ClickUrl)
	} else {
		res.ResCode = 0
		res.AdList = make([]adserver.AdInfo, 0, 1)
		log.Printf("appId=%d slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s resNum=0\n",
			req.AppId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion)
	}

	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)
}
