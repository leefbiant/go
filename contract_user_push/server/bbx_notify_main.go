package main

import (
	"bbexgo/bbexutil"
	"bbexgo/config"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"time"
)

var (
	r_mq_queue string
	bbx_notify = config.Get("contract.bbx_notify_addr")
)

func SyncBBx(msg string) {
	glog.Info("SyncBBx data:", msg)
	bbexutil.HttpPost(msg, bbx_notify)
	return
}

func Handle() {
	glog.Info("MsgPushHandle start ....")
	done := make(chan struct{})
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	msgchan := make(chan string, 10)

	go bbexutil.WatchMq(done, msgchan, r_mq_queue)
	for {
		select {
		case msg := <-msgchan:
			SyncBBx(msg)
		case <-ticker.C:
			glog.Info("MsgPushHandle live .....\n")
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
	r_mq_queue = config.Get("mq.queue_bbx_notify")
	go Handle()
	for {
		time.Sleep(30 * time.Second)
	}
}
