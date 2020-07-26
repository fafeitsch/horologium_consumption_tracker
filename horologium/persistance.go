package horologium

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-yaml"
	"io"
	"time"
)

// Reads the yaml file provided by the reader and returns a series struct.
// In case of parsing errors, an error is returned.
func LoadFromReader(reader io.Reader) (*Series, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read reader: %v", err)
	}
	series := seriesDto{}
	err = yaml.Unmarshal(buf.Bytes(), &series)
	if err != nil {
		formatError := yaml.FormatError(err, true, true)
		return nil, fmt.Errorf("unmarshalling yaml failed: " + formatError)
	}
	return series.mapToDomain()
}

type seriesDto struct {
	Name     string
	Plans    []pricingPlanDto
	Readings []meterReadingDto
}

func (s *seriesDto) mapToDomain() (*Series, error) {
	plans := make([]PricingPlan, 0, len(s.Plans))
	for index, plan := range s.Plans {
		domainPlan, err := plan.mapToDomain()
		if err != nil {
			return nil, fmt.Errorf("could not parse plan %d: %v", index, err)
		}
		plans = append(plans, *domainPlan)
	}
	readings := make([]MeterReading, 0, len(s.Readings))
	for index, reading := range s.Readings {
		domainReading, err := reading.mapToDomain()
		if err != nil {
			return nil, fmt.Errorf("could not parse reading %d: %v", index, err)
		}
		readings = append(readings, *domainReading)
	}
	return &Series{Name: s.Name, PricingPlans: plans, MeterReadings: readings}, nil
}

type pricingPlanDto struct {
	Name      string
	BasePrice float64 `json:"basePrice"`
	UnitPrice float64 `json:"unitPrice"`
	ValidFrom string  `json:"validFrom"`
	ValidTo   *string `json:"validTo"`
}

func (p *pricingPlanDto) mapToDomain() (*PricingPlan, error) {
	validFrom, err := time.Parse(DateFormat, p.ValidFrom)
	if err != nil {
		return nil, fmt.Errorf("could not parse validFrom date: %v", err)
	}
	var validTo *time.Time
	if p.ValidTo != nil {
		validToVal, err := time.Parse(DateFormat, *p.ValidTo)
		if err != nil && len(*p.ValidTo) > 1 {
			return nil, fmt.Errorf("could not parse validTo date: %v", err)
		}
		validTo = &validToVal
	}
	return &PricingPlan{ValidFrom: &validFrom, ValidTo: validTo, Name: p.Name, BasePrice: p.BasePrice, UnitPrice: p.UnitPrice}, nil
}

type meterReadingDto struct {
	Count float64
	Date  string
}

func (m *meterReadingDto) mapToDomain() (*MeterReading, error) {
	date, err := time.Parse(DateFormat, m.Date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}
	return &MeterReading{Date: date, Count: m.Count}, nil
}
