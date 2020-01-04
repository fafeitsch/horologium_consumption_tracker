package orm

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSeriesServiceImpl_CRUD(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	service := NewSeriesService(db)

	testify := assert.New(t)

	written := make([]domain.Series, 0, 10)
	for i := 0; i < 10; i++ {
		series := domain.Series{Name: fmt.Sprintf("Series %d", i)}
		err := service.Save(&series)
		written = append(written, series)
		testify.NoError(err, "error while saving series %v", series)
	}

	got, err := service.QueryAll()
	testify.NoError(err)
	testify.Equal(written, got, "queried series differ from expected ones")

	err = service.Delete(written[4].Id)
	testify.NoError(err, "deleting should not throw an error, but did")
	got, err = service.QueryAll()
	testify.Equal(9, len(got), "after deletion, there should be one series missing, but was not")
}

func TestSeriesServiceImpl_DeleteZero(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	service := NewSeriesService(db)

	err := service.Delete(0)
	assert.EqualError(t, err, "cannot delete entity with id 0", "Id = 0 is not allowed for deletion.")
}

func Test_seriesServiceImpl_QueryById(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()
	service := NewSeriesService(db)

	series1 := domain.Series{Name: "Series ABC"}
	series2 := domain.Series{Name: "Series XYZ"}
	err := service.Save(&series1)
	require.NoError(t, err)
	err = service.Save(&series2)
	require.NoError(t, err)

	got1, err := service.QueryById(series1.Id)
	assert.NoError(t, err, "no error while querying expected")
	assert.Equal(t, series1, *got1, "queried series 1 is not as expected")

	got2, err := service.QueryById(series2.Id)
	assert.NoError(t, err, "no error while querying expected")
	assert.Equal(t, series2, *got2, "queried series 2 is not as expected")
}
