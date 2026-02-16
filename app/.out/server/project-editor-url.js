/**
 *
 *
 * Get Project Editor URL...
 *
 *
 */
export class GetProjectEditorUrl {
    main;
    constructor(main) {
        this.main = main;
    }
    getUrl() {
        const url = "/editor";
        this.main.protocol(url);
    }
}
