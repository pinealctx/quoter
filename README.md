# quoter

Get real-time Cryptocurrency quotes via CoinMarketCap.

#### Get it

```shell
go get -u github.com/pinealctx/opensea-go
```

#### Use it

```go
package main

import (
	"github.com/pinealctx/quoter"
	"log"
	"time"
)

func main() {
	var p, err = quoter.New()
	if err != nil {
		log.Fatalf("new.price.latest.error:%+v", err)
		return
	}
	for i := 0; i < 2; i++ {
		log.Printf("1ETH = %fUSD, 1ETH = %fCNY", p.Price(quoter.USD, quoter.ETH), p.Price(quoter.CNY, quoter.ETH))
		time.Sleep(time.Second * 5)
	}
}
```