import { TestBed, inject } from '@angular/core/testing';

import { PollyService } from './polly.service';
import { HttpClient } from 'selenium-webdriver/http';

describe('PollyService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PollyService, HttpClient ]
    });
  });

  it('should be created', inject([PollyService], (service: PollyService) => {
    expect(service).toBeTruthy();
  }));
});
