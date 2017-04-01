package exchange






type Order struct {
    Id string
    Type int
    Amount float64
    Price float64
    Time int64

    DealAmount float64
    AvgPrice float64
    Status int
    Fee float64
}

