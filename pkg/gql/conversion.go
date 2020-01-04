package gql

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"strconv"
)

func toQLSeries(series *domain.Series) *Series {
	if series == nil {
		return nil
	}
	return &Series{
		ID:   strconv.Itoa(int(series.Id)),
		Name: series.Name,
	}
}
