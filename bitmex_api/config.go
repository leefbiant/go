package bitmax_api

type Config struct {
	// Rest api endpoint url. eg: http://www.bitmax.com/
	Endpoint string

	// Rest websocket api endpoint url. eg: ws://192.168.80.113:10442/
	WSEndpoint string

	// The user's api key provided by OKEx.
	ApiKey string
	// The user's secret key provided by OKEx. The secret key used to sign your request data.
	SecretKey string

	IsPrint bool
}
