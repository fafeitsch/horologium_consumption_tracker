package horologium

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInterpolateValueAtDate(t *testing.T) {
	readings := MeterReadings{
		{Date: CreateDate(2019, 9, 12), Count: 45},
		{Date: CreateDate(2019, 9, 23), Count: 80},
		{Date: CreateDate(2019, 9, 27), Count: 134},
		{Date: CreateDate(2019, 9, 30), Count: 178},
	}
	tests := []struct {
		date string
		want float64
	}{
		{date: "2019-09-27", want: 134},
		{date: "2019-09-30", want: 178},
		{date: "2020-01-02", want: 178},
		{date: "2019-09-26", want: 120.5},
	}
	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			date, _ := time.Parse(DateFormat, tt.date)
			got := readings.interpolateValueAtDate(date)
			assert.Equal(t, tt.want, got, "interpolated value is wrong")
		})
	}
}

func TestLastReadingBefore(t *testing.T) {
	readings := MeterReadings{
		{Date: CreateDate(2019, 9, 12)},
		{Date: CreateDate(2019, 9, 23)},
		{Date: CreateDate(2019, 9, 27)},
		{Date: CreateDate(2019, 9, 30)},
	}
	t.Run("simple", func(t *testing.T) {
		got := readings.lastReadingBefore(CreateDate(2019, 9, 26))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := readings.lastReadingBefore(CreateDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all smaller", func(t *testing.T) {
		got := readings.lastReadingBefore(CreateDate(2019, 10, 1))
		assert.Equal(t, readings[3], got, "last reading calculated incorrectly")
	})
}

func TestFirstReadingAfter(t *testing.T) {
	readings := MeterReadings{
		{Date: CreateDate(2019, 9, 12)},
		{Date: CreateDate(2019, 9, 23)},
		{Date: CreateDate(2019, 9, 27)},
		{Date: CreateDate(2019, 9, 30)},
	}
	cpy := make(MeterReadings, len(readings))
	copy(cpy, readings)
	t.Run("simple", func(t *testing.T) {
		got := cpy.firstReadingAfter(CreateDate(2019, 9, 17))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := cpy.firstReadingAfter(CreateDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all bigger", func(t *testing.T) {
		got := cpy.firstReadingAfter(CreateDate(2019, 9, 1))
		assert.Equal(t, readings[0], got, "last reading calculated incorrectly")
	})
}

func TestMeterReadings_Consumption_Order(t *testing.T) {
	m1 := MeterReading{Date: CreateDate(2019, 4, 5), Count: 1500}
	m2 := MeterReading{Date: CreateDate(2019, 4, 10), Count: 2000}
	m3 := MeterReading{Date: CreateDate(2019, 4, 15), Count: 2250}
	readings := MeterReadings{m1, m3, m2}
	// The consumption will be calcutlated wrongly because of wrong order
	// This test is to check that the order is not changed
	_ = readings.Consumption(CreateDate(2019, 4, 9), CreateDate(2019, 4, 11))
	assert.Equal(t, MeterReadings{m1, m3, m2}, readings, "order has changed, that is not allowed")
}

func ExampleMeterReadings_Consumption() {
	m1 := MeterReading{Date: CreateDate(2019, 4, 5), Count: 1500}
	m2 := MeterReading{Date: CreateDate(2019, 4, 10), Count: 2000}
	m3 := MeterReading{Date: CreateDate(2019, 4, 15), Count: 2250}
	readings := MeterReadings{m1, m2, m3}
	consumption := readings.Consumption(CreateDate(2019, 4, 9), CreateDate(2019, 4, 11))
	fmt.Printf("%.2f", consumption)
	// Output: 150.00
}

func ExampleMeterReadings_Sort() {
	m1 := MeterReading{Date: CreateDate(2019, 4, 5)}
	m2 := MeterReading{Date: CreateDate(2019, 4, 10)}
	m3 := MeterReading{Date: CreateDate(2019, 4, 15)}
	readings := MeterReadings{m2, m3, m1}
	fmt.Printf("%s,%s,%s\n", readings[0].Date.Format(DateFormat), readings[1].Date.Format(DateFormat), readings[2].Date.Format(DateFormat))
	readings.Sort()
	fmt.Printf("%s,%s,%s", readings[0].Date.Format(DateFormat), readings[1].Date.Format(DateFormat), readings[2].Date.Format(DateFormat))
	// Output: 2019-04-10,2019-04-15,2019-04-05
	// 2019-04-05,2019-04-10,2019-04-15
}

func ExampleCreateDate() {
	date := CreateDate(2019, 10, 12)
	formatted := date.Format(time.RFC1123)
	fmt.Printf("%s", formatted)
	// Output: Sat, 12 Oct 2019 00:00:00 UTC
}
