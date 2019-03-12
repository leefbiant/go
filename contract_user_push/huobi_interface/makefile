all:
	go build -o huobi_depth depth.go huobi_util.go
	go build -o huobi_trade trade.go huobi_util.go
	go build -o huobi_ticker ticker.go huobi_util.go
	go build -o huobi_index index.go
clean:
	rm -rf huobi_depth huobi_trade huobi_ticker huobi_index
