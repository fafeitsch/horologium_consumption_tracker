package gql

import (
	"context"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestMutationResolver_CreateSeries(t *testing.T) {
	service := new(mockSeriesService)
	series := domain.Series{Name: "A new series", Id: 6}
	service.On("Save").Return(6, nil)
	resolver := NewResolver(service)

	newSeries := NewSeriesInput{Name: series.Name}
	got, err := resolver.Mutation().CreateSeries(context.Background(), newSeries)
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, got.Name, newSeries.Name, "the names should be equal")
	assert.Equal(t, strconv.Itoa(int(series.Id)), got.ID, "an ID must be given")
}

func TestMutationResolver_DeleteSeries(t *testing.T) {
	service := new(mockSeriesService)
	service.On("Delete", uint(44)).Return(nil)
	resolver := NewResolver(service)

	id, err := resolver.Mutation().DeleteSeries(context.Background(), 44)
	assert.Equal(t, 44, id, "reported id should be 44")
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, 1, len(service.Calls), "exactly on call should have happened on the service")
}

func TestQueryResolver_Series(t *testing.T) {
	service := new(mockSeriesService)
	series := domain.Series{Id: 55, Name: "Water"}
	service.On("QueryById", uint(55)).Return(&series, nil)
	resolver := NewResolver(service)

	got, err := resolver.Query().Series(context.Background(), 55)
	assert.NoError(t, err, "no error expected")
	compareSeries(t, series, got, "got series differs from expected")
}

func TestQueryResolver_AllSeries(t *testing.T) {
	service := new(mockSeriesService)
	series := []domain.Series{{Id: 25, Name: "Power"}, {Id: 33, Name: "Water"}}
	service.On("QueryAll").Return(series, nil)
	resolver := NewResolver(service)

	got, err := resolver.Query().AllSeries(context.Background())
	assert.NoError(t, err, "no error expected")
	require.Equal(t, len(series), len(got), "number of series not correct")
	for index, s := range series {
		compareSeries(t, s, got[index], fmt.Sprintf("series at index %d", index))
	}
}

func compareSeries(t *testing.T, s domain.Series, got *Series, msg string) {
	assert.Equal(t, strconv.Itoa(int(s.Id)), got.ID, "id of %d")
	assert.Equal(t, s.Name, got.Name, "name of %s", msg)
}

type mockSeriesService struct {
	mock.Mock
}

func (m *mockSeriesService) Save(series *domain.Series) error {
	idToSet := m.Called().Int(0)
	err := m.Called().Error(1)
	if err != nil {
		return err
	}
	series.Id = uint(idToSet)
	return nil
}

func (m *mockSeriesService) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *mockSeriesService) QueryById(id uint) (*domain.Series, error) {
	args := m.Called(id).Get(0)
	err := m.Called(id).Error(1)
	if err != nil {
		return nil, err
	}
	return args.(*domain.Series), nil
}

func (m *mockSeriesService) QueryAll() ([]domain.Series, error) {
	args := m.Called().Get(0)
	err := m.Called().Error(1)
	if err != nil {
		return nil, err
	}
	return args.([]domain.Series), nil
}
