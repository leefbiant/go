package common

import (
	"bbexgo/config"
	"bbexgo/help"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 通用数据定义
type ServerQueueStruct struct {
	Type         int    `json:"type"`         // 消息类型id
	SubType      int    `json:"subType"`      // 订阅类型
	Exchange     int    `json:"exchange"`     // 交易所ID
	Symbol       string `json:"symbol"`       // 合约币种
	ContractType int    `json:"contractType"` // 合约类型
	ContractName string `json:"contractName"` // 完整合约名
	Time         int64  `json:"time"`         // 时间戳
	Msg          string `json:"msg"`          // 信息内容 json结构
}

// QueueSpotStruct 限价点位
type QueueSpotStruct struct {
	Price float64 `json:"price"` // 当前价格
	Side  string  `json:"side"`  // 价格方向
}

// QueueChangeStruct 涨跌幅
type QueueChangeStruct struct {
	Change     float64 `json:"change"`     // 涨跌幅
	StartTime  int64   `json:"startTime"`  // 开始时间
	StartPrice float64 `json:"startPrice"` // 开始价格
	EndTime    int64   `json:"endTime"`    // 结束时间
	EndPrice   float64 `json:"endPrice"`   // 结束价格
}

// QueuePremiumStruct 期现溢价
type QueuePremiumStruct struct {
	Premium     float64 `json:"premium"`     // 溢价率
	IndexPrice  float64 `json:"indexPrice"`  // 指数价格
	FuturePrice float64 `json:"futurePrice"` // 期货价格
}

// QueueBigOrderStruct 大额挂单
type QueueBigOrderStruct struct {
	Side         string  `json:"side"`         // 挂单方向（buy,sell）
	Price        float64 `json:"price"`        // 挂单价格
	Size         int64   `json:"size"`         // 挂单数量(张)
	AmountSymbol float64 `json:"amountSymbol"` // 挂单总量(币)
}

// QueueBigTradeStruct 大额交易
type QueueBigTradeStruct struct {
	Price        float64 `json:"price"`        // 成交价格
	AmountUSD    float64 `json:"amountUSD"`    // 成交总额(USD)
	Size         int64   `json:"size"`         // 成交张数
	AmountSymbol float64 `json:"amountSymbol"` // 挂单总量(币)
	Side         string  `json:"side"`         // 挂单方向（buy,sell）
}

// QueueBlastingOrderStruct 爆仓订单
type QueueBlastingOrderStruct struct {
	AmountUSD    float64 `json:"amountUSD"`    // 爆仓金额(USD)
	AmountSymbol float64 `json:"amountSymbol"` // 委托单量(币)
	Type         string  `json:"type"`         // 交易类型
	Price        float64 `json:"price"`        // 委托价格
	Size         int64   `json:"size"`         // 爆仓张数
}

// QueueDeliveryDateStruct 交割时间
type QueueDeliveryDateStruct struct {
	Dif          int64 `json:"dif"`          // 距离价格时间差值（秒）
	DeliveryTime int64 `json:"deliveryTime"` // 交割时间
}

type QueueEveryDayMarketStruct struct {
	TraceSymbol         string  `json:"traceSymbol"`         // 追踪币种
	IndexPrice          float64 `json:"indexPrice"`          // 指数价格
	Change24Houer       float64 `json:"change24Hour"`        // 24小时波动
	ShortChange         int     `json:"shortChange"`         // 每日超过1%的次数
	ContractHold        float64 `json:"contractHold"`        // 合约持仓量
	ContractHoldChange  float64 `json:"contractHoldChange"`  // 合约持仓相比昨日涨跌幅
	ContractTrade       float64 `json:"contractTrade"`       // 合约交易量
	ContractTradeChange float64 `json:"contractTradeChange"` // 合约交易量相比昨日涨跌幅
	ContractRatioMore   float64 `json:"contractRatioMore"`   // 合约做多占比
	COntractRatioEmpty  float64 `json:"contractRatioEmpty"`  // 合约做空占比
}

// 消息队列中push消息通用结构
type PushData struct {
	// UID          string `json:"uid"`          // 用户id
	Type         int    `json:"type"`         // 消息类型id
	SubType      int    `json:"subType"`      // 订阅类型
	Exchange     int    `json:"exchange"`     // 交易所ID
	Symbol       string `json:"symbol"`       // 合约币种
	ContractType int    `json:"contractType"` // 合约类型
	ContractName int    `json:"contractName"` // 完整合约名
	Time         int64  `json:"time"`         // 时间戳
	Msg          string `json:"msg"`          // 信息内容
}

// push数据结构, 继承monitor.PushData 并添加openID
type PushSruct struct {
	ServerQueueStruct
	UID            string
	Openid         string
	PushStatusCode int64                        // 推送后状态码
	PushErrorMsg   string                       // 推送错误信息
	WechatMsgId    string                       // 推送后微信返回的ID
	Bpush          bool                         // 是否推送
	Pusg_msg       map[string]map[string]string // 最终推送的微信数据
	Url            string
	Title          string
}

//////////////////////////////////////////////////////////
// push 用户结构
type UserSubInfo struct {
	Id           string
	Openid       string
	LastPushTime uint64
	NextPushTime uint64
	Val1         float64
	Val2         float64
	Night_push   string
}

type BBXNotify struct {
	Symbol       string `json:"symbol"`       // 合约币种
	ContractName string `json:"contractName"` // 完整合约名
	Type         int    `json:"type"`         // 消息类型id
	Time         int64  `json:"time"`         // 时间戳
	Msg          string `json:"msg"`          // 信息内容 json结构
}

type EventCache struct {
	EventInfo map[string]*UserSubInfo
}

// push 用户缓存
type PushUserCache struct {
	L        *sync.Mutex
	Db       *sql.DB
	Cahce    map[string]*EventCache
	Queue    string
	sub_type int
}

func (this *PushUserCache) Init(queue string, sub_type int) error {
	constr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", config.Get("mysql.default.username"),
		config.Get("mysql.default.password"), config.Get("mysql.default.host"), config.Get("mysql.default.port"),
		config.Get("mysql.default.database"), config.Get("mysql.default.charset"))

	db, err := sql.Open("mysql", constr)
	if err != nil {
		glog.Fatal(err)
		return err
	}
	this.Db = db
	this.L = new(sync.Mutex)
	this.Queue = queue
	this.sub_type = sub_type
	return nil
}

func (this *PushUserCache) GetCache(key string) *EventCache {
	this.L.Lock()
	defer this.L.Unlock()
	obj, ok := this.Cahce[key]
	if ok {
		// event_cache := new(EventCache)
		// event_cache.EventInfo = make(map[string]*UserSubInfo)
		// event_cache.EventInfo = obj.EventInfo
		// return event_cache
		return obj
	}
	// log.Error("not find cahce for key:", key)
	return nil
}

func (this *PushUserCache) UpdateCache() {
	t_now := time.Now().Unix()
	sql := fmt.Sprintf("select t1.id, t1.openid, t2.exchange, t2.symbol, t2.contract_type, t2.sub_type, t2.val1, t2.val2, t1.night_push, t1.wechat_push, t1.report_push	from user_subs as t2 join users as t1 on t2.uid = t1.id where t1.subscribe_status = 1  and t2.status = 1 and sub_type = %d and t2.off_remind < %d", this.sub_type, t_now)

	if -1 == this.sub_type {
		sql = fmt.Sprintf("select t1.id, t1.openid, 'exchange', 'symbol', 1, 1, 0, 0, t1.night_push, t1.wechat_push, t1.report_push	from users as t1 where t1.subscribe_status = 1 and t1.report_push = 1", t_now)
	}
	glog.Info("sql:", sql)
	list, err := help.QueryFormMysql(this.Db, sql)
	if err != nil {
		glog.Error("QueryFormMysql failed:", err)
		return
	}
	cache := make(map[string]*EventCache, 0)
	for _, v := range *list {
		// 是否开启PUSH
		if v["wechat_push"] != "1" {
			continue
		}

		key := fmt.Sprintf("%s|%s|%s|%s", v["sub_type"], v["exchange"], strings.ToUpper(v["symbol"]), v["contract_type"])
		if -1 == this.sub_type {
			key = fmt.Sprintf("%s", v["openid"])
		}

		_, ok := cache[key]
		if !ok {
			new_obj := new(EventCache)
			new_obj.EventInfo = make(map[string]*UserSubInfo)
			cache[key] = new_obj
		}
		Val1, _ := strconv.ParseFloat(v["val1"], 64)
		Val2, _ := strconv.ParseFloat(v["val2"], 64)

		_, ok = cache[key].EventInfo[v["openid"]]
		if !ok {
			cache[key].EventInfo[v["openid"]] = &UserSubInfo{
				Id:           v["id"],
				Openid:       v["openid"],
				LastPushTime: 0,
				NextPushTime: 0,
				Val1:         Val1,
				Val2:         Val2,
				Night_push:   v["night_push"],
			}
		}
		// log.Info("add key %s user:%s", key, v["openid"])
	}
	// 更新缓存
	this.L.Lock()
	defer this.L.Unlock()
	for k, v := range this.Cahce {
		obj, ok := cache[k]
		if !ok {
			glog.Info("not find key:", k)
			continue
		}
		for openid, info := range v.EventInfo {
			cache_user, ok := obj.EventInfo[openid]
			if !ok {
				glog.Error("not find user:", openid)
				continue
			}
			// 更新当前的push时间
			cache_user.LastPushTime = info.LastPushTime
			cache_user.NextPushTime = info.NextPushTime
			//glog.Info("update cache for key:", k, " openid:", openid)
		}
	}
	this.Cahce = cache
	glog.Info("update cache:", len(cache))
	// for k, v := range cache {
	// 	glog.Info("key:", k, " num:", len(v.EventInfo))
	// }
}

func (this *PushUserCache) DisEnableSub(uid uint64, sub_type int, exchange string,
	symbol string, contract_type int) error {

	stmt, err := this.Db.Prepare("update user_subs set status = 0 where uid =? and sub_type = ? and exchange = ? and symbol = ? and contract_type =?")
	if err != nil {
		glog.Error("Prepare failed:", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(uid, sub_type, exchange, symbol, contract_type)
	if err != nil {
		glog.Error("Prepare Exec:", err)
		return err
	}
	num, err := res.RowsAffected()
	if num != 1 {
		glog.Error("Prepare Exec:", err)
		return err
	}
	return nil
}

func (this *PushUserCache) UpdateUserNextNotify(uid uint64, sub_type int, exchange string,
	symbol string, contract_type int, next_notify_time uint64) error {

	stmt, err := this.Db.Prepare("update user_subs set off_remind = ? where uid =? and sub_type = ? and exchange = ? and symbol = ? and contract_type =?")
	if err != nil {
		glog.Error("Prepare failed:", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(next_notify_time, uid, sub_type, exchange, symbol, contract_type)
	if err != nil {
		glog.Error("Prepare Exec:", err)
		return err
	}
	num, err := res.RowsAffected()
	if num != 1 {
		glog.Error("Prepare Exec:", err)
		return err
	}
	return nil
}

func UpdateCacheHandle(cache *PushUserCache) {
	done := make(chan struct{})
	notify := make(chan int)
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	cache.UpdateCache()
	go WatchUpdateChacheSignal(done, cache.Queue, notify)
	for {
		select {
		case <-done:
			glog.Info("recv sig 2 update cache")
			done = make(chan struct{})
			time.Sleep(1 * time.Second)
			go WatchUpdateChacheSignal(done, cache.Queue, notify)

		case <-notify:
			cache.UpdateCache()

		case t := <-ticker.C:
			t = t
			cache.UpdateCache()
			break
		}
	}
}

func WatchUpdateChacheSignal(done chan struct{}, queue string, notify chan int) {
	defer close(done)
	mq_addr := config.Get("mq.addr")

	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Error("amqp failed:%v", err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		glog.Error("Channel failed:%v", err)
		return
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"sub_notify", // name
		"fanout",     // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		glog.Error("QueueDeclare failed:%v", err)
		return
	}

	ch.QueueBind(queue, "", "sub_notify", false, nil)

	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		glog.Error("Consume failed:", err)
		return
	}
	for d := range msgs {
		glog.Error("Received a message: ", d.Body)
		notify <- 1
	}
}
