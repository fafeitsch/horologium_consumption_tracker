package gql

import (
	"errors"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/consumption"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_toQlStatisticsNil(t *testing.T) {
	assert.Nil(t, toQlStatistics(nil), "nil as input should get nil as output")
}

func Test_toQlStatistics(t *testing.T) {
	stats := &consumption.Statistics{
		ValidFrom:   util.FormatDate(2019, 9, 23),
		ValidTo:     util.FormatDate(2019, 11, 15),
		Costs:       444,
		Consumption: 666,
	}
	got := toQlStatistics(stats)
	assert.Equal(t, stats.ValidFrom.Format(util.DateFormat), got.ValidFrom, "valid_from is wrong")
	assert.Equal(t, stats.ValidTo.Format(util.DateFormat), got.ValidTo, "valid_to is wrong")
	assert.Equal(t, stats.Consumption, got.Consumption, "consumption is wrong")
	assert.Equal(t, stats.Costs, got.Costs, "costs are wrong")
}

func Test_parseDates(t *testing.T) {
	date1 := "2019-04-04"
	date2 := "2019-08-31"
	date3 := "20-20-20"
	tests := []struct {
		dates   []*string
		want    []*time.Time
		wantErr error
	}{
		{dates: []*string{&date1, nil, &date2}, want: []*time.Time{util.FormatDatePtr(2019, 4, 4), nil, util.FormatDatePtr(2019, 8, 31)}, wantErr: nil},
		{dates: []*string{&date1, &date3}, want: nil, wantErr: errors.New("could not parse \"20-20-20\" as format YYYY-MM-DD")},
	}
	for index, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", index), func(t *testing.T) {
			got, err := parseDates(tt.dates...)
			assert.Equal(t, tt.want, got, "converted times are not correct")
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error(), "errors not equal")
			} else {
				assert.NoError(t, err, "no error expected")
			}
		})
	}
}
