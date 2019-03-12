package main

import (
	"bbexgo/config"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"reflect"
	"time"
	"unicode"
)

type BaseData struct {
	DataType       string `json:"dataType"`
	Exchange       string `json:"exchange"`
	ContractIndex  int    `json:"contractIndex"`
	ContractName   string `json:"contractName"`
	ContractSymbol string `json:"contractSymbol"`
	Data           string `json:"data"`
}

type TradeObj []struct {
	Timestamp uint64  `json:"timestamp"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
}

type IndexObj []struct {
	Timestamp uint64  `json:"timestamp"`
	Price     float64 `json:"price"`
}

type TickerObj []struct {
	Timestamp uint64  `json:"timestamp"`
	BuyPrice  float64 `json:"buyPrice"`
	BuySize   float64 `json:"buySize"`
	SellPrice float64 `json:"sellPrice"`
	SellSize  float64 `json:"sellSize"`
}

type DepthObj []struct {
	Sells struct {
		Action string `json:"action"`
		List   []struct {
			Idx   string  `json:"idx"`
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Type  string  `json:"type"`
		} `json:"list"`
	} `json:"sells"`
	Buys struct {
		Action string `json:"action"`
		List   []struct {
			Idx   string  `json:"idx"`
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Type  string  `json:"type"`
		} `json:"list"`
	} `json:"buys"`
}

/////////////////////////

type BlastingObj []struct {
	TradeType string  `json:"tradeType"`
	AmountUSD float64 `json:"amountUSD"`
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Count     float64 `json:"count"`
}

type responseStruct struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func ReflectInterface(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}
}

func WriteMq(c chan struct{}, data_chan chan string, queuq_name string) {
	defer close(c)
	mq_addr := config.Get("mq.addr")
	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Error("amqp.Dial failed:", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		glog.Error("conn.Channel failed:", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"interface_chan", // name
		"direct",         // durable
		true,             // delete when unused
		false,            // exclusive
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		glog.Error("ExchangeDeclare failed:", err, " mq_addr:", mq_addr, " queue:", queuq_name)
		return
	}

	_, err = ch.QueueDeclare(
		queuq_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		glog.Error("QueueDeclare.Dial failed:", err, " mq_addr:", mq_addr, " queue:", queuq_name)
		return
	}

	for {
		select {
		case push_obj := <-data_chan:
			body := string(push_obj)
			err = ch.Publish(
				"interface_chan", // exchange
				queuq_name,       // routing key
				false,            // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			if err != nil {
				glog.Info("Publish 2 queue:", queuq_name, " failed err:", err)
				return
			}
			// glog.Info("send msg 2 queue:", queuq_name)
		}
	}
	glog.Info("WriteMq exit ....")
}

func DataDoneHandle(queue_name string, msgchan chan string) {
	glog.Info("DataDoneHandle start queue:", queue_name)
	done := make(chan struct{})
	go WriteMq(done, msgchan, queue_name)
	for {
		select {
		case <-done:
			done = make(chan struct{})
			glog.Info("DataDoneHandle WriteMq exit start")
			time.Sleep(1 * time.Second)
			go WriteMq(done, msgchan, queue_name)
		}
	}
}
