import { Component, OnInit, OnDestroy } from '@angular/core';
import { VgAPI } from 'videogular2/core';
import { PollyService } from '@app/shared/polly/polly.service';
import { Subscription } from '../../../node_modules/rxjs';

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

  encodingText: string;
  encodingVoice: string;

  progressbarMode: string;
  progressbarValue: number;

  constructor(private pollyservice: PollyService) {
    this.encodingTextSubscription = pollyservice.encodingTextUpdate$.subscribe(encodingText => {
      this.encodingText = encodingText;
    });
    this.encodingVoiceSubscription = pollyservice.encodingVoiceUpdate$.subscribe(encodingVoice => {
      this.encodingVoice = encodingVoice;
    });
   }

  ngOnInit() {
    this.sources = new Array<Object>();
    this.sources.push({
      src: 'https://s3-us-west-2.amazonaws.com/uk.ac.ncl.atnu.sot.output/16f4e11e-658b-4c9a-a1e8-e24cc54b4fc8.mp3',
      type: 'audio/mpeg'
    });
    this.progressbarMode = 'determinate';
    this.progressbarValue = 100;
  }

  ngOnDestroy() {
    this.encodingTextSubscription.unsubscribe();
    this.encodingVoiceSubscription.unsubscribe();
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

    this.pollyservice.getEncoding(this.encodingText).subscribe(url => {
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
