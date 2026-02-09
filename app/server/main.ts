import { GetTimeStream } from "./get-time-stream.js";
import { GetClientCount } from "./get-client-count.js";
import { ProjectEditor } from "../project-editor.js";

export class Main {
    private timeStream: GetTimeStream;
    private clientCount: GetClientCount;
    private projectEditor: ProjectEditor;

    constructor() {
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
    public protocol(url: string): string {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/${url}`;
        return wsUrl;
    }

    /**
     * 
     * Connect
     * 
     */
    private async connect(): Promise<void> {
        await this.timeStream.connect();
        await this.clientCount.connect();
    }

    /**
     * Cleanup
     */
    private cleanup(): void {
        window.addEventListener('beforeunload', () => {
            if(this.timeStream.ws) this.timeStream.ws.close();
            if(this.clientCount.ws) this.clientCount.ws.close();
        });
    }
}