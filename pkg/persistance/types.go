package orm

import (
	"github.com/jinzhu/gorm"
	"time"
)

type seriesEntity struct {
	gorm.Model
	Name string
}

type meterReadingEntity struct {
	gorm.Model
	Count    float64
	Date     time.Time
	Series   seriesEntity
	SeriesID int
}

type pricingPlanEntity struct {
	gorm.Model
	BasePrice float64
	UnitPrice float64
	ValidFrom *time.Time
	ValidTo   *time.Time
	Series    seriesEntity
	SeriesID  int
}
