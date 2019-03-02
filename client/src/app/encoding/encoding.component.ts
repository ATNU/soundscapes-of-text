import { Component, OnInit, OnDestroy } from '@angular/core';
import { VgAPI } from 'videogular2/core';
import { PollyService } from '@app/shared/polly/polly.service';
import { Subscription } from '../../../node_modules/rxjs';
import { PollySelection } from '@app/shared/polly/polly-selection';

@Component({
  selector: 'app-encoding',
  templateUrl: './encoding.component.html',
  styleUrls: ['./encoding.component.scss']
})
export class EncodingComponent implements OnInit, OnDestroy {

  api: VgAPI;
  sources: Array<Object>;

  encodingTextSubscription: Subscription;
  encodingVoiceSubscription: Subscription;
  encodingSelectionsSubscription: Subscription;

  encodingText: string;
  encodingVoice: string;
  encodingSelections: PollySelection[];

  progressbarMode: string;
  progressbarValue: number;
  generatingEncoding: boolean;

  constructor(private pollyservice: PollyService) {
    this.encodingTextSubscription = pollyservice.encodingTextUpdate$.subscribe(encodingText => {
      this.encodingText = encodingText;
    });
    this.encodingVoiceSubscription = pollyservice.encodingVoiceUpdate$.subscribe(encodingVoice => {
      this.encodingVoice = encodingVoice;
    });
    this.encodingSelectionsSubscription = pollyservice.encodingSelections$.subscribe(encodingSelections => {
      this.encodingSelections = encodingSelections;
    });
   }

  ngOnInit() {
    this.sources = new Array<Object>();
    this.sources.push({
      src: 'https://s3-eu-west-1.amazonaws.com/atnu.soundscapes/bfeeb015-d93f-4961-ae77-61dbdd700e0d.mp3',
      type: 'audio/mpeg'
    });
    this.progressbarMode = 'determinate';
    this.progressbarValue = 100;
    this.generatingEncoding = true;
  }

  ngOnDestroy() {
    this.encodingTextSubscription.unsubscribe();
    this.encodingVoiceSubscription.unsubscribe();
    this.encodingSelectionsSubscription.unsubscribe();
  }

  createEncoding() {
    console.log(this.encodingSelections);
  }

  /**
   * Subscribe to key events when audio player is loaded
   * @param VgAPI videogular plauyer
   */
  onPlayerReady(api: VgAPI) {
    this.api = api;
    this.api.getDefaultMedia().subscriptions.ended.subscribe(
      () => {
          this.api.getDefaultMedia().currentTime = 0;
      }
    );
  }

  /**
   * Update the audio src of media player to newly generated tts encoding
   * and set current playback time to  0
   * @param encoding source of media
   */
  updatePlayerSource() {
    this.progressbarMode = 'indeterminate';
    this.generatingEncoding = true;

    this.encodingSelections.sort((ls, rs): number => {
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

    this.encodingSelections.forEach(selection => {

      text = [
        text.slice(0, (selection.caretStart + offset)),
        selection.ssml,
        text.slice(selection.caretEnd + offset)]
        .join('');
      offset = offset + selection.litter;
    });

    // clean up invalid break tags
    text = text.replace(new RegExp('</break>', 'g'), '');


    this.pollyservice.getEncoding(text).subscribe(url => {
      console.log('URL', url);
      this.api.pause();
      this.sources = [];
      this.sources.push({
        src: url,
        type: 'audio/mpeg'
      });
      this.api.currentTime = 0;
      this.progressbarMode = 'determinate';
      this.progressbarValue = 100;
      this.generatingEncoding = false;
    });
  }


  /**
   * Download the current encoding and related SSML
   */
  downloadEncoding() {
    this.pollyservice.getDownload(this.encodingText).subscribe(url => {
      window.open(url as string);
    });
  }

}
