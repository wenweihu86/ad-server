package adserver

import (
	"io"
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
	"sort"
	"utils"
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

type IpDict struct {
	ipPairs utils.IpPairs
	ipLocationMap map[utils.IpPair]*LocationInfo
}

var LocationDict *IpDict

func init() {
	LocationDict = &IpDict{
		ipPairs: make(utils.IpPairs, 0),
		ipLocationMap: make(map[utils.IpPair]*LocationInfo),
	}
}

func LoadLocationDict(blockFileName, locationFileName string) {
	dictFile, err := os.Open(blockFileName)
	if err != nil {
		fmt.Printf("open file error, name=%s\n", blockFileName)
	}
	defer dictFile.Close()

	geoLocationMap := loadGeoLocation(locationFileName)
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
			fmt.Printf("invalid format, blockFileName=%s, line=%s\n",
				blockFileName, lineString)
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
		LocationDict.ipPairs = append(LocationDict.ipPairs, ipPair)
		LocationDict.ipLocationMap[ipPair] = locationInfo
	}

	sort.Sort(LocationDict.ipPairs)
	fmt.Printf("read dict success, blockFileName=%s locationFileName=%s\n", blockFileName, locationFileName)
	fmt.Printf("location dict size=%d\n", len(LocationDict.ipPairs))
}

func SearchLocationByIp(ipString string) *LocationInfo {
	ip := utils.StringIpToUint(ipString)
	ipPairs := LocationDict.ipPairs
	size := len(ipPairs)
	if size == 0 || ip < ipPairs[0].BeginIp || ip > ipPairs[size - 1].EndIp {
		return nil
	}
	left := 0
	right := size - 1
	for left <= right {
		mid := (left + right) / 2
		if ip >= ipPairs[mid].BeginIp && ip <= ipPairs[mid].EndIp {
			return LocationDict.ipLocationMap[ipPairs[mid]];
		} else if ip < ipPairs[mid].BeginIp {
			right = mid - 1
		} else if ip > ipPairs[mid].EndIp {
			left = mid + 1
		}
	}
	return nil
}

func loadGeoLocation(fileName string) map[uint64]*GeoLocationInfo {
	dictFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("open file error, name=%s\n", fileName)
		return nil
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
			fmt.Printf("invalid format, file=%s, line=%s\n",
				fileName, lineString)
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
	fmt.Printf("load dict success, file=%s, lineNum=%d\n", fileName, lineNum)
	return geoLocationMap
}
