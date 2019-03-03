import { Component, OnInit, OnDestroy } from '@angular/core';
import { VgAPI } from 'videogular2/core';
import { PollyService } from '@app/shared/polly/polly.service';
import { Subscription, Observable, interval } from 'rxjs';
import { PollySelection } from '@app/shared/polly/polly-selection';
import { HttpClient } from '@angular/common/http';

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

  generatingEncoding: boolean;

  constructor(private pollyservice: PollyService, private http: HttpClient) {
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
    this.encodingSelections = new Array<PollySelection>();
    this.generatingEncoding = true;
  }

  ngOnDestroy() {
    this.encodingTextSubscription.unsubscribe();
    this.encodingVoiceSubscription.unsubscribe();
    this.encodingSelectionsSubscription.unsubscribe();
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

      let tag: string;

      if (selection.tag.name === 'break') {
        tag = selection.ssml.replace('>', '/>');
        tag = tag.replace('</break>', '');
      } else {
        tag = selection.ssml;
      }

      text = [
        text.slice(0, (selection.caretStart + offset)),
        tag,
        text.slice(selection.caretEnd + offset)]
        .join('');
      offset = offset + selection.litter;
    });

    this.pollyservice.getEncoding(text).subscribe(url => {

      const source = interval(5000);

      const subscribe = source.subscribe(val => {
        this.http.get(url.toString()).subscribe(
          data => {
            console.log('success', data);
            this.generatingEncoding = false;
            subscribe.unsubscribe();
          },
          response => {
            if (response.status === 200) {
              this.api.pause();
              this.sources = [];
              this.sources.push({
                src: url,
                type: 'audio/mpeg'
              });
              this.api.currentTime = 0;
              this.generatingEncoding = false;
              subscribe.unsubscribe();
            }
          }
        );
      });
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
