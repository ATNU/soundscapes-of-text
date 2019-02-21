import { Component, OnInit, OnDestroy } from '@angular/core';
import { TextSelectEvent } from './text-select.directive';
import { Subscription } from '../../../node_modules/rxjs';
import { PollyTag } from '@app/shared/polly/polly-tag';
import { PollyService } from '@app/shared/polly/polly.service';
import { PollySelection } from '@app/shared/polly/polly-selection';
import { SafeHtml, DomSanitizer } from '@angular/platform-browser';

interface SelectionRectangle {
  left: number;
  top: number;
  width: number;
  height: number;
}

@Component({
  selector: 'app-control',
  templateUrl: './control.component.html',
  styleUrls: ['./control.component.scss']
})
export class ControlComponent implements OnInit, OnDestroy {

  public hostRectangle: SelectionRectangle | null;
  private selectedText: string;

  paintText: SafeHtml = '';

  breakSliderMode: string;
  breakSliderStep: number;
  breakSliderMin: number;
  breakSliderMax: number;
  breakValue: number;

  emphasisMode: string;

  volumeDefaultCheck: boolean;
  rateDefaultCheck: boolean;
  pitchDefaultCheck: boolean;
  volumeValue: number;
  rateValue: number;
  pitchValue: number;

  lastSelection: PollySelection;
  selections = Array<PollySelection>();

  encodingTag: PollyTag;
  encodingTagSubscription: Subscription;
  encodingText: string;
  encodingTextSubscription: Subscription;

  constructor(private pollyservice: PollyService, private sanitizer: DomSanitizer) {
    this.hostRectangle = null;
    this.selectedText = '';
    this.encodingTagSubscription = pollyservice.encodingTagUpdate$.subscribe(encodingTag => {
      this.encodingTag = encodingTag;
    });
    this.encodingTextSubscription = pollyservice.encodingTextUpdate$.subscribe(encodingText => {
      this.encodingText = encodingText;
    });
  }

  ngOnInit() {

    this.breakSliderMode = 's';
    this.breakSliderStep = 1;
    this.breakSliderMin = 0;
    this.breakSliderMax = 10;
    this.breakValue = 0;

    this.emphasisMode = 'moderate';

    this.volumeDefaultCheck = true;
    this.rateDefaultCheck = true;
    this.pitchDefaultCheck = true;
    this.volumeValue = 0;
    this.rateValue = 100;
    this.pitchValue = 0;

    this.encodingTag = new PollyTag('', '', '', '');
  }

  ngOnDestroy() {
    this.encodingTagSubscription.unsubscribe();
  }

  // ---
  // PUBLIC METHODS.
  // ---

  // I render the rectangles emitted by the [textSelect] directive.
  public renderRectangles(event: TextSelectEvent): void {

    console.group('Text Select Event');
    console.log('Text:', event.text);
    console.log('Viewport Rectangle:', event.viewportRectangle);
    console.log('Host Rectangle:', event.hostRectangle);
    console.groupEnd();

    // If a new selection has been created, the viewport and host rectangles will
    // exist. Or, if a selection is being removed, the rectangles will be null.
    if (event.hostRectangle) {

      this.hostRectangle = event.hostRectangle;
      this.selectedText = event.text;

    } else {

      this.hostRectangle = null;
      this.selectedText = '';

    }

  }


  // I share the selected text with friends :)
  public shareSelection(): void {

    console.group('Shared Text');
    console.log(this.selectedText);
    console.groupEnd();

    // Now that we've shared the text, let's clear the current selection.
    document.getSelection().removeAllRanges();
    // CAUTION: In modern browsers, the above call triggers a 'selectionchange'
    // event, which implicitly calls our renderRectangles() callback. However,
    // in IE, the above call doesn't appear to trigger the 'selectionchange'
    // event. As such, we need to remove the host rectangle explicitly.
    this.hostRectangle = null;
    this.selectedText = '';

  }

  /**
   * Create a new local record of the user selected tag
   * Ensure that this tag does not overlap or override any prexisting tags
   * @param oField
   */
  addTag(oField: any) {
    if (window.getSelection().toString() !== '' || oField.selectionStart === '0') {

      let selection: PollySelection;
      if (oField.selectionStart === undefined && oField.selectionEnd === undefined) {
        let dif = 0;
        if (this.selections.length > 0) {
          dif = this.lastSelection.caretEnd;
        }
        selection = new PollySelection((window.getSelection().anchorOffset + dif),
          window.getSelection().focusOffset + dif, window.getSelection().toString());

        const posSelect = this.encodingText.substring(selection.caretStart, selection.caretEnd);

        if (selection.range !== posSelect && this.selections.length > 0) {
          console.log('In');
          selection.caretStart = selection.caretStart - this.lastSelection.caretEnd;
          selection.caretEnd = selection.caretEnd - this.lastSelection.caretEnd;
          console.log('last', this.lastSelection.caretStart, this.lastSelection.caretEnd);
          console.log(selection.caretStart, selection.caretEnd);
          console.log('end');
        }
      } else {
        selection = new PollySelection(oField.selectionStart, oField.selectionEnd
          , window.getSelection().toString());
      }

      if (selection.caretEnd < selection.caretStart) {
        const temp = selection.caretStart;
        selection.caretStart = selection.caretEnd;
        selection.caretEnd = temp;
      }
      console.log('Final values', selection.caretStart, selection.caretEnd);
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
          this.lastSelection = selection;
          return;
        }
      });
      if (!error) {
        this.selections.push(selection);
        this.lastSelection = selection;
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
    this.selections.sort((ls, rs): number => {
      if (ls.caretStart > rs.caretStart) {
        console.log(ls.caretStart, 'selectedLarger', rs.caretStart);
        return 1;
      }
      if (ls.caretStart < rs.caretStart) {
        console.log(ls.caretStart, 'selectedLess', rs.caretStart);
        return -1;
      }
      console.log(this.selections);
      return 0;
    });

    let p = this.encodingText;
    let css = this.encodingText;

    let litter = 0;
    let csslitter = 0;
    this.selections.forEach(selection => {

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

  /**
   * Update the encoding tag to be used
   * @param any event of button click
   * @param string name of ssml tag
   */
  updateEncodingTag(event: any, name: string) {
    console.log(name);
    if (name === 'break') {
      const pre = '<break time="' + this.breakValue + this.breakSliderMode + '"/>';
      this.encodingTag = new PollyTag(name, 'e07575', pre, '');
    }
    if (name === 'emphasis') {
      const pre = '<emphasis level="' + this.emphasisMode + '">';
      const post = '</emphasis>';
      this.encodingTag = new PollyTag(name, '9175e0', pre, post);
    }
    if (name === 'prosody') {
      let pre = '<prosody';
      if (this.volumeDefaultCheck === false) {
        pre = pre + ' volume="' + this.volumeValue + '"';
      }
      if (this.rateDefaultCheck === false) {
        pre = pre + ' rate="' + this.rateValue + '"';
      }
      if (this.pitchDefaultCheck === false) {
        pre = pre + ' pitch="' + this.pitchValue + '"';
      }
      pre = pre + '>';
      const post = '</prosody>';
      this.encodingTag = new PollyTag(name, 'c9e075', pre, post);
    }
    this.pollyservice.updateTag(this.encodingTag);
  }

  /**
   * Returns if the provided brush is currently the selected brush
   * Useful for CSS manipulation
   * @param string brush id
   * @returns string boolean if selected
   */
  isBrushSelected(brush: string): string {
    if (this.encodingTag.name === brush) {
      return 'Selected';
    }
    return 'Select';
  }

  /**
   * Set the min and max value for break slider based upon which mode is selected (s/ms)
   * @param any event
   */
  breakModeSwitch(event: any) {
    if (event.value === 's') {
      this.breakSliderStep = 1;
      this.breakSliderMin = 0;
      this.breakSliderMax = 10;
    } else if (event.value === 'ms') {
      this.breakSliderStep = 100;
      this.breakSliderMin = 10;
      this.breakSliderMax = 10000;
    }
  }

}
