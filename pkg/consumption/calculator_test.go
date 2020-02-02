package consumption

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLastReadingBefore(t *testing.T) {
	readings := []domain.MeterReading{
		{Date: util.FormatDate(2019, 9, 12), Id: 5},
		{Date: util.FormatDate(2019, 9, 23), Id: 6},
		{Date: util.FormatDate(2019, 9, 27), Id: 7},
		{Date: util.FormatDate(2019, 9, 30), Id: 8},
	}
	t.Run("simple", func(t *testing.T) {
		got := lastReadingBefore(util.FormatDate(2019, 9, 26), readings)
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := lastReadingBefore(util.FormatDate(2019, 9, 23), readings)
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all smaller", func(t *testing.T) {
		got := lastReadingBefore(util.FormatDate(2019, 10, 1), readings)
		assert.Equal(t, readings[3], got, "last reading calculated incorrectly")
	})
}

func TestFirstReadingAfter(t *testing.T) {
	readings := []domain.MeterReading{
		{Date: util.FormatDate(2019, 9, 12), Id: 5},
		{Date: util.FormatDate(2019, 9, 23), Id: 6},
		{Date: util.FormatDate(2019, 9, 27), Id: 7},
		{Date: util.FormatDate(2019, 9, 30), Id: 8},
	}
	cpy := make([]domain.MeterReading, len(readings))
	copy(cpy, readings)
	t.Run("simple", func(t *testing.T) {
		got := firstReadingAfter(util.FormatDate(2019, 9, 17), cpy)
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := firstReadingAfter(util.FormatDate(2019, 9, 23), cpy)
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all bigger", func(t *testing.T) {
		got := firstReadingAfter(util.FormatDate(2019, 9, 1), cpy)
		assert.Equal(t, readings[0], got, "last reading calculated incorrectly")
	})
}

func TestCalculate_Simple(t *testing.T) {
	plan := domain.PricingPlan{
		ValidFrom: util.FormatDatePtr(2019, 1, 1),
		ValidTo:   util.FormatDatePtr(2019, 12, 31),
		BasePrice: 10.8,
		UnitPrice: 2.3,
	}
	firstReading := domain.MeterReading{
		Count: 125,
		Date:  util.FormatDate(2019, 4, 12),
	}
	secondReading := domain.MeterReading{
		Count: 335,
		Date:  util.FormatDate(2019, 6, 13),
	}
	params := Parameters{
		Start:    util.FormatDate(2019, 4, 15),
		End:      util.FormatDate(2019, 5, 31),
		readings: []domain.MeterReading{firstReading, secondReading},
		plans:    []domain.PricingPlan{plan},
	}
	costs := Costs(params)
	consumption := Consumption(params)
	assert.Equal(t, 380.429518, costs, "calcuated costs not correct")
	assert.Equal(t, 155.80645161290326, consumption, "calculated consumption not correct")
}
