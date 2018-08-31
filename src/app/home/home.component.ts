import { Component, OnInit } from '@angular/core';
import { finalize } from 'rxjs/operators';

import { PollyService } from '@app/shared/polly/polly.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  providers: [ PollyService ]
})
export class HomeComponent implements OnInit {

  quote: string;
  isLoading: boolean;

  constructor(private pollyService: PollyService) { }

  ngOnInit() {
    this.isLoading = true;
  }

}
