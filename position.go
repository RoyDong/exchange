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
    LeverRate int
}


func NewPosition(typ, rate int) *Position {
    return &Position{
        Type: typ,
        LeverRate: rate,
    }
}

func (p *Position) Change(amount, price float64) {
    p.Amount += amount
    p.Money += amount * price
    if p.Amount > 0 && p.Money > 0 {
        p.AvgPrice = p.Amount / p.Money
    } else {
        p.AvgPrice = 0
    }
}

