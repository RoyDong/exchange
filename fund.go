package exchange

import (
    "time"
    "sync"
)

type Fund struct {
    //货币类型
    currency int

    amount, money, fee float64

    frozenAmount, frozenMoney float64

    deposit float64

    bills []Bill

    locker *sync.RWMutex
}


func NewFund(currency int) *Fund {
    f := &Fund{}
    f.currency = currency
    f.bills = make([]Bill, 0)
    f.locker = &sync.RWMutex{}

    return f
}

func (f *Fund) Change(amount, money, fee float64) {
    f.locker.Lock()
    defer f.locker.Unlock()

    changed := false
    if amount != 0 {
        f.amount += amount
        changed = true
    }
    if money != 0 {
        f.money += money
        changed = true
    }
    if fee != 0 {
        f.fee += fee
        changed = true
    }
    if changed {
        f.bills = append(f.bills, Bill{amount, money, fee,time.Now()})
    }
}

func (f *Fund) ChangeFrozon(amount, money float64) {
    f.locker.Lock()
    defer f.locker.Unlock()

    if amount != 0 {
        f.frozenAmount += amount
    }
    if money != 0 {
        f.frozenMoney += money
    }
}

func (f *Fund) ClearStatements() {
    f.bills = make([]Bill, 0)
}

func (f *Fund) Amount() float64 {
    f.locker.RLock()
    defer f.locker.RUnlock()
    return f.amount
}

func (f *Fund) Money() float64 {
    f.locker.RLock()
    defer f.locker.RUnlock()
    return f.money
}

func (f *Fund) Fee() float64 {
    f.locker.RLock()
    defer f.locker.RUnlock()
    return f.fee
}

func (f *Fund) Currency() int {
    return f.currency
}

/*
清单
 */
type Bill struct {
    Amount, Money, Fee float64
    CreateTime time.Time
}
