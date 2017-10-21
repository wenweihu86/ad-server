package adserver

import (
	"io"
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
)

type AdUnitInfo struct {
	UnitId uint32
	CreativeId uint32
}

type AdCreativeInfo struct {
	CreativeId uint32
	Title string
	Description string
	AppPackageName string // for native app only
	IconImageUrl string
	MainImageUrl string
	ClickUrl string
}

type AdDictInfo struct {
	AdUnitMap map[uint32]AdUnitInfo
	AdCreativeMap map[uint32]AdCreativeInfo
	LocationUnitMap map[string][]uint32
}

var AdDict *AdDictInfo

func init()  {
	AdDict = &AdDictInfo{
		AdUnitMap: make(map[uint32]AdUnitInfo),
		AdCreativeMap: make(map[uint32]AdCreativeInfo),
		LocationUnitMap: make(map[string][]uint32),
	}
}

func LoadAdDict(dictFileName string) {
	dictFile, err := os.Open(dictFileName)
	if err != nil {
		AdServerLog.Error(fmt.Sprintf(
			"open file error, name=%s\n", dictFileName))
		panic(-1)
	}
	defer dictFile.Close()

	br := bufio.NewReader(dictFile)
	lineNum := 0
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineString := string(line)
		lines := strings.Split(lineString, "\t")
		level, _ := strconv.Atoi(lines[0])
		if level == 1 {
			// ad unit info
			unitId, _ := strconv.ParseUint(lines[1], 10, 32)
			creativeId, _ := strconv.ParseUint(lines[2], 10, 32)
			adUnit := AdUnitInfo{
				UnitId: uint32(unitId),
				CreativeId: uint32(creativeId),
			}
			AdDict.AdUnitMap[adUnit.UnitId] = adUnit
			lineNum++
			AdServerLog.Debug(fmt.Sprintf(
				"read ad unit info, unitId=%d creativeId=%d\n",
				unitId, creativeId))
		} else if level == 2 {
			// creative info
			creativeId, _ := strconv.ParseUint(lines[1], 10, 32)
			title := lines[2]
			description := lines[3]
			packageName := lines[4]
			iconImageUrl := lines[5]
			mainImageUrl := lines[6]
			clickUrl := lines[7]
			adCreative := AdCreativeInfo{
				CreativeId: uint32(creativeId),
				Title: title,
				Description: description,
				AppPackageName: packageName,
				IconImageUrl: iconImageUrl,
				MainImageUrl: mainImageUrl,
				ClickUrl: clickUrl,
			}
			AdDict.AdCreativeMap[adCreative.CreativeId] = adCreative
			lineNum++
			AdServerLog.Debug(fmt.Sprintf(
				"read ad creative info, creativeId=%d " +
				"title=%s description=%s package=%s iconImageUrl=%s " +
				"mainImageUrl=%s clickUrl=%s\n",
				creativeId, title, description, packageName,
				iconImageUrl, mainImageUrl, clickUrl))
		} else if level == 3 {
			// location target
			unitId, _ := strconv.ParseUint(lines[1], 10, 32)
			country := strings.ToLower(lines[2])
			city := strings.ToLower(lines[3])
			key := country + "_" + city
			unitIdList, exist := AdDict.LocationUnitMap[key]
			if !exist {
				unitIdList = make([]uint32, 0)
			}
			unitIdList = append(unitIdList, uint32(unitId))
			AdDict.LocationUnitMap[key] = unitIdList
			AdServerLog.Debug(fmt.Sprintf(
				"read location target info, unitId=%d country=%s city=%s\n",
				unitId, country, city))
			lineNum++
		}
	}
	AdServerLog.Info(fmt.Sprintf(
		"read ad info file success, lineNum=%d\n", lineNum))
}
