package storage

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"time"
)

type Series struct {
	Name     string
	Plans    []PricingPlan
	Readings []MeterReading
}

type PricingPlan struct {
	Name      string
	BasePrice float64
	UnitPrice float64
	ValidFrom string
	ValidTo   *string
}

func (p *PricingPlan) mapToDomain() (*domain.PricingPlan, error) {
	validFrom, err := time.Parse(util.DateFormat, p.ValidFrom)
	if err != nil {
		return nil, fmt.Errorf("could not parse validFrom date: %v", err)
	}
	var validTo *time.Time
	if p.ValidTo != nil {
		validToVal, err := time.Parse(util.DateFormat, *p.ValidTo)
		if err != nil && len(*p.ValidTo) > 1 {
			return nil, fmt.Errorf("could not parse validTo date: %v", err)
		}
		validTo = &validToVal
	}
	return &domain.PricingPlan{ValidFrom: &validFrom, ValidTo: validTo, Name: p.Name, BasePrice: p.BasePrice, UnitPrice: p.UnitPrice}, nil
}

type MeterReading struct {
	Count float64
	Date  string
}

func (m *MeterReading) mapToDomain() (*domain.MeterReading, error) {
	date, err := time.Parse(util.DateFormat, m.Date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}
	return &domain.MeterReading{Date: date, Count: m.Count}, nil
}
