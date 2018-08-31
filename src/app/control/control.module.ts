import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlComponent } from './control.component';
import { MaterialModule } from '@app/material/material.module';

@NgModule({
  imports: [
    CommonModule,
    MaterialModule
  ],
  declarations: [ControlComponent],
  exports: [
    ControlComponent
  ]
})
export class ControlModule { }
