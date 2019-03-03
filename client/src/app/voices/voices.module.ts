import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '@app/shared';
import { VoicesComponent } from './voices.component';
import { MaterialModule } from '@app/material/material.module';

@NgModule({
  imports: [
    CommonModule,
    MaterialModule,
    SharedModule
  ],
  declarations: [VoicesComponent],
  exports: [
    VoicesComponent
  ]
})
export class VoicesModule { }
