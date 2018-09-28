import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
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

  selectedTextPreset: TextPreset;
  textPresets: TextPreset[] = Array<TextPreset>();

  step = 0;
  paintText: SafeHtml = 'Nothing here yet...';

  encodingText = '';
  encodingTextSubscription: Subscription;

  encodingTag: PollyTag;
  encodingTagSubscription: Subscription;

  selections = Array<PollySelection>();

  constructor(private pollyservice: PollyService, private sanitizer: DomSanitizer) {
    this.encodingTextSubscription = pollyservice.encodingTextUpdate$.subscribe(encodingText => {
      // this.encodingText = encodingText;
    });
    this.encodingTagSubscription = pollyservice.encodingTagUpdate$.subscribe(encodingTag => {
      this.encodingTag = encodingTag;
    });
  }

  ngOnInit() {
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

  }

  ngOnDestroy() {
    this.encodingTextSubscription.unsubscribe();
    this.encodingTagSubscription.unsubscribe();
  }

  nextStep() {
    this.step++;
  }

  prevStep() {
    this.step--;
  }

  updateText(oField: any) {
    this.addTags();
    // this.pollyservice.updateText(this.encodingText); - will remove ssml?
    console.log(this.encodingText);
  }

  onChange(deviceValue: TextPreset) {
    this.selectedTextPreset = deviceValue;
    this.encodingText = this.selectedTextPreset.text;
    this.addTags();
  }

  /**
   * Create a new local record of the user selected tag
   * Ensure that this tag does not overlap or override any prexisting tags
   * @param oField
   */
  addTag(oField: any) {
    if (window.getSelection().toString() !== '' || oField.selectionStart === '0') {

      let selection: PollySelection;
      if (oField.selectionStart === undefined && oField.selectionEnd === undefined ) {
        let dif = 0;
        if (this.selections.length > 0) {
          dif = this.selections[this.selections.length - 1].caretEnd;
        }

        // dif needs to figure out if the selection is before
        // else minus it

        selection = new PollySelection((window.getSelection().anchorOffset + dif),
        window.getSelection().focusOffset + dif
          , window.getSelection().toString());
      } else {
        selection = new PollySelection(oField.selectionStart, oField.selectionEnd
          , window.getSelection().toString());
      }

      if (selection.caretEnd < selection.caretStart) {
        const temp = selection.caretStart;
        selection.caretStart = selection.caretEnd;
        selection.caretEnd = temp;
      }
      console.log(selection.caretStart, selection.caretEnd);
      selection.ssml = this.encodingTag.wrap(window.getSelection().toString());
      selection.css = this.encodingTag.paint(window.getSelection().toString());
      selection.litter = this.encodingTag.litter();
      selection.csslitter = this.encodingTag.csslitter();

      let error = false;
      this.selections.forEach(idx => {
        console.log(idx.caretStart, idx.caretEnd, selection.caretStart, selection.caretEnd);
        if (selection.overrides(idx) || selection.overlaps(idx)) {
          console.log('Error');
          error = true;
          return;
        }
      });
      if (!error) {
        this.selections.push(selection);
        this.addTags();
        console.log(this.selections);
      }
    }
  }

  /**
   * Add ssml tags to input text in preparation for sending
   * to AWS Polly
   * @returns string input text with ssml tags included
   */
  addTags(): string {
    let p = this.encodingText;
    let css = this.encodingText;

    let litter = 0;
    let csslitter = 0;
    this.selections.forEach(selection => {
      // Within this loop:
      // - Build the encoding ssml text
      // - Build the safehtml css text to show colour
      p = p.substring(0, selection.caretStart + litter) +
        selection.ssml + p.substring(selection.caretEnd + litter);

      litter = litter + selection.litter;

      css = css.substring(0, selection.caretStart + csslitter) +
      selection.css + css.substring(selection.caretEnd + csslitter);

      csslitter = csslitter + selection.csslitter;

    });
    this.paintText = this.sanitizer.bypassSecurityTrustHtml(css);

    localStorage.setItem('encodingText', this.encodingText);

    this.pollyservice.updateText('<speak>' + p + '</speak>');
    return p;
  }

  resetPaint() {
    this.selections = [];
    this.paintText = this.encodingText;
  }
}
