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

func (r *mutationResolver) CreatePricingPlan(ctx context.Context, plan *NewPricingPlanInput) (*PricingPlan, error) {
	start, err := time.Parse(util.DateFormat, plan.ValidFrom)
	if err != nil {
		return nil, fmt.Errorf("could not parse the validFrom date as RFC3339: %v", err)
	}
	var end *time.Time
	if plan.ValidTo != nil {
		tmp, err := time.Parse(util.DateFormat, *plan.ValidTo)
		if err != nil {
			return nil, fmt.Errorf("could not parse the validTo date as RFC3339: %v", err)
		}
		end = &tmp
	}
	series, err := r.seriesService.QueryById(uint(plan.SeriesID))
	if err != nil {
		return nil, fmt.Errorf("could not find a series with ID %d", plan.SeriesID)
	}
	newPlan := domain.PricingPlan{
		Name:      plan.Name,
		BasePrice: plan.BasePrice,
		UnitPrice: plan.UnitPrice,
		ValidFrom: &start,
		ValidTo:   end,
		Series:    series,
	}
	err = r.planService.Save(&newPlan)
	if err != nil {
		return nil, fmt.Errorf("the pricing plan could not be saved: %v", err)
	}
	return toQLPricingPlan(&newPlan), nil
}

func (r *mutationResolver) CreateMeterReading(ctx context.Context, reading *MeterReadingInput) (*MeterReading, error) {
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

func (r *queryResolver) Statistics(ctx context.Context, seriesId int, startString string, endString string, granularity Granularity) ([]*Statistics, error) {
	readings, err := r.meterService.QueryForSeries(uint(seriesId))
	if err != nil {
		return []*Statistics{}, err
	}
	plans, err := r.planService.QueryForSeries(uint(seriesId))
	if err != nil {
		return []*Statistics{}, err
	}
	start, err := time.Parse(util.DateFormat, startString)
	if err != nil {
		return []*Statistics{}, fmt.Errorf("could not parse the start date \"%s\": %v", startString, err)
	}
	end, err := time.Parse(util.DateFormat, endString)
	if err != nil {
		return []*Statistics{}, fmt.Errorf("could not parse the end date \"%s\": %v", endString, err)
	}
	if start.After(end) {
		return []*Statistics{}, fmt.Errorf("the start date \"%s\" is after the end date \"%s\"", startString, endString)
	}
	if granularity != GranularityMonthly {
		return []*Statistics{}, fmt.Errorf("other granularity than monthly is currently not supported")
	}
	result := make([]*Statistics, 0, 0)
	monthStart := start
	for monthStart != end {
		addedStart := monthStart.AddDate(0, 1, 0)
		month := addedStart.Month()
		year := addedStart.Year()
		monthEnd := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		if end.Before(monthEnd) {
			monthEnd = end
		}
		params := consumption.Parameters{
			Start:    monthStart,
			End:      monthEnd,
			Readings: readings,
			Plans:    plans,
		}
		costs, cons := consumption.Costs(params)
		stats := &Statistics{
			ValidFrom:   monthStart.Format(util.DateFormat),
			ValidTo:     monthEnd.Format(util.DateFormat),
			Costs:       costs,
			Consumption: cons,
		}
		result = append(result, stats)
		monthStart = monthEnd
	}
	return result, nil
}
