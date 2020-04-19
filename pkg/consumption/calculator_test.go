package consumption

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func testData() Parameters {
	plan1 := domain.PricingPlan{
		ValidFrom: util.FormatDatePtr(2019, 1, 1),
		ValidTo:   util.FormatDatePtr(2019, 7, 31),
		BasePrice: 10.8,
		UnitPrice: 2.3,
	}
	plan2 := domain.PricingPlan{
		ValidFrom: util.FormatDatePtr(2019, 8, 1),
		ValidTo:   util.FormatDatePtr(2019, 9, 30),
		BasePrice: 11.2,
		UnitPrice: 2.7,
	}
	plan3 := domain.PricingPlan{
		ValidFrom: util.FormatDatePtr(2019, 10, 1),
		ValidTo:   util.FormatDatePtr(2019, 12, 31),
		BasePrice: 11.9,
		UnitPrice: 3.4,
	}
	zeroReading := domain.MeterReading{
		Count: 85,
		Date:  util.FormatDate(2019, 1, 1),
	}
	firstReading := domain.MeterReading{
		Count: 125,
		Date:  util.FormatDate(2019, 4, 12),
	}
	secondReading := domain.MeterReading{
		Count: 335,
		Date:  util.FormatDate(2019, 6, 13),
	}
	thirdReading := domain.MeterReading{
		Count: 400,
		Date:  util.FormatDate(2019, 7, 1),
	}
	forthReading := domain.MeterReading{
		Date:  util.FormatDate(2019, 10, 10),
		Count: 652,
	}
	fifthReading := domain.MeterReading{
		Count: 932,
		Date:  util.FormatDate(2019, 12, 31),
	}
	return Parameters{
		Readings: []domain.MeterReading{zeroReading, firstReading, secondReading, thirdReading, forthReading, fifthReading},
		Plans:    []domain.PricingPlan{plan1, plan2, plan3},
	}
}

func TestCalculate_Simple(t *testing.T) {
	params := testData()
	tests := []struct {
		name            string
		start           time.Time
		end             *time.Time
		wantCosts       float64
		wantConsumption float64
	}{
		{
			name:            "two months (simple)",
			start:           util.FormatDate(2019, 4, 15),
			end:             util.FormatDatePtr(2019, 5, 31),
			wantCosts:       379.9548387096775,
			wantConsumption: 155.80645161290326,
		}, {
			name:            "year",
			start:           util.FormatDate(2019, 1, 1),
			end:             util.FormatDatePtr(2019, 12, 31),
			wantCosts:       2475.380198019802,
			wantConsumption: 847,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params.Start = tt.start
			params.End = *tt.end
			costs, accumulatedConsumption := Costs(params)
			consumption := params.Readings.Consumption(tt.start, *tt.end)
			assert.Equal(t, tt.wantCosts, costs, "calculated costs not correct")
			assert.Equal(t, tt.wantConsumption, consumption, "calculated consumption not correct")
			assert.Equal(t, tt.wantConsumption, accumulatedConsumption, "accumulated calculated consumption not correct")
		})
	}
}

func TestMonthsBetween(t *testing.T) {
	tests := []struct {
		start string
		end   string
		want  int
	}{
		{start: "2020-05-02", end: "2020-05-03", want: 1},
		{start: "2020-01-01", end: "2020-07-01", want: 6},
		{start: "2020-01-01", end: "2020-07-02", want: 7},
		{start: "2019-09-15", end: "2021-10-01", want: 25},
		{start: "2019-11-23", end: "2020-02-22", want: 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("start: %s, end: %s", tt.start, tt.end), func(t *testing.T) {
			start, _ := time.Parse(util.DateFormat, tt.start)
			end, _ := time.Parse(util.DateFormat, tt.end)
			got := monthsBetween(start, end)
			assert.Equal(t, tt.want, got, "months between two dates differ")
		})
	}
}

func TestMonthlyCosts(t *testing.T) {
	params := testData()
	params.Start, _ = time.Parse(util.DateFormat, "2019-01-01")
	params.End, _ = time.Parse(util.DateFormat, "2019-03-24")
	got := MonthlyCosts(params)
	assert.Equal(t, 3, len(got), "there should be twelve months in the statistic")
	wantJanuary := Statistics{
		ValidFrom:   util.FormatDate(2019, 1, 1),
		ValidTo:     util.FormatDate(2019, 2, 1),
		Costs:       39.03762376237624,
		Consumption: 12.277227722772281,
	}
	assertStats(t, wantJanuary, got[0], "January")
	wantFebruary := Statistics{
		ValidFrom:   util.FormatDate(2019, 2, 1),
		ValidTo:     util.FormatDate(2019, 3, 1),
		Costs:       36.304950495049496,
		Consumption: 11.089108910891085,
	}
	assertStats(t, wantFebruary, got[1], "February")
	wantMarch := Statistics{
		ValidFrom:   util.FormatDate(2019, 3, 1),
		ValidTo:     util.FormatDate(2019, 3, 24),
		Costs:       31.750495049504952,
		Consumption: 9.10891089108911,
	}
	assertStats(t, wantMarch, got[2], "March")
}

func assertStats(t *testing.T, want Statistics, got Statistics, msg string) {
	assert.Equal(t, want.ValidFrom, got.ValidFrom, "valid_from of %s is not correct", msg)
	assert.Equal(t, want.ValidTo, got.ValidTo, "valid_to of %s is not correct", msg)
	assert.Equal(t, want.Consumption, got.Consumption, "consumption of %s is not correct", msg)
	assert.Equal(t, want.Costs, got.Costs, "costs of %s is not correct", msg)
}
