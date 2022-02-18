# quoter

[![Go Reference](https://pkg.go.dev/badge/github.com/pinealctx/quoter.svg)](https://pkg.go.dev/github.com/pinealctx/quoter)
[![golangci-lint](https://github.com/pinealctx/quoter/actions/workflows/ci.yml/badge.svg)](https://github.com/pinealctx/quoter/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinealctx/quoter)](https://goreportcard.com/report/github.com/pinealctx/quoter)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

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