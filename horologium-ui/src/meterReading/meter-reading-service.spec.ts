import { TestBed } from '@angular/core/testing';

import { MeterReadingService } from './meter-reading.service';

describe('MeterReadingServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: MeterReadingService = TestBed.get(MeterReadingService);
    expect(service).toBeTruthy();
  });
});
