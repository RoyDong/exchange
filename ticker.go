package exchange

import "time"

type Ticker struct {
    High float64
    Low  float64
    Ask float64
    Bid  float64
    Last float64
    Vol  float64
    CreateTime time.Time
}
