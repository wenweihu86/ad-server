package adserver

import (
	"io"
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
	"sort"
	"ad-server/utils"
	"time"
)

type GeoLocationInfo struct {
	LocId uint64
	Country string
	City string
}

type LocationInfo struct {
	BeginIp uint32
	EndIp uint32
	Country string
	City string
}

type IpDataInfo struct {
	ipPairs utils.IpPairs
	ipLocationMap map[utils.IpPair]*LocationInfo
}

type IpDict struct {
	IpDataArray []*IpDataInfo
	CurrentIndex uint32
	BlockLastModifiedTime int64
	LocationLastModifiedTime int64
}

var LocationDict *IpDict

func init() {
	LocationDict = &IpDict{
		IpDataArray: make([]*IpDataInfo, 2, 2),
		CurrentIndex: 0,
		BlockLastModifiedTime: 0,
		LocationLastModifiedTime: 0,
	}
	for i := 0; i < 2; i++ {
		LocationDict.IpDataArray[i] = NewIpDataInfo()
	}
}

// 初始化之后首次加载Ip字典信息
func (ipDict *IpDict) Load() error {
	ipDataInfo,err := LoadLocationDict(GlobalConfObject.GeoBlockFileName,
		GlobalConfObject.GeoLocationFileName)
	if err != nil {
		return err
	}
	ipDict.IpDataArray[ipDict.CurrentIndex] = ipDataInfo
	blockFileStat, _ := os.Stat(GlobalConfObject.GeoBlockFileName)
    locationFileStat, _ := os.Stat(GlobalConfObject.GeoLocationFileName)
	ipDict.BlockLastModifiedTime = blockFileStat.ModTime().Unix()
	ipDict.LocationLastModifiedTime = locationFileStat.ModTime().Unix()
	return nil
}

func NewIpDataInfo() *IpDataInfo {
	return &IpDataInfo{
		ipPairs: make(utils.IpPairs, 0),
		ipLocationMap: make(map[utils.IpPair]*LocationInfo),
	}
}

func LoadLocationDict(blockFileName, locationFileName string) (*IpDataInfo, error) {
	dictFile, err := os.Open(blockFileName)
	if err != nil {
		AdServerLog.Error(fmt.Sprintf("open file error, name=%s\n", blockFileName))
		return nil, err
	}
	defer dictFile.Close()

	ipDataInfo := NewIpDataInfo()
	geoLocationMap, err := loadGeoLocation(locationFileName)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(dictFile)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineString := string(line)
		if len(lineString) == 0 || lineString[0] == '#' {
			continue
		}
		lines := strings.Split(lineString, ",")
		if len(lines) != 3 {
			AdServerLog.Warn(fmt.Sprintf(
				"invalid format, blockFileName=%s, line=%s\n",
				blockFileName, lineString))
			continue
		}

		tmpInt, _ := strconv.ParseUint(strings.Trim(lines[0], "\""), 10 ,32)
		beginIp := uint32(tmpInt)
		tmpInt, _ = strconv.ParseUint(strings.Trim(lines[1], "\""), 10 ,32)
		endIp := uint32(tmpInt)
		locId, _ := strconv.ParseUint(strings.Trim(lines[2], "\""), 10 ,64)
		geoLocationInfo, exist := geoLocationMap[locId]
		if !exist {
			//fmt.Printf("geoLocationInfo not found, locId=%d\n", locId)
			continue
		}
		locationInfo := &LocationInfo{
			BeginIp: beginIp,
			EndIp: endIp,
			Country: geoLocationInfo.Country,
			City: geoLocationInfo.City,
		}
		ipPair := utils.IpPair{
			BeginIp: beginIp,
			EndIp: endIp,
		}
		ipDataInfo.ipPairs = append(ipDataInfo.ipPairs, ipPair)
		ipDataInfo.ipLocationMap[ipPair] = locationInfo
	}

	sort.Sort(ipDataInfo.ipPairs)
	AdServerLog.Info(fmt.Sprintf(
		"read dict success, blockFileName=%s locationFileName=%s\n",
		blockFileName, locationFileName))
	AdServerLog.Info(fmt.Sprintf(
		"location dict size=%d\n", len(ipDataInfo.ipPairs)))

	return ipDataInfo, nil
}

// 启动定时器，用于定期重新加载Ip字典信息
func (locationDict *IpDict) StartReloadTimer() {
	duration := int64(time.Second) * GlobalConfObject.IpFileReloadInterval
	t := time.NewTicker(time.Duration(duration))
	go func() {
		for t1 := range t.C {
			AdServerLog.Debug("IpDict reload timer execute")
			blockFileStat, _ := os.Stat(GlobalConfObject.GeoBlockFileName)
			locationFileStat, _ := os.Stat(GlobalConfObject.GeoLocationFileName)
			blockCurrentModifiedTime := blockFileStat.ModTime().Unix()
			locationCurrentModifiedTime := locationFileStat.ModTime().Unix()
			// 如果文件有更新，则重新加载广告内容
			if blockCurrentModifiedTime > locationDict.BlockLastModifiedTime || locationCurrentModifiedTime > locationDict.LocationLastModifiedTime {
				AdServerLog.Info(fmt.Sprintf("start reload ad info dict at %s",
					t1.Format("2006-01-02 03:04:05")))
				_, err := LoadLocationDict(
					GlobalConfObject.GeoBlockFileName,
					GlobalConfObject.GeoLocationFileName)
				if err != nil {
					continue
				}
				nextIndex := 1 - locationDict.CurrentIndex
				locationDict.CurrentIndex = nextIndex
				locationDict.BlockLastModifiedTime = blockCurrentModifiedTime
				locationDict.LocationLastModifiedTime = locationCurrentModifiedTime
			}
		}
	}()
}

// 获取当前可用的Ip字典信息
func (ipDict *IpDict) GetCurrentIpData() *IpDataInfo {
	return ipDict.IpDataArray[ipDict.CurrentIndex]
}

func (ipDataInfo *IpDataInfo) SearchLocationByIp(ipString string) *LocationInfo {
	ip := utils.StringIpToUint(ipString)
	ipPairs := ipDataInfo.ipPairs
	size := len(ipPairs)
	if size == 0 || ip < ipPairs[0].BeginIp || ip > ipPairs[size - 1].EndIp {
		return nil
	}
	left := 0
	right := size - 1
	for left <= right {
		mid := (left + right) / 2
		if ip >= ipPairs[mid].BeginIp && ip <= ipPairs[mid].EndIp {
			return ipDataInfo.ipLocationMap[ipPairs[mid]];
		} else if ip < ipPairs[mid].BeginIp {
			right = mid - 1
		} else if ip > ipPairs[mid].EndIp {
			left = mid + 1
		}
	}
	return nil
}

func loadGeoLocation(fileName string) (map[uint64]*GeoLocationInfo, error) {
	dictFile, err := os.Open(fileName)
	if err != nil {
		AdServerLog.Error(fmt.Sprintf(
		   "open file error, name=%s\n", fileName))
		return nil, err
	}
	defer dictFile.Close()

	geoLocationMap := make(map[uint64]*GeoLocationInfo)
	br := bufio.NewReader(dictFile)
	lineNum := 0
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineString := string(line)
		if len(lineString) == 0 || lineString[0] == '#' {
			continue
		}
		lines := strings.Split(lineString, ",")
		if len(lines) != 9 {
			AdServerLog.Warn(fmt.Sprintf(
				"invalid format, file=%s, line=%s\n",
				fileName, lineString))
			continue
		}
		locId, _ := strconv.ParseUint(lines[0], 10, 64)
		geoLocationInfo := &GeoLocationInfo{
			LocId: locId,
			Country: strings.Trim(lines[1], "\""),
			City: strings.Trim(lines[3], "\""),
		}
		geoLocationMap[locId] = geoLocationInfo
		lineNum++
	}
	AdServerLog.Info(fmt.Sprintf(
		"load dict success, file=%s, lineNum=%d\n",
		fileName, lineNum))
	return geoLocationMap, nil
}
