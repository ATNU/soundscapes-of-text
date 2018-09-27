export class PollyTag {

    private _name: string;
    private _color: string;
    private _preTag: string;
    private _postTag: string;

    constructor(name: string, color: string, preTag: string, postTag: string) {
        this._name = name;
        this._color = color;
        this._preTag = preTag;
        this._postTag = postTag;
    }

    public get name(): string {
        return this._name;
    }

    public get color(): string {
        return this._color;
    }

    public get preTag(): string {
        return this._preTag;
    }

    public get postTag(): string {
        return this._postTag;
    }

    /**
     * Wrap the body of text in ssml tags
     * @param string body of text
     * @returns string wrapped body of text
     */
    public wrap(body: string): string {
        return this.preTag + body + this._postTag;
    }

    /**
     * Paints the body of text in a span with inline css
     * @param string body of text
     * @returns string painted body of text
     */
    public paint(body: string): string {
        return '<span style="background-color:#' + this.color + ';">' + body + '</span>';
    }

    /**
     * The length of the surrounding tags
     * @returns number of characters tag consists of
     */
    public litter(): number {
        return this._preTag.length + this._postTag.length;
    }

    public csslitter(): number {
        const len = '<span style="background-color:#' + this.color + ';">' + '</span>';
        return len.length;
    }
}
