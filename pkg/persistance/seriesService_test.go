package orm

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeriesServiceImpl_CRUD(t *testing.T) {
	db, _ := createInMemoryDb()
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
	db, _ := createInMemoryDb()
	defer func() { _ = db.Close() }()

	service := NewSeriesService(db)

	err := service.Delete(0)
	assert.EqualError(t, err, "cannot delete entity with id 0", "Id = 0 is not allowed for deletion.")
}
