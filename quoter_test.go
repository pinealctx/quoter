package quoter

import (
	"testing"
	"time"
)

func TestPriceLatest_Price(t *testing.T) {
	var p, err = New()
	if err != nil {
		t.Errorf("new.price.latest.error:%+v", err)
		return
	}
	for i := 0; i < 2; i++ {
		t.Logf("1ETH = %fUSD, 1ETH = %fCNY", p.Price(USD, ETH), p.Price(CNY, ETH))
		time.Sleep(time.Second * 5)
	}
}
