package domain

import "time"

type Series struct {
	Id   uint
	Name string
}

type MeterReading struct {
	Id    *int
	Count float64
	Date  time.Time
}

type PricingPlan struct {
	Id        *int
	BasePrice float64
	UnitPrice float64
	ValidFrom *time.Time
	ValidTo   *time.Time
}
