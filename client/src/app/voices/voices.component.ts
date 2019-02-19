import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { PollyVoice } from '@app/shared/polly/polly-voice';
import { PollyLanguage } from '@app/shared/polly/polly-language';
import { PollyService } from '@app/shared/polly/polly.service';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';


@Component({
  selector: 'app-voices',
  templateUrl: './Voices.component.html',
  styleUrls: ['./Voices.component.scss']
})
export class VoicesComponent implements OnInit, OnDestroy {

  voicesFormGroup: FormGroup;

  selectedVoice: string;
  selectedLanguage: PollyLanguage;
  languages: PollyLanguage[] = Array<PollyLanguage>();

  pollyVoices: PollyVoice[] = Array<PollyVoice>();

  encodingVoice: string;
  encodingVoiceSubscription: Subscription;

  constructor(private pollyservice: PollyService, private formBuilder: FormBuilder) {
    this.encodingVoiceSubscription = pollyservice.encodingVoiceUpdate$.subscribe(encodingVoice => {
      this.encodingVoice = encodingVoice;
    });
  }

  ngOnInit() {

    this.voicesFormGroup = this.formBuilder.group({
      textPreset: [''],
      textInput: ['', Validators.required]
    });

    this.selectedVoice = 'Emma';
    this.pollyservice.updateVoice('Emma');

    this.languages.push(new PollyLanguage('en-GB', 'English'));
    this.languages.push(new PollyLanguage('en-US', 'American'));

    this.getVoices(new PollyLanguage('en-GB', 'English British'));
  }

  ngOnDestroy() {
    this.encodingVoiceSubscription.unsubscribe();
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
}
