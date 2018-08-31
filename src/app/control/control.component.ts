import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from '../../../node_modules/rxjs';
import { PollyVoice } from '@app/shared/polly/polly-voice';
import { PollyLanguage } from '@app/shared/polly/polly-language';
import { PollyTag } from '@app/shared/polly/polly-tag';
import { PollyService } from '@app/shared/polly/polly.service';


@Component({
  selector: 'app-control',
  templateUrl: './control.component.html',
  styleUrls: ['./control.component.scss']
})
export class ControlComponent implements OnInit, OnDestroy {

  selectedVoice: string;


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

  selectedLanguage: PollyLanguage;
  languages: PollyLanguage[] = Array<PollyLanguage>();

  pollyVoices: PollyVoice[] = Array<PollyVoice>();

  encodingVoice: string;
  encodingVoiceSubscription: Subscription;
  encodingTag: PollyTag;
  encodingTagSubscription: Subscription;

  constructor(private pollyservice: PollyService) {
    this.encodingVoiceSubscription = pollyservice.encodingVoiceUpdate$.subscribe(encodingVoice => {
      this.encodingVoice = encodingVoice;
    });
    this.encodingTagSubscription = pollyservice.encodingTagUpdate$.subscribe(encodingTag => {
      this.encodingTag = encodingTag;
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

    this.languages.push(new PollyLanguage('en-GB', 'English'));
    this.languages.push(new PollyLanguage('en-US', 'American'));

    this.getVoices(new PollyLanguage('en-GB', 'English British'));
    this.encodingTag = new PollyTag('', '', '', '');
  }

  ngOnDestroy() {
    this.encodingVoiceSubscription.unsubscribe();
    this.encodingTagSubscription.unsubscribe();
  }

  /**
   * Update the encoding tag to be used
   * @param any event of button click
   * @param string name of ssml tag
   */
  updateEncodingTag(event: any, name: string) {
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
        pre = pre + ' rate="' + this.rateValue  + '"';
      }
      if (this.pitchDefaultCheck === false) {
        pre = pre + ' pitch="' + this.pitchValue  + '"';
      }
      pre = pre + '>';
      const post = '</prosody>';
      this.encodingTag = new PollyTag(name, 'c9e075', pre, post);
    }
    this.pollyservice.updateTag(this.encodingTag);
  }

  /**
   * Retreieve all AWS Polly voices available for selected language
   * @param PollyLanguage language selected by user
   */
  getVoices(language: PollyLanguage) {
    this.pollyVoices = [];
    this.pollyservice.getVoices(language).subscribe(voices =>
      voices.forEach(element => {
        this.pollyVoices.push(element);
      }));
  }

  /**
   * Return if the provided voice is currently the selected voice
   * Useful for CSS manipulation
   * @param PollyVoice voice to query
   * @returns string boolean if selected
   */
  isVoiceSelected(voice: PollyVoice): string {
    if (this.selectedVoice === voice.Id) {
      return 'selectedTrue';
    }
    return 'selectedFalse';
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
   * Update the selected voice
   * @param PollyVoice selected voice
   */
  updateVoice(voice: PollyVoice) {
    this.selectedVoice = voice.Id;
    this.pollyservice.updateVoice(voice.Id);
  }

  /**
   * Retreive a demo audio clip of the selected voice
   * @param PollyVoice voice selected by user
   */
  play(voice: PollyVoice) {
    this.pollyservice.getDemo(voice);
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
