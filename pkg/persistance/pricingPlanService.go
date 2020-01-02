package orm

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
)

type PricingPlanService interface {
	Create(new *domain.PricingPlan) error
	Delete(id int) error
	QueryById(id int) (*domain.PricingPlan, *error)
}

func NewPricingPlanService(db *gorm.DB) PricingPlanService {
	return &PricingPlanServiceImpl{db: db}
}

type PricingPlanServiceImpl struct {
	db *gorm.DB
}

func (p *PricingPlanServiceImpl) Create(new *domain.PricingPlan) error {
	return nil
}

func (p *PricingPlanServiceImpl) Delete(id int) error {
	return nil
}

func (p *PricingPlanServiceImpl) QueryById(id int) (*domain.PricingPlan, *error) {
	return nil, nil
}
