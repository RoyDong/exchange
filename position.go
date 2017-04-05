package exchange

const (
    //仓位 方向
    PositionLong = 1
    PositionShort = 2
)


type Position struct {
    Type int
    Amount float64
    AvgPrice float64
    Money float64
    Deposit float64
    Leverage int
}


func NewPosition(typ, lever int, amount, price, money float64) *Position {
    p := &Position{
        Type: typ,
        Amount: amount,
        AvgPrice: price,
        Money: money,
        Deposit: 0,
        Leverage: lever,
    }

    if p.Amount > 0 {
        if p.AvgPrice > 0 {
            p.Money = p.Amount * p.AvgPrice
        } else if p.Money > 0 {
            p.AvgPrice = p.Amount / p.Money
        }

        p.Deposit = p.Money / float64(p.Leverage)
    }

    return p
}

