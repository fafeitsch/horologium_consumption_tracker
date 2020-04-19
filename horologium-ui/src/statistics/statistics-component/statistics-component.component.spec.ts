import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StatisticsComponentComponent } from './statistics-component.component';

describe('StatisticsComponentComponent', () => {
  let component: StatisticsComponentComponent;
  let fixture: ComponentFixture<StatisticsComponentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StatisticsComponentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatisticsComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
