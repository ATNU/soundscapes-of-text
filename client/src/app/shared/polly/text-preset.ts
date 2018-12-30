export class TextPreset {

    private _name: string;
    private _text: string;


    constructor(name: string, text: string) {
        this._name = name;
        this._text = text;
    }

    public get name(): string {
        return this._name;
    }
    public set name(value: string) {
        this._name = value;
    }

    public get text(): string {
        return this._text;
    }
    public set text(value: string) {
        this._text = value;
    }
}
