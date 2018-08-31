import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TextComponent } from '@app/text/text.component';
import { MaterialModule } from '@app/material/material.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';


@NgModule({
  imports: [
    CommonModule,
    MaterialModule,
    FormsModule,
    ReactiveFormsModule
  ],
  declarations: [TextComponent],
  exports: [
    TextComponent
  ]
})
export class TextModule { }
