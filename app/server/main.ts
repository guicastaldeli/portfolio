import { GetTimeStream } from "./get-time-stream.js";
import { GetClientCount } from "./get-client-count.js";
import { ProjectEditor } from "../project-editor.js";
import { Env, GetEnv } from "../env.js";

export class Main {
    private timeStream: GetTimeStream;
    private clientCount: GetClientCount;
    private projectEditor: ProjectEditor;

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
    private async start(): Promise<void> {
        await Env.Ready();
        await this.connect();
        this.cleanup();
    }

    /**
     * 
     * Protocol
     * 
     */
    public protocol(url: string): string {
        const base = GetEnv('SERVER_URL') ?? '';
        const val = `${base}/${url}`;
        return val;
    }

    /**
     * 
     * Connect
     * 
     */
    private async connect(): Promise<void> {
        await this.timeStream.connect();
        await this.clientCount.connect();
        await this.projectEditor.init();
    }

    /**
     * Cleanup
     */
    private cleanup(): void {
        document.addEventListener('beforeunload', () => {
            if(this.timeStream.ws) this.timeStream.ws.close();
            if(this.clientCount.ws) this.clientCount.ws.close();
        });
    }
}