import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlComponent, ClearTagsDialogComponent } from './control.component';
import { MaterialModule } from '@app/material/material.module';
import { TextSelectDirective } from './text-select.directive';

@NgModule({
  imports: [
    CommonModule,
    MaterialModule
  ],
  declarations: [
    ControlComponent,
    ClearTagsDialogComponent,
    TextSelectDirective
   ],
  exports: [
    ControlComponent,
    TextSelectDirective
  ]
})
export class ControlModule { }
