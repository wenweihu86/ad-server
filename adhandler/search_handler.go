package adhandler

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"time"
	"math/rand"
	"github.com/wenweihu86/ad-server/adserver"
	"strings"
	"fmt"
	"bytes"
	"encoding/base64"
	"github.com/satori/go.uuid"
)

func SearchHandler(ctx *fasthttp.RequestCtx) {
	beginTime := time.Now().Nanosecond()
	args := ctx.QueryArgs()
	req := new(adserver.Request)
	// slot_id
	slotId, err := args.GetUint("slot_id")
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return
	}
	req.SlotId = uint32(slotId)

	// ad_num
	reqAdNum, err := args.GetUint("ad_num")
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return
	}
	req.ReqAdNum = uint32(reqAdNum)

	// ip
	if args.Has("ip") {
		req.Ip = string(args.Peek("ip"))
	}

	// device_id
	if args.Has("device_id") {
		req.DeviceId = string(args.Peek("device_id"))
	}

	// os
	os, err := args.GetUint("os")
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return
	}
	req.Os = uint32(os)

	// os_version
	if args.Has("os_version") {
		req.OsVersion = string(args.Peek("os_version"))
	}

	// searchId
	req.SearchId = uuid.NewV4().String()

	adData := adserver.AdDictObject.GetCurrentAdData()
	// search by request ip
	var unitIdList1 []uint32
	var exist1 bool
	ipDataInfo := adserver.LocationDict.GetCurrentIpData()
	locationInfo := ipDataInfo.SearchLocationByIp(req.Ip)
	if locationInfo != nil {
		country := locationInfo.Country
		city := locationInfo.City
		adserver.AdServerLog.Debug(fmt.Sprintf(
			"ip=%s country=%s city=%s\n", req.Ip, country, city))
		key := strings.ToLower(country) + "_" + strings.ToLower(city)
		unitIdList1, exist1 = adData.LocationUnitMap[key]
	}
	// search by CN_ALL
	key := "cn_all"
	unitIdList2, exist2 := adData.LocationUnitMap[key]
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
	adList := make([]adserver.AdInfo, 0, 1)
	unitIdMap := make(map[int] bool)
	var unitIdsStr, creativeIdsStr string
	resAdNum := 0 
	if unitIdList != nil && req.ReqAdNum >= 1 {
		unitNum = len(unitIdList)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < unitNum && i < int(req.ReqAdNum); i++ {
			randIndex := random.Intn(unitNum)
			if unitIdMap[randIndex] {
				i--
				continue
			}
			resAdNum++
			unitIdMap[randIndex] = true
			unitId := unitIdList[randIndex]
			unitInfo := adData.AdUnitMap[unitId]
			adCreative := adData.AdCreativeMap[unitInfo.CreativeId]
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
			adInfo.ConversionTrackUrl = buildConversionTrackUrl(req, adInfo)
			adList = append(adList, adInfo)
			if i == unitNum - 1 || i == int(req.ReqAdNum) - 1 {
				unitIdsStr += fmt.Sprint(adInfo.UnitId)
				creativeIdsStr += fmt.Sprint(adInfo.CreativeId)
			} else {
				unitIdsStr += fmt.Sprint(adInfo.UnitId) + ","
				creativeIdsStr += fmt.Sprint(adInfo.CreativeId) + ","
			}
		}
		res.ResCode = 0
		res.AdList = adList
	} else {
		res.ResCode = 0
		res.AdList = make([]adserver.AdInfo, 0, 1)
	}
	elapseTimeMs := (time.Now().Nanosecond() - beginTime) / 1000000
    adserver.SearchLog.Info(fmt.Sprintf(
			"searchId=%s elpaseMs=%d slotId=%d reqAdNum=%d resAdNum=%d " +
				"iP=%s deviceId=%s oS=%d osVersion=%s unitId=%s creativeId=%s\n",
			req.SearchId, elapseTimeMs, req.SlotId, req.ReqAdNum, resAdNum,
			req.Ip, req.DeviceId, req.Os, req.OsVersion, unitIdsStr, creativeIdsStr))

	resBytes, err := json.Marshal(res)
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return 
	}
	ctx.SetBody(resBytes)
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

func buildConversionTrackUrl(req *adserver.Request, adInfo adserver.AdInfo) string {
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
	conversionTrackUrl := fmt.Sprintf("%s?i=%s",
		adserver.GlobalConfObject.ConversionTrackUrlPrefix, paramEncoded)
	return conversionTrackUrl
}

