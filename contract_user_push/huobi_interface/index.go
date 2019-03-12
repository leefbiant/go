// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bbexgo/config"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type IndexObj struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol     string  `json:"symbol"`
		IndexPrice float64 `json:"index_price"`
		IndexTs    uint64  `json:"index_ts"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type IndexDataObj struct {
	Timestamp uint64  `json:"timestamp"`
	Price     float64 `json:"price"`
}

type IndexObjInside struct {
	DataType       string `json:"dataType"`
	Exchange       string `json:"exchange"`
	ContractIndex  int    `json:"contractIndex"`
	ContractName   string `json:"contractName"`
	ContractSymbol string `json:"contractSymbol"`
	Data           string `json:"data"`
}

func httpPost(param string) {
	if param == "" {
		glog.Error("httpPost input null")
		return
	}
	data := fmt.Sprintf("data=%s", param)
	url := config.Get("contract_sdk_agent.addr")
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		glog.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("ReadAll err:", err)
		return
	}
	glog.Info("req:", data, " res:", string(body))
}

func GetIndex(name string) {
	url := fmt.Sprintf("https://api.hbdm.com/api/v1/contract_index?symbol=%s", name)
	// glog.Info("url:", url)
	resp, err := http.Get(url)
	if err != nil {
		glog.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("ReadAll err:", err)
		return
	}
	var msg IndexObj
	err = json.Unmarshal([]byte(body), &msg)
	if err != nil {
		glog.Error("Unmarshal failed:", err)
		return
	}

	var obj IndexObjInside = IndexObjInside{
		DataType:       "index",
		Exchange:       "HUOBI",
		ContractIndex:  0,
		ContractName:   "",
		ContractSymbol: name,
	}

	index_array := make([]IndexDataObj, 0)
	var index_obj IndexDataObj
	// 买价
	for _, v := range msg.Data {
		index_obj.Price = v.IndexPrice
		index_obj.Timestamp = v.IndexTs / 1000
		index_array = append(index_array, index_obj)
	}
	data_str, _ := json.Marshal(index_array)
	obj.Data = string(data_str)
	json_data, err := json.Marshal(obj)
	if err != nil {
		glog.Error("Marshal failed:", err)
		return
	}
	httpPost(string(json_data))
}

func WatchHandle(name string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			go GetIndex(name)
		}
	}
}

func LogInit() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "50")
	flag.Set("flushInterval", "1")
	flag.Set("v", "3")
}

func main() {
	defer glog.Flush()
	LogInit()
	flag.Parse()
	go WatchHandle("BTC")
	go WatchHandle("ETH")

	for {
		time.Sleep(1 * time.Second)
	}
}
