import { GetTimeStream } from "./get-time-stream.js";
import { GetClientCount } from "./get-client-count.js";
import { ProjectEditor } from "../project-editor.js";
import { Env, GetEnv } from "../env.js";
export class Main {
    timeStream;
    clientCount;
    projectEditor;
    constructor() {
        this.timeStream = new GetTimeStream(this);
        this.clientCount = new GetClientCount(this);
        this.projectEditor = new ProjectEditor(this);
        this.start();
    }
    /**
     *
     * Start
     *
     */
    async start() {
        await Env.Ready();
        await this.connect();
        this.cleanup();
    }
    /**
     *
     * Protocol
     *
     */
    protocol(url) {
        const base = GetEnv('SERVER_URL') ?? '';
        const val = `${base}/${url}`;
        return val;
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
