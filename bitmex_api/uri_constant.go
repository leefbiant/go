package bitmax_api

const (
	// for user
	USER   = "/api/v1/user"
	WALLET = "/api/v1/user/wallet"

	// for position
	POSITION       = "/api/v1/position"
	MODIFYMARGIN   = "/api/v1/position/transferMargin"
	MODIFYLEVERAGE = "/api/v1/position/leverage"

	// for order
	ORDER     = "/api/v1/order"
	ORDEREXEC = "/api/v1/execution"

	// stst
	GETSTATS       = "/api/v1/stats"
	GETUSER_MARGIN = "/api/v1/user/margin"

	GET_INSTRUMENT = "/api/v1/instrument/active"
)
