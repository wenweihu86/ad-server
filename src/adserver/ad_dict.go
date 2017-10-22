package adserver

import (
	"io"
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
	"time"
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

type AdDataInfo struct {
	AdUnitMap map[uint32]AdUnitInfo
	AdCreativeMap map[uint32]AdCreativeInfo
	LocationUnitMap map[string][]uint32
}

func NewAdDataInfo() *AdDataInfo {
	return &AdDataInfo{
		AdUnitMap: make(map[uint32]AdUnitInfo),
		AdCreativeMap: make(map[uint32]AdCreativeInfo),
		LocationUnitMap: make(map[string][]uint32),
	}
}

type AdDict struct {
	AdDataArray []*AdDataInfo
	CurrentIndex uint32
	FileName string
	LastModifiedTime int64
}

var AdDictObject *AdDict

// 初始化广告对象
func NewAdDict(dictFileName string) *AdDict {
	adDictObject := &AdDict{
		AdDataArray: make([]*AdDataInfo, 2, 2),
		CurrentIndex: 0,
		FileName: dictFileName,
		LastModifiedTime: 0,
	}
	for i := 0; i < 2; i++ {
		adDictObject.AdDataArray[i] = NewAdDataInfo()
	}
	return adDictObject
}

// 初始化之后首次加载广告信息
func (ad *AdDict) Load() {
	adDataInfo := ad.loadAdDict()
	ad.AdDataArray[ad.CurrentIndex] = adDataInfo
	fileStat, _ := os.Stat(ad.FileName)
	ad.LastModifiedTime = fileStat.ModTime().Unix()
}

// 启动定时器，用于定期重新加载广告信息
func (ad *AdDict) StartReloadTimer() {
	duration := int64(time.Second) * GlobalConfObject.AdFileReloadInterval
	t := time.NewTicker(time.Duration(duration))
	go func() {
		for t1 := range t.C {
			AdServerLog.Debug("AdDict reload timer execute")
			fileStat, _ := os.Stat(ad.FileName)
			currentModifiedTime := fileStat.ModTime().Unix()
			// 如果文件有更新，则重新加载广告内容
			if currentModifiedTime > ad.LastModifiedTime {
				AdServerLog.Info(fmt.Sprintf("start reload ad info dict at %s",
					t1.Format("2006-01-02 03:04:05")))
				adDataInfo := ad.loadAdDict()
				nextIndex := 1 - ad.CurrentIndex
				ad.AdDataArray[nextIndex] = adDataInfo
				ad.CurrentIndex = nextIndex
				ad.LastModifiedTime = currentModifiedTime
			}
		}
	}()
}

// 获取当前可用的广告信息
func (ad *AdDict) GetCurrentAdData() *AdDataInfo {
	return ad.AdDataArray[ad.CurrentIndex]
}

func (ad *AdDict) loadAdDict() *AdDataInfo {
	dictFile, err := os.Open(ad.FileName)
	if err != nil {
		AdServerLog.Error(fmt.Sprintf(
			"open file error, name=%s\n", ad.FileName))
		panic(-1)
	}
	defer dictFile.Close()

	adDataInfo := NewAdDataInfo()
	lineNum := 0
	br := bufio.NewReader(dictFile)
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
			adDataInfo.AdUnitMap[adUnit.UnitId] = adUnit
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
			adDataInfo.AdCreativeMap[adCreative.CreativeId] = adCreative
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
			unitIdList, exist := adDataInfo.LocationUnitMap[key]
			if !exist {
				unitIdList = make([]uint32, 0)
			}
			unitIdList = append(unitIdList, uint32(unitId))
			adDataInfo.LocationUnitMap[key] = unitIdList
			AdServerLog.Debug(fmt.Sprintf(
				"read location target info, unitId=%d country=%s city=%s\n",
				unitId, country, city))
			lineNum++
		} else {
			panic(1)
		}
	}
	AdServerLog.Info(fmt.Sprintf(
		"read ad info file success, lineNum=%d\n", lineNum))
	return adDataInfo
}
