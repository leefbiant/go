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

type DataCache struct {
	Data common.QueueBigOrderStruct
	Time uint64
}

var (
	r_mq_queue string
	w_mq_queue string

	user_push_msg_chan = make(chan common.PushSruct, 10)
	data_cache         = make(map[string]DataCache)
)

func DonePush(cache *common.PushUserCache, c chan common.ServerQueueStruct,
	user_push_msg_chan *chan common.PushSruct) {
	push_interval := uint64(bbexutil.String2Int(config.Get("push_interval.bigorder"), 300))
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueueBigOrderStruct
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
			if obj.Price*obj.AmountSymbol < 100000 {
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
			Side := fmt.Sprintf("空单")
			if strings.ToUpper(obj.Side) == "BUY" {
				Side = fmt.Sprintf("多单")
			}

			template := "【委托类型】%s\r\n【委托价格】%.2f\r\n【委托张数】%d张\r\n【折算价值】%.2f %s"

			push_msg.Msg = fmt.Sprintf(template, Side, obj.Price, obj.Size, obj.AmountSymbol, queue_obj.Symbol)

			// 检查是否重复消息
			now_t := uint64(time.Now().Unix())
			cache_key := fmt.Sprintf("%s|%.2f", keys, obj.Price)
			cache_obj, ok := data_cache[cache_key]
			if !ok {
				data_cache[cache_key] = DataCache{Data: obj, Time: now_t}
			} else {
				if now_t-cache_obj.Time < 300 || obj.Size < cache_obj.Data.Size {
					glog.Error("mabye recv a Repeated msg please check key:", cache_key, " data:", queue_obj)
					continue
				}
			}
			// end check
			// sync bbx if bbx notify
			if bbexutil.GetExchangeName(queue_obj.Exchange) == "BBX" {
				bbx_msg := fmt.Sprintf("【大额挂单】%s盘口出现数量为%d张的挂单，挂单类型：%s", queue_obj.ContractName, obj.Size, Side)
				nodify_obj := common.BBXNotify{
					Symbol:       queue_obj.Symbol,
					ContractName: queue_obj.ContractName,
					Type:         queue_obj.SubType,
					Time:         queue_obj.Time,
					Msg:          bbx_msg,
				}
				bbexutil.SyncBBx(nodify_obj)
			}
			push_msg.Url = bbexutil.GetContractUrl(queue_obj)

			user_map := cache.GetCache(keys)
			if user_map == nil {
				continue
			}

			for _, user_info := range user_map.EventInfo {
				if now < user_info.NextPushTime {
					glog.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
					continue
				}
				if math.Abs(float64(obj.AmountSymbol)*obj.Price) > math.Abs(user_info.Val1*10000) {

					push_msg.Bpush = true
					// 不推送 但是 存储
					if user_info.Night_push == "0" && hour < 8 || IsDelay {
						push_msg.Bpush = false
					}

					push_msg.Openid = user_info.Openid
					push_msg.UID = user_info.Id
					bbexutil.MakePushContent(&push_msg)
					*user_push_msg_chan <- push_msg

					// doing
					glog.Info("send push 2 id:", user_info.Id, " openid:", user_info.Openid, " msg:", push_msg, " channle len:", cap(*user_push_msg_chan))
					user_info.LastPushTime = now
					user_info.NextPushTime = user_info.LastPushTime + push_interval
				}
			}

			// 删除老的数据
			for k, cache_obj := range data_cache {
				if now_t-cache_obj.Time > 600 {
					delete(data_cache, k)
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
	msgchan := make(chan common.ServerQueueStruct, 10)

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
				msgchan := make(chan common.ServerQueueStruct, 10)
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
	r_mq_queue = config.Get("mq.queue_bigorder")
	var cache common.PushUserCache
	err := cache.Init("sub_notify_bigorder", 4)
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
