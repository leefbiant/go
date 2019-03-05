package bitmax_api

import (
	"time"
)

type ApiKey struct {
	Key       string `json:"ApiKey"`
	SecretKey string `json:"SecretKey"`
}

type Positions []struct {
	Account       float64 `json:"account"`
	Symbol        string  `json:"symbol"`
	Currency      string  `json:"currency"`
	Underlying    string  `json:"underlying"`
	QuoteCurrency string  `json:"quoteCurrency"`
	Commission    float64 `json:"commission"` // 佣金
	/* 	InitMarginReq        float64   `json:"initMarginReq"`
	   	MaintMarginReq       float64   `json:"maintMarginReq"`
	   	RiskLimit            float64   `json:"riskLimit"` */
	Leverage float64 `json:"leverage"` // 当前杠杆倍数
	/* 	CrossMargin          bool      `json:"crossMargin"`
	   	DeleveragePercentile float64   `json:"deleveragePercentile"`
	   	RebalancedPnl        float64   `json:"rebalancedPnl"`
	   	PrevRealisedPnl      float64   `json:"prevRealisedPnl"`
	   	PrevUnrealisedPnl    float64   `json:"prevUnrealisedPnl"` */
	PrevClosePrice   float64   `json:"prevClosePrice"`   // 标记价格
	OpeningTimestamp time.Time `json:"openingTimestamp"` // 开仓时间
	OpeningQty       float64   `json:"openingQty"`       // 开仓杠杆数
	/* 	OpeningCost          float64   `json:"openingCost"` */
	OpeningComm float64 `json:"openingComm"` // 已实现盈亏

	/* OpenOrderBuyQty      float64   `json:"openOrderBuyQty"`
	OpenOrderBuyCost     float64   `json:"openOrderBuyCost"`
	OpenOrderBuyPremium  float64   `json:"openOrderBuyPremium"`
	OpenOrderSellQty     float64   `json:"openOrderSellQty"`
	OpenOrderSellCost    float64   `json:"openOrderSellCost"`
	OpenOrderSellPremium float64   `json:"openOrderSellPremium"`
	ExecBuyQty           float64   `json:"execBuyQty"`
	ExecBuyCost          float64   `json:"execBuyCost"`
	ExecSellQty          float64   `json:"execSellQty"`
	ExecSellCost         float64   `json:"execSellCost"`
	ExecQty              float64   `json:"execQty"`
	ExecCost             float64   `json:"execCost"`
	ExecComm             float64   `json:"execComm"`
	CurrentTimestamp     time.Time `json:"currentTimestamp"`
	*/
	CurrentQty float64 `json:"currentQty"`
	/*
		CurrentCost          float64   `json:"currentCost"`
		CurrentComm          float64   `json:"currentComm"`
		RealisedCost         float64   `json:"realisedCost"`
		UnrealisedCost       float64   `json:"unrealisedCost"`
		GrossOpenCost        float64   `json:"grossOpenCost"`
		GrossOpenPremium     float64   `json:"grossOpenPremium"`
		GrossExecCost        float64   `json:"grossExecCost"`
		IsOpen               bool      `json:"isOpen"`
		MarkPrice            float64   `json:"markPrice"`
		MarkValue            float64   `json:"markValue"`
		RiskValue            float64   `json:"riskValue"`
		HomeNotional         float64   `json:"homeNotional"`
		ForeignNotional      float64   `json:"foreignNotional"`
		PosState             string    `json:"posState"`
		PosCost              float64   `json:"posCost"`
		PosCost2             float64   `json:"posCost2"`
		PosCross             float64   `json:"posCross"`
		PosInit              float64   `json:"posInit"`
		PosComm              float64   `json:"posComm"`
		PosLoss              float64   `json:"posLoss"`
		PosMargin            float64   `json:"posMargin"`
		PosMaint             float64   `json:"posMaint"`
		PosAllowance         float64   `json:"posAllowance"`
		TaxableMargin        float64   `json:"taxableMargin"`
		InitMargin           float64   `json:"initMargin"` */

	MaintMargin float64 `json:"maintMargin"`

	/* SessionMargin        float64   `json:"sessionMargin"`
	TargetExcessMargin   float64   `json:"targetExcessMargin"`
	VarMargin            float64   `json:"varMargin"` */

	RealisedGrossPnl float64 `json:"realisedGrossPnl"` // 已实现盈亏
	/* 	RealisedTax          float64   `json:"realisedTax"`
	   	RealisedPnl          float64   `json:"realisedPnl"` */
	UnrealisedGrossPnl float64 `json:"unrealisedGrossPnl"` // 未实现盈亏

	/* LongBankrupt         float64   `json:"longBankrupt"`
	ShortBankrupt        float64   `json:"shortBankrupt"`
	TaxBase              float64   `json:"taxBase"`
	IndicativeTaxRate    float64   `json:"indicativeTaxRate"`
	IndicativeTax        float64   `json:"indicativeTax"`
	UnrealisedTax        float64   `json:"unrealisedTax"`
	UnrealisedPnl        float64   `json:"unrealisedPnl"`
	*/
	UnrealisedPnlPcnt    float64   `json:"unrealisedPnlPcnt"`
	UnrealisedRoePcnt    float64   `json:"unrealisedRoePcnt"`
	/*
	SimpleQty            float64   `json:"simpleQty"`
	SimpleCost           float64   `json:"simpleCost"`
	SimpleValue          float64   `json:"simpleValue"`
	SimplePnl            float64   `json:"simplePnl"`
	SimplePnlPcnt        float64   `json:"simplePnlPcnt"` */
	AvgCostPrice float64 `json:"avgCostPrice"`

	/* 	AvgEntryPrice        float64   `json:"avgEntryPrice"`
	   	BreakEvenPrice       float64   `json:"breakEvenPrice"` */
	MarginCallPrice float64 `json:"marginCallPrice"`
	/* 	LiquidationPrice     float64   `json:"liquidationPrice"`
	   	BankruptPrice        float64   `json:"bankruptPrice"` */

	Timestamp time.Time `json:"timestamp"`
	LastPrice float64   `json:"lastPrice"`
	LastValue float64   `json:"lastValue"`
	// OrderTime int64     `json:"ordertime"`
}

type ContractInfo[] struct {
	// Account          float64       `json:"account"`
	Symbol           string    `json:"symbol"`
	Currency         string    `json:"currency"`
	CurrentQty       float64   `json:"currentQty"`
	MarkPrice        float64   `json:"markPrice"`
	Leverage         float64   `json:"leverage"`
}

type Orders []struct {
	OrderID string `json:"orderID"`
	ClOrdID string `json:"clOrdID"`

	/* 	ClOrdLinkID           string    `json:"clOrdLinkID"` */

	Account float64 `json:"account"`
	Symbol  string  `json:"symbol"`
	Side    string  `json:"side"`

	SimpleOrderQty        float64   `json:"simpleOrderQty"`

	OrderQty float64 `json:"orderQty"`
	Price    float64 `json:"price"`

	/* 	DisplayQty            float64   `json:"displayQty"`
	   	StopPx                float64   `json:"stopPx"`
	   	PegOffsetValue        float64   `json:"pegOffsetValue"`
	   	PegPriceType          string    `json:"pegPriceType"` */

	Currency string `json:"currency"`

	/* 	SettlCurrency         string    `json:"settlCurrency"` */

	OrdType string `json:"ordType"`
	Commission    float64 `json:"commission"` 

	/* 	TimeInForce           string    `json:"timeInForce"`
	   	ExecInst              string    `json:"execInst"`
	   	ContingencyType       string    `json:"contingencyType"`
	   	ExDestination         string    `json:"exDestination"`
	*/
	OrdStatus string `json:"ordStatus"`
	
/* 		Triggered             string    `json:"triggered"`
		WorkingIndicator      bool      `json:"workingIndicator"`
		OrdRejReason          string    `json:"ordRejReason"`
		SimpleLeavesQty       float64   `json:"simpleLeavesQty"`*/ 

		LeavesQty             float64   `json:"leavesQty"`

/* 		SimpleCumQty          float64   `json:"simpleCumQty"` */

		CumQty                float64   `json:"cumQty"`

/* 		AvgPx                 float64   `json:"avgPx"`
		MultiLegReportingType string    `json:"multiLegReportingType"`
		Text                  string    `json:"text"`  */

	TransactTime time.Time `json:"transactTime"`
	Timestamp    time.Time `json:"timestamp"`
	OrderTime    int64     `json:"ordertime"`
}

type NewOrder struct {
	OrderID string `json:"orderID"`
	/* 	ClOrdID               string    `json:"clOrdID"`
	   	ClOrdLinkID           string    `json:"clOrdLinkID"` */
	Account float64 `json:"account"`
	Symbol  string  `json:"symbol"`
	Side    string  `json:"side"`
	/* 	SimpleOrderQty        float64   `json:"simpleOrderQty"` */
	OrderQty float64 `json:"orderQty"`
	Price    float64 `json:"price"`
	/* 	DisplayQty            float64   `json:"displayQty"`
	   	StopPx                float64   `json:"stopPx"`
	   	PegOffsetValue        float64   `json:"pegOffsetValue"`
	   	PegPriceType          string    `json:"pegPriceType"` */
	Currency      string `json:"currency"`
	SettlCurrency string `json:"settlCurrency"`
	OrdType       string `json:"ordType"`
	/* 	TimeInForce           string    `json:"timeInForce"`
	   	ExecInst              string    `json:"execInst"`
	   	ContingencyType       string    `json:"contingencyType"`
	   	ExDestination         string    `json:"exDestination"` */
	OrdStatus string `json:"ordStatus"`
	/* 	Triggered             string    `json:"triggered"`
	   	WorkingIndicator      bool      `json:"workingIndicator"`
	   	OrdRejReason          string    `json:"ordRejReason"`
	   	SimpleLeavesQty       float64   `json:"simpleLeavesQty"`
	   	LeavesQty             float64   `json:"leavesQty"`
	   	SimpleCumQty          float64   `json:"simpleCumQty"`
	   	CumQty                float64   `json:"cumQty"`
	   	AvgPx                 float64   `json:"avgPx"`
	   	MultiLegReportingType string    `json:"multiLegReportingType"`
	   	Text                  string    `json:"text"` */
	TransactTime time.Time `json:"transactTime"`
	Timestamp    time.Time `json:"timestamp"`
	OrderTime    int64     `json:"ordertime"`
}

type Stats []struct {
	RootSymbol   string `json:"rootSymbol"`
	Currency     string `json:"currency"`
	Volume24H    int    `json:"volume24h"`
	Turnover24H  int    `json:"turnover24h"`
	OpenInterest int    `json:"openInterest"`
	OpenValue    int    `json:"openValue"`
}

type DelOrder []struct {
	OrderID string `json:"orderID"`
	/* 	ClOrdID               string    `json:"clOrdID"`
	   	ClOrdLinkID           string    `json:"clOrdLinkID"` */
	Account float64 `json:"account"`
	Symbol  string  `json:"symbol"`
	Side    string  `json:"side"`
	/* 	SimpleOrderQty        float64   `json:"simpleOrderQty"` */
	OrderQty float64 `json:"orderQty"`
	Price    float64 `json:"price"`
	/* 	DisplayQty            float64   `json:"displayQty"`
	   	StopPx                float64   `json:"stopPx"`
	   	PegOffsetValue        float64   `json:"pegOffsetValue"`
	   	PegPriceType          string    `json:"pegPriceType"` */
	Currency string `json:"currency"`
	/* 	SettlCurrency         string    `json:"settlCurrency"` */
	OrdType string `json:"ordType"`
	/* 	TimeInForce           string    `json:"timeInForce"`
	   	ExecInst              string    `json:"execInst"`
	   	ContingencyType       string    `json:"contingencyType"`
	   	ExDestination         string    `json:"exDestination"` */
	OrdStatus string `json:"ordStatus"`
	/* 	Triggered             string    `json:"triggered"`
	   	WorkingIndicator      bool      `json:"workingIndicator"`
	   	OrdRejReason          string    `json:"ordRejReason"`
	   	SimpleLeavesQty       float64   `json:"simpleLeavesQty"`
	   	LeavesQty             float64   `json:"leavesQty"`
	   	SimpleCumQty          float64   `json:"simpleCumQty"`
	   	CumQty                float64   `json:"cumQty"`
	   	AvgPx                 float64   `json:"avgPx"`
	   	MultiLegReportingType string    `json:"multiLegReportingType"` */
	Text         string    `json:"text"`
	TransactTime time.Time `json:"transactTime"`
	Timestamp    time.Time `json:"timestamp"`
}

type ExchangeInfo []struct {
	Symbol          string    `json:"symbol"`
	RootSymbol      string    `json:"rootSymbol"`
	State           string    `json:"state"`
	Expiry          time.Time `json:"expiry"`
	ExpiryTime      int64 `json:"expiry_time"`
	MaxOrderQty     float64       `json:"maxOrderQty"`
	MaxPrice        float64       `json:"maxPrice"`
	InitMargin      float64   `json:"initMargin"`
	RiskLimit       float64     `json:"riskLimit"`
	RiskStep        float64     `json:"riskStep"`
	QuoteCurrency   string     `json:"quoteCurrency"`
	Leveles         []int     `json:"leveles"`
}

type UserMargin struct {
	Account   int    `json:"account"`
	Currency  string `json:"currency"`
	RiskLimit int    `json:"riskLimit"`
	/* 	PrevState          string    `json:"prevState"`
	   	State              string    `json:"state"`
	   	Action             string    `json:"action"` */
	Amount float64 `json:"amount"`
	/* 	PendingCredit      float64   `json:"pendingCredit"`
	   	PendingDebit       float64   `json:"pendingDebit"`
	   	ConfirmedDebit     float64   `json:"confirmedDebit"`
	   	PrevRealisedPnl    float64   `json:"prevRealisedPnl"`
	   	PrevUnrealisedPnl  float64   `json:"prevUnrealisedPnl"`
	   	GrossComm          float64   `json:"grossComm"`
	   	GrossOpenCost      float64   `json:"grossOpenCost"`
	   	GrossOpenPremium   float64   `json:"grossOpenPremium"`
	   	GrossExecCost      float64   `json:"grossExecCost"`
	   	GrossMarkValue     float64   `json:"grossMarkValue"`
	   	RiskValue          float64   `json:"riskValue"`
	   	TaxableMargin      float64   `json:"taxableMargin"`
	   	InitMargin         float64   `json:"initMargin"` */
	MaintMargin float64 `json:"maintMargin"` // 仓位保证金
	/* 	SessionMargin      float64   `json:"sessionMargin"`
	   	TargetExcessMargin float64   `json:"targetExcessMargin"`
	   	VarMargin          float64   `json:"varMargin"`
	   	RealisedPnl        float64   `json:"realisedPnl"`
	   	UnrealisedPnl      float64   `json:"unrealisedPnl"`
	   	IndicativeTax      float64   `json:"indicativeTax"`
	   	UnrealisedProfit   float64   `json:"unrealisedProfit"`
	   	SyntheticMargin    float64   `json:"syntheticMargin"` */
	WalletBalance float64 `json:"walletBalance"` // 钱包余额
	MarginBalance float64 `json:"marginBalance"` // 保证金余额
	/* 	MarginBalancePcnt  float64   `json:"marginBalancePcnt"` */
	MarginLeverage float64 `json:"marginLeverage"` // 已使用杠杆
	/* 	MarginUsedPcnt     float64   `json:"marginUsedPcnt"`
	   	ExcessMargin       float64   `json:"excessMargin"`
	   	ExcessMarginPcnt   float64   `json:"excessMarginPcnt"` */
	AvailableMargin float64 `json:"availableMargin"` // 可用保证金
	/* 	WithdrawableMargin float64   `json:"withdrawableMargin"`
	   	Timestamp          time.Time `json:"timestamp"`
	   	GrossLastValue     float64   `json:"grossLastValue"`
		   Commission         float64   `json:"commission"` */
	WalletBalanceUsd float64 `json:"walletBalanceUsd"`
	XbtPrice         float64 `json:"xbtprice"`
	ExchangeInfoList ExchangeInfo `json:"exchangelist"`
}

type Execution []struct {
	/* 	ExecID                string    `json:"execID"` */
	OrderID string `json:"orderID"`
	/* 	ClOrdID               string    `json:"clOrdID"`
	   	ClOrdLinkID           string    `json:"clOrdLinkID"` */
	Account float64 `json:"account"`
	Symbol  string  `json:"symbol"`
	Side    string  `json:"side"`
	/* 	LastQty               float64   `json:"lastQty"`
	   	LastPx                float64   `json:"lastPx"`
	   	UnderlyingLastPx      float64   `json:"underlyingLastPx"`
	   	LastMkt               string    `json:"lastMkt"`
	   	LastLiquidityInd      string    `json:"lastLiquidityInd"`
	   	SimpleOrderQty        float64   `json:"simpleOrderQty"` */
	OrderQty float64 `json:"orderQty"`
	Price    float64 `json:"price"`
	/* 	DisplayQty            float64   `json:"displayQty"`
	   	StopPx                float64   `json:"stopPx"`
	   	PegOffsetValue        float64   `json:"pegOffsetValue"`
	   	PegPriceType          string    `json:"pegPriceType"` */
	Currency string `json:"currency"`
	/* 	SettlCurrency         string    `json:"settlCurrency"` */
	ExecType string `json:"execType"`
	OrdType  string `json:"ordType"`
	/* 	TimeInForce           string    `json:"timeInForce"`
	   	ExecInst              string    `json:"execInst"`
	   	ContingencyType       string    `json:"contingencyType"`
	   	ExDestination         string    `json:"exDestination"` */
	OrdStatus string `json:"ordStatus"`
	/* 	Triggered             string    `json:"triggered"`
	   	WorkingIndicator      bool      `json:"workingIndicator"`
	   	OrdRejReason          string    `json:"ordRejReason"`
	   	SimpleLeavesQty       float64   `json:"simpleLeavesQty"`
	   	LeavesQty             float64   `json:"leavesQty"`
	   	SimpleCumQty          float64   `json:"simpleCumQty"`
	   	CumQty                float64   `json:"cumQty"`
	   	AvgPx                 float64   `json:"avgPx"`
	   	Commission            float64   `json:"commission"`
	   	TradePublishIndicator string    `json:"tradePublishIndicator"`
	   	MultiLegReportingType string    `json:"multiLegReportingType"`
	   	Text                  string    `json:"text"`
	   	TrdMatchID            string    `json:"trdMatchID"`
	   	ExecCost              float64   `json:"execCost"`
	   	ExecComm              float64   `json:"execComm"`
	   	HomeNotional          float64   `json:"homeNotional"`
	   	ForeignNotional       float64   `json:"foreignNotional"` */
	TransactTime time.Time `json:"transactTime"`
	Timestamp    time.Time `json:"timestamp"`
}

type ModifyOrder struct {
	Account       float64 `json:"account"`
	Symbol        string  `json:"symbol"`
	Currency      string  `json:"currency"`
	Underlying    string  `json:"underlying"`
	QuoteCurrency string  `json:"quoteCurrency"`
	/* 	Commission           float64       `json:"commission"`
	   	InitMarginReq        float64       `json:"initMarginReq"`
	   	MaintMarginReq       float64       `json:"maintMarginReq"`
	   	RiskLimit            float64       `json:"riskLimit"` */
	Leverage float64 `json:"leverage"`
	/* 	CrossMargin          bool      `json:"crossMargin"`
	   	DeleveragePercentile float64       `json:"deleveragePercentile"`
	   	RebalancedPnl        float64       `json:"rebalancedPnl"`
	   	PrevRealisedPnl      float64       `json:"prevRealisedPnl"`
	   	PrevUnrealisedPnl    float64       `json:"prevUnrealisedPnl"` */
	PrevClosePrice   float64   `json:"prevClosePrice"`
	OpeningTimestamp time.Time `json:"openingTimestamp"`
	/* 	OpeningQty           float64       `json:"openingQty"`
	   	OpeningCost          float64       `json:"openingCost"`
	   	OpeningComm          float64       `json:"openingComm"`
	   	OpenOrderBuyQty      float64       `json:"openOrderBuyQty"`
	   	OpenOrderBuyCost     float64       `json:"openOrderBuyCost"`
	   	OpenOrderBuyPremium  float64       `json:"openOrderBuyPremium"`
	   	OpenOrderSellQty     float64       `json:"openOrderSellQty"`
	   	OpenOrderSellCost    float64       `json:"openOrderSellCost"`
	   	OpenOrderSellPremium float64       `json:"openOrderSellPremium"`
	   	ExecBuyQty           float64       `json:"execBuyQty"`
	   	ExecBuyCost          float64       `json:"execBuyCost"`
	   	ExecSellQty          float64       `json:"execSellQty"`
	   	ExecSellCost         float64       `json:"execSellCost"`
	   	ExecQty              float64       `json:"execQty"`
	   	ExecCost             float64       `json:"execCost"`
	   	ExecComm             float64       `json:"execComm"`
	   	CurrentTimestamp     time.Time `json:"currentTimestamp"`
	   	CurrentQty           float64       `json:"currentQty"`
	   	CurrentCost          float64       `json:"currentCost"`
	   	CurrentComm          float64       `json:"currentComm"`
	   	RealisedCost         float64       `json:"realisedCost"`
	   	UnrealisedCost       float64       `json:"unrealisedCost"`
	   	GrossOpenCost        float64       `json:"grossOpenCost"`
	   	GrossOpenPremium     float64       `json:"grossOpenPremium"`
	   	GrossExecCost        float64       `json:"grossExecCost"`
	   	IsOpen               bool      `json:"isOpen"`
	   	MarkPrice            float64       `json:"markPrice"`
	   	MarkValue            float64       `json:"markValue"`
	   	RiskValue            float64       `json:"riskValue"`
	   	HomeNotional         float64       `json:"homeNotional"`
	   	ForeignNotional      float64       `json:"foreignNotional"`
	   	PosState             string    `json:"posState"`
	   	PosCost              float64       `json:"posCost"`
	   	PosCost2             float64       `json:"posCost2"`
	   	PosCross             float64       `json:"posCross"`
	   	PosInit              float64       `json:"posInit"`
	   	PosComm              float64       `json:"posComm"`
	   	PosLoss              float64       `json:"posLoss"`
	   	PosMargin            float64       `json:"posMargin"`
	   	PosMaint             float64       `json:"posMaint"`
	   	PosAllowance         float64       `json:"posAllowance"`
	   	TaxableMargin        float64       `json:"taxableMargin"`
	   	InitMargin           float64       `json:"initMargin"` */
	MaintMargin float64 `json:"maintMargin"`
	/* 	SessionMargin        float64       `json:"sessionMargin"`
	   	TargetExcessMargin   float64       `json:"targetExcessMargin"`
	   	VarMargin            float64       `json:"varMargin"`
	   	RealisedGrossPnl     float64       `json:"realisedGrossPnl"`
	   	RealisedTax          float64       `json:"realisedTax"` */
	RealisedPnl        float64 `json:"realisedPnl"`
	UnrealisedGrossPnl float64 `json:"unrealisedGrossPnl"`
	/* 	LongBankrupt         float64       `json:"longBankrupt"`
	   	ShortBankrupt        float64       `json:"shortBankrupt"`
	   	TaxBase              float64       `json:"taxBase"`
	   	IndicativeTaxRate    float64       `json:"indicativeTaxRate"`
	   	IndicativeTax        float64       `json:"indicativeTax"`
	   	UnrealisedTax        float64       `json:"unrealisedTax"`
	   	UnrealisedPnl        float64       `json:"unrealisedPnl"` */
	   	UnrealisedPnlPcnt    float64       `json:"unrealisedPnlPcnt"`
		UnrealisedRoePcnt    float64       `json:"unrealisedRoePcnt"`
		   /*
	   	SimpleQty            float64       `json:"simpleQty"`
	   	SimpleCost           float64       `json:"simpleCost"`
	   	SimpleValue          float64       `json:"simpleValue"`
	   	SimplePnl            float64       `json:"simplePnl"`
	   	SimplePnlPcnt        float64       `json:"simplePnlPcnt"` */
	AvgCostPrice float64 `json:"avgCostPrice"`
	/* 	AvgEntryPrice        float64       `json:"avgEntryPrice"`
	   	BreakEvenPrice       float64       `json:"breakEvenPrice"` */
	MarginCallPrice float64 `json:"marginCallPrice"`
	/* 	LiquidationPrice     float64       `json:"liquidationPrice"`
	   	BankruptPrice        float64       `json:"bankruptPrice"` */
	Timestamp time.Time `json:"timestamp"`
	/* 	LastPrice            float64       `json:"lastPrice"`
	   	LastValue            float64       `json:"lastValue"` */
}

type InstrumentObj []struct {
	Symbol                         string    `json:"symbol"`
	RootSymbol                     string    `json:"rootSymbol"`
	State                          string    `json:"state"`
	Typ                            string    `json:"typ"`
	Listing                        time.Time `json:"listing"`
	Front                          time.Time `json:"front"`
	Expiry                         time.Time `json:"expiry"`
	Settle                         time.Time `json:"settle"`

/* 	Relistfloat64erval             time.Time `json:"relistfloat64erval"`
	InverseLeg                     string    `json:"inverseLeg"`
	SellLeg                        string    `json:"sellLeg"`
	BuyLeg                         string    `json:"buyLeg"`
	OptionStrikePcnt               float64   `json:"optionStrikePcnt"`
	OptionStrikeRound              float64   `json:"optionStrikeRound"`
	OptionStrikePrice              float64   `json:"optionStrikePrice"`
	OptionMultiplier               float64   `json:"optionMultiplier"`
	PositionCurrency               string    `json:"positionCurrency"`
	Underlying                     string    `json:"underlying"`
	QuoteCurrency                  string    `json:"quoteCurrency"`
	UnderlyingSymbol               string    `json:"underlyingSymbol"`
	Reference                      string    `json:"reference"`
	ReferenceSymbol                string    `json:"referenceSymbol"`
	Calcfloat64erval               time.Time `json:"calcfloat64erval"`
	Publishfloat64erval            time.Time `json:"publishfloat64erval"`
	PublishTime                    time.Time `json:"publishTime"` */

	MaxOrderQty                    float64   `json:"maxOrderQty"`
	MaxPrice                       float64   `json:"maxPrice"`
	LotSize                        float64   `json:"lotSize"`

/* 	TickSize                       float64   `json:"tickSize"`
	Multiplier                     float64   `json:"multiplier"`
	SettlCurrency                  string    `json:"settlCurrency"`
	UnderlyingToPositionMultiplier float64   `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   float64   `json:"underlyingToSettleMultiplier"`
	QuoteToSettleMultiplier        float64   `json:"quoteToSettleMultiplier"`
	IsQuanto                       bool      `json:"isQuanto"`
	IsInverse                      bool      `json:"isInverse"` */

	InitMargin                     float64   `json:"initMargin"`
	Mafloat64Margin                float64   `json:"mafloat64Margin"`
	RiskLimit                      float64   `json:"riskLimit"`
	RiskStep                       float64   `json:"riskStep"`

	/* Limit                          float64   `json:"limit"`
	Capped                         bool      `json:"capped"`
	Taxed                          bool      `json:"taxed"`
	Deleverage                     bool      `json:"deleverage"`
	MakerFee                       float64   `json:"makerFee"`
	TakerFee                       float64   `json:"takerFee"`
	SettlementFee                  float64   `json:"settlementFee"`
	InsuranceFee                   float64   `json:"insuranceFee"`
	FundingBaseSymbol              string    `json:"fundingBaseSymbol"`
	FundingQuoteSymbol             string    `json:"fundingQuoteSymbol"`
	FundingPremiumSymbol           string    `json:"fundingPremiumSymbol"`
	FundingTimestamp               time.Time `json:"fundingTimestamp"`
	Fundingfloat64erval            time.Time `json:"fundingfloat64erval"`
	FundingRate                    float64   `json:"fundingRate"`
	IndicativeFundingRate          float64   `json:"indicativeFundingRate"`
	RebalanceTimestamp             time.Time `json:"rebalanceTimestamp"`
	Rebalancefloat64erval          time.Time `json:"rebalancefloat64erval"`
	OpeningTimestamp               time.Time `json:"openingTimestamp"`
	ClosingTimestamp               time.Time `json:"closingTimestamp"`
	Sessionfloat64erval            time.Time `json:"sessionfloat64erval"`
	PrevClosePrice                 float64   `json:"prevClosePrice"`
	LimitDownPrice                 float64   `json:"limitDownPrice"`
	LimitUpPrice                   float64   `json:"limitUpPrice"`
	BankruptLimitDownPrice         float64   `json:"bankruptLimitDownPrice"`
	BankruptLimitUpPrice           float64   `json:"bankruptLimitUpPrice"`
	PrevTotalVolume                float64   `json:"prevTotalVolume"`
	TotalVolume                    float64   `json:"totalVolume"`
	Volume                         float64   `json:"volume"`
	Volume24H                      float64   `json:"volume24h"`
	PrevTotalTurnover              float64   `json:"prevTotalTurnover"`
	TotalTurnover                  float64   `json:"totalTurnover"`
	Turnover                       float64   `json:"turnover"`
	Turnover24H                    float64   `json:"turnover24h"`
	HomeNotional24H                float64   `json:"homeNotional24h"`
	ForeignNotional24H             float64   `json:"foreignNotional24h"`
	PrevPrice24H                   float64   `json:"prevPrice24h"`
	Vwap                           float64   `json:"vwap"`
	HighPrice                      float64   `json:"highPrice"`
	LowPrice                       float64   `json:"lowPrice"`
	LastPrice                      float64   `json:"lastPrice"`
	LastPriceProtected             float64   `json:"lastPriceProtected"`
	LastTickDirection              string    `json:"lastTickDirection"`
	LastChangePcnt                 float64   `json:"lastChangePcnt"`
	BidPrice                       float64   `json:"bidPrice"`
	MidPrice                       float64   `json:"midPrice"`
	AskPrice                       float64   `json:"askPrice"`
	ImpactBidPrice                 float64   `json:"impactBidPrice"`
	ImpactMidPrice                 float64   `json:"impactMidPrice"`
	ImpactAskPrice                 float64   `json:"impactAskPrice"`
	HasLiquidity                   bool      `json:"hasLiquidity"`
	Openfloat64erest               float64   `json:"openfloat64erest"`
	OpenValue                      float64   `json:"openValue"`
	FairMethod                     string    `json:"fairMethod"`
	FairBasisRate                  float64   `json:"fairBasisRate"`
	FairBasis                      float64   `json:"fairBasis"`
	FairPrice                      float64   `json:"fairPrice"`
	MarkMethod                     string    `json:"markMethod"`
	MarkPrice                      float64   `json:"markPrice"`
	IndicativeTaxRate              float64   `json:"indicativeTaxRate"`
	IndicativeSettlePrice          float64   `json:"indicativeSettlePrice"`
	OptionUnderlyingPrice          float64   `json:"optionUnderlyingPrice"`
	SettledPrice                   float64   `json:"settledPrice"`
	Timestamp                      time.Time `json:"timestamp"` */

}