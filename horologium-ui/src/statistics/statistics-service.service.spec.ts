import { TestBed } from '@angular/core/testing';

import { StatisticsService } from './statistics.service';

describe('StatisticsServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: StatisticsService = TestBed.get(StatisticsService);
    expect(service).toBeTruthy();
  });
});
