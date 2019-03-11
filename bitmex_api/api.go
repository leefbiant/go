package bitmax_api

import (
// "fmt"
)

var (
	Leveles1   = []int{0, 1}
	Leveles20  = []int{0, 1, 2, 3, 5, 10, 20}
	Leveles30  = []int{0, 1, 2, 3, 5, 10, 20, 33}
	Leveles50  = []int{0, 1, 2, 3, 5, 10, 20, 25, 50}
	Leveles100 = []int{0, 1, 2, 3, 5, 10, 20, 25, 50, 100}
)

func GetLevLes(initMargin float64) []int {
	value := 1 / initMargin
	if value >= 100 {
		return Leveles100
	}
	if value >= 50 {
		return Leveles50
	}
	if value >= 30 {
		return Leveles30
	}
	if value >= 20 {
		return Leveles30
	}
	return Leveles1
}

func GetUserMargin(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string) (UserMargin, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	opt := make(map[string]interface{})
	res, err := client.GetUserMargin(opt)
	return res, err
}

func GetPosition(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (Positions, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = true
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.GetPosition(opt)
	return res, err
}

func GetContractInfo(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (ContractInfo, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.GetContractInfo(opt)
	return res, err
}

func GetOrders(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (Orders, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.GetOrders(opt)
	return res, err
}

func CancelOredr(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (DelOrder, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.DelOrders(opt)
	return res, err
}

func PutOrders(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (NewOrder, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.PutOrders(opt)
	return res, err
}

func TransModifyOrderMargin(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (ModifyOrder, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.ModifyPostionMargin(opt)
	return res, err
}

func ModifyLeverAge(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (ModifyOrder, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.ModifyLeverAge(opt)
	return res, err
}

func GetInstrument(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (ExchangeInfo, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.GetInstrument(opt)
	return res, err
}

func ModifyEntrustOrder(ApiKey string, SecretKey string, Endpoint string, ProxyAddr string,
	opt map[string]interface{}) (NewOrder, error) {
	var config Config
	config.Endpoint = Endpoint
	config.ApiKey = ApiKey
	config.SecretKey = SecretKey

	config.IsPrint = false
	client := NewClient(config, ProxyAddr)
	defer client.Close()

	res, err := client.ModifyEntrustOrder(opt)
	return res, err
}
