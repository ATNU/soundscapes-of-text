import { Injectable, OnInit } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';

import { SnackbarComponent } from '@app/home/snackbar/snackbar.component';
import { PollyVoice } from '@app/shared/polly/polly-voice';
import { PollyLanguage } from '@app/shared/polly/polly-language';
import { PollyTag } from '@app/shared/polly/polly-tag';
import { environment } from '@env/environment';
import { PollySelection } from './polly-selection';

@Injectable({
  providedIn: 'root'
})
export class PollyService implements OnInit {

  private voicesUrl = '/voices';
  private demoUrl = '/demo';
  private genUrl = '/generate';

  private encodingText = new Subject<string>();
  private encodingVoice = new Subject<string>();
  private encodingTag = new Subject<PollyTag>();
  private encodingSelections = new Subject<PollySelection[]>();

  private encodingVoiceValue: string;

  encodingTextUpdate$ = this.encodingText.asObservable();
  encodingVoiceUpdate$ = this.encodingVoice.asObservable();
  encodingTagUpdate$ = this.encodingTag.asObservable();
  encodingSelections$ = this.encodingSelections.asObservable();

  constructor(
    private http: HttpClient,
    public snackBar: MatSnackBar
  ) {}

  ngOnInit() {
    this.encodingText.next(localStorage.getItem('encodingText'));
  }

  /**
   * Provide observable with newly received encoding text
   * @param string text to encode
   */
  updateText(encodingText: string) {
    localStorage.setItem('encodingText', encodingText);
    this.encodingText.next(encodingText);
  }

  /**
   * Provide observable with newly received encoding voice
   * @param string voice to use for encoding
   */
  updateVoice(encodingVoice: string) {
    this.encodingVoice.next(encodingVoice);
    this.encodingVoiceValue = encodingVoice;
  }

  /**
   * Provide observable with newly received encoding tag
   * @param PollyTag ssml tag to be painted
   */
  updateTag(encodingTag: PollyTag) {
    this.encodingTag.next(encodingTag);
  }

  /**
   * Provide observable with newly received encoding tag
   * @param PollyTag ssml tag to be painted
   */
  updateSelections(selections: PollySelection[]) {
    this.encodingSelections.next(selections);
  }

  /**
   * Retrieve available languages from AWS Polly in specified language
   *
   * @param  PollyLanguage language to search for
   * @returns Observable<PollyVoice[]> Returned array of PollyVoice
   */
  getVoices(language: PollyLanguage): Observable<PollyVoice[]> {
    const url = `${this.voicesUrl}/${language.code}`;
    return this.http.get<PollyVoice[]>(url)
        .pipe(
            tap(polly => this.log(`Retrieved voices`)),
            catchError(this.handleError('getVoices', []))
        );
  }

  /**
   * Retrieve and play a demo (.mp3) from AWS Polly of specified voice
   * @param PollyVoice queried voice
   */
  getDemo(voice: PollyVoice) {
    const audio = new Audio();
    // This overrides any http prefix intercepting unfortunately
    audio.src = `${environment.serverUrl}${this.demoUrl}/${voice.Id}`;
    audio.load();
    audio.play();
  }

  /**
   * Start an asynchronous tts encoding task and pipe response URL
   * @param string text to encode
   */
  getEncoding(text: string): Observable<Object> {
    if (text === '') {
      // Show snackbar
      return;
    }
    return this.http.post(this.genUrl, text, { params:
      new HttpParams().set('voice', this.encodingVoiceValue), responseType: 'text' })
        .pipe(
            tap(polly => this.log(`fetched generation`)),
            catchError(this.handleError('getEncoding', []))
        );
  }

  getDownload(text: string): Observable<Object> {
    if (text === '') {
      // Show snackbar
      return;
    }
    return this.http.post(this.genUrl, text, { params:
      new HttpParams().set('voice', this.encodingVoiceValue), responseType: 'text' })
        .pipe(
            tap(polly => this.log(`fetched generation`)),
            catchError(this.handleError('getEncoding', []))
        );
  }

  /**
   * Simple logger to console
   * @param message
   */
  private log(message: string) {
    console.log(message);
  }

  /**
   * Simple error handler
   */
  private handleError<T>(operation: string = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      this.snackBar.openFromComponent(SnackbarComponent, {
        data: 'Unable to generate encoding, please try later',
        duration: 3000,
      });

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }


}
