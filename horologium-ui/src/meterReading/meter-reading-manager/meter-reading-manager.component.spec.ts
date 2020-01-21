import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MeterReadingManagerComponent } from './meter-reading-manager.component';

describe('MeterReadingManagerComponent', () => {
  let component: MeterReadingManagerComponent;
  let fixture: ComponentFixture<MeterReadingManagerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MeterReadingManagerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MeterReadingManagerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
