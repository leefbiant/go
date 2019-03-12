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

type TradeObj struct {
	Ch   string `json:"ch"`
	Ts   uint64 `json:"ts"`
	Tick struct {
		ID   int    `json:"id"`
		Ts   uint64 `json:"ts"`
		Data []struct {
			Amount    int     `json:"amount"`
			Ts        uint64  `json:"ts"`
			ID        int64   `json:"id"`
			Price     float64 `json:"price"`
			Direction string  `json:"direction"`
		} `json:"data"`
	} `json:"tick"`
}

type TradeDataObj struct {
	Timestamp uint64  `json:"timestamp"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
}

type TradeObjInside struct {
	DataType       string  `json:"dataType"`
	Exchange       string  `json:"exchange"`
	ContractIndex  int     `json:"contractIndex"`
	ContractName   string  `json:"contractName"`
	ContractSymbol string  `json:"contractSymbol"`
	Data           string  `json:"data"`
	SinglePriceUSD float64 `json:"singlePriceUSD"` // 单张价格（USD）
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
		glog.Error("ReadAll err:", err)
		return
	}
	glog.Info("req:", data, " res:", string(body))
}

func WsChannel(name string, id int, req string, done chan struct{}) {
	defer close(done)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	glog.Info("connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("dial:", err)
		return
	}
	defer c.Close()
	c.SetReadDeadline(time.Now().Add(30 * time.Second))

	msg := fmt.Sprintf("{\"sub\": \"%s\", \"id\": \"id1\"}", req)
	c.WriteMessage(websocket.TextMessage, []byte(msg))

	ticker := time.NewTicker(180 * time.Second)
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
				glog.Info("error timeout for name:", name, " id:", id, " err:", err)
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

			if msg[0:7] == "{\"ping\"" {
				ts := msg[8:21]
				pong := fmt.Sprintf("{\"pong\":%s}", ts)
				err = c.WriteMessage(websocket.TextMessage, []byte(pong))
				if err != nil {
					glog.Error("WriteMessage:", err)
					return
				}
			} else {
				var s TradeObj
				err = json.Unmarshal([]byte(msg), &s)
				if err != nil {
					glog.Error("Unmarshal err:", err)
					return
				}
				var obj TradeObjInside = TradeObjInside{
					DataType:       "trade",
					Exchange:       "HUOBI",
					ContractIndex:  id,
					ContractName:   GetContractName(name, id),
					ContractSymbol: name,
					SinglePriceUSD: GetSinglePriceUSD(name),
				}

				trade_array := make([]TradeDataObj, 0)
				var trade_obj TradeDataObj
				// 买价
				for _, v := range s.Tick.Data {
					trade_obj.Size = v.Amount
					trade_obj.Price = v.Price
					trade_obj.Type = v.Direction
					trade_obj.Timestamp = v.Ts / 1000
					trade_array = append(trade_array, trade_obj)
				}
				data_str, _ := json.Marshal(trade_array)
				obj.Data = string(data_str)
				json_data, err := json.Marshal(obj)
				if err != nil {
					glog.Error("Marshal failed:", err)
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
	defer glog.Flush()
	Init()
	flag.Parse()
	go WatchHandle("BTC", 0, "market.BTC_CW.trade.detail") // 当周
	go WatchHandle("BTC", 1, "market.BTC_NW.trade.detail") // 次周
	go WatchHandle("BTC", 2, "market.BTC_CQ.trade.detail") // 季度

	go WatchHandle("ETH", 0, "market.ETH_CW.trade.detail")
	go WatchHandle("ETH", 1, "market.ETH_NW.trade.detail")
	go WatchHandle("ETH", 2, "market.ETH_CQ.trade.detail")

	for {
		time.Sleep(1 * time.Second)
	}
}
