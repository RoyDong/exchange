package exchange

import (
    "sync"
)

type Fund struct {
    //货币类型
    Currency int

    /*
    期货: 多头仓位, 空头仓位
    现货：多头仓位表示仓位，空头仓位忽略
     */
    LongAmount, ShortAmount float64

    Money, Fee float64

    Deposit float64
}


