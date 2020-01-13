import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PricingPlanEditorComponent } from './pricing-plan-editor.component';

describe('PricingPlanEditorComponent', () => {
  let component: PricingPlanEditorComponent;
  let fixture: ComponentFixture<PricingPlanEditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PricingPlanEditorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PricingPlanEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
