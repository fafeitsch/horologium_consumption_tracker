package orm

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMeterReadingServiceImpl_CRUD(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	powerSeries := &domain.Series{Name: "Power"}
	waterSeries := &domain.Series{Name: "Water"}

	time1, _ := time.Parse(util.DateFormat, "2018-01-31")
	time2, _ := time.Parse(util.DateFormat, "2018-02-28")
	time3, _ := time.Parse(util.DateFormat, "2018-03-31")
	meter1 := domain.MeterReading{
		Count:  1000,
		Date:   time1,
		Series: powerSeries,
	}
	meter2 := domain.MeterReading{
		Count:  1000,
		Date:   time1,
		Series: waterSeries,
	}
	meter3 := domain.MeterReading{
		Count:  2000,
		Date:   time2,
		Series: powerSeries,
	}
	meter4 := domain.MeterReading{
		Count:  3000,
		Date:   time3,
		Series: powerSeries,
	}

	service := NewMeterReadingService(db)
	err := service.Save(&meter1)
	require.NoError(t, err)
	err = service.Save(&meter2)
	require.NoError(t, err)
	err = service.Save(&meter3)
	require.NoError(t, err)
	err = service.Save(&meter4)
	require.NoError(t, err)
	want := []domain.MeterReading{meter1, meter3, meter4}

	got, err := service.QueryForSeries(powerSeries.Id)
	require.NoError(t, err, "no error while quering expected")
	assert.Equal(t, want, got, "queried meter readings are not correct")

	err = service.Delete(meter3.Id)
	require.NoError(t, err, "no error while deleting expected")

	got, err = service.QueryForSeries(powerSeries.Id)
	require.NoError(t, err, "no error while quering expected")
	assert.Equal(t, 2, len(got), "after deletion, result set should be one smaller than before")
}

func TestMeterReadingServiceImpl_DeleteZero(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	service := NewMeterReadingService(db)

	err := service.Delete(0)
	assert.EqualError(t, err, "cannot delete entity with id 0", "Id = 0 is not allowed for deletion.")
}

func compareMeterReadings(t *testing.T, want domain.MeterReading, got domain.MeterReading, msg string) {
	assert.Equal(t, want.Series, got.Series, "series of %s wrong", msg)
	assert.Equal(t, want.Id, got.Id, "id of %s wrong", msg)
	assert.Equal(t, want.Date, got.Date, "date of %s wrong", msg)
	assert.Equal(t, want.Count, got.Count, "count of %s wrong", msg)
}

func TestMeterReadingServiceImpl_QueryOpenInterval(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()
	service := NewMeterReadingService(db)

	april, _ := time.Parse(util.DateFormat, "2019-04-28")
	may, _ := time.Parse(util.DateFormat, "2019-05-29")
	june1, _ := time.Parse(util.DateFormat, "2019-06-10")
	june2, _ := time.Parse(util.DateFormat, "2019-06-28")
	july, _ := time.Parse(util.DateFormat, "2019-07-02")
	august, _ := time.Parse(util.DateFormat, "2019-08-01")
	dates := []time.Time{april, may, june1, june2, july, august}

	water := domain.Series{Name: "Water"}
	power := domain.Series{Name: "Power"}
	for _, date := range dates {
		waterReading := domain.MeterReading{Count: 636, Date: date, Series: &water}
		powerReading := domain.MeterReading{Count: 700, Date: date, Series: &power}
		_ = service.Save(&powerReading)
		_ = service.Save(&waterReading)
	}
	waterReadings, _ := service.QueryForSeries(water.Id)
	powerReadings, _ := service.QueryForSeries(power.Id)
	require.Equal(t, 6, len(waterReadings), "the water readings were not saved correctly")
	require.Equal(t, 6, len(powerReadings), "the power readings were not saved correctly")

	firstJune, _ := time.Parse(util.DateFormat, "2019-06-01")
	lastJune, _ := time.Parse(util.DateFormat, "2019-06-30")
	got, err := service.QueryOpenInterval(water.Id, firstJune, lastJune)
	require.NoError(t, err, "no error expected")
	require.Equal(t, 4, len(got), "number of got readings not correct")
	for index, want := range waterReadings[1:5] {
		compareMeterReadings(t, want, got[index], fmt.Sprintf("meter reading at index %d", index))
	}

	got, err = service.QueryOpenInterval(power.Id, april, august)
	require.NoError(t, err, "no error expected")
	require.Equal(t, 6, len(got), "number of got readings not correct")
	for index, want := range powerReadings {
		compareMeterReadings(t, want, got[index], fmt.Sprintf("meter reading at index %d", index))
	}
}

func TestMeterReadingServiceImpl_QueryById(t *testing.T) {
	db, _ := CreateInMemoryDb()
	service := NewMeterReadingService(db)
	series := domain.Series{Name: "Test Series"}
	r1 := domain.MeterReading{Count: 666, Date: util.FormatDate(2020, 4, 26), Series: &series}
	r2 := domain.MeterReading{Count: 777, Date: util.FormatDate(2020, 6, 26), Series: &series}
	_ = service.Save(&r1)
	_ = service.Save(&r2)
	t.Run("Query by Id 1", func(t *testing.T) {
		queried1, err := service.QueryById(r1.Id)
		require.NoError(t, err, "no error expected by query")
		assert.Equal(t, r1, queried1, "got meter reading differs from wanted")
	})
	t.Run("Query by Id 2", func(t *testing.T) {
		queried2, err := service.QueryById(r2.Id)
		require.NoError(t, err, "no error expected while querying")
		assert.Equal(t, r2, queried2, "got meter reading differs from wanted")
	})
	t.Run("Id not found", func(t *testing.T) {
		_, err := service.QueryById(uint(55))
		require.EqualError(t, err, "could not query meter reading with id 55: record not found", "error message is wrong")
	})
}
