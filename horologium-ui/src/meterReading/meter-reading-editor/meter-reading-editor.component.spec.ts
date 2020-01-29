import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MeterReadingEditorComponent } from './meter-reading-editor.component';

describe('MeterReadingEditorComponent', () => {
  let component: MeterReadingEditorComponent;
  let fixture: ComponentFixture<MeterReadingEditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MeterReadingEditorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MeterReadingEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
