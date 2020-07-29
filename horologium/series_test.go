package horologium

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func testData() *Series {
	plan0 := PricingPlan{
		ValidFrom: formatDatePtr(2018, 1, 1),
		ValidTo:   formatDatePtr(2019, 1, 1),
		BasePrice: 100,
		UnitPrice: 100,
	}
	plan1 := PricingPlan{
		ValidFrom: formatDatePtr(2019, 1, 1),
		ValidTo:   formatDatePtr(2019, 8, 1),
		BasePrice: 10.8,
		UnitPrice: 2.3,
	}
	plan2 := PricingPlan{
		ValidFrom: formatDatePtr(2019, 8, 1),
		ValidTo:   formatDatePtr(2019, 10, 1),
		BasePrice: 11.2,
		UnitPrice: 2.7,
	}
	plan3 := PricingPlan{
		ValidFrom: formatDatePtr(2019, 10, 1),
		ValidTo:   formatDatePtr(2019, 12, 31),
		BasePrice: 11.9,
		UnitPrice: 3.4,
	}
	zeroReading := MeterReading{
		Count: 85,
		Date:  CreateDate(2019, 1, 1),
	}
	firstReading := MeterReading{
		Count: 125,
		Date:  CreateDate(2019, 4, 12),
	}
	secondReading := MeterReading{
		Count: 335,
		Date:  CreateDate(2019, 6, 13),
	}
	thirdReading := MeterReading{
		Count: 400,
		Date:  CreateDate(2019, 7, 1),
	}
	forthReading := MeterReading{
		Date:  CreateDate(2019, 10, 10),
		Count: 652,
	}
	fifthReading := MeterReading{
		Count: 932,
		Date:  CreateDate(2019, 12, 31),
	}
	return &Series{
		PricingPlans:  []PricingPlan{plan0, plan1, plan2, plan3},
		MeterReadings: []MeterReading{zeroReading, firstReading, secondReading, thirdReading, forthReading, fifthReading},
	}
}

func formatDatePtr(year int, month int, day int) *time.Time {
	result := CreateDate(year, month, day)
	return &result
}

func TestCalculate_Simple(t *testing.T) {
	series := testData()
	tests := []struct {
		name            string
		start           time.Time
		end             time.Time
		wantCosts       float64
		wantConsumption float64
	}{
		{
			name:            "two months (simple)",
			start:           CreateDate(2019, 4, 15),
			end:             CreateDate(2019, 5, 31),
			wantCosts:       379.9548387096775,
			wantConsumption: 155.80645161290326,
		}, {
			name:            "year",
			start:           CreateDate(2019, 1, 1),
			end:             CreateDate(2019, 12, 31),
			wantCosts:       2475.380198019802,
			wantConsumption: 847,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			costs, accumulatedConsumption := series.CostsAndConsumption(tt.start, tt.end)
			consumption := series.MeterReadings.Consumption(tt.start, tt.end)
			assert.Equal(t, tt.wantCosts, costs, "calculated costs not correct")
			assert.Equal(t, tt.wantConsumption, consumption, "calculated consumption not correct")
			assert.Equal(t, tt.wantConsumption, accumulatedConsumption, "accumulated calculated consumption not correct")
		})
	}
}

func ExampleSeries_CostsAndConsumption() {
	pricingPlanEnd := CreateDate(2019, 6, 1)
	p1 := PricingPlan{ValidFrom: nil, ValidTo: &pricingPlanEnd, BasePrice: 10, UnitPrice: .20}
	p2 := PricingPlan{ValidFrom: &pricingPlanEnd, ValidTo: nil, BasePrice: 10, UnitPrice: .30}
	plans := PricingPlans{p1, p2}
	// Simple calculation, we consume constantly 100 units per day:
	m1 := MeterReading{Date: CreateDate(2019, 5, 1), Count: 1000}
	m2 := MeterReading{Date: CreateDate(2019, 6, 1), Count: 4100}
	m3 := MeterReading{Date: CreateDate(2019, 7, 1), Count: 7100}
	readings := MeterReadings{m1, m2, m3}
	series := Series{PricingPlans: plans, MeterReadings: readings}
	cost, consumption := series.CostsAndConsumption(CreateDate(2019, 5, 15), CreateDate(2019, 6, 15))
	// Base price of first pricing plan is not included because the the month is already half over.
	// This gives: 31 days between 2019-05-15 and 2019-06-15 → 31 * 100 = 3100 consumption
	// First 16 days in pricing plan 1: 16 * 100 * 0.2 = 320
	// Last 25 days in pricing plan 2: 15 * 100 * 0.3 = 450
	// 320 + 450 + base price of pricing plan 2 = 780
	fmt.Printf("Cost: %.2f, Consumption: %.2f", cost, consumption)
	// Output: Cost: 780.00, Consumption: 3100.00
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
			start, _ := time.Parse(DateFormat, tt.start)
			end, _ := time.Parse(DateFormat, tt.end)
			got := monthsBetween(start, end)
			assert.Equal(t, tt.want, got, "months between two dates differ")
		})
	}
}

func TestMonthlyCosts(t *testing.T) {
	series := testData()
	got := series.MonthlyStatistics(CreateDate(2019, 1, 1), CreateDate(2019, 3, 24))
	assert.Equal(t, 3, len(got), "there should be twelve months in the statistic")
	wantJanuary := Statistics{
		ValidFrom:   CreateDate(2019, 1, 1),
		ValidTo:     CreateDate(2019, 2, 1),
		Costs:       39.03762376237624,
		Consumption: 12.277227722772281,
	}
	assertStats(t, wantJanuary, got[0], "January")
	wantFebruary := Statistics{
		ValidFrom:   CreateDate(2019, 2, 1),
		ValidTo:     CreateDate(2019, 3, 1),
		Costs:       36.304950495049496,
		Consumption: 11.089108910891085,
	}
	assertStats(t, wantFebruary, got[1], "February")
	wantMarch := Statistics{
		ValidFrom:   CreateDate(2019, 3, 1),
		ValidTo:     CreateDate(2019, 3, 24),
		Costs:       31.750495049504952,
		Consumption: 9.10891089108911,
	}
	assertStats(t, wantMarch, got[2], "March")
}

func ExampleSeries_MonthlyStatistics() {
	pricingPlanEnd := CreateDate(2019, 6, 1)
	p1 := PricingPlan{ValidFrom: nil, ValidTo: &pricingPlanEnd, BasePrice: 10, UnitPrice: .20}
	p2 := PricingPlan{ValidFrom: &pricingPlanEnd, ValidTo: nil, BasePrice: 10, UnitPrice: .30}
	plans := PricingPlans{p1, p2}
	// Simple calculation, we consume constantly 100 units per day:
	m1 := MeterReading{Date: CreateDate(2019, 5, 1), Count: 1000}
	m2 := MeterReading{Date: CreateDate(2019, 6, 1), Count: 4100}
	m3 := MeterReading{Date: CreateDate(2019, 7, 1), Count: 7100}
	readings := MeterReadings{m1, m2, m3}
	series := Series{PricingPlans: plans, MeterReadings: readings}
	statistics := series.MonthlyStatistics(CreateDate(2019, 5, 1), CreateDate(2019, 7, 1))
	for _, stat := range statistics {
		fmt.Printf("%s – %s, Consumption: %.2f, Cost: %.2f\n", stat.ValidFrom.Format(DateFormat), stat.ValidTo.Format(DateFormat), stat.Consumption, stat.Costs)
	}
	// Output: 2019-05-01 – 2019-06-01, Consumption: 3100.00, Cost: 630.00
	// 2019-06-01 – 2019-07-01, Consumption: 3000.00, Cost: 910.00
}

func assertStats(t *testing.T, want Statistics, got Statistics, msg string) {
	assert.Equal(t, want.ValidFrom, got.ValidFrom, "valid_from of %s is not correct", msg)
	assert.Equal(t, want.ValidTo, got.ValidTo, "valid_to of %s is not correct", msg)
	assert.Equal(t, want.Consumption, got.Consumption, "consumption of %s is not correct", msg)
	assert.Equal(t, want.Costs, got.Costs, "costs of %s is not correct", msg)
}
