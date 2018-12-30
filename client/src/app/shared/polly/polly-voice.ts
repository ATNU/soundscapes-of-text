import { PollyLanguage } from '@app/shared/polly/polly-language';

/**
 * @description Class to represent AWS Polly Voice
 */
export class PollyVoice {

    private _Gender: string;
    private _Id: string;
    private _LanguageCode: string;
    private _LanguageName: string;
    private _Name: string;

    constructor(gender: string, id: string, languageCode: string, languageName: string, name: string) {
        this._Gender = gender;
        this._Id = id;
        this._LanguageCode = languageCode;
        this._LanguageName = languageName;
        this._Name = name;
    }

    public get Gender(): string {
        return this._Gender;
    }

    public set Gender(gender: string) {
        this._Gender = gender;
    }

    public get Id(): string {
        return this._Id;
    }

    public set Id(value: string) {
        this._Id = value;
    }

    public get LanguageCode(): string {
        return this._LanguageCode;
    }

    public set LanguageCode(language: string) {
        this._LanguageCode = language;
    }

    public get LanguageName(): string {
        return this._LanguageName;
    }

    public set LanguageName(language: string) {
        this._LanguageName = language;
    }

    public get Name(): string {
        return this._Name;
    }

    public set name(name: string) {
        this._Name = name;
    }

}
