package adhandler

import (
	"net/http"
	"strconv"
	"adserver"
	"encoding/base64"
	"fmt"
	"net/url"
)

// 展现监控handler
func ImpressionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form["i"]) == 0 {
		w.Write([]byte("{\"status\": 1}"))
		return
	}
	i := r.Form["i"][0]
	queryStringBytes, err := base64.StdEncoding.DecodeString(i)
	if err != nil {
		w.Write([]byte("{\"status\": 1}"))
		return
	}
	queryString := string(queryStringBytes)
	paramMap, _ := url.ParseQuery(queryString)

	// search_id
	var searchId string
	if searchIds, exist := paramMap["search_id"]; exist {
		searchId = searchIds[0]
	}

	// slot_id
	var slotId uint32
	if slotIds, exist := paramMap["slot_id"]; exist {
		tmpInt, _ := strconv.ParseUint(slotIds[0], 10, 32)
		slotId = uint32(tmpInt)
	}

	// ip
	var ip string
	if ips, exist := paramMap["ip"]; exist {
		ip = ips[0]
	}

	// device_id
	var deviceId string
	if deviceIds, exist := paramMap["device_id"]; exist {
		deviceId = deviceIds[0]
	}

	// os
	var os uint32
	if osString, exist := paramMap["os"]; exist {
		tmpInt, _ := strconv.ParseUint(osString[0], 10, 32)
		os = uint32(tmpInt)
	}

	// os_version
	var osVersion string
	if osVersions, exist := paramMap["os_version"]; exist {
		osVersion = osVersions[0]
	}

	// unit_id
	var unitId uint32
	if unitIdString, exist := paramMap["unit_id"]; exist {
		tmpInt, _ := strconv.ParseUint(unitIdString[0], 10, 32)
		unitId = uint32(tmpInt)
	}

	// creative_id
	var creativeId uint32
	if creativeIdString, exist := paramMap["creative_id"]; exist {
		tmp, _ := strconv.ParseUint(creativeIdString[0], 10, 32)
		creativeId = uint32(tmp)
	}

    adserver.ImpressionLog.Info(fmt.Sprintf(
    	"searchId=%s slotId=%d ip=%s deviceId=%s os=%d osVersion=%s unit_id=%d creativeId=%d",
		searchId, slotId, ip, deviceId, os, osVersion, unitId, creativeId))
	res := "{\"status\": 0}"
	w.Write([]byte(res))
}