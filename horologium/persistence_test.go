package horologium

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestLoadFromReader(t *testing.T) {
	file, _ := ioutil.ReadFile("../test-resources/series/powerSeries.yml")
	buffer := bytes.NewBuffer(file)
	got, err := LoadFromReader(buffer)
	require.NoError(t, err, "no error expected")
	assert.Equal(t, "A pseudo power consumption for testing", got.Name, "name not correct")
	require.Equal(t, 3, len(got.PricingPlans))
	assert.Equal(t, "2019", got.PricingPlans[1].Name)
	require.Equal(t, 3, len(got.MeterReadings))
	assert.Equal(t, 1299.23, got.MeterReadings[2].Count)
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
		wantValidTo   *time.Time
		wantErr       error
	}{
		{name: "success", validFrom: "2018-07-10", validTo: &to, wantValidFrom: FormatDate(2018, 7, 10), wantValidTo: FormatDatePtr(2018, 8, 13), wantErr: nil},
		{name: "success without validTo", validFrom: "2018-07-10", validTo: nil, wantValidFrom: FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: nil},
		{name: "cannot parse validFrom", validFrom: "notADate", validTo: nil, wantValidFrom: FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validFrom date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
		{name: "cannot parse validFrom", validFrom: "2018-07-10", validTo: &notADate, wantValidFrom: FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validTo date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
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
				assert.Equal(t, tt.wantValidTo, got.ValidTo, "validTo is wrong")
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
		{name: "success", date: "2018-05-31", wantTime: FormatDate(2018, 5, 31), wantErr: nil},
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