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

func httpPost(id int, param string) {
	glog.Info("start task id:", id)
	defer glog.Info("end task id:", id)
	data := fmt.Sprintf("data=%s", param)
	url := config.Get("contract_sdk_agent.addr")
	conn := &http.Client{
		Timeout: 5 * time.Second,
	}

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		go func() {
			resp, err := conn.Post(url,
				"application/x-www-form-urlencoded",
				strings.NewReader(data))
			if err != nil {
				glog.Error(err)
				return
			}

			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				glog.Error("ReadAll err:", err)
				return
			}
		}()
	}
}

func Init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "200")
	flag.Set("v", "3")
}

func main() {
	Init()
	defer glog.Flush()
	flag.Parse()

	base_data := BaseData{
		DataType:       "index",
		Exchange:       "HUOBI",
		ContractIndex:  1,
		ContractName:   "BTC_CQ",
		ContractSymbol: "BTC",
	}

	index_obj := make(IndexObj, 100)
	for i := 0; i < 100; i++ {
		index_obj[i].Timestamp = uint64(time.Now().Unix())
		index_obj[i].Price = 4800.1 + float64(i)
	}
	obj_data, _ := json.Marshal(index_obj)
	base_data.Data = string(obj_data)
	data, _ := json.Marshal(base_data)

	for i := 0; i < 1000; i++ {
		go httpPost(i, string(data))
	}
	for {
		time.Sleep(30 * time.Second)
	}
	fmt.Println("test")
}
