import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PricingPlanTableComponent } from './pricing-plan-table.component';

describe('PricingPlanTableComponent', () => {
  let component: PricingPlanTableComponent;
  let fixture: ComponentFixture<PricingPlanTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PricingPlanTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PricingPlanTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
