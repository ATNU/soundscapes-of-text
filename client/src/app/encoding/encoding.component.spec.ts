import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EncodingComponent } from './encoding.component';
import { VideogularModule } from '@app/videogular/videogular.module';

describe('EncodingComponent', () => {
  let component: EncodingComponent;
  let fixture: ComponentFixture<EncodingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EncodingComponent, VideogularModule ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EncodingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
