import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EncodingComponent } from '@app/encoding/encoding.component';
import { VideogularModule } from '@app/videogular/videogular.module';
import { MaterialModule } from '@app/material/material.module';

@NgModule({
  imports: [
    CommonModule,
    VideogularModule,
    MaterialModule
  ],
  declarations: [EncodingComponent],
  exports: [
    EncodingComponent
  ]
})
export class EncodingModule { }
