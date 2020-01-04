package gql

import (
	"context"
	"github.com/fafeitsch/Horologium/pkg/domain"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

func NewResolver(seriesService orm.SeriesService) Resolver {
	if seriesService == nil {
		panic("the series seriesService is nil")
	}
	return &resolverImpl{seriesService: seriesService}
}

type resolverImpl struct {
	seriesService orm.SeriesService
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
