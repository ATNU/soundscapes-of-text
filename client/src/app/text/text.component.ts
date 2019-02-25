import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { PollyService } from '@app/shared/polly/polly.service';
import { PollySelection } from '@app/shared/polly/polly-selection';
import { TextPreset } from '@app/shared/polly/text-preset';

import { Subscription } from '../../../node_modules/rxjs';
import { PollyTag } from '@app/shared/polly/polly-tag';
import { DomSanitizer, SafeHtml } from '../../../node_modules/@angular/platform-browser';

@Component({
  selector: 'app-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class TextComponent implements OnInit, OnDestroy {

  textFormGroup: FormGroup;

  selectedTextPreset: TextPreset;
  textPresets: TextPreset[] = Array<TextPreset>();

  encodingText = '';
  encodingTextSubscription: Subscription;

  encodingTag: PollyTag;
  encodingTagSubscription: Subscription;

  constructor(private pollyservice: PollyService, private sanitizer: DomSanitizer, private formBuilder: FormBuilder) {
    this.encodingTextSubscription = pollyservice.encodingTextUpdate$.subscribe(encodingText => {
      // this.encodingText = encodingText;
    });
    this.encodingTagSubscription = pollyservice.encodingTagUpdate$.subscribe(encodingTag => {
      this.encodingTag = encodingTag;
    });
  }

  ngOnInit() {
    this.textPresets.push(new TextPreset('Test', 'abc def ghi jkl mno pqr stu vwx yz'));
    this.textPresets.push(new TextPreset('Thomas Nashe', 'Brightness falls from the air, Queens have died young and fair, Dust hath closed Helen\'s eye. I am sick, I must die: Lord, have mercy on us.'));
    this.textPresets.push(new TextPreset('Thomas Nashe 2', 'towards Venice we progrest, and tooke Roterdam in our waie, that was cleane out of our waie: there we met with aged learnings chiefe ornament, that abundant and superingenious clarke, Erasmus, as also with merrie Sir Thomas Moore, our Countriman, who was come purposelie ouer a little before vs, to visite the said graue father Erasmus: what talke, what conference wee had then, it were here superfluous to rehearse, but this I can assure you, Erasmus in all his speeches seemed so much to mislike the indiscretion of Princes in preferring of parasites and fooles, that he decreed with himselfe to swim with the stream, and write a booke forthwith in commendation of follie.'));
    this.textPresets.push(new TextPreset('Thomas Nashe 3', 'What drugs, what sorceries, what oils, what waters, what ointments do our curious dames use to enlarge their withered beauties? Their lips are as lavishly red as if they used to kiss an ochreman every morning, and their cheeks sugar-candied and cherry-blushed so sweetly, after the colour of a new Lord Mayor\'s posts, as if the pageant of their wedlock holiday were hard at the door; so that if a painter were to draw any of their counterfeits on table, he needs no more but wet his pencil and dab it on their cheeks and he shall have vermilion and white enough to furnish out his work, though he leave his tar-box at home behind him. . . .'));

    this.encodingText = localStorage.getItem('encodingText');
    this.encodingTag = new PollyTag('', '', '', '');
    /*
    if (JSON.parse(localStorage.getItem('selections')) != null) {
      this.selections = JSON.parse(localStorage.getItem('selections'));
    }
    */

   this.textFormGroup = this.formBuilder.group({
    textPreset: [''],
    textInput: ['', Validators.required]
  });

  }

  ngOnDestroy() {
    this.encodingTextSubscription.unsubscribe();
    this.encodingTagSubscription.unsubscribe();
  }

  updateText(oField: any) {
    this.pollyservice.updateText(this.encodingText);
  }

  onChange(deviceValue: TextPreset) {
    this.selectedTextPreset = deviceValue;
    if (this.selectedTextPreset) {
      this.encodingText = this.selectedTextPreset.text;
      this.pollyservice.updateText(this.encodingText);
    }
  }
}
