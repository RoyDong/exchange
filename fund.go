package exchange

import (
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

func (f *Fund) AddBill(bill Bill) {
    f.locker.Lock()
    defer f.locker.Unlock()

    changed := false
    if bill.Amount != 0 {
        f.amount += bill.Amount
        changed = true
    }
    if bill.Money != 0 {
        f.money += bill.Money
        changed = true
    }
    if bill.Fee != 0 {
        f.fee += bill.Fee
        changed = true
    }
    if changed {
        f.bills = append(f.bills, bill)
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

func (f *Fund) ClearBills() {
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
