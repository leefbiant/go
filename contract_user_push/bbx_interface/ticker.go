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

var (
	addr            = flag.String("addr", "api.bbx.com", "http service address")
	watch_trade_map map[int]TickerWatchObj
)

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
	glog.Info("httpPost data:", data)
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

func WsChannel(done chan struct{}) {
	defer close(done)
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/v1/ifcontract/realTime"}
	glog.Info("connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("dial:", err)
		return
	}
	defer c.Close()

	msg := fmt.Sprintf("{\"action\":\"subscribe\",\"args\":[\"Ticker\"]}")
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
			var obj TickerObj
			err = json.Unmarshal([]byte(msg), &obj)
			if err != nil {
				glog.Error("Unmarshal err:", err)
				return
			}

			data, ok := watch_trade_map[obj.Data.ContractID]
			if !ok {
				continue
			}
			{

				// for ticker
				var inside_obj TradeObjInside = TradeObjInside{
					DataType:       "ticker",
					Exchange:       "BBX",
					ContractIndex:  GetContractType(data.Contract_name),
					ContractName:   GetContractName(data.Contract_name),
					ContractSymbol: data.Currency_name,
				}

				array := make([]TickerDataObj, 0)
				var data_obj TickerDataObj
				data_obj.Timestamp = uint64(obj.Data.Timestamp)
				data_obj.BuyPrice = String2Float(obj.Data.LastPrice)
				data_obj.SellPrice = String2Float(obj.Data.LastPrice)
				data_obj.BuySize = uint64(String2Int(obj.Data.Volume))
				data_obj.SellSize = uint64(String2Int(obj.Data.Volume))
				data_obj.RiseFallRate = String2Float(obj.Data.RiseFallRate)

				array = append(array, data_obj)

				data_str, _ := json.Marshal(array)
				inside_obj.Data = string(data_str)
				json_data, err := json.Marshal(inside_obj)
				if err != nil {
					glog.Error("Marshal failed:", err)
					return
				}
				httpPost(string(json_data))
			}

			{
				// for index
				var inside_obj InterfaceInside = InterfaceInside{
					DataType:       "index",
					Exchange:       "BBX",
					ContractIndex:  GetContractType(data.Contract_name),
					ContractName:   GetContractName(data.Contract_name),
					ContractSymbol: data.Currency_name,
				}

				array := make([]IndexDataObj, 0)
				var data_obj IndexDataObj
				data_obj.Timestamp = uint64(obj.Data.Timestamp)
				data_obj.Price = String2Float(obj.Data.IndexPrice)

				array = append(array, data_obj)

				data_str, _ := json.Marshal(array)
				inside_obj.Data = string(data_str)
				json_data, err := json.Marshal(inside_obj)
				if err != nil {
					glog.Error("Marshal failed:", err)
					return
				}
				httpPost(string(json_data))
			}
		}
	}
}

func WatchHandle() {
	done := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go WsChannel(done)

	for {
		select {
		case <-done:
			glog.Error("recv done....")
			done = make(chan struct{})
			go WsChannel(done)
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

	watch_trade_map = make(map[int]TickerWatchObj)
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
		watch_trade_map[v.Contract.ContractID] = TickerWatchObj{
			ContractID:    v.Contract.ContractID,
			Currency_name: v.Contract.BaseCoin,
			Contract_name: v.Contract.DisplayName,
		}
	}
	go WatchHandle()
	for {
		time.Sleep(1 * time.Second)
	}
}
