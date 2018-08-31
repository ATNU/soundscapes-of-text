import { Component, OnInit } from '@angular/core';

import { environment } from '@env/environment';

export interface Aim {
  name: string;
  number: string;
}

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit {

  version: string = environment.version;
  aims: Aim[] = [
    {
      name: 'To explore ways of highlighting performance cues in a text',
      number: 'looks_one',
    },
    {
      name: 'To consider whether revisions are made according to sound qualities',
      number: 'looks_two',
    },
    {
      name: 'To enable readers to explore the relationship between sound and meaning',
      number: 'looks_3',
    },
    {
      name: 'To make clear the ways in which reader’ oral choices determine meaning',
      number: 'looks_4',
    }
  ];

  constructor() { }

  ngOnInit() { }

}
