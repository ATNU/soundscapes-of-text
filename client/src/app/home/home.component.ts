import { Component, OnInit, ViewChild } from '@angular/core';
import { finalize } from 'rxjs/operators';

import { PollyService } from '@app/shared/polly/polly.service';
import { TextComponent } from '@app/text/text.component';
import { VoicesComponent } from '@app/voices/voices.component';
import { EncodingComponent } from '@app/encoding/encoding.component';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  providers: [ PollyService ]
})
export class HomeComponent implements OnInit {

  @ViewChild(TextComponent) textComponent: TextComponent;
  @ViewChild(VoicesComponent) voicesComponent: VoicesComponent;
  @ViewChild(EncodingComponent) encodingComponent: EncodingComponent;

  get frmStepOne() {
      return this.textComponent ? this.textComponent.textFormGroup : null;
  }

  get frmStepThree() {
    return this.voicesComponent ? this.voicesComponent.voicesFormGroup : null;
}

  constructor(private pollyService: PollyService) { }

  ngOnInit() { }

  public onStepChange(event: any): void {
    if (event.selectedIndex === 3) {
      console.log('Fetch Encoding');

      this.encodingComponent.updatePlayerSource();
    }
  }
}
