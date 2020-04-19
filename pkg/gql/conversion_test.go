package gql

import (
	"github.com/fafeitsch/Horologium/pkg/consumption"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
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
