var MessageType;
(function (MessageType) {
    MessageType["PROJECT_CREATED"] = "project_created";
    MessageType["PROJECT_UPDATED"] = "project_updated";
    MessageType["PROJECT_DELETED"] = "project_deleted";
})(MessageType || (MessageType = {}));
export class GetProjectHandler {
    ws = null;
    onProjectCreated;
    onProjectUpdated;
    onProjectDeleted;
    /**
     *
     * Connect
     *
     */
    async connect() {
        return new Promise((res, rej) => {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/ws`;
            this.ws = new WebSocket(wsUrl);
            this.ws.onopen = () => {
                this.subscribe('projects');
                res();
            };
            this.ws.onerror = (err) => {
                rej(err);
            };
            this.ws.onmessage = (e) => {
                try {
                    const message = JSON.parse(e.data);
                    this.handleMessage(message);
                }
                catch (err) {
                    console.error('Failed to parse message', err);
                }
            };
            this.ws.onclose = () => {
                console.log('WS disconnected');
                setTimeout(() => this.connect(), 3000);
            };
        });
    }
    /**
     *
     * Subscribe
     *
     */
    subscribe(channel) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'subscribe',
                channel: channel
            }));
        }
    }
    /**
     * Handle Message
     */
    handleMessage(message) {
        switch (message.type) {
            case MessageType.PROJECT_CREATED:
                if (this.onProjectCreated)
                    this.onProjectCreated(message.data);
                break;
            case MessageType.PROJECT_UPDATED:
                if (this.onProjectUpdated)
                    this.onProjectUpdated(message.data);
                break;
            case MessageType.PROJECT_DELETED:
                if (this.onProjectDeleted)
                    this.onProjectDeleted(message.data);
                break;
        }
    }
    setOnProjectCreated(cb) {
        this.onProjectCreated = cb;
    }
    setOnProjectUpdated(cb) {
        this.onProjectUpdated = cb;
    }
    setOnProjectDeleted(cb) {
        this.onProjectDeleted = cb;
    }
    /**
     *
     * Disconnect
     *
     */
    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}
