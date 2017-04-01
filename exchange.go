package exchange


import (
    "github.com/roydong/gmvc"
    //"math"
    "time"
    "fmt"
    "container/list"
)


/*
Amount 是指交易物品的单位数量(stock, contract, bitcoin...)
Price  交易的单价(单个物品需要多少的计价货币) Price = Money / Amount
Money  购买所用的货币数量(rmb, usd, btc ...), 需要事先定好以什么作为计价货币, 价格可以通过汇率转换

 */


const (

    //仓位 方向
    PositionLong = 1
    PositionShort = 2

    //交易 下单方式
    OpenLongPosition = 1
    OpenShortPosition = 2
    CloseLongPosition = 3
    CloseShortPosition = 4

    //order status
    OrderCreated   = 0
    OrderPartial   = 1  //部分成交
    OrderComplete  = 2
    OrderCanceled  = -1
    OrderCanceling = 4


    //货币种类
    CurrencyRMB = 1
    CurrencyUSD = 2
    CurrencyBTC = 3
)

type IExchange interface {

    /*
    pos    开仓类型
    amount 数量
    price 价格  price = 0 市价, price = -1 对手价
     */
    Trade(pos int, amount, price float64) int64

    Order(id int64) Order

    CancelOrder(id int64) bool

    //GetTicker() Ticker

    //GetTrade() Trade

    GetTrades() []Trade

    GetDepth() ([][]float64, [][]float64)

    GetFund() *Fund
}

type Unit struct {
    Amount, Price float64
}

type Exchange struct {
    IExchange

    //交易所名称
    name string

    //账户金额
    fund *Fund

    lastAsks, lastBids [][]float64

    trades *list.List
}

func NewExchange(name string) *Exchange {
    e := &Exchange{
        name: name,
        currency: "cny",
        trades: list.New(),
    }

    switch e.name {
    case "okcoin":
        //e.IExchange = NewOKCoin()
    case "huobi":
        //e.IExchange = NewHuobi()
    case "okfuture_quarter":
        e.IExchange = NewOKFuture("quarter", 20)
        e.currency = "usd"
    case "okfuture_thisweek":
        e.IExchange = NewOKFuture("this_week", 20)
        e.currency = "usd"

    default:
        gmvc.Logger.Fatalln("invalid exchange " + e.name)
    }

    return e
}

func (e *Exchange) Name() string {
    return e.name
}

/*
获取盘口买单对应深度的价格和数量
 */
func (e *Exchange) GetBidDepth(depth int) (float64, float64) {
    var price, amount float64
    for i, bid := range e.lastBids {
        amount += bid[1]
        if i + 1 >= depth {
            price = bid[0]
            break
        }
    }
    return price, amount
}

/*
获取盘口卖单对应深度的价格和深度
 */
func (e *Exchange) GetAskDepth(depth int) (float64, float64) {
    var price, amount float64
    for i, ask := range e.lastAsks {
        amount += ask[1]
        if i >= depth {
            price = ask[0]
            break
        }
    }
    return price, amount
}

func (e *Exchange) SyncDepth() {
    lastAsks, lastBids := e.GetDepth()
    if len(lastAsks) > 0 && len(lastBids) > 0 {
        e.lastAsks = lastAsks
        e.lastBids = lastBids
    }
}

func (e *Exchange) SyncBalance() {
    b := e.GetBalance()
    e.amount, e.money = b.Amount, b.Money
}

func (e *Exchange) Balance() (float64, float64) {
    return e.amount, e.money
}

/*
根据最近的交易计算出价格的平均数
 */
func (e *Exchange) calcMa() {
    trades := e.IExchange.GetTrades()
    var trade Trade
    if e.trades.Len() > 1 {
        trade, _ = e.trades.Front().Value.(Trade)
    }

    newTrade := false
    for _, t := range trades {
        if t.Id > trade.Id {
            e.trades.PushFront(t)
            if e.trades.Len() > 600 {
                e.trades.Remove(e.trades.Back())
            }
            newTrade = true
        }
    }

    if !newTrade {
        return
    }

    var n, sum float64
    for el := e.trades.Back(); el != nil; el = el.Prev() {
        trade, _ := el.Value.(Trade)
        if trade.Price > 0 {
            sum = sum + trade.Price * trade.Amount
            n = n + trade.Amount
        }
    }
    e.ma = sum / n
}

/*
price = 0 对手价
 */
func (e *Exchange) Trade(position int, amount, price float64) Order {
    var order Order
    var id int64
    for i := 0; i < 2; i++ {
        id = e.IExchange.Trade(position, amount, price)
        if id > 0 {
            break
        }
    }
    if id > 0 {
        //每隔0.5s读取一次，最多等待3次
        for i := 0; i < 5; i++ {
            time.Sleep(200 * time.Millisecond)
            order = e.IExchange.Order(id)
            if order.Status == 2 {
                break
            }
        }

        //如果订单没有完全成交
        if order.Status != 2 {
            canceled := false

            //重试两次，如果都失败中断程序
            for i := 0; i < 2; i++ {
                canceled = e.IExchange.CancelOrder(id)
            }

            if !canceled {
                gmvc.Logger.Println(fmt.Sprintf("cancel order failed %v order id = %v", e.Name(), id))
            }

            //更新order info
            for i := 0; i < 2; i++ {
                order = e.IExchange.Order(id)
                if order.Id > 0 {
                    break
                }
            }

            if order.Id == 0 {
                gmvc.Logger.Println(fmt.Sprintf("update order info failed %v order id = %v", e.Name(), id))
            }
        }
    } else {
        gmvc.Logger.Println(fmt.Sprintf("make order failed %v", e.Name()))
    }
    return order
}

/*
不断以对手价下单直到交易完amount数量
 */
func (e *Exchange) TradeAll(position int, amount float64) Order {
    var order Order
    var money float64
    for order.DealAmount < amount {
        o := e.Trade(position, amount - order.DealAmount, 0)
        if order.Id <= 0 {
            order.Id = o.Id
        }
        money += o.DealAmount * o.AvgPrice
        order.DealAmount += o.DealAmount
    }
    order.AvgPrice = money / order.DealAmount
    return order
}

