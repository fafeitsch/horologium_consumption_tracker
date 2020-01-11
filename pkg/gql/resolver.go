package gql

import (
	"context"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

func NewResolver(seriesService orm.SeriesService, planService orm.PricingPlanService) Resolver {
	if seriesService == nil {
		panic("the series service is nil")
	}
	if planService == nil {
		panic("the plan service is nil")
	}
	return &resolverImpl{seriesService: seriesService, planService: planService}
}

type resolverImpl struct {
	seriesService orm.SeriesService
	planService   orm.PricingPlanService
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
	start, err := time.Parse(orm.DateFormat, plan.ValidFrom)
	if err != nil {
		return nil, fmt.Errorf("could not parse the validFrom date as RFC3339: %v", err)
	}
	end := new(time.Time)
	if plan.ValidTo != nil {
		tmp, err := time.Parse(orm.DateFormat, *plan.ValidTo)
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
