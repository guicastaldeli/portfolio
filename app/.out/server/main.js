import { GetTimeStream } from "./get-time-stream.js";
import { GetClientCount } from "./get-client-count.js";
import { ProjectEditor } from "../project-editor.js";
import window from "../window.js";
export class Main {
    timeStream;
    clientCount;
    projectEditor;
    constructor() {
        window.vars.APP_ENV = 'prod';
        window.init();
        this.timeStream = new GetTimeStream(this);
        this.clientCount = new GetClientCount(this);
        this.projectEditor = new ProjectEditor(this);
        this.connect();
        this.cleanup();
    }
    /**
     *
     * Protocol
     *
     */
    protocol(url) {
        const wsUrl = `${window.vars.SERVER_URL}/${url}`;
        return wsUrl;
    }
    /**
     *
     * Connect
     *
     */
    async connect() {
        await this.timeStream.connect();
        await this.clientCount.connect();
        await this.projectEditor.init();
    }
    /**
     * Cleanup
     */
    cleanup() {
        document.addEventListener('beforeunload', () => {
            if (this.timeStream.ws)
                this.timeStream.ws.close();
            if (this.clientCount.ws)
                this.clientCount.ws.close();
        });
    }
}
