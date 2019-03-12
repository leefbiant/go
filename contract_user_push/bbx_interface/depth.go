// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bbexgo/config"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	addr = flag.String("addr", "api.bbx.com", "http service address")
)

func ReBuildDepthObj(contract_name string, base_price float64, obj_slice DepthDataSlice) DepthDataSlice {
	depth_slice := make(DepthDataSlice, 0)
	for _, obj := range obj_slice {
		if math.Abs(obj.Price-base_price)/base_price*100 <= 10.0 {
			depth_slice = append(depth_slice, obj)
		}
	}
	return depth_slice
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func httpPost(param string) {
	data := fmt.Sprintf("data=%s", param)
	// glog.Info("httpPost data:", data)
	url := config.Get("contract_sdk_agent.addr")
	conn := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := conn.Post(url,
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
	glog.Info(string(body))
}

func WsChannel(currency_name string, contract_name string, id int, done chan struct{}) {
	defer close(done)
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/v1/ifcontract/realTime"}
	glog.Info("connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("dial:", err)
		return
	}
	defer c.Close()
	// msg := fmt.Sprintf("{\"action\":\"subscribe\",\"args\":[\"Depth:%d\"]}", id)
	msg := fmt.Sprintf("{\"action\":\"subscribe\",\"args\":[\"Depth:%d\"]}", id)
	c.WriteMessage(websocket.TextMessage, []byte(msg))
	var inside_obj_base InterfaceInside = InterfaceInside{
		DataType:       "depth",
		Exchange:       "BBX",
		ContractIndex:  GetContractType(contract_name),
		ContractName:   GetContractName(contract_name),
		ContractSymbol: currency_name,
	}

	var inside_obj DepthInsideObj

	for {
		var msg string
		messageType, message, err := c.ReadMessage()
		if err != nil {
			glog.Error("read:", err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			msg = string(message)
			break
		case websocket.BinaryMessage:
			text, err := GzipDecode(message)
			if err != nil {
				glog.Error("err:", err)
			} else {
				msg = string(text)
			}
			break
		}
		// glog.Info("recv: ", msg)
		if msg[11:15] == "ping" {
			pong := fmt.Sprintf("{\"group\":\"System\",\"data\":\"pong\"}")
			glog.Info("send pong: ", pong)
			err = c.WriteMessage(websocket.TextMessage, []byte(pong))
			if err != nil {
				glog.Error("WriteMessage:", err)
				return
			}
		} else {
			var obj DepthObj
			err = json.Unmarshal([]byte(msg), &obj)
			if err != nil {
				glog.Error("Unmarshal err:", err)
				return
			}

			if obj.Data.Way == 1 {
				inside_obj.Buys.List = inside_obj.Buys.List[0:0]
			} else if obj.Data.Way == 2 {
				inside_obj.Sells.List = inside_obj.Sells.List[0:0]
			}
			for _, v := range obj.Data.Depths {
				Type := "buy"
				if obj.Data.Way == 1 {
					Type = "buy"
					inside_obj.Buys.List = append(inside_obj.Buys.List, DepthData{
						Idx:   v.Price,
						Price: String2Float(v.Price),
						Size:  uint(String2Int(v.Vol)),
						Type:  Type,
					})
					inside_obj.Buys.Action = "partial"
				} else if obj.Data.Way == 2 {
					Type := "sell"
					inside_obj.Sells.List = append(inside_obj.Sells.List, DepthData{
						Idx:   v.Price,
						Price: String2Float(v.Price),
						Size:  uint(String2Int(v.Vol)),
						Type:  Type,
					})
					inside_obj.Sells.bset = true
					inside_obj.Sells.Action = "partial"
				}
			}
			if !inside_obj.Sells.bset || inside_obj.Buys.bset {
				continue
			}
			if len(inside_obj.Buys.List) <= 0 || len(inside_obj.Sells.List) <= 0 {
				continue
			}

			sort.Sort(inside_obj.Sells.List)
			sort.Sort(inside_obj.Buys.List)
			if len(inside_obj.Buys.List) > 0 {
				inside_obj_base.SinglePriceUSD = GetSinglePriceUSD(contract_name, inside_obj.Buys.List[0].Price)
			}

			base_price := (inside_obj.Buys.List[0].Price + inside_obj.Sells.List[len(inside_obj.Sells.List)-1].Price) / 2

			inside_obj.Sells.List = ReBuildDepthObj(contract_name, base_price, inside_obj.Sells.List)
			inside_obj.Buys.List = ReBuildDepthObj(contract_name, base_price, inside_obj.Buys.List)

			array := make([]DepthInsideObj, 0)
			array = append(array, inside_obj)
			data_str, _ := json.Marshal(array)
			inside_obj_base.Data = string(data_str)
			json_data, err := json.Marshal(inside_obj_base)
			if err != nil {
				glog.Error("Marshal failed:", err)
				return
			}
			httpPost(string(json_data))
			inside_obj.Sells.bset = false
			inside_obj.Buys.bset = false
		}
	}
}

func WatchHandle(currency_name string, contract_name string, id int) {
	done := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go WsChannel(currency_name, contract_name, id, done)

	for {
		select {
		case <-done:
			glog.Error("recv done....")
			done = make(chan struct{})
			go WsChannel(currency_name, contract_name, id, done)
		case t := <-ticker.C:
			t = t
			break
		}
	}
}

func Init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "50")
	flag.Set("v", "3")

	InitContractMap()
}

func main() {
	defer glog.Flush()
	Init()
	flag.Parse()
	err, contracts := GetContracts()
	if err != nil {
		glog.Error("GetContracts failed:", err)
		return
	}

	if contracts.Errno != "OK" || len(contracts.Data.Contracts) < 1 {
		glog.Error("not find any Contracts :")
		return
	}

	for _, v := range contracts.Data.Contracts {
		_, ok := trade_map[v.Contract.DisplayName]
		if !ok {
			continue
		}
		go WatchHandle(v.Contract.BaseCoin, v.Contract.DisplayName, v.Contract.ContractID)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
