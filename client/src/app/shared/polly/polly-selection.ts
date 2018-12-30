export class PollySelection {

    private _caretStart: number;
    private _caretEnd: number;
    private _range: string;
    private _ssml: string;
    private _css: string;
    private _litter: number;
    private _csslitter: number;

    constructor(caretStart: number, caretEnd: number, range: string) {
        this._caretStart = caretStart;
        this._caretEnd = caretEnd;
        this._range = range;
    }

    /** Checks if this selection will overlap with the provided selection
     * @param arg0 PollySelection to compare
     * @returns boolean
     */
    public overlaps(arg0: PollySelection): boolean {
        return ((this.caretStart < arg0._caretStart) && (this.caretEnd > arg0.caretStart)
            && (this.caretEnd < arg0.caretEnd))
                || ((this.caretStart > arg0.caretStart) && (this.caretStart < arg0.caretEnd)
                    && (this.caretEnd > arg0.caretEnd));
    }

    /**
     * Checks if this selection will override the provided selection
     * @param arg0 PollySelection to compare
     * @returns boolean
     */
    public overrides(arg0: PollySelection): boolean {
        return ((this.caretStart < arg0.caretStart) && (this.caretEnd > arg0.caretEnd))
            || ((this.caretStart > arg0.caretStart) && (this.caretEnd < arg0.caretEnd));
    }

    public reinvoke(text: string): string {
        return '';
    }

    public get caretStart(): number {
        return this._caretStart;
    }

    public get caretEnd(): number {
        return this._caretEnd;
    }

    public set caretStart(start: number)  {
        this._caretStart = start;
    }

    public set caretEnd(end: number) {
        this._caretEnd = end;
    }

    public get range(): string {
        return this._range;
    }

    public get ssml(): string {
        return this._ssml;
    }

    public set ssml(ssml: string) {
        this._ssml = ssml;
    }

    public get css(): string {
        return this._css;
    }

    public set css(css: string) {
        this._css = css;
    }

    public get litter(): number {
        return this._litter;
    }

    public set litter(litter: number) {
        this._litter = litter;
    }

    public get csslitter(): number {
        return this._csslitter;
    }

    public set csslitter(csslitter: number) {
        this._csslitter = csslitter;
    }
}
