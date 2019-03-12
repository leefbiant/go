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
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TickerObj struct {
	Ch   string `json:"ch"`
	Ts   uint64 `json:"ts"`
	Tick struct {
		ID     int     `json:"id"`
		Mrid   int     `json:"mrid"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
		Amount float64 `json:"amount"`
		Vol    int     `json:"vol"`
		Count  int     `json:"count"`
	} `json:"tick"`
}

type TickerDataObj struct {
	Timestamp uint64  `json:"timestamp"`
	BuyPrice  float64 `json:"buyPrice"`
	BuySize   uint64  `json:"buySize"`
	SellPrice float64 `json:"sellPrice"`
	SellSize  uint64  `json:"sellSize"`
}

type TickerObjInside struct {
	DataType       string `json:"dataType"`
	Exchange       string `json:"exchange"`
	ContractIndex  int    `json:"contractIndex"`
	ContractName   string `json:"contractName"`
	ContractSymbol string `json:"contractSymbol"`
	Data           string `json:"data"`
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
		glog.Error("ReadAll err:%v", err)
		return
	}
	glog.Info("req:", data, " res:", string(body))
}

func WsChannel(name string, id int, req string, done chan struct{}) {
	defer close(done)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	glog.Error("connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("dial:", err)
		return
	}
	defer c.Close()
	c.SetReadDeadline(time.Now().Add(180 * time.Second))

	msg := fmt.Sprintf("{\"sub\": \"%s\", \"id\": \"id6\"}", req)
	// glog.Info("%s", msg)
	c.WriteMessage(websocket.TextMessage, []byte(msg))
	ticker := time.NewTicker(30 * time.Second)
	last_live := time.Now().Unix()
	for {
		select {
		case _ = <-ticker.C:
			if time.Now().Unix()-last_live > 300 {
				glog.Info("error timeout for name:", name, " id:", id)
				return
			}
		default:
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

			last_live = time.Now().Unix()
			// glog.Info("recv: ", msg)
			if msg[0:7] == "{\"ping\"" {
				ts := msg[8:21]
				pong := fmt.Sprintf("{\"pong\":%s}", ts)
				err = c.WriteMessage(websocket.TextMessage, []byte(pong))
				if err != nil {
					glog.Error("WriteMessage:", err)
					return
				}
			} else {
				var s TickerObj
				err = json.Unmarshal([]byte(msg), &s)
				if err != nil {
					glog.Error("Unmarshal err:%v", err)
					return
				}
				var obj TickerObjInside = TickerObjInside{
					DataType:       "ticker",
					Exchange:       "huobi",
					ContractIndex:  id,
					ContractName:   GetContractName(name, id),
					ContractSymbol: name,
				}

				if 0 == s.Tick.Close || 0 == s.Tick.Vol {
					continue
				}

				ticker_array := make([]TickerDataObj, 0)
				var ticker_obj TickerDataObj
				ticker_obj.Timestamp = s.Ts / 1000
				ticker_obj.BuyPrice = s.Tick.Close
				ticker_obj.SellPrice = s.Tick.Close
				ticker_obj.BuySize = uint64(s.Tick.Vol)
				ticker_obj.SellSize = uint64(s.Tick.Vol)

				ticker_array = append(ticker_array, ticker_obj)

				data_str, _ := json.Marshal(ticker_array)
				obj.Data = string(data_str)
				json_data, err := json.Marshal(obj)
				if err != nil {
					glog.Error("Marshal failed:%v", err)
					return
				}
				httpPost(string(json_data))
				c.SetReadDeadline(time.Now().Add(180 * time.Second))
			}
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

func Init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "50")
	flag.Set("v", "3")
}

func main() {
	Init()
	flag.Parse()
	// glog.SetFlags(0)
	go WatchHandle("BTC", 0, "market.BTC_CW.kline.1min")
	go WatchHandle("BTC", 1, "market.BTC_NW.kline.1min")
	go WatchHandle("BTC", 2, "market.BTC_CQ.kline.1min")

	go WatchHandle("ETH", 0, "market.ETH_CW.kline.1min")
	go WatchHandle("ETH", 1, "market.ETH_NW.kline.1min")
	go WatchHandle("ETH", 2, "market.ETH_CQ.kline.1min")

	for {
		time.Sleep(30 * time.Second)
	}
}
