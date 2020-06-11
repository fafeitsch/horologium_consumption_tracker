package storage

import (
	"errors"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
		{name: "success", validFrom: "2018-07-10", validTo: &to, wantValidFrom: util.FormatDate(2018, 7, 10), wantValidTo: util.FormatDatePtr(2018, 8, 13), wantErr: nil},
		{name: "success without validTo", validFrom: "2018-07-10", validTo: nil, wantValidFrom: util.FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: nil},
		{name: "cannot parse validFrom", validFrom: "notADate", validTo: nil, wantValidFrom: util.FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validFrom date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
		{name: "cannot parse validFrom", validFrom: "2018-07-10", validTo: &notADate, wantValidFrom: util.FormatDate(2018, 7, 10), wantValidTo: nil, wantErr: errors.New("could not parse validTo date: parsing time \"notADate\" as \"2006-01-02\": cannot parse \"notADate\" as \"2006\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := PricingPlan{UnitPrice: 34.23, BasePrice: 1023.12, ValidTo: tt.validTo, ValidFrom: tt.validFrom, Name: "Unit Test Plan"}
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
		{name: "success", date: "2018-05-31", wantTime: util.FormatDate(2018, 5, 31), wantErr: nil},
		{name: "wrong date", date: "not a date", wantErr: errors.New("could not parse date: parsing time \"not a date\" as \"2006-01-02\": cannot parse \"not a date\" as \"2006\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reading := MeterReading{Date: tt.date, Count: 34.3}
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
