package exchange

import "time"

const (
    //order status
    OrderStatusCreated   = 0
    OrderStatusPartial   = 1  //部分成交
    OrderStatusComplete  = 2
    OrderStatusCanceled  = -1
    OrderStatusCanceling = 4
)

type Order struct {
    Id int64
    Type int
    Amount float64
    Price float64

    //成交的数量
    DealAmount float64
    //已使用了多少货币
    DealMoney float64
    //成交均价
    AvgPrice float64
    //手续费
    Fee float64
    //订单状态
    Status int

    CreateTime time.Time
}

