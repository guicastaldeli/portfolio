import { GetTimeStream } from "./get-time-stream.js";
import { GetClientCount } from "./get-client-count.js";
import { ProjectEditor } from "../project-editor.js";
import window from "../window.js";

export class Main {
    private timeStream: GetTimeStream;
    private clientCount: GetClientCount;
    private projectEditor: ProjectEditor;

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
    public protocol(url: string): string {
        const wsUrl = `${window.vars.SERVER_URL}/${url}`;
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