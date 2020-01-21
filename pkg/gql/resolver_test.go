package gql

import (
	"context"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMutationResolver_CreateSeries(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := domain.Series{Name: "A new series", Id: 6}
	seriesService.On("Save").Return(6, nil)
	resolver := NewResolver(seriesService, planService, readingService)

	newSeries := NewSeriesInput{Name: series.Name}
	got, err := resolver.Mutation().CreateSeries(context.Background(), newSeries)
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, got.Name, newSeries.Name, "the names should be equal")
	assert.Equal(t, int(series.Id), got.ID, "an ID must be given")
}

func TestMutationResolver_DeleteSeries(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	seriesService.On("Delete", uint(44)).Return(nil)
	resolver := NewResolver(seriesService, planService, readingService)

	id, err := resolver.Mutation().DeleteSeries(context.Background(), 44)
	assert.Equal(t, 44, id, "reported id should be 44")
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, 1, len(seriesService.Calls), "exactly on call should have happened on the service")
}

func TestMutationResolver_CreatePricingPlan(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := domain.Series{Id: 27, Name: "Power"}
	seriesService.On("QueryById", uint(27)).Return(&series, nil)
	validFrom1, _ := time.Parse(orm.DateFormat, "2018-01-01")
	validTo1, _ := time.Parse(orm.DateFormat, "2018-12-31")
	planService.On("Save").Return(677, nil)
	resolver := NewResolver(seriesService, planService, readingService)
	validToStr := validTo1.Format(orm.DateFormat)

	t.Run("with both dates", func(t *testing.T) {
		newPlan := NewPricingPlanInput{Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: validFrom1.Format(orm.DateFormat), ValidTo: &validToStr, SeriesID: 27}
		got, err := resolver.Mutation().CreatePricingPlan(context.Background(), &newPlan)
		assert.NoError(t, err, "no error expected")
		plan := domain.PricingPlan{Id: 677, Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: &validFrom1, ValidTo: &validTo1, Series: &series}
		comparePlans(t, plan, got, "created plan")
	})
	t.Run("with only start date", func(t *testing.T) {
		newPlan := NewPricingPlanInput{Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: validFrom1.Format(orm.DateFormat), ValidTo: nil, SeriesID: 27}
		got, err := resolver.Mutation().CreatePricingPlan(context.Background(), &newPlan)
		assert.NoError(t, err, "no error expected")
		plan := domain.PricingPlan{Id: 677, Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: &validFrom1, ValidTo: nil, Series: &series}
		comparePlans(t, plan, got, "created plan")
	})
}

func Test_mutationResolver_CreateMeterReading(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := domain.Series{
		Id:   62,
		Name: "Chocolate Consumption",
	}
	date, _ := time.Parse(orm.DateFormat, "2020-01-20")
	seriesService.On("QueryById", uint(62)).Return(&series, nil)
	readingService.On("Save").Return(82, nil)

	newReading := MeterReadingInput{
		Count: 53.2, Date: date.Format(orm.DateFormat), SeriesID: 62,
	}
	resolver := NewResolver(seriesService, planService, readingService)
	got, err := resolver.Mutation().CreateMeterReading(context.Background(), &newReading)
	assert.NoError(t, err, "no error expected")
	reading := domain.MeterReading{Id: 82, Date: date, Series: &series, Count: 53.2}
	compareMeterReadings(t, reading, got, "created reading")
}

func compareMeterReadings(t *testing.T, expected domain.MeterReading, got *MeterReading, msg string) {
	assert.Equal(t, expected.Date.Format(orm.DateFormat), got.Date, "date of %s is wrong", msg)
	assert.Equal(t, expected.Count, got.Count, "count of %s is wrong", msg)
	assert.Equal(t, expected.Series.Id, uint(got.SeriesID), "seriesId of %s is wrong", msg)
	assert.Equal(t, expected.Id, uint(got.ID), "id of %s is wrong", msg)
}

func TestQueryResolver_Series(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := domain.Series{Id: 55, Name: "Water"}
	seriesService.On("QueryById", uint(55)).Return(&series, nil)
	resolver := NewResolver(seriesService, planService, readingService)

	got, err := resolver.Query().Series(context.Background(), 55)
	assert.NoError(t, err, "no error expected")
	compareSeries(t, series, got, "got series differs from expected")
}

func compareSeries(t *testing.T, s domain.Series, got *Series, msg string) {
	assert.Equal(t, int(s.Id), got.ID, "id of %d")
	assert.Equal(t, s.Name, got.Name, "name of %s", msg)
}

func TestQueryResolver_AllSeries(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := []domain.Series{{Id: 25, Name: "Power"}, {Id: 33, Name: "Water"}}
	seriesService.On("QueryAll").Return(series, nil)
	resolver := NewResolver(seriesService, planService, readingService)

	got, err := resolver.Query().AllSeries(context.Background())
	assert.NoError(t, err, "no error expected")
	require.Equal(t, len(series), len(got), "number of series not correct")
	for index, s := range series {
		compareSeries(t, s, got[index], fmt.Sprintf("series at index %d", index))
	}
}

func TestQueryResolver_PricingPlans(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	validFrom1, _ := time.Parse(orm.DateFormat, "2018-01-01")
	validTo1, _ := time.Parse(orm.DateFormat, "2018-12-31")
	validFrom2, _ := time.Parse(orm.DateFormat, "2019-01-01")
	series := domain.Series{Id: 25, Name: "Power"}
	plans := []domain.PricingPlan{
		{Id: 5, Name: "Year 2018", BasePrice: 12, UnitPrice: 0.34, ValidFrom: &validFrom1, ValidTo: &validTo1, Series: &series},
		{Id: 6, Name: "Year 2019", BasePrice: 13, UnitPrice: 0.35, ValidFrom: &validFrom2, Series: &series},
	}
	planService.On("QueryForSeries", uint(25)).Return(plans, nil)
	resolver := NewResolver(seriesService, planService, readingService)

	got, err := resolver.Query().PricingPlans(context.Background(), 25)
	assert.NoError(t, err, "no error expected")
	require.Equal(t, 2, len(got), "number of pricing plans not correct")
	for index, actual := range got {
		comparePlans(t, plans[index], actual, fmt.Sprintf("plan at index %d", index))
	}
}

func comparePlans(t *testing.T, p domain.PricingPlan, got *PricingPlan, msg string) {
	assert.Equal(t, int(p.Id), got.ID, "id of %s", msg)
	assert.Equal(t, p.Name, got.Name, "name of %s", msg)
	assert.Equal(t, p.BasePrice, got.BasePrice, "base price of %s", msg)
	assert.Equal(t, p.UnitPrice, got.UnitPrice, "unit price of %s", msg)
	validFrom := p.ValidFrom.Format(orm.DateFormat)
	assert.Equal(t, validFrom, got.ValidFrom, "validFrom of %s", msg)
	if p.ValidTo != nil {
		validTo := p.ValidTo.Format(orm.DateFormat)
		assert.Equal(t, validTo, *got.ValidTo, "validTo of %s", msg)
	} else {
		assert.Nil(t, got.ValidTo, "validTo of %s", msg)
	}
}

func TestQueryResolver_MeterReadings(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	now := time.Now()
	series := domain.Series{Name: "Beer", Id: 93}
	reading1 := domain.MeterReading{
		Id:     62,
		Count:  25.2,
		Date:   now,
		Series: &series,
	}
	reading2 := domain.MeterReading{
		Id:     43,
		Count:  74.2,
		Date:   now.Add(24 * time.Hour),
		Series: &series,
	}
	readings := []domain.MeterReading{reading1, reading2}
	readingService.On("QueryForSeries", uint(93)).Return(readings, nil)

	resolver := NewResolver(seriesService, planService, readingService)
	got, err := resolver.Query().MeterReadings(context.Background(), &MeterReadingQuery{SeriesID: 93})
	assert.NoError(t, err, "no error expected")
	require.Equal(t, len(readings), len(got), "number of got meter readings is not correct")
	for index, gotReading := range got {
		compareMeterReadings(t, readings[index], gotReading, fmt.Sprintf("reading at %d", index))
	}
}

func createMockServices() (*mockSeriesService, *mockPricingPlanService, *mockReadingService) {
	return new(mockSeriesService), new(mockPricingPlanService), new(mockReadingService)
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

type mockPricingPlanService struct {
	mock.Mock
}

func (m mockPricingPlanService) Save(d *domain.PricingPlan) error {
	idToSet := m.Called().Int(0)
	err := m.Called().Error(1)
	if err != nil {
		return err
	}
	d.Id = uint(idToSet)
	return nil
}

func (m mockPricingPlanService) Delete(uint) error {
	panic("implement me")
}

func (m mockPricingPlanService) QueryAll() ([]domain.PricingPlan, error) {
	panic("implement me")
}

func (m mockPricingPlanService) QueryForSeries(id uint) ([]domain.PricingPlan, error) {
	args := m.Called(id).Get(0)
	err := m.Called(id).Error(1)
	if err != nil {
		return nil, err
	}
	return args.([]domain.PricingPlan), nil
}

type mockReadingService struct {
	mock.Mock
}

func (m mockReadingService) Save(reading *domain.MeterReading) error {
	idToSet := m.Called().Int(0)
	err := m.Called().Error(1)
	if err != nil {
		return err
	}
	reading.Id = uint(idToSet)
	return nil
}

func (m mockReadingService) Delete(uint) error {
	panic("implement me")
}

func (m mockReadingService) QueryForSeries(seriesId uint) ([]domain.MeterReading, error) {
	args := m.Called(seriesId).Get(0)
	err := m.Called(seriesId).Error(1)
	if err != nil {
		return nil, err
	}
	return args.([]domain.MeterReading), nil
}
