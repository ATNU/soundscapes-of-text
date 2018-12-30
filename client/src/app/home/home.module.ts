import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';
import { FlexLayoutModule } from '@angular/flex-layout';

import { CoreModule } from '@app/core';
import { SharedModule } from '@app/shared';
import { MaterialModule } from '@app/material.module';
import { HomeRoutingModule } from './home-routing.module';
import { HomeComponent } from './home.component';

import { TextModule } from '@app/text/text.module';
import { ControlModule } from '@app/control/control.module';
import { EncodingModule } from '@app/encoding/encoding.module';
import { SnackbarComponent } from './snackbar/snackbar.component';
import { PagenotfoundComponent } from './pagenotfound/pagenotfound.component';


@NgModule({
  imports: [
    CommonModule,
    TranslateModule,
    CoreModule,
    SharedModule,
    FlexLayoutModule,
    MaterialModule,
    HomeRoutingModule,
    TextModule,
    ControlModule,
    EncodingModule
  ],
  declarations: [
    HomeComponent,
    SnackbarComponent,
    PagenotfoundComponent
  ],
  providers: [
  ],
  entryComponents: [
    SnackbarComponent
  ]
})
export class HomeModule { }
