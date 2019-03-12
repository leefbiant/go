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

type OrderData struct {
	Price float64
	Size  int
}

type DepthObj struct {
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		Mrid    int           `json:"mrid"`
		ID      int           `json:"id"`
		Bids    []interface{} `json:"bids"` // 买入
		Asks    []interface{} `json:"asks"` // 卖出
		Ts      int64         `json:"ts"`
		Version int           `json:"version"`
		Ch      string        `json:"ch"`
	} `json:"tick"`
}

type DepthData struct {
	Idx   string  `json:"idx"`
	Price float64 `json:"price"`
	Size  uint    `json:"size"`
	Type  string  `json:"type"`
}

type DepthDataSlice []DepthData

type DepthDataObj struct {
	Sells struct {
		Action string         `json:"action"`
		List   DepthDataSlice `json:"list"`
	} `json:"sells"`
	Buys struct {
		Action string         `json:"action"`
		List   DepthDataSlice `json:"list"`
	} `json:"buys"`
}

func (s DepthDataSlice) Len() int           { return len(s) }
func (s DepthDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s DepthDataSlice) Less(i, j int) bool { return s[i].Price > s[j].Price }

type DepthObjInside struct {
	DataType       string  `json:"dataType"`
	Exchange       string  `json:"exchange"`
	ContractIndex  int     `json:"contractIndex"`
	ContractName   string  `json:"contractName"`
	ContractSymbol string  `json:"contractSymbol"`
	Data           string  `json:"data"`
	SinglePriceUSD float64 `json:"singlePriceUSD"` // 单张价格（USD）
}

func ReBuildDepthObj(base_price float64, obj_slice DepthDataSlice) DepthDataSlice {
	depth_slice := make(DepthDataSlice, 0)
	for _, obj := range obj_slice {
		if math.Abs(obj.Price-base_price)/base_price*100 <= 10.0 {
			depth_slice = append(depth_slice, obj)
		}
	}
	return depth_slice
}

var addr = flag.String("addr", "www.hbdm.com", "http service address")

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
	if param == "" {
		glog.Error("httpPost input null")
		return
	}

	data := fmt.Sprintf("data=%s", param)
	url := config.Get("contract_sdk_agent.addr")
	// glog.Info("post data:", data)

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
	glog.Info("req:", len(data), " res:", string(body))
}

func WsChannel(name string, id int, req string, done chan struct{}) {
	defer close(done)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	glog.Info("connecting to ", u.String())
	last_msg_t := time.Now().Unix()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("dial:", err)
		return
	}
	defer c.Close()
	c.SetReadDeadline(time.Now().Add(30 * time.Second))

	msg := fmt.Sprintf("{\"sub\": \"%s\", \"id\": \"id5\"}", req)
	glog.Info("WsChannel:", msg)
	c.WriteMessage(websocket.TextMessage, []byte(msg))
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
		if msg[0:7] == "{\"ping\"" {
			ts := msg[8:21]
			pong := fmt.Sprintf("{\"pong\":%s}", ts)
			err = c.WriteMessage(websocket.TextMessage, []byte(pong))
			if err != nil {
				glog.Error("WriteMessage:", err)
				return
			}
			// glog.Info("recv: ", msg)
		} else {

			now_t := time.Now().Unix()
			if last_msg_t == now_t {
				continue
			}

			var s DepthObj
			err = json.Unmarshal([]byte(msg), &s)
			if err != nil {
				glog.Error("Unmarshal err:", err)
				return
			}
			var obj DepthObjInside = DepthObjInside{
				DataType:       "depth",
				Exchange:       "HUOBI",
				ContractIndex:  id,
				ContractName:   GetContractName(name, id),
				ContractSymbol: name,
				SinglePriceUSD: GetSinglePriceUSD(name),
			}

			var depobj DepthDataObj

			// // 买价
			for _, v := range s.Tick.Bids {
				vs := v.([]interface{})
				depobj.Buys.Action = "partial"
				dep_data := DepthData{
					Idx:   fmt.Sprintf("%f", vs[0].(float64)),
					Price: vs[0].(float64),
					Size:  uint(vs[1].(float64)),
					Type:  "buy",
				}
				depobj.Buys.List = append(depobj.Buys.List, dep_data)
			}

			// 卖价
			for _, v := range s.Tick.Asks {
				vs := v.([]interface{})
				depobj.Sells.Action = "partial"
				dep_data := DepthData{
					Idx:   fmt.Sprintf("%f", vs[0].(float64)),
					Price: vs[0].(float64),
					Size:  uint(vs[1].(float64)),
					Type:  "sell",
				}
				depobj.Sells.List = append(depobj.Sells.List, dep_data)
			}

			if len(depobj.Buys.List) <= 0 || len(depobj.Sells.List) <= 0 {
				continue
			}

			sort.Sort(depobj.Sells.List)
			sort.Sort(depobj.Buys.List)

			base_price := (depobj.Buys.List[0].Price + depobj.Sells.List[len(depobj.Sells.List)-1].Price) / 2

			depobj.Sells.List = ReBuildDepthObj(base_price, depobj.Sells.List)
			depobj.Buys.List = ReBuildDepthObj(base_price, depobj.Buys.List)

			dep_array := make([]DepthDataObj, 0)
			dep_array = append(dep_array, depobj)
			data_str, _ := json.Marshal(dep_array)
			obj.Data = string(data_str)
			json_data, err := json.Marshal(obj)
			if err != nil {
				glog.Error("Marshal failed:", err)
				return
			}
			httpPost(string(json_data))
			last_msg_t = now_t
			c.SetReadDeadline(time.Now().Add(30 * time.Second))
		}
	}
}

func WatchHandle(name string, id int, req string) {
	done := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go WsChannel(name, id, req, done)

	for {
		select {
		case <-done:
			glog.Error("recv done....")
			done = make(chan struct{})
			go WsChannel(name, id, req, done)
		case t := <-ticker.C:
			t = t
			break
		}
	}
}

func LogInit() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "200")
	flag.Set("flushInterval", "1")
	flag.Set("v", "3")
}

func main() {
	defer glog.Flush()
	LogInit()
	flag.Parse()
	go WatchHandle("BTC", 0, "market.BTC_CW.depth.step0")
	go WatchHandle("BTC", 1, "market.BTC_NW.depth.step0")
	go WatchHandle("BTC", 2, "market.BTC_CQ.depth.step0")

	go WatchHandle("ETH", 0, "market.ETH_CW.depth.step0")
	go WatchHandle("ETH", 1, "market.ETH_NW.depth.step0")
	go WatchHandle("ETH", 2, "market.ETH_CQ.depth.step0")

	for {
		time.Sleep(1 * time.Second)
	}
}
