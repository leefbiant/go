package bitmax_api

func (client *Client) GetUserMargin(opt map[string]interface{}) (UserMargin, error) {
	var res UserMargin
	_, err := client.Request(GET, GETUSER_MARGIN, opt, nil, &res)
	return res, err
}
