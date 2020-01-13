import { TestBed } from '@angular/core/testing';

import { PricingPlanService } from './pricing-plan.service';

describe('PricingPlanServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: PricingPlanService = TestBed.get(PricingPlanService);
    expect(service).toBeTruthy();
  });
});
