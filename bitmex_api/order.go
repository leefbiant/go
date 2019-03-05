package bitmax_api

func (client *Client) GetOrders(opt map[string]interface{}) (Orders, error) {
	var res Orders
	_, err := client.Request(GET, ORDER, opt, nil, &res)
	return res, err
}

func (client *Client) PutOrders(opt map[string]interface{}) (NewOrder, error) {
	var res NewOrder
	_, err := client.Request(POST, ORDER, opt, nil, &res)
	return res, err
}

func (client *Client) ModifyEntrustOrder(opt map[string]interface{}) (NewOrder, error) {
	var res NewOrder
	_, err := client.Request(PUT, ORDER, opt, nil, &res)
	return res, err
}

func (client *Client) DelOrders(opt map[string]interface{}) (DelOrder, error) {
	var res DelOrder
	_, err := client.Request(DELETE, ORDER, opt, nil, &res)
	return res, err
}

func (client *Client) OrderExecution(opt map[string]interface{}) (Execution, error) {
	var res Execution
	_, err := client.Request(GET, ORDEREXEC, opt, nil, &res)
	return res, err
}
