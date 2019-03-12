package main

import (
	"bbexgo/bbexutil"
	"bbexgo/common"
	"bbexgo/config"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"math"
	"strings"
	"time"
)

var (
	r_mq_queue string
	w_mq_queue string

	user_push_msg_chan = make(chan common.PushSruct, 10)
)

func DonePush(cache *common.PushUserCache, c chan common.ServerQueueStruct,
	user_push_msg_chan *chan common.PushSruct) {
	push_interval := uint64(bbexutil.String2Int(config.Get("push_interval.change"), 300))
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueueChangeStruct
			err := json.Unmarshal([]byte(queue_obj.Msg), &obj)
			if err != nil {
				glog.Error("Unmarshal failed:", err)
				continue
			}

			Exchange := bbexutil.GetExchangeName(queue_obj.Exchange)
			if Exchange == "" {
				glog.Error("not find GetExchangeName for id:", queue_obj.Exchange)
				continue
			}
			keys := fmt.Sprintf("%d|%s|%s|%d", queue_obj.SubType, Exchange, queue_obj.Symbol, queue_obj.ContractType)
			user_map := cache.GetCache(keys)
			if user_map == nil {
				continue
			}
			now := uint64(time.Now().Unix())
			push_msg := common.PushSruct{
				ServerQueueStruct: queue_obj,
				Bpush:             true,
			}

			IsDelay := false
			if time.Now().Unix()-queue_obj.Time > 60*30 {
				IsDelay = true
			}

			hour := time.Now().Hour()

			side := "上涨"
			push_msg.Title = fmt.Sprintf("%s急涨提醒\n", strings.ToUpper(queue_obj.Symbol))
			if obj.Change < 0 {
				side = "下跌"
				push_msg.Title = fmt.Sprintf("%s急跌提醒\n", strings.ToUpper(queue_obj.Symbol))
			}
			side = fmt.Sprintf("%s%.2f%%", side, math.Abs(obj.Change))

			template := "【涨跌情况】%s\r\n【当前价格】%g"
			push_msg.Msg = fmt.Sprintf(template, side, obj.EndPrice)

			for _, user_info := range user_map.EventInfo {
				if now < user_info.NextPushTime {
					glog.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
					continue
				}

				if 0 == user_info.Val1 {
					user_info.Val1 = 1
				}

				user_set := user_info.Val2 / user_info.Val1                              // 用户设置的每分钟涨跌幅
				envent_set := obj.Change / (float64((obj.EndTime - obj.StartTime)) / 60) // 数据中每分钟涨跌幅

				glog.Info("user:", user_info.Openid, " set event:", math.Abs(user_set), " data event:", math.Abs(envent_set))
				if math.Abs(envent_set) > math.Abs(user_set) {

					push_msg.Bpush = true
					// 不推送 但是 存储
					if user_info.Night_push == "0" && hour < 8 || IsDelay {
						push_msg.Bpush = false
					}

					push_msg.Openid = user_info.Openid
					push_msg.UID = user_info.Id
					push_msg.Url = bbexutil.GetContractUrl(queue_obj)
					bbexutil.MakePushContent(&push_msg)
					*user_push_msg_chan <- push_msg

					// doing
					glog.Info("send push 2 id:", user_info.Id, " openid:", user_info.Openid, " msg:", push_msg, " channle len:", cap(*user_push_msg_chan))

					if push_msg.Bpush {
						user_info.LastPushTime = now
						user_info.NextPushTime = user_info.LastPushTime + push_interval
						glog.Info("user:", user_info.Openid, " key:", keys, " send push next time:", user_info.NextPushTime)
					}
				}
			}
		}
	}
	return
}

func MsgPushHandle(cache *common.PushUserCache) {
	glog.Info("MsgPushHandle start ....")
	done := make(chan struct{})
	send_mq_chan := make(chan struct{})

	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	cache.UpdateCache()
	msgchan := make(chan common.ServerQueueStruct, 100)

	event_chan := make(map[int]chan common.ServerQueueStruct)

	go bbexutil.WatchPushMsg(done, msgchan, r_mq_queue)
	go bbexutil.WritePushMsg(send_mq_chan, &user_push_msg_chan, w_mq_queue)
	for {
		select {
		case <-done:
			cache.UpdateCache()
			done = make(chan struct{})
			time.Sleep(1 * time.Second)
			go bbexutil.WatchPushMsg(done, msgchan, r_mq_queue)

		// 写push到Mq异常
		case <-send_mq_chan:
			time.Sleep(1 * time.Second)
			send_mq_chan = make(chan struct{})
			go bbexutil.WritePushMsg(send_mq_chan, &user_push_msg_chan, w_mq_queue)

		case msg := <-msgchan:
			c, ok := event_chan[msg.SubType]
			if !ok {
				msgchan := make(chan common.ServerQueueStruct, 100)
				event_chan[msg.SubType] = msgchan
				glog.Info("create a new chan for type:", msg.SubType)
				go DonePush(cache, msgchan, &user_push_msg_chan)
			} else {
				c <- msg
			}

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
	w_mq_queue = config.Get("mq.queue_user_push")
	r_mq_queue = config.Get("mq.queue_change")
	var cache common.PushUserCache
	err := cache.Init("sub_notify_change", 2)
	if err != nil {
		glog.Error("InitDB failed:%v", err)
		return
	}
	// 更新用户数据缓存
	go common.UpdateCacheHandle(&cache)
	go MsgPushHandle(&cache)
	for {
		time.Sleep(30 * time.Second)
	}
}
