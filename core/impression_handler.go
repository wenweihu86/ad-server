package core

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"encoding/base64"
	"fmt"
	"net/url"
)

// 展现监控handler
func ImpressionHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	if !args.Has("i") {
		ctx.SetBody([]byte("{\"status\": 1}"))
		return
    }
	argsValueBytes := args.Peek("i")
	queryStringBytes, err := base64.URLEncoding.DecodeString(string(argsValueBytes))
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return
	}
	queryString := string(queryStringBytes)
	paramMap, err := url.ParseQuery(queryString)
	if err != nil {
		ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
		return 
	}

	// search_id
	var searchId string
	if searchIds, exist := paramMap["search_id"]; exist {
		searchId = searchIds[0]
	}

	// slot_id
	var slotId uint32
	if slotIds, exist := paramMap["slot_id"]; exist {
		tmpInt, err := strconv.ParseUint(slotIds[0], 10, 32)
		if err != nil {
			ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
			return 
		}
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
		tmpInt, err := strconv.ParseUint(osString[0], 10, 32)
		if err != nil {
			ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
			return 
		}
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
		tmpInt, err := strconv.ParseUint(unitIdString[0], 10, 32)
		if err != nil {
			ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
			return 
		}
		unitId = uint32(tmpInt)
	}

	// creative_id
	var creativeId uint32
	if creativeIdString, exist := paramMap["creative_id"]; exist {
		tmp, err := strconv.ParseUint(creativeIdString[0], 10, 32)
		if err != nil {
			ctx.SetBody([]byte("{\"status\": 1," + "\"error\":" + err.Error() + "}"))
			return 
		}
		creativeId = uint32(tmp)
	}

    ImpressionLog.Info(fmt.Sprintf(
    	"impression=1 searchId=%s slotId=%d ip=%s deviceId=%s os=%d osVersion=%s unit_id=%d creativeId=%d",
		searchId, slotId, ip, deviceId, os, osVersion, unitId, creativeId))
	res := "{\"status\": 0}"
	ctx.SetBody([]byte(res))
}