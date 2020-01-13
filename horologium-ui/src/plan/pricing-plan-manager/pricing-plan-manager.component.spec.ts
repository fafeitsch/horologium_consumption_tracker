import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PricingPlanManagerComponent } from './pricing-plan-manager.component';

describe('PricingPlanManagerComponent', () => {
  let component: PricingPlanManagerComponent;
  let fixture: ComponentFixture<PricingPlanManagerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PricingPlanManagerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PricingPlanManagerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
