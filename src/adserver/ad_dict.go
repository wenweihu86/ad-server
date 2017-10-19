package adserver

import (
	"io"
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
)

type AdUnit struct {
	UnitId uint
	CreativeId uint
}

type AdCreative struct {
	CreativeId uint
	Title string
	IconImageUrl string
	ClickUrl string
}

var AdUnits = make([]AdUnit, 0, 10000)
var AdCreativeMap = make(map[uint]AdCreative, 10000)

func ReadAdDict() {
	dictFileName := "./data/ad_info.txt"
	dictFile, err := os.Open(dictFileName)
	if err != nil {
		fmt.Printf("open file error, name=%s\n", dictFileName)
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
			adUnit := AdUnit{
				UnitId: uint(unitId),
				CreativeId: uint(creativeId),
			}
			AdUnits = append(AdUnits, adUnit)
			lineNum++
			fmt.Printf("read ad unit info, unitId=%d creativeId=%d\n",
				unitId, creativeId)
		} else if level == 2 {
			// creative info
			creativeId, _ := strconv.ParseUint(lines[1], 10, 32)
			title := lines[2]
			iconImageUrl := lines[3]
			clickUrl := lines[4]
			adCreative := AdCreative{
				CreativeId: uint(creativeId),
				Title: title,
				IconImageUrl: iconImageUrl,
				ClickUrl: clickUrl,
			}
			AdCreativeMap[adCreative.CreativeId] = adCreative
			lineNum++
			fmt.Printf("read ad creative info, creativeId=%d title=%s iconImageUrl=%s clickUrl=%s\n",
				creativeId, title, iconImageUrl, clickUrl)
		}
	}
	fmt.Printf("read ad info file success, lineNum=%d\n", lineNum)
}
