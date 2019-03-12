package main

import (
	"bbexgo/config"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"strings"
	// "time"
)

var (
	queue_map map[string]chan string
)

type ResObj struct {
	Code       int64
	Msg        string
	ReturnData interface{}
}

func HttpChannle(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) == 2 {
		w.Write([]byte("path error."))
	}
	// r.ParseForm()
	r.ParseMultipartForm(2 << 20)
	params := make(map[string]string)
	for k, v := range r.Form {
		params[k] = strings.Join(v, "")
	}

	defer func() {
		if err := recover(); err != nil {
			glog.Error(err)
			var res []byte
			res, _ = json.Marshal(&responseStruct{
				Code: -1,
				Msg:  fmt.Sprintf("system error: [%s]", err),
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	}()

	obj := &ResObj{}
	obj.PushMsg(params)
	result := *obj
	if result.ReturnData == nil {
		result.ReturnData = ""
	}

	responbyte, _ := json.Marshal(&responseStruct{
		Code: result.Code,
		Msg:  result.Msg,
		Data: result.ReturnData,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(responbyte)
}

func (this *ResObj) PushMsg(params map[string]string) {
	this.Code = -100
	this.Msg = fmt.Sprintf("not find data from params")
	req_data, ok := params["data"]
	if !ok {
		glog.Error("not find data in req")
		return
	}
	var base_data BaseData
	err := json.Unmarshal([]byte(req_data), &base_data)
	if err != nil {
		glog.Error("Unmarshal failed:", err)
		return
	}
	switch base_data.DataType {
	case "trade":
		var trade_obj TradeObj
		err = json.Unmarshal([]byte(base_data.Data), &trade_obj)
		if err != nil {
			glog.Error("Unmarshal trade failed:", err, " data:", base_data.Data)
			this.Msg = fmt.Sprintf("Unmarshal trade failed:%v", err)
			return
		}
		break
	case "index":
		var obj IndexObj
		err = json.Unmarshal([]byte(base_data.Data), &obj)
		if err != nil {
			glog.Error("Unmarshal index failed:", err, " data:", base_data.Data)
			this.Msg = fmt.Sprintf("Unmarshal index failed:%v", err)
			return
		}
	case "ticker":
		var obj TickerObj
		err = json.Unmarshal([]byte(base_data.Data), &obj)
		if err != nil {
			glog.Error("Unmarshal ticker failed:", err, " data:", base_data.Data)
			this.Msg = fmt.Sprintf("Unmarshal ticker failed:%v", err)
			return
		}
	case "depth":
		var obj DepthObj
		err = json.Unmarshal([]byte(base_data.Data), &obj)
		if err != nil {
			glog.Error("Unmarshal depth failed:", err, " data:", base_data.Data)
			this.Msg = fmt.Sprintf("Unmarshal depth failed:%v", err)
			return
		}
	case "blasting":
		var obj BlastingObj
		err = json.Unmarshal([]byte(base_data.Data), &obj)
		if err != nil {
			glog.Error("Unmarshal blasting failed:", err, " data:", base_data.Data)
			this.Msg = fmt.Sprintf("Unmarshal blasting failed:%v", err)
			return
		}
	}

	msg_chan, ok := queue_map[base_data.DataType]
	if !ok {
		glog.Info("not find msg chan:", base_data.DataType)
		this.Msg = fmt.Sprintf("not find data type:%s", base_data.DataType)
		return
	}
	glog.Info("sucess recv DataType:", base_data.DataType, " Exchange:", base_data.Exchange, " ContractName:", base_data.ContractName)
	this.Code = 0
	this.Msg = "OK"
	msg_chan <- req_data
}

func Init() {
	queue_map = make(map[string]chan string)
	queue_map["trade"] = make(chan string, 100)
	queue_map["ticker"] = make(chan string, 100)
	queue_map["depth"] = make(chan string, 100)
	queue_map["index"] = make(chan string, 100)
	queue_map["blasting"] = make(chan string, 100)

	for k, v := range queue_map {
		go DataDoneHandle(k, v)
	}

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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sdk_agent...."))
	})
	mux.HandleFunc("/contract_sdk_agent/PushMsg", HttpChannle)

	port := config.Get("sdk_agent.listen_port")
	server := &http.Server{
		Addr: port,
		// ReadTimeout:  300 * time.Second,
		// WriteTimeout: time.Second * 300,
		Handler: mux,
	}
	glog.Info(server.ListenAndServe())
}
