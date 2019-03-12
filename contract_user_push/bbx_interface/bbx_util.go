package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ContractAttr struct {
	Type           int
	Name           string
	SinglePriceUSD float64
}

var (
	//single_map = make(map[string]ContractAttr)
	trade_map = make(map[string]ContractAttr)
)

type InterfaceInside struct {
	DataType       string  `json:"dataType"`
	Exchange       string  `json:"exchange"`
	ContractIndex  int     `json:"contractIndex"`
	ContractName   string  `json:"contractName"`
	ContractSymbol string  `json:"contractSymbol"`
	Data           string  `json:"data"`
	SinglePriceUSD float64 `json:"singlePriceUSD"` // 单张价格（USD）
}

/////////////////////////////////////

type Contracts struct {
	Errno   string `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		Contracts []struct {
			Contract struct {
				ContractID   int    `json:"contract_id"`
				IndexID      int    `json:"index_id"`
				Name         string `json:"name"`
				DisplayName  string `json:"display_name"`
				BaseCoin     string `json:"base_coin"`
				ContractSize string `json:"contract_size"`
			} `json:"contract"`
		} `json:"contracts"`
	} `json:"data"`
}

type TradeObj struct {
	Group string `json:"group"`
	Data  []struct {
		TradeID       int64     `json:"trade_id"`
		ContractID    int       `json:"contract_id"`
		SellAccountID int64     `json:"sell_account_id"`
		BuyAccountID  int64     `json:"buy_account_id"`
		SellOrderID   int64     `json:"sell_order_id"`
		BuyOrderID    int64     `json:"buy_order_id"`
		DealPrice     string    `json:"deal_price"`
		DealVol       string    `json:"deal_vol"`
		Fluctuation   string    `json:"fluctuation"`
		MakeFee       string    `json:"make_fee"`
		TakeFee       string    `json:"take_fee"`
		CreatedAt     time.Time `json:"created_at"`
		Way           int       `json:"way"`
		Type          int       `json:"type"`
	} `json:"data"`
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

//////////////////////////////////////////////

type TickerWatchObj struct {
	ContractID    int
	Currency_name string // 币种
	Contract_name string // 合约名字
}

type TickerObj struct {
	Group string `json:"group"`
	Data  struct {
		LastPrice    string `json:"last_price"`
		AvgPrice     string `json:"avg_price"`
		Volume       string `json:"volume"`
		TotalVolume  string `json:"total_volume"`
		Timestamp    int    `json:"timestamp"`
		ContractID   int    `json:"contract_id"`
		PremiumIndex string `json:"premium_index"`
		FundingRate  string `json:"funding_rate"`
		IndexPrice   string `json:"index_price"`
		RiseFallRate string `json:"rise_fall_rate"`
	} `json:"data"`
}

type TickerDataObj struct {
	Timestamp    uint64  `json:"timestamp"`
	BuyPrice     float64 `json:"buyPrice"`
	BuySize      uint64  `json:"buySize"`
	SellPrice    float64 `json:"sellPrice"`
	SellSize     uint64  `json:"sellSize"`
	RiseFallRate float64 `json:"riseFallRate"`
}

type TickerObjInside struct {
	DataType       string `json:"dataType"`
	Exchange       string `json:"exchange"`
	ContractIndex  int    `json:"contractIndex"`
	ContractName   string `json:"contractName"`
	ContractSymbol string `json:"contractSymbol"`
	Data           string `json:"data"`
}

////////////////////////////////////
type IndexDataObj struct {
	Timestamp uint64  `json:"timestamp"`
	Price     float64 `json:"price"`
}

/////////////////////////////////////
type DepthObj struct {
	Group string `json:"group"`
	Data  struct {
		Way    int `json:"way"`
		Depths []struct {
			Price string `json:"price"`
			Vol   string `json:"vol"`
		} `json:"depths"`
	} `json:"data"`
}

type DepthData struct {
	Idx   string  `json:"idx"`
	Price float64 `json:"price"`
	Size  uint    `json:"size"`
	Type  string  `json:"type"`
}

type DepthDataSlice []DepthData

type DepthInsideObj struct {
	Sells struct {
		bset   bool           `json:""`
		Action string         `json:"action"`
		List   DepthDataSlice `json:"list"`
	} `json:"sells"`
	Buys struct {
		bset   bool           `json:""`
		Action string         `json:"action"`
		List   DepthDataSlice `json:"list"`
	} `json:"buys"`
}

func (s DepthDataSlice) Len() int           { return len(s) }
func (s DepthDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s DepthDataSlice) Less(i, j int) bool { return s[i].Price > s[j].Price }

func String2Float(value string) float64 {
	if s, err := strconv.ParseFloat(value, 64); err == nil {
		return s
	}
	glog.Error("err value:", value)
	return 0.0
}

func String2Int(value string) int64 {
	if s, err := strconv.ParseInt(value, 10, 64); err == nil {
		return s
	}
	glog.Error("err value:", value)
	return 0
}

func GetContracts() (error, *Contracts) {
	url := "https://api.bbx.com/v1/ifcontract/contracts"
	conn := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := conn.Get(url)
	if err != nil {
		glog.Error(err)
		return err, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("ReadAll err:", err)
		return err, nil
	}

	var obj Contracts
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		glog.Error("Unmarshal failed:", err)
		return err, nil
	}
	for _, v := range obj.Data.Contracts {
		obj_set, ok := trade_map[v.Contract.DisplayName]
		if !ok {
			continue
		}
		obj_set.Name = v.Contract.Name
		obj_set.SinglePriceUSD = String2Float(v.Contract.ContractSize)
		trade_map[v.Contract.DisplayName] = obj_set
	}
	glog.Info(obj)
	return nil, &obj
}

func InitContractMap() {
	trade_map["BTCUSDT永续合约"] = ContractAttr{Type: 0}
	trade_map["ETHUSDT永续合约"] = ContractAttr{Type: 0}
	trade_map["EOSUSDT永续合约"] = ContractAttr{Type: 0}
	trade_map["XRPUSDT永续合约"] = ContractAttr{Type: 0}
	trade_map["BTCUSD反向永续"] = ContractAttr{Type: 1}
	trade_map["ETHUSD反向永续"] = ContractAttr{Type: 1}
	trade_map["EOSUSDT反向永续"] = ContractAttr{Type: 1}
}

func GetSinglePriceUSD(name string, price float64) float64 {
	obj, ok := trade_map[name]
	if !ok {
		glog.Error("not find data for Contract name:", name)
		return 0.0
	}
	if obj.Type == 1 {
		return obj.SinglePriceUSD
	}
	return price * obj.SinglePriceUSD
}

func GetContractType(name string) int {
	obj, ok := trade_map[name]
	if !ok {
		glog.Error("not find data for Contract name:", name)
		return 0
	}
	return obj.Type
}

func GetContractName(name string) string {
	obj, ok := trade_map[name]
	if !ok {
		glog.Error("not find data for Contract name:", name)
		return ""
	}
	return obj.Name
}
