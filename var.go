package quoter

import (
	"net/http"
)

var (
	baseURL                = "https://api.coinmarketcap.com"
	wsURL                  = "wss://stream.coinmarketcap.com/price/latest"
	wsHeader               = genHeader()
	coinMarketSymbolIDHash = map[Symbol]int{
		BTC: 1,
		ETH: 1027,
		USD: 2781,
		CNY: 2787,
	}
	coinMarketIDSymbolHash = genIDSymbolHash()
)

func genHeader() http.Header {
	var h = http.Header{}
	h.Set("Accept-Encoding", "gzip, deflate, br")
	h.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	h.Set("Cache-Control", "no-cache")
	h.Set("Host", "stream.coinmarketcap.com")
	h.Set("Origin", "https://coinmarketcap.com")
	h.Set("Pragma", "no-cache")
	h.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	return h
}

func genIDSymbolHash() map[int]Symbol {
	var hash = make(map[int]Symbol)
	for k, v := range coinMarketSymbolIDHash {
		hash[v] = k
	}
	return hash
}
