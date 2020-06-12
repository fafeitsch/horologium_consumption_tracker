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

func (s *Series) mapToDomain() (*domain.Series, error) {
	plans := make([]domain.PricingPlan, 0, len(s.Plans))
	for index, plan := range s.Plans {
		domainPlan, err := plan.mapToDomain()
		if err != nil {
			return nil, fmt.Errorf("could not parse plan %d: %v", index, err)
		}
		plans = append(plans, *domainPlan)
	}
	readings := make([]domain.MeterReading, 0, len(s.Readings))
	for index, reading := range s.Readings {
		domainReading, err := reading.mapToDomain()
		if err != nil {
			return nil, fmt.Errorf("could not parse reading %d: %v", index, err)
		}
		readings = append(readings, *domainReading)
	}
	return &domain.Series{Name: s.Name, PricingPlans: plans, MeterReadings: readings}, nil
}

type PricingPlan struct {
	Name      string
	BasePrice float64 `json:"basePrice"`
	UnitPrice float64 `json:"unitPrice"`
	ValidFrom string  `json:"validFrom"`
	ValidTo   *string `json:"validTo"`
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
