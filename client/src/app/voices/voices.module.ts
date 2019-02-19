import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { VoicesComponent } from './voices.component';
import { MaterialModule } from '@app/material/material.module';

@NgModule({
  imports: [
    CommonModule,
    MaterialModule
  ],
  declarations: [VoicesComponent],
  exports: [
    VoicesComponent
  ]
})
export class VoicesModule { }
