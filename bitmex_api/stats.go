package bitmax_api

func (client *Client) GetStats() (Stats, error) {
	var res Stats
	_, err := client.Request(GET, GETSTATS, nil, nil, &res)
	return res, err
}
