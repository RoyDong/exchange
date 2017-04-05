package okcoin

import (
    "strings"
    "fmt"
    "github.com/roydong/gtools"
    "github.com/roydong/exchange/utils"
    "github.com/roydong/exchange"
    "time"
)

type OKFuture struct {
    httpHost   string
    apiKey    string
    apiSecret string

    contractType string
    leverRate int
}


func NewOKFuture(httpHost, apiKey, apiSecret, contractType string, leverRate int) *OKFuture {
    ok := &OKFuture{}
    ok.httpHost  = httpHost
    ok.apiKey = apiKey
    ok.apiSecret = apiSecret
    ok.contractType = contractType
    ok.leverRate = leverRate
    return ok
}

func (ok *OKFuture) Name() string {
    return "okfuture/" + ok.contractType
}

func (ok *OKFuture) Trade(position int, amount, price float64) int64 {
    params := map[string]interface{}{
        "symbol": "btc_usd",
        "contract_type": ok.contractType,
        "type": position,
        "amount": fmt.Sprintf("%.0f", amount),
        "price": price,
        "match_price": 0,
        "lever_rate": ok.leverRate,
    }
    if price == 0 {
        params["match_price"] = 1
    }
    rs := ok.callHttp("future_trade.do", nil, params)
    if rs == nil {
        return 0
    }
    id, _ := rs.Int64("order_id")
    return id
}

func (ok *OKFuture) GetOrder(id int64) exchange.Order {
    params := map[string]interface{}{
        "symbol": "btc_usd",
        "contract_type": ok.contractType,
        "order_id": id,
    }

    order := exchange.Order{}
    rs := ok.callHttp("future_order_info.do", nil, params)
    if rs == nil {
        return order
    }

    rst := rs.Tree("orders.0")
    if rst == nil {
        return order
    }

    order.Id = id
    order.Amount, _ = rst.Float("amount")
    order.Price, _ = rst.Float("price")
    order.DealAmount, _ = rst.Float("deal_amount")
    order.AvgPrice, _ = rst.Float("price_avg")
    order.Status, _ = rst.Int("status")
    sec, _ := rst.Int64("create_date")
    order.CreateTime = time.Unix(sec, 0)

    return order
}

func (ok *OKFuture) CancelOrder(id int64) bool {
    params := map[string]interface{} {
        "symbol": "btc_usd",
        "contract_type": ok.contractType,
        "order_id": id,
    }

    rs := ok.callHttp("future_cancel.do", nil, params)
    if rs == nil {
        return false
    }

    if result, _ := rs.String("result"); result == "true" {
        return true
    }

    return false
}

func (ok *OKFuture) GetTrades() []exchange.Trade {
    params := map[string]interface{}{
        "symbol": "btc_usd",
        "contract_type": ok.contractType,
    }

    rs := ok.callHttp("future_trades.do", params, nil)
    if rs == nil {
        return nil
    }

    n := rs.NodeNum("")
    trades := make([]exchange.Trade, 0, n)
    for i := 0; i < n; i++ {
        trade := exchange.Trade{}
        trade.Id, _ = rs.Int64(fmt.Sprintf("%d.tid", i))
        trade.Price, _ = rs.Float(fmt.Sprintf("%d.price", i))
        trade.Amount, _ = rs.Float(fmt.Sprintf("%d.amount", i))
        sec, _ := rs.Int64(fmt.Sprintf("%d.date_ms", i))
        trade.CreateTime = time.Unix(sec, 0)
        trade.Type, _ = rs.String(fmt.Sprintf("%d.type", i))
        trades = append(trades, trade)
    }
    return trades
}

func (ok *OKFuture) GetTicker() exchange.Ticker {
    t := exchange.Ticker{}
    q := map[string]interface{}{"symbol": "btc_usd", "contract_type": ok.contractType}
    rs := ok.callHttp("future_ticker.do", q, nil)
    if rs == nil {
        return t
    }

    rst := rs.Tree("ticker")
    t.High, _ = rst.Float("high")
    t.Low,  _ = rst.Float("low")
    t.Ask, _ = rst.Float("sell")
    t.Bid,  _ = rst.Float("buy")
    t.Last, _ = rst.Float("last")
    t.Vol,  _ = rst.Float("vol")
    sec, _ := rs.Int64("date")
    t.CreateTime = time.Unix(sec, 0)

    return t
}

func (ok *OKFuture) GetDepth() ([]exchange.SmallBill, []exchange.SmallBill) {
    query := map[string]interface{}{
        "symbol": "btc_usd",
        "size": 50,
        "merge": 0,
        "contract_type": ok.contractType,
    }

    rs := ok.callHttp("future_depth.do", query, nil)
    if rs == nil {
        return nil, nil
    }

    var l int
    ask := make([]exchange.SmallBill, 0, l)
    l = rs.NodeNum("asks")
    for i := l - 1; i >= 0; i-- {
        price, _ := rs.Float(fmt.Sprintf("asks.%v.0", i))
        amount, _ := rs.Float(fmt.Sprintf("asks.%v.1", i))
        ask = append(ask, exchange.SmallBill{amount, price})
    }

    bid := make([]exchange.SmallBill, 0, l)
    l = rs.NodeNum("bids")
    for i := 0; i < l; i++ {
        price, _ := rs.Float(fmt.Sprintf("bids.%v.0", i))
        amount, _ := rs.Float(fmt.Sprintf("bids.%v.1", i))
        bid = append(bid, exchange.SmallBill{amount, price})
    }

    return ask, bid
}

func (ok *OKFuture) GetIndex() float64 {
    q := map[string]interface{}{"symbol": "btc_usd"}
    rs := ok.callHttp("future_index.do", q, nil)
    idx, _ := rs.Float("future_index")
    return idx
}


func (ok *OKFuture) GetBalance() (int, float64) {
    rs := ok.callHttp("future_userinfo.do", nil, map[string]interface{}{})
    if rs == nil {
        return 0, 0
    }

    btcInfo := rs.Tree("info.btc")
    if btcInfo == nil {
        return 0, 0
    }

    amount, _ := btcInfo.Float("account_rights")
    deposit, _ := btcInfo.Float("keep_deposit")

    return exchange.CurrencyBTC, amount - deposit
}


func (ok *OKFuture) callHttp(api string, query, params map[string]interface{}) *gtools.Tree {
    if params != nil {
        params["api_key"] = ok.apiKey
        params["sign"] = strings.ToUpper(utils.CreateSignature(params, ok.apiSecret))
    }

    tree := utils.CallRest(ok.httpHost + api, query, params)
    if tree == nil {
        return nil
    }

    if _, has := tree.Int64("error_code"); has {
        return nil
    }

    return tree
}

