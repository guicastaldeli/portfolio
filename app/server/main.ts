import { GetTimeStream } from "./get-time-stream";
import { GetClientCount } from "./get-client-count";

export class Main {
    private timeStream: GetTimeStream;
    private clientCount: GetClientCount;

    constructor() {
        this.timeStream = new GetTimeStream(this);
        this.clientCount = new GetClientCount(this);

        this.connect();
        this.cleanup();
    }

    /**
     * 
     * Protocol
     * 
     */
    public protocol(url: string) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/${url}`;
        return wsUrl;
    }

    /**
     * 
     * Connect
     * 
     */
    private connect(): void {
        this.timeStream.connect();
        this.clientCount.connect();
    }

    /**
     * Cleanup
     */
    private cleanup() {
        window.addEventListener('beforeunload', () => {
            if(this.timeStream.ws) this.timeStream.ws.close();
            if(this.clientCount.ws) this.clientCount.ws.close();
        });
    }
}