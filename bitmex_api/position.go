package bitmax_api

func (client *Client) GetPosition(opt map[string]interface{}) (Positions, error) {
	var res Positions
	_, err := client.Request(GET, POSITION, opt, nil, &res)
	return res, err
}

func (client *Client) GetContractInfo(opt map[string]interface{}) (ContractInfo, error) {
	var res ContractInfo
	_, err := client.Request(GET, POSITION, opt, nil, &res)
	return res, err
}

func (client *Client) ModifyPostionMargin(opt map[string]interface{}) (ModifyOrder, error) {
	var res ModifyOrder
	_, err := client.Request(POST, MODIFYMARGIN, opt, nil, &res)
	return res, err
}

func (client *Client) ModifyLeverAge(opt map[string]interface{}) (ModifyOrder, error) {
	var res ModifyOrder
	_, err := client.Request(POST, MODIFYLEVERAGE, opt, nil, &res)
	return res, err
}

func (client *Client) GetInstrument(opt map[string]interface{}) (ExchangeInfo, error) {
	var res ExchangeInfo
	_, err := client.Request(GET, GET_INSTRUMENT, opt, nil, &res)
	return res, err
}
