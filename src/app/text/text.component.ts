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
      this.encodingText = encodingText;
    });
    this.encodingTagSubscription = pollyservice.encodingTagUpdate$.subscribe(encodingTag => {
      this.encodingTag = encodingTag;
    });
  }

  ngOnInit() {
    this.textPresets.push(new TextPreset('Thomas Nashe', 'This is a Thomas Nashe Text'));
    this.textPresets.push(new TextPreset('Thomas Nashe 2', 'This is a Thomas Nashe Text For Part 2'));
    this.textPresets.push(new TextPreset('Thomas Nashe 3', 'This is a Thomas Nashe Text For the Final Part'));

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
    this.addTag(oField);
    this.pollyservice.updateText(this.encodingText);
  }

  onChange(deviceValue: TextPreset) {
    this.selectedTextPreset = deviceValue;
    this.encodingText = this.selectedTextPreset.text;
  }

  /**
   * Create a new local record of the user selected tag
   * Ensure that this tag does not overlap or override any prexisting tags
   * @param oField
   */
  addTag(oField: any) {
    if (window.getSelection().toString() !== '' || oField.selectionStart === '0') {
      // Because caret start and caret end are undefined
      const selection = new PollySelection(oField.selectionStart, oField.selectionEnd
        , window.getSelection().toString());

      console.log('Ofield ' + oField);
      selection.ssml = this.encodingTag.wrap(window.getSelection().toString());
      selection.css = this.encodingTag.paint(window.getSelection().toString());
      selection.litter = this.encodingTag.litter();

      this.selections.forEach(idx => {
        if (selection.overrides(idx) || selection.overlaps(idx)) {
          console.log('Error');
          return;
        }
      });
      this.selections.push(selection);
      this.addTags();
      console.log(this.selections);
    }
  }

  /**
   * Add ssml tags to input text in preparation for sending
   * to AWS Polly
   * @returns string input text with ssml tags included
   */
  addTags(): string {
    let p = this.encodingText;
    let litter = 0;
    this.selections.forEach(selection => {
      // Within this loop:
      // - Build the encoding ssml text
      // - Build the safehtml css text to show colour
      p = p.substring(0, selection.caretStart + litter) +
        selection.css + p.substring(selection.caretEnd + litter);

      litter = litter + selection.litter;
    });
    this.paintText = this.sanitizer.bypassSecurityTrustHtml(p);

    localStorage.setItem('encodingText', this.encodingText);

    return p;
  }

  resetPaint() {
    this.selections = [];
    this.paintText = this.encodingText;
  }
}
