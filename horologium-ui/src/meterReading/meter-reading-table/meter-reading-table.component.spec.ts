import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MeterReadingTableComponent } from './meter-reading-table.component';

describe('MeterReadingTableComponent', () => {
  let component: MeterReadingTableComponent;
  let fixture: ComponentFixture<MeterReadingTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MeterReadingTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MeterReadingTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
