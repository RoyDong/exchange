package exchange

import "time"

const (
    TradeTypeBuy = 1
    TradeTypeSell = 2
)


type Trade struct {
    Id string
    Amount float64
    Price float64
    Type string
    CreateTime time.Time
}

