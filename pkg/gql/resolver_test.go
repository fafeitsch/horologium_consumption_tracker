package gql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
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
	validFrom1, _ := time.Parse(util.DateFormat, "2018-01-01")
	validTo1, _ := time.Parse(util.DateFormat, "2018-12-31")
	planService.On("Save").Return(677, nil)
	resolver := NewResolver(seriesService, planService, readingService)
	validToStr := validTo1.Format(util.DateFormat)

	t.Run("with both dates", func(t *testing.T) {
		newPlan := PricingPlanInput{Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: validFrom1.Format(util.DateFormat), ValidTo: &validToStr, SeriesID: 27}
		got, err := resolver.Mutation().CreatePricingPlan(context.Background(), newPlan)
		assert.NoError(t, err, "no error expected")
		plan := domain.PricingPlan{Id: 677, Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: &validFrom1, ValidTo: &validTo1, Series: &series}
		comparePlans(t, plan, got, "created plan")
	})
	t.Run("with only start date", func(t *testing.T) {
		newPlan := PricingPlanInput{Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: validFrom1.Format(util.DateFormat), ValidTo: nil, SeriesID: 27}
		got, err := resolver.Mutation().CreatePricingPlan(context.Background(), newPlan)
		assert.NoError(t, err, "no error expected")
		plan := domain.PricingPlan{Id: 677, Name: "Power 2020", BasePrice: 40, UnitPrice: 23, ValidFrom: &validFrom1, ValidTo: nil, Series: &series}
		comparePlans(t, plan, got, "created plan")
	})
}

func TestMutationResolver_ModifyPricingPlan(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	resolver := NewResolver(seriesService, planService, readingService)
	series := domain.Series{Id: 5}
	plan := domain.PricingPlan{Id: 77, ValidFrom: util.FormatDatePtr(2019, 6, 20), BasePrice: 44, UnitPrice: 2.3, Series: &series}

	planService.On("QueryById", plan.Id).Return(&plan, nil)
	planService.On("Save").Return(int(plan.Id), nil)
	planService.On("QueryById", uint(34)).Return(nil, fmt.Errorf("not found"))
	t.Run("success", func(t *testing.T) {
		validTo := "2020-12-04"
		change := PricingPlanChange{
			ValidFrom: "2020-06-20",
			ValidTo:   &validTo,
			BasePrice: 3,
			UnitPrice: 5,
			ID:        77,
		}
		got, err := resolver.Mutation().ModifyPricingPlan(context.Background(), change)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, got.ValidFrom, change.ValidFrom, "the validFrom was not changed")
		assert.Equal(t, got.ValidTo, change.ValidTo, "the validTo was not changed")
		assert.Equal(t, got.BasePrice, change.BasePrice, "the basePrice was not changed")
		assert.Equal(t, got.UnitPrice, change.UnitPrice, "the unitPrice was not changed")
	})
	t.Run("incorrect date", func(t *testing.T) {
		change := PricingPlanChange{
			ValidFrom: "not a date",
			BasePrice: 3,
			UnitPrice: 5,
		}
		_, err := resolver.Mutation().ModifyPricingPlan(context.Background(), change)
		require.EqualError(t, err, "could not parse \"not a date\" as format YYYY-MM-DD", "error message not correct")
	})
	t.Run("not found", func(t *testing.T) {
		_, err := resolver.Mutation().ModifyPricingPlan(context.Background(), PricingPlanChange{ID: 34, ValidFrom: "2020-03-04"})
		require.EqualError(t, err, "could not find pricing plan with id 34: not found", "error message not correct")
	})
}

func Test_mutationResolver_CreateMeterReading(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	series := domain.Series{
		Id:   62,
		Name: "Chocolate Consumption",
	}
	date, _ := time.Parse(util.DateFormat, "2020-01-20")
	seriesService.On("QueryById", uint(62)).Return(&series, nil)
	readingService.On("Save").Return(82, nil)

	newReading := MeterReadingInput{
		Count: 53.2, Date: date.Format(util.DateFormat), SeriesID: 62,
	}
	resolver := NewResolver(seriesService, planService, readingService)
	got, err := resolver.Mutation().CreateMeterReading(context.Background(), newReading)
	assert.NoError(t, err, "no error expected")
	reading := domain.MeterReading{Id: 82, Date: date, Series: &series, Count: 53.2}
	compareMeterReadings(t, reading, got, "created reading")
}

func compareMeterReadings(t *testing.T, expected domain.MeterReading, got *MeterReading, msg string) {
	assert.Equal(t, expected.Date.Format(util.DateFormat), got.Date, "date of %s is wrong", msg)
	assert.Equal(t, expected.Count, got.Count, "count of %s is wrong", msg)
	assert.Equal(t, expected.Series.Id, uint(got.SeriesID), "seriesId of %s is wrong", msg)
	assert.Equal(t, expected.Id, uint(got.ID), "id of %s is wrong", msg)
}

func TestMutationResolver_ModifyMeterReading(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	resolver := NewResolver(seriesService, planService, readingService)
	series := domain.Series{Id: 5}
	reading := domain.MeterReading{Id: 77, Date: util.FormatDate(2019, 9, 23), Count: 242, Series: &series}
	readingService.On("QueryById", reading.Id).Return(&reading, nil)
	readingService.On("Save").Return(int(reading.Id), nil)
	readingService.On("QueryById", uint(34)).Return(nil, fmt.Errorf("not found"))
	t.Run("success", func(t *testing.T) {
		change := MeterReadingChange{
			Count: 700,
			Date:  "2019-10-23",
			ID:    77,
		}
		got, err := resolver.Mutation().ModifyMeterReading(context.Background(), change)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, got.Count, change.Count, "the count was not changed")
		assert.Equal(t, got.Date, change.Date, "the date was not changed")
	})
	t.Run("incorrect date", func(t *testing.T) {
		change := MeterReadingChange{
			Count: 700,
			Date:  "not a date",
			ID:    77,
		}
		_, err := resolver.Mutation().ModifyMeterReading(context.Background(), change)
		require.EqualError(t, err, "could not format date: parsing time \"not a date\" as \"2006-01-02\": cannot parse \"not a date\" as \"2006\"", "error message not correct")
	})
	t.Run("not found", func(t *testing.T) {
		_, err := resolver.Mutation().ModifyMeterReading(context.Background(), MeterReadingChange{ID: 34})
		require.EqualError(t, err, "could not find meter reading with id 34: not found", "error message not correct")
	})
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
	validFrom1, _ := time.Parse(util.DateFormat, "2018-01-01")
	validTo1, _ := time.Parse(util.DateFormat, "2018-12-31")
	validFrom2, _ := time.Parse(util.DateFormat, "2019-01-01")
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
	validFrom := p.ValidFrom.Format(util.DateFormat)
	assert.Equal(t, validFrom, got.ValidFrom, "validFrom of %s", msg)
	if p.ValidTo != nil {
		validTo := p.ValidTo.Format(util.DateFormat)
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
	got, err := resolver.Query().MeterReadings(context.Background(), 93)
	assert.NoError(t, err, "no error expected")
	require.Equal(t, len(readings), len(got), "number of got meter readings is not correct")
	for index, gotReading := range got {
		compareMeterReadings(t, readings[index], gotReading, fmt.Sprintf("reading at %d", index))
	}
}

func TestQueryResolver_MonthlyStatisticsSuccess(t *testing.T) {
	seriesService, planService, readingService := createMockServices()
	plans := []domain.PricingPlan{{
		Id:        66,
		BasePrice: 10.8,
		UnitPrice: 3.5,
		ValidFrom: util.FormatDatePtr(2020, 6, 1),
		ValidTo:   util.FormatDatePtr(2020, 12, 1),
	}}
	planService.On("QueryForSeries", uint(15)).Return(plans, nil)
	reading1 := domain.MeterReading{
		Count: 25.2,
		Date:  util.FormatDate(2020, 8, 1),
	}
	reading2 := domain.MeterReading{
		Count: 74.2,
		Date:  util.FormatDate(2020, 10, 15),
	}
	start := util.FormatDate(2020, 9, 1)
	end := util.FormatDate(2020, 10, 15)
	readings := []domain.MeterReading{reading1, reading2}
	readingService.On("QueryOpenInterval", uint(15), start, end).Return(readings, nil)
	resolver := NewResolver(seriesService, planService, readingService)
	got, err := resolver.Query().MonthlyStatistics(context.Background(), 15, "2020-09-01", "2020-10-15")
	require.NoError(t, err, "no error expected")
	require.Equal(t, 2, len(got), "two statistics expected")
	require.Equal(t, 19.599999999999994, got[0].Consumption, "consumption of first statistic is wrong")
	require.Equal(t, 79.39999999999998, got[0].Costs, "costs of first statistic are wrong")
	require.Equal(t, "2020-09-01", got[0].ValidFrom, "valid_from of first statistic is wrong")
	require.Equal(t, "2020-10-01", got[0].ValidTo, "valid_to of first statistic is wrong")
}

func TestQueryResolver_MonthlyStatisticsErrors(t *testing.T) {
	testcases := []struct {
		start         string
		end           string
		planError     error
		readingsError error
		want          string
	}{
		{start: "2019-01-01", end: "2020-01-01", planError: errors.New("plans could not be loaded"), readingsError: nil, want: "plans could not be loaded"},
		{start: "2019-01-01", end: "2020-01-01", planError: nil, readingsError: errors.New("readings could not be loaded"), want: "readings could not be loaded"},
		{start: "nodate", end: "2020-01-01", planError: nil, readingsError: nil, want: "could not parse \"nodate\" as format YYYY-MM-DD"},
		{start: "2019-01-01", end: "2018-12-31", planError: nil, readingsError: nil, want: "the start date \"2019-01-01\" is after the end date \"2018-12-31\""},
	}
	for _, tt := range testcases {
		t.Run(tt.want, func(t *testing.T) {
			seriesService, planService, readingService := createMockServices()
			planService.On("QueryForSeries", uint(25)).Return([]domain.PricingPlan{}, tt.planError)
			readingService.On("QueryOpenInterval", uint(25), mock.Anything, mock.Anything).Return([]domain.MeterReading{}, tt.readingsError)
			resolver := NewResolver(seriesService, planService, readingService)
			got, err := resolver.Query().MonthlyStatistics(context.Background(), 25, tt.start, tt.end)
			assert.Equal(t, 0, len(got), "there should be no object returned")
			assert.EqualError(t, err, tt.want, "the error message is wrong")
		})
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

func (m mockPricingPlanService) QueryForTime(uint, time.Time) (*domain.PricingPlan, error) {
	panic("implement me")
}

func (m mockPricingPlanService) QueryById(id uint) (*domain.PricingPlan, error) {
	args := m.Called(id).Get(0)
	err := m.Called(id).Error(1)
	if err != nil {
		return nil, err
	}
	return args.(*domain.PricingPlan), nil
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

func (m mockReadingService) QueryOpenInterval(seriesId uint, start time.Time, end time.Time) ([]domain.MeterReading, error) {
	args := m.Called(seriesId, start, end).Get(0)
	err := m.Called(seriesId, start, end).Error(1)
	if err != nil {
		return nil, err
	}
	return args.([]domain.MeterReading), nil
}

func (m mockReadingService) QueryById(u uint) (*domain.MeterReading, error) {
	args := m.Called(u).Get(0)
	err := m.Called(u).Error(1)
	if err != nil {
		return nil, err
	}
	return args.(*domain.MeterReading), nil
}
