package quoter

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pinealctx/neptune/jsonx"
	"github.com/pinealctx/neptune/ulog"
	"github.com/pinealctx/restgo"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Quoter struct {
	cryptoIDs []int
	prices    *sync.Map
	socket    *Socket
	request   string
	running   *atomic.Bool
	cli       *restgo.Client
}

func New(symbols ...Symbol) (*Quoter, error) {
	var l = len(symbols)
	if l == 0 {
		symbols = []Symbol{ETH, BTC, USD, CNY}
	}
	var cryptoIDs = make([]int, 0, len(symbols))
	for _, s := range symbols {
		id, ok := coinMarketSymbolIDHash[s]
		if ok {
			cryptoIDs = append(cryptoIDs, id)
		}
	}
	var req = NewSubscribeRequest(cryptoIDs)
	var request, err = jsonx.JSONFastMarshal(req)
	if err != nil {
		return nil, fmt.Errorf("json.marshal.error:%+v", err)
	}
	var p = &Quoter{
		cryptoIDs: cryptoIDs,
		prices:    &sync.Map{},
		request:   string(request),
		running:   atomic.NewBool(true),
		cli:       restgo.New(restgo.WithBaseURL(baseURL)),
	}
	err = p.latest()
	if err != nil {
		return nil, fmt.Errorf("fetch.latest.error:%+v", err)
	}
	go p.loopLatest()
	go p.loopListen()
	return p, nil
}

func (q *Quoter) Close() {
	if !q.running.Load() {
		return
	}
	q.running.Store(false)
	if q.socket != nil {
		var err = q.socket.Close()
		if err != nil {
			ulog.Error("socket.close.error", zap.Error(err))
		}
	}
}

func (q *Quoter) Price(target, symbol Symbol) float64 {
	var v = q.price(symbol)
	if target == USD || target == "" {
		return v
	}
	var t = q.price(target)
	if t == 0 {
		return 0
	}
	return v / t
}

func (q *Quoter) price(symbol Symbol) float64 {
	var v, ok = q.prices.Load(symbol)
	if !ok {
		return 0
	}
	return v.(float64)
}

func (q *Quoter) latest() error {
	var usdID = coinMarketSymbolIDHash[USD]
	var idList = make([]string, len(q.cryptoIDs))
	for i, id := range q.cryptoIDs {
		idList[i] = strconv.Itoa(id)
	}
	var rsp, err = q.cli.Get(context.Background(),
		"/data-api/v3/cryptocurrency/quote/latest",
		restgo.NewURLQueryParam("id", strings.Join(idList, ",")),
		restgo.NewURLQueryParam("convertId", strconv.Itoa(usdID)),
	)
	if err != nil {
		return fmt.Errorf("get.latest.error:%+v", err)
	}
	var response LatestResponse
	err = rsp.JSONUnmarshal(&response)
	if err != nil {
		return fmt.Errorf("json.unmarshal.error:%+v", err)
	}
	if response.Status == nil || response.Status.ErrorCode != 0 {
		return fmt.Errorf("latest.fail:%s", response.Status.ErrorMessage)
	}
	for _, data := range response.Data {
		for _, quote := range data.Quotes {
			if quote.Name == usdID {
				q.store(data.ID, quote.Price)
				break
			}
		}
	}
	return nil
}

func (q *Quoter) loopLatest() {
	var ticker = time.NewTicker(time.Hour)
	for q.running.Load() {
		<-ticker.C
		ticker.Reset(time.Hour)
		var err = q.latest()
		if err != nil {
			ulog.Error("fetch.latest.error", zap.Error(err))
		}
	}
}

func (q *Quoter) loopListen() {
	var count int
	for q.running.Load() {
		count++
		ulog.Error("socket.connecting...", zap.Int("count", count))
		var err error
		q.socket, err = NewSocket(wsURL, wsHeader)
		if err != nil {
			ulog.Error("new.socket.error", zap.Error(err))
			continue
		}
		err = q.socket.SendStr(q.request)
		if err != nil {
			ulog.Error("send.message.error", zap.Error(err))
			continue
		}
		go func() {
			var msg Message
			for {
				msg = <-q.socket.Message()
				switch msg.Typ {
				case websocket.TextMessage:
					q.procPrice(msg.Content)
				default:
					ulog.Info("unknown.msg.type", zap.Int("typ", msg.Typ))
				}
			}
		}()
		err = <-q.socket.Wait()
		if err != nil {
			ulog.Error("socket.error", zap.Error(err))
		}
	}
}

func (q *Quoter) procPrice(data []byte) {
	var rsp SubscribeResponse
	var err = jsonx.JSONFastUnmarshal(data, &rsp)
	if err != nil {
		ulog.Error("json.unmarshal.error", zap.String("d", string(data)), zap.Error(err))
		return
	}
	if rsp.ID != "price" {
		return
	}
	var cr = rsp.Data.CR
	if cr.ID == 0 || cr.Price == 0 {
		return
	}
	q.store(cr.ID, cr.Price)
}

func (q *Quoter) store(id int, price float64) {
	var symbol, ok = coinMarketIDSymbolHash[id]
	if !ok {
		return
	}
	q.prices.Store(symbol, price)
}
