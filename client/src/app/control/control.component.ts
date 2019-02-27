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
  private selectionBounds: [number, number];

  paintText: SafeHtml = '';
  currentTag: string;

  breakConfig: {
    breakSliderMode: string;
    breakSliderStep: number;
    breakSliderMin: number;
    breakSliderMax: number;
    breakValue: number;
  };

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
      if (this.paintText === '') {
        this.addTags();
      }
    });
  }

  ngOnInit() {

    this.currentTag = null;

    this.breakConfig = {
      breakSliderMode: 's',
      breakSliderStep: 1,
      breakSliderMin: 0,
      breakSliderMax: 10,
      breakValue: 1
    };

    this.breakSliderMode = 's';
    this.breakSliderStep = 1;
    this.breakSliderMin = 0;
    this.breakSliderMax = 10;
    this.breakValue = 1;

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

  onBreakChange(event: any) {
    this.updateEncodingTag(null, this.currentTag);
  }

  // ---
  // PUBLIC METHODS.
  // ---

  // I render the rectangles emitted by the [textSelect] directive.
  public renderRectangles(event: TextSelectEvent): void {

    // If a new selection has been created, the viewport and host rectangles will
    // exist. Or, if a selection is being removed, the rectangles will be null.
    if (event.hostRectangle) {

      this.hostRectangle = event.hostRectangle;
      this.selectedText = event.text;
      this.selectionBounds = [event.start, event.end];

    } else {

      this.hostRectangle = null;
      this.selectedText = '';

    }

  }

  public cancel(): void {

    // Now that we've shared the text, let's clear the current selection.
    document.getSelection().removeAllRanges();
    // CAUTION: In modern browsers, the above call triggers a 'selectionchange'
    // event, which implicitly calls our renderRectangles() callback. However,
    // in IE, the above call doesn't appear to trigger the 'selectionchange'
    // event. As such, we need to remove the host rectangle explicitly.
    this.hostRectangle = null;
    this.selectedText = '';
    this.selectionBounds = [null, null];

  }

  /**
   * Create a new local record of the user selected tag
   * Ensure that this tag does not overlap or override any prexisting tags
   */
  addTag() {
    if (this.selectionBounds && this.selectedText) {

      let selection: PollySelection;

      const start = this.selectionBounds[0],
          end = this.selectionBounds[1],
          range = (end - start).toString();

      selection = new PollySelection(start, end, range);
      selection.tag = this.encodingTag;
      selection.ssml = this.encodingTag.wrap(this.selectedText);
      selection.css = this.encodingTag.paint(this.selectedText);
      selection.litter = this.encodingTag.litter();
      selection.csslitter = this.encodingTag.csslitter();

      let error = false;
      this.selections.forEach(idx => {
        if (selection.overrides(idx) || selection.overlaps(idx)) {
          console.log('Error');
          console.log(idx);
          error = true;
          this.lastSelection = selection;
          return;
        }
      });
      if (!error) {
        this.selections.push(selection);
        this.lastSelection = selection;
        this.addTags();
      }

      // Now that we've shared the text, let's clear the current selection.
      document.getSelection().removeAllRanges();
      // CAUTION: In modern browsers, the above call triggers a 'selectionchange'
      // event, which implicitly calls our renderRectangles() callback. However,
      // in IE, the above call doesn't appear to trigger the 'selectionchange'
      // event. As such, we need to remove the host rectangle explicitly.
      this.hostRectangle = null;
      this.selectedText = '';
      this.selectionBounds = [null, null];
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
        return 1;
      }
      if (ls.caretStart < rs.caretStart) {
        return -1;
      }
      return 0;
    });

    let text = String(this.encodingText).replace(/<[^>]+>/gm, '');
    text = '<speak>' + text + '</speak>';
    let offset = 7;

    this.selections.forEach(selection => {

      text = [
        text.slice(0, (selection.caretStart + offset)),
        selection.ssml,
        text.slice(selection.caretEnd + offset)]
        .join('');
      offset = offset + selection.litter;
    });

    this.paintText = this.sanitizer.bypassSecurityTrustHtml(text);

    this.pollyservice.updateText(text);
    localStorage.setItem('encodingText', this.encodingText);

    return this.encodingText;
  }

  /**
   * Update the encoding tag to be used
   * @param any event of button click
   * @param string name of ssml tag
   */
  updateEncodingTag(event: any, name: string) {
    this.currentTag = name;
    console.log(this.breakConfig);
    if (name === 'break') {
      const pre = '<break time="' + this.breakConfig.breakValue + this.breakConfig.breakSliderMode + '">';
      const post = '</break>';
      this.encodingTag = new PollyTag(name, 'break', pre, post);
    }
    if (name === 'emphasis') {
      const pre = '<emphasis level="' + this.emphasisMode + '">';
      const post = '</emphasis>';
      this.encodingTag = new PollyTag(name, 'emphasis', pre, post);
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
      this.encodingTag = new PollyTag(name, 'prosody', pre, post);
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
