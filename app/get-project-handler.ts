import { WebSocketMessage } from "./types.js";
import window from "./window.js";

enum MessageType {
    PROJECT_CREATED = 'project_created',
    PROJECT_UPDATED = 'project_updated',
    PROJECT_DELETED = 'project_deleted'
}

export class GetProjectHandler {
    private ws: WebSocket | null = null;
    private onProjectCreated?: (data: any) => void;
    private onProjectUpdated?: (data: any) => void;
    private onProjectDeleted?: (data: any) => void;
    
    /**
     * 
     * Connect
     * 
     */
    public async connect(): Promise<void> {
        return new Promise((res, rej) => {
            const wsUrl = window.vars.SERVER_WS;

            this.ws = new WebSocket(wsUrl);
            this.ws.onopen = () => {
                this.subscribe('projects');
                res();
            }
            this.ws.onerror = (err) => {
                rej(err);
            }
            this.ws.onmessage = (e) => {
                try {
                    const message: WebSocketMessage = JSON.parse(e.data);
                    this.handleMessage(message);
                } catch(err) {
                    console.error('Failed to parse message', err);
                }
            }
            this.ws.onclose = () => {
                console.log('WS disconnected');
                setTimeout(() => this.connect(), 3000);
            }
        });
    }
    
    /**
     * 
     * Subscribe
     * 
     */
    private subscribe(channel: string): void {
        if(this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'subscribe',
                channel: channel
            }));
        }
    }

    /**
     * Handle Message
     */
    private handleMessage(message: WebSocketMessage): void {
        switch(message.type) {
            case MessageType.PROJECT_CREATED:
                if(this.onProjectCreated) this.onProjectCreated(message.data);
                break;
            case MessageType.PROJECT_UPDATED:
                if(this.onProjectUpdated) this.onProjectUpdated(message.data);
                break;
            case MessageType.PROJECT_DELETED:
                if(this.onProjectDeleted) this.onProjectDeleted(message.data);
                break;
        }
    }

    public setOnProjectCreated(cb: (data: any) => void): void {
        this.onProjectCreated = cb;
    }

    public setOnProjectUpdated(cb: (data: any) => void): void {
        this.onProjectUpdated = cb;
    }

    public setOnProjectDeleted(cb: (data: any) => void): void {
        this.onProjectDeleted = cb;
    }

    /**
     * 
     * Disconnect
     * 
     */
    public disconnect(): void {
        if(this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}