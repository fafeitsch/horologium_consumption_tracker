package consumption

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculate_Simple(t *testing.T) {

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
	params := Parameters{
		readings: []domain.MeterReading{zeroReading, firstReading, secondReading, thirdReading, forthReading, fifthReading},
		plans:    []domain.PricingPlan{plan1, plan2, plan3},
	}
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
			wantCosts:       2497.3801980198014,
			wantConsumption: 847,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params.Start = tt.start
			params.End = *tt.end
			costs := Costs(params)
			consumption := Consumption(params)
			assert.Equal(t, tt.wantCosts, costs, "calculated costs not correct")
			assert.Equal(t, tt.wantConsumption, consumption, "calculated consumption not correct")
		})
	}
}
