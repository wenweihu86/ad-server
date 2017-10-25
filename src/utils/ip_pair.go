package utils

import (
	"strings"
	"strconv"
)

type IpPair struct {
	BeginIp uint32
	EndIp uint32
}

type IpPairs []IpPair

func (p IpPairs) Len() int  {
	return len(p)
}

func (p IpPairs) Less(i, j int) bool {
	if p[i].BeginIp == p[j].BeginIp {
		return p[i].EndIp < p[j].EndIp
	} else {
		return p[i].BeginIp < p[j].BeginIp
	}
}

func (p IpPairs) Swap(i, j int)  {
	p[i], p[j] = p[j], p[i]
}

func StringIpToUint(ipstring string) uint32 {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt uint32 = 0
	var pos uint32 = 24
	for _, ipSeg := range ipSegs {
		tempInt, err := strconv.ParseUint(ipSeg, 10, 32)
		if err !=nil {
			panic(-1)
		}
		intSeg := (uint32(tempInt)) << pos
		ipInt = ipInt | intSeg
		pos -= 8
	}
	return ipInt
}
