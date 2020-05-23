package gql

import (
	"context"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/consumption"
	"github.com/fafeitsch/Horologium/pkg/domain"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
	"github.com/fafeitsch/Horologium/pkg/util"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

func NewResolver(seriesService orm.SeriesService, planService orm.PricingPlanService, meterService orm.MeterReadingService) Resolver {
	if seriesService == nil {
		panic("the series service is nil")
	}
	if planService == nil {
		panic("the plan service is nil")
	}
	if meterService == nil {
		panic("the meter service is nil")
	}
	return &resolverImpl{seriesService: seriesService, planService: planService, meterService: meterService}
}

type resolverImpl struct {
	seriesService orm.SeriesService
	planService   orm.PricingPlanService
	meterService  orm.MeterReadingService
}

func (r *resolverImpl) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *resolverImpl) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *resolverImpl }

func (r *mutationResolver) CreateSeries(ctx context.Context, series NewSeriesInput) (*Series, error) {
	newSeries := domain.Series{
		Name: series.Name,
	}
	err := r.seriesService.Save(&newSeries)
	return toQLSeries(&newSeries), err
}

func (r *mutationResolver) DeleteSeries(ctx context.Context, id int) (int, error) {
	return id, r.seriesService.Delete(uint(id))
}

func (r *mutationResolver) CreatePricingPlan(ctx context.Context, plan PricingPlanInput) (*PricingPlan, error) {
	dates, err := parseDates(&plan.ValidFrom, plan.ValidTo)
	if err != nil {
		return nil, fmt.Errorf("could not parse date as format YYYY-MM-DD: %v", err)
	}
	series, err := r.seriesService.QueryById(uint(plan.SeriesID))
	if err != nil {
		return nil, fmt.Errorf("could not find a series with ID %d", plan.SeriesID)
	}
	newPlan := domain.PricingPlan{
		Name:      plan.Name,
		BasePrice: plan.BasePrice,
		UnitPrice: plan.UnitPrice,
		ValidFrom: dates[0],
		ValidTo:   dates[1],
		Series:    series,
	}
	err = r.planService.Save(&newPlan)
	if err != nil {
		return nil, fmt.Errorf("the pricing plan could not be saved: %v", err)
	}
	return toQLPricingPlan(&newPlan), nil
}

func (r *mutationResolver) ModifyPricingPlan(ctx context.Context, plan PricingPlanChange) (*PricingPlan, error) {
	dates, err := parseDates(&plan.ValidFrom, plan.ValidTo)
	if err != nil {
		return nil, err
	}
	existingPlan, err := r.planService.QueryById(uint(plan.ID))
	if err != nil {
		return nil, fmt.Errorf("could not find pricing plan with id %d: %v", plan.ID, err)
	}
	existingPlan.Name = plan.Name
	existingPlan.UnitPrice = plan.UnitPrice
	existingPlan.BasePrice = plan.BasePrice
	existingPlan.ValidFrom = dates[0]
	existingPlan.ValidTo = dates[1]
	err = r.planService.Save(existingPlan)
	if err != nil {
		return nil, fmt.Errorf("the pricing plan could not be saved: %v", err)
	}
	return toQLPricingPlan(existingPlan), nil
}

func (r *mutationResolver) CreateMeterReading(ctx context.Context, reading MeterReadingInput) (*MeterReading, error) {
	date, err := time.Parse(util.DateFormat, reading.Date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date \"%s\" as \"%s\"", reading.Date, util.DateFormat)
	}
	series, err := r.seriesService.QueryById(uint(reading.SeriesID))
	if err != nil {
		return nil, fmt.Errorf("could not find a series with ID %d", reading.SeriesID)
	}
	newReading := domain.MeterReading{
		Count:  reading.Count,
		Date:   date,
		Series: series,
	}
	err = r.meterService.Save(&newReading)
	if err != nil {
		return nil, fmt.Errorf("the meter reading could not be saved: %v", err)
	}
	return toQlMeterReading(&newReading), nil
}

func (r *mutationResolver) ModifyMeterReading(ctx context.Context, input MeterReadingChange) (*MeterReading, error) {
	existing, err := r.meterService.QueryById(uint(input.ID))
	if err != nil {
		return nil, fmt.Errorf("could not find meter reading with id %d: %v", input.ID, err)
	}
	date, err := time.Parse(util.DateFormat, input.Date)
	if err != nil {
		return nil, fmt.Errorf("could not format date: %v", err)
	}
	existing.Date = date
	existing.Count = input.Count
	err = r.meterService.Save(existing)
	return toQlMeterReading(existing), err
}

type queryResolver struct{ *resolverImpl }

func (r *queryResolver) AllSeries(ctx context.Context) ([]*Series, error) {
	dbResult, err := r.seriesService.QueryAll()
	if err != nil {
		return []*Series{}, err
	}
	result := make([]*Series, 0, len(dbResult))
	for _, res := range dbResult {
		result = append(result, toQLSeries(&res))
	}
	return result, nil
}

func (r *queryResolver) Series(ctx context.Context, id int) (*Series, error) {
	dbResult, err := r.seriesService.QueryById(uint(id))
	return toQLSeries(dbResult), err
}

func (r *queryResolver) PricingPlans(ctx context.Context, seriesID int) ([]*PricingPlan, error) {
	dbResult, err := r.planService.QueryForSeries(uint(seriesID))
	if err != nil {
		return []*PricingPlan{}, err
	}
	result := make([]*PricingPlan, 0, len(dbResult))
	for _, res := range dbResult {
		result = append(result, toQLPricingPlan(&res))
	}
	return result, nil
}

func (r *queryResolver) MeterReadings(ctx context.Context, id int) ([]*MeterReading, error) {
	dbResult, err := r.meterService.QueryForSeries(uint(id))
	if err != nil {
		return []*MeterReading{}, err
	}
	result := make([]*MeterReading, 0, len(dbResult))
	for _, res := range dbResult {
		result = append(result, toQlMeterReading(&res))
	}
	return result, nil
}

func (r *queryResolver) MonthlyStatistics(ctx context.Context, seriesId int, startString string, endString string) ([]*Statistics, error) {
	plans, err := r.planService.QueryForSeries(uint(seriesId))
	if err != nil {
		return []*Statistics{}, err
	}
	dates, err := parseDates(&startString, &endString)
	if err != nil {
		return nil, err
	}
	start, end := *dates[0], *dates[1]
	if start.After(end) {
		return []*Statistics{}, fmt.Errorf("the start date \"%s\" is after the end date \"%s\"", startString, endString)
	}
	readings, err := r.meterService.QueryOpenInterval(uint(seriesId), start, end)
	if err != nil {
		return []*Statistics{}, err
	}
	params := consumption.Parameters{
		Start:    start,
		End:      end,
		Readings: readings,
		Plans:    plans,
	}
	stats := consumption.MonthlyCosts(params)
	result := make([]*Statistics, 0, len(stats))
	for _, stat := range stats {
		result = append(result, toQlStatistics(&stat))
	}
	return result, nil
}
