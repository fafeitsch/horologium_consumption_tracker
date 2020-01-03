package domain

import "time"

type Series struct {
	Id   uint
	Name string
}

type MeterReading struct {
	Id     uint
	Count  float64
	Date   time.Time
	Series *Series
}

type PricingPlan struct {
	Id        uint
	Name      string
	BasePrice float64
	UnitPrice float64
	ValidFrom *time.Time
	ValidTo   *time.Time
	Series    *Series
}
