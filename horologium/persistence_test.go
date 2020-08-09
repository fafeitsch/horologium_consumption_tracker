package horologium

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"
	"time"
)

func ExampleLoadFromReader() {
	file :=
		`name: "A pseudo power consumption for testing"
consumptionFormat: "%.2f kWh"
currencyFormat: "%.2f Euro"
plans:
  - {name: 2018, basePrice: 1241.34, unitPrice: 26.32, validFrom: "2018-01-01", validTo: "2018-12-31"}
  - {name: 2019, basePrice: 1341.12, unitPrice: 27.28, validFrom: "2019-01-01", validTo: "2019-12-31"}
  - {name: 2020, basePrice: 1400.28, unitPrice: 26.56, validFrom: "2020-01-01", validTo: "2020-12-31"}
readings:
  - {date: 2018-02-01, count: 1223.34}
  - {date: 2018-03-01, count: 1256.93}
  - {date: 2018-04-01, count: 1299.23}`

	reader := strings.NewReader(file)
	got, err := LoadFromReader(reader)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}
	fmt.Printf("Name: %s\n", got.Name)
	fmt.Printf("Number of plans: %d\n", len(got.PricingPlans))
	fmt.Printf("Name of first plan: %s\n", got.PricingPlans[0].Name)
	fmt.Printf("Number of readings: %d\n", len(got.MeterReadings))
	fmt.Printf("Count of third reading: %.2f\n", got.MeterReadings[2].Count)
	fmt.Printf("Consumption format: %s\n", got.ConsumptionFormat)
	fmt.Printf("Currency format: %s\n", got.CurrencyFormat)
	// Output: Name: A pseudo power consumption for testing
	// Number of plans: 3
	// Name of first plan: 2018
	// Number of readings: 3
	// Count of third reading: 1299.23
	// Consumption format: %.2f kWh
	// Currency format: %.2f Euro
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestLoadFromReader_ReaderError(t *testing.T) {
	got, err := LoadFromReader(&errReader{})
	assert.EqualError(t, err, "could not read reader: test error", "error message wrong")
	assert.Nil(t, got, "result should be nil in case of an error")
}

func TestLoadFromReader_YamlError(t *testing.T) {
	reader := strings.NewReader("I'm not { a valid yaml")
	got, err := LoadFromReader(reader)
	assert.EqualError(t, err, "unmarshalling yaml failed: String node doesn't MapNode", "error message wrong")
	assert.Nil(t, got, "result should be nil in case of an error")
}

//noinspection GoNilness
func TestSeries_MapToDomain(t *testing.T) {
	to := "2020-02-29"
	plan1 := pricingPlanDto{
		Name:      "Year 2000",
		BasePrice: 1202.23,
		UnitPrice: 19.2,
		ValidFrom: "2000-01-01",
		ValidTo:   &to,
	}
	plan2 := pricingPlanDto{
		Name:      "To Infinity",
		BasePrice: 1823.12,
		UnitPrice: 27.23,
		ValidFrom: "2020-03-01",
		ValidTo:   nil,
	}
	reading1 := meterReadingDto{
		Count: 3242,
		Date:  "2008-10-14",
	}
	reading2 := meterReadingDto{
		Count: 6585,
		Date:  "2013-05-23",
	}
	tests := []struct {
		name           string
		wantPlanErr    bool
		wantReadingErr bool
		wantErrMessage string
	}{
		{name: "wrong plan", wantPlanErr: true, wantReadingErr: false, wantErrMessage: "could not parse plan 0: could not parse validFrom date: parsing time \"10.04.2004\" as \"2006-01-02\": cannot parse \"4.2004\" as \"2006\""},
		{name: "wrong reading", wantPlanErr: false, wantReadingErr: true, wantErrMessage: "could not parse reading 0: could not parse date: parsing time \"----\" as \"2006-01-02\": cannot parse \"----\" as \"2006\""},
		{name: "success", wantPlanErr: false, wantReadingErr: false, wantErrMessage: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var plans []pricingPlanDto
			if tt.wantPlanErr {
				plans = []pricingPlanDto{{ValidFrom: "10.04.2004"}}
			} else {
				plans = []pricingPlanDto{plan1, plan2}
			}
			var readings []meterReadingDto
			if tt.wantReadingErr {
				readings = []meterReadingDto{{Date: "----"}}
			} else {
				readings = []meterReadingDto{reading1, reading2}
			}
			series := seriesDto{
				Name:     "Test Series",
				Plans:    plans,
				Readings: readings,
			}
			got, err := series.mapToDomain()
			if tt.wantPlanErr || tt.wantReadingErr {
				assert.Nil(t, got, "result should be nil in case of an error")
				assert.EqualError(t, err, tt.wantErrMessage, "error message wrong")
			} else {
				require.NotNil(t, got, "got must not be nil")
				assert.Equal(t, series.Name, got.Name, "name is wrong")
				require.Equal(t, len(series.Plans), len(got.PricingPlans), "number of pricing plans is wrong")
				assert.Equal(t, series.Plans[1].UnitPrice, got.PricingPlans[1].UnitPrice, "UnitPrice of second plan wrong")
				require.Equal(t, len(series.Readings), len(got.MeterReadings), "number or meter readings is wrong")
				assert.Equal(t, series.Readings[1].Count, got.MeterReadings[1].Count, "count of second meter reading is wrong")
			}
		})
	}
}

func TestPricingPlan_MapToDomain(t *testing.T) {
	to := "2018-08-13"
	notADate := "notADate"
	tests := []struct {
		name          string
		validFrom     string
		validTo       *string
		wantValidFrom time.Time
		wantValidTo   *string
		wantErr       error
	}{
		{name: "success", validFrom: "2018-07-10", validTo: &to, wantValidFrom: CreateDate(2018, 7, 10), wantValidTo: &to, wantErr: nil},
		{name: "success without validTo", validFrom: "2018-07-10", validTo: nil, wantValidFrom: CreateDate(2018, 7, 10), wantValidTo: nil, wantErr: nil},
		{name: "cannot parse validFrom", validFrom: "notADate", validTo: nil, wantValidFrom: CreateDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validFrom date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
		{name: "cannot parse validFrom", validFrom: "2018-07-10", validTo: &notADate, wantValidFrom: CreateDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validTo date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := pricingPlanDto{UnitPrice: 34.23, BasePrice: 1023.12, ValidTo: tt.validTo, ValidFrom: tt.validFrom, Name: "Unit Test Plan"}
			got, err := plan.mapToDomain()
			if err != nil || tt.wantErr != nil {
				assert.Nil(t, got, "got should be nil in case of an error")
				assert.EqualError(t, err, tt.wantErr.Error(), "error is wrong")
			} else {
				assert.Equal(t, plan.Name, got.Name, "name is wrong")
				assert.Equal(t, tt.wantValidFrom, *got.ValidFrom, "validFrom is wrong")
				var wantValidTo *time.Time
				if tt.wantValidTo != nil {
					t, _ := time.Parse(DateFormat, *tt.wantValidTo)
					wantValidTo = &t
				}
				assert.Equal(t, wantValidTo, got.ValidTo, "validTo is wrong")
				assert.Equal(t, 1023.12, got.BasePrice, "basePrice  is wrong")
				assert.Equal(t, 34.23, got.UnitPrice, "unitPrice is wrong")
			}
		})
	}
}

func TestMeterReading_MapToDomain(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		wantTime time.Time
		wantErr  error
	}{
		{name: "success", date: "2018-05-31", wantTime: CreateDate(2018, 5, 31), wantErr: nil},
		{name: "wrong date", date: "not a date", wantErr: errors.New("could not parse date: parsing time \"not a date\" as \"2006-01-02\": cannot parse \"not a date\" as \"2006\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reading := meterReadingDto{Date: tt.date, Count: 34.3}
			got, err := reading.mapToDomain()
			if tt.wantErr != nil || err != nil {
				assert.Nil(t, got, "got should be nil if an error occurred")
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Equal(t, reading.Count, got.Count, "count is wrong")
				assert.Equal(t, tt.wantTime, got.Date, "date is wrong")
			}
		})
	}
}
