package exchange

import "time"

type Fund struct {
    //货币类型
    CurrencyType int

    amount, money float64

    FrozenAmount, FrozenMoney float64

    Deposit float64

    Statments []*FundStmt
}


func (f *Fund) NewFundStmt() {

}


/*
流水帐
 */
type FundStmt struct {
    //金额变动
    Amount, Money float64

    //时间点
    CreateTime time.Time
}