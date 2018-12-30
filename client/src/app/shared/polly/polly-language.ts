/**
 * @description Class to represent AWS Polly Language types
 */
export class PollyLanguage {

    private _code: string;
    private _name: string;

    constructor(code: string, name: string) {
        this._code = code;
        this._name = name;
    }

    public get code() {
        return this._code;
    }

    public set code(code: string) {
        this.code = code;
    }

    public get name(): string {
        return this._name;
    }

    public set name(name: string) {
        this._name = name;
    }
}
