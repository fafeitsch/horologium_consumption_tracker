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

	got, err = service.QueryForSeries(waterSeries.Id)
	require.NoError(t, err, "no error while querying expected")
	assert.Equal(t, 1, len(got), "number of got plans incorrect")
	assert.Equal(t, waterPlan1, got[0], "got pricing plan not correct")

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

func comparePrincingPlans(t *testing.T, got, want *domain.PricingPlan) {
	assert.Equal(t, got.Name, want.Name, "name is different")
	assert.Equal(t, got.ValidTo, want.ValidTo, "validTo is different")
	assert.Equal(t, got.ValidFrom, want.ValidFrom, "validFrom is different")
	assert.Equal(t, got.Series, want.Series, "series in different")
}

func TestPricingPlanServiceImpl_QueryForTime(t *testing.T) {
	db, _ := CreateInMemoryDb()
	defer func() { _ = db.Close() }()

	april, _ := time.Parse(DateFormat, "2019-04-01")
	lastApril, _ := time.Parse(DateFormat, "2019-04-30")
	may, _ := time.Parse(DateFormat, "2019-05-01")
	lastMay, _ := time.Parse(DateFormat, "2019-05-31")
	june, _ := time.Parse(DateFormat, "2019-06-01")
	lastJune, _ := time.Parse(DateFormat, "2019-06-30")
	readingTime, _ := time.Parse(DateFormat, "2019-05-30")

	service := NewPricingPlanService(db)

	water := domain.Series{Name: "Water", Id: 4}
	power := domain.Series{Name: "Power", Id: 5}

	aprilWaterPlan := domain.PricingPlan{Name: "April", Series: &water, ValidFrom: &april, ValidTo: &lastApril}
	mayWaterPlan := domain.PricingPlan{Name: "May", Series: &water, ValidFrom: &may, ValidTo: &lastMay}
	mayPowerPlan := domain.PricingPlan{Name: "May Power", Series: &power, ValidFrom: &may, ValidTo: &lastMay}
	juneWaterPlan := domain.PricingPlan{Name: "June", Series: &water, ValidFrom: &june, ValidTo: &lastJune}
	_ = service.Save(&aprilWaterPlan)
	_ = service.Save(&mayWaterPlan)
	_ = service.Save(&mayPowerPlan)
	_ = service.Save(&juneWaterPlan)

	saved, _ := service.QueryForSeries(4)
	require.Equal(t, 3, len(saved), "saving the plans did not succeed")
	saved, _ = service.QueryForSeries(5)
	require.Equal(t, 1, len(saved), "saving the plans did not succeed")

	got, err := service.QueryForTime(4, readingTime)
	require.NoError(t, err, "no error expected")
	comparePrincingPlans(t, got, &mayWaterPlan)
}
