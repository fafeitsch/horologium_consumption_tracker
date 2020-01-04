package orm

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
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

	time1, _ := time.Parse(time.RFC3339, "2018-01-31")
	time2, _ := time.Parse(time.RFC3339, "2018-02-28")
	time3, _ := time.Parse(time.RFC3339, "2018-03-31")
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
