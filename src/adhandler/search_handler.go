package adhandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"math/rand"
	"adserver"
	"strings"
	"fmt"
	"bytes"
	"encoding/base64"
	"github.com/satori/go.uuid"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	req := new(adserver.Request)
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

	// searchId
	req.SearchId = uuid.NewV4().String()

	// search by request ip
	adDict := adserver.AdDict
	var unitIdList1 []uint32
	var exist1 bool
	locationInfo := adserver.SearchLocationByIp(req.Ip)
	if locationInfo != nil {
		country := locationInfo.Country
		city := locationInfo.City
		adserver.AdServerLog.Debug(fmt.Sprintf(
			"ip=%s country=%s city=%s\n", req.Ip, country, city))
		key := strings.ToLower(country) + "_" + strings.ToLower(city)
		unitIdList1, exist1 = adDict.LocationUnitMap[key]
	}
	// search by CN_ALL
	key := "cn_all"
	unitIdList2, exist2 := adDict.LocationUnitMap[key]
	// merge two unit id list
	unitNum := 0
	if exist1 {
		unitNum += len(unitIdList1)
	}
	if exist2 {
		unitNum += len(unitIdList2)
	}
	unitIdList := make([]uint32, unitNum)
	if exist1 && unitIdList1 != nil {
		copy(unitIdList, unitIdList1)
	}
	if exist2 && unitIdList2 != nil {
		copy(unitIdList, unitIdList2)
	}

	// select one from unit id list
	var res = &adserver.Response{}
	if unitIdList != nil {
		unitNum = len(unitIdList)
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
		adInfo.ImpressionTrackUrl = buildImpressionTrackUrl(req, adInfo)
		adInfo.ClickTrackUrl = buildClickTrackUrl(req, adInfo)
		adList := make([]adserver.AdInfo, 0, 1)
		adList = append(adList, adInfo)
		res.ResCode = 0
		res.AdList = adList
		adserver.SearchLog.Info(fmt.Sprintf(
			"searchId=%s slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s " +
			"unitId=%d creativeId=%d\n",
			req.SearchId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion,
			adInfo.UnitId, adInfo.CreativeId))
	} else {
		res.ResCode = 0
		res.AdList = make([]adserver.AdInfo, 0, 1)
		adserver.SearchLog.Info(fmt.Sprintf(
			"searchId=%s slotId=%d adNum=%d iP=%s deviceId=%s oS=%d osVersion=%s resNum=0\n",
			req.SearchId, req.SlotId, req.AdNum, req.Ip, req.DeviceId, req.Os, req.OsVersion))
	}

	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)
}

func buildImpressionTrackUrl(req *adserver.Request, adInfo adserver.AdInfo) string {
	var paramBuf bytes.Buffer
	paramBuf.WriteString(fmt.Sprintf("search_id=%s", req.SearchId))
	paramBuf.WriteString(fmt.Sprintf("&slot_id=%d", req.SlotId))
	paramBuf.WriteString(fmt.Sprintf("&ip=%s", req.Ip))
	paramBuf.WriteString(fmt.Sprintf("&device_id=%s", req.DeviceId))
	paramBuf.WriteString(fmt.Sprintf("&os=%d", req.Os))
	paramBuf.WriteString(fmt.Sprintf("&os_version=%s", req.OsVersion))
	paramBuf.WriteString(fmt.Sprintf("&unit_id=%d", adInfo.UnitId))
	paramBuf.WriteString(fmt.Sprintf("&creative_id=%d", adInfo.CreativeId))
	paramEncoded := base64.StdEncoding.EncodeToString(paramBuf.Bytes())
	impressionTrackUrl := fmt.Sprintf("%s?i=%s",
		adserver.GlobalConfObject.ImpressionTrackUrlPrefix, paramEncoded)
	return impressionTrackUrl
}

func buildClickTrackUrl(req *adserver.Request, adInfo adserver.AdInfo) string {
	var paramBuf bytes.Buffer
	paramBuf.WriteString(fmt.Sprintf("search_id=%s", req.SearchId))
	paramBuf.WriteString(fmt.Sprintf("&slot_id=%d", req.SlotId))
	paramBuf.WriteString(fmt.Sprintf("&ip=%s", req.Ip))
	paramBuf.WriteString(fmt.Sprintf("&device_id=%s", req.DeviceId))
	paramBuf.WriteString(fmt.Sprintf("&os=%d", req.Os))
	paramBuf.WriteString(fmt.Sprintf("&os_version=%s", req.OsVersion))
	paramBuf.WriteString(fmt.Sprintf("&unit_id=%d", adInfo.UnitId))
	paramBuf.WriteString(fmt.Sprintf("&creative_id=%d", adInfo.CreativeId))
	paramBuf.WriteString(fmt.Sprintf("&click_url=%s", adInfo.ClickUrl))
	paramEncoded := base64.StdEncoding.EncodeToString(paramBuf.Bytes())
	impressionTrackUrl := fmt.Sprintf("%s?i=%s",
		adserver.GlobalConfObject.ClickTrackUrlPrefix, paramEncoded)
	return impressionTrackUrl
}

