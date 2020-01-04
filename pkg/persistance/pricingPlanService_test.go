package orm

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPricingPlanService_CRUD(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	powerSeries := &domain.Series{Name: "Power"}
	waterSeries := &domain.Series{Name: "Water"}

	time1, _ := time.Parse(time.RFC3339, "2018-01-01")
	time2, _ := time.Parse(time.RFC3339, "2018-12-31")
	time3, _ := time.Parse(time.RFC3339, "2020-01-01")
	powerPlan1 := domain.PricingPlan{
		Name:      "Power Plan 1",
		BasePrice: 20,
		UnitPrice: 0.5,
		ValidFrom: &time1,
		ValidTo:   &time2,
		Series:    powerSeries,
	}

	powerPlan2 := domain.PricingPlan{
		Name:      "Power Plan 2",
		BasePrice: 21,
		UnitPrice: 0.6,
		ValidFrom: &time3,
		Series:    powerSeries,
	}

	waterPlan1 := domain.PricingPlan{
		Name:      "Water Plan 1",
		BasePrice: 10.5,
		UnitPrice: 0.1,
		ValidFrom: &time3,
		Series:    waterSeries,
	}

	service := NewPricingPlanService(db)
	err := service.Save(&powerPlan1)
	require.NoError(t, err)
	err = service.Save(&powerPlan2)
	require.NoError(t, err)
	err = service.Save(&waterPlan1)
	require.NoError(t, err)
	plans := []domain.PricingPlan{powerPlan1, powerPlan2, waterPlan1}

	got, err := service.QueryAll()
	require.NoError(t, err, "no error while querying expected.")
	assert.Equal(t, plans, got, "actual and wanted pricing plans differ")

	err = service.Delete(plans[1].Id)
	require.NoError(t, err, "no error while deleting expected")
	got, err = service.QueryAll()
	require.NoError(t, err, "no error while querying expected")
	assert.Equal(t, 2, len(got), "after deletion there should be one less pricing plan")
}

func TestPricingPlanServiceImpl_DeleteZero(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	service := NewPricingPlanService(db)

	err := service.Delete(0)
	assert.EqualError(t, err, "cannot delete entity with id 0", "Id = 0 is not allowed for deletion.")
}
