package main

import (
	"fmt"
	"github.com/golang/glog"
)

var (
	single_map = make(map[string]float64)
)

func GetContractName(name string, id int) string {
	ext := ""
	switch id {
	case 0:
		ext = "_CW"
		break
	case 1:
		ext = "_NW"
		break
	case 2:
		ext = "_CQ"
		break
	}
	return fmt.Sprintf("%s%s", name, ext)
}

func InitSinglePrice() {
	single_map["BTC"] = 100
	single_map["ETH"] = 10
	single_map["EOS"] = 10
	single_map["ETC"] = 10
}

func GetSinglePriceUSD(name string) float64 {
	if len(single_map) <= 0 {
		InitSinglePrice()
	}
	obj, ok := single_map[name]
	if !ok {
		glog.Error("not find data for Contract name:", name)
		return 0.0
	}
	return obj
}
