package exchange


import (
    "time"
)


/*
Amount 是指交易物品的单位数量(stock, contract, bitcoin...)
Price  交易的单价(单个物品需要多少的计价货币) Price = Money / Amount
Money  购买所用的货币数量(rmb, usd, btc, ltc ...), 需要事先定好以什么作为计价货币, 价格可以通过汇率转换

 */


const (

    //交易 下单类型
    OpenLongPosition = 1
    OpenShortPosition = 2
    CloseLongPosition = 3
    CloseShortPosition = 4


    //货币种类
    CurrencyRMB = 1
    CurrencyUSD = 2
    CurrencyBTC = 3
)

/*
清单
 */
type Bill struct {
    Amount, Money, Fee float64
    CreateTime time.Time
}

/*
小清单
 */
type SmallBill struct {
    Amount, Prince float64
}

type IExchange interface {

    Name() string

    /*
    pos    开仓类型
    amount 数量
    price 价格  price = 0 市价, price = -1 对手价
     */
    Trade(pos int, amount, price float64) (int64, error)

    CancelOrder(id int64) error

    GetOrder(id int64) (Order, error)

    GetTicker() (Ticker, error)

    GetTrades() ([]Trade, error)

    GetDepth() ([]SmallBill, []SmallBill, error)

    /*
    return currency type, money
     */
    GetBalance() (int, float64, error)

    /*
    获取当前所持有的仓位信息，如果是现货则用一个杠杆为1的多头Position表示: Position.Leverage = 1 Position.Type = LongPosition
     */
    GetPositions() ([]*Position, error)
}

