export class GetTimeStream {
    main;
    el = null;
    ws = null;
    constructor(main) {
        this.main = main;
        this.el = document.querySelector('.main #time-stream #t-time');
    }
    /**
     * Connect
     */
    async connect() {
        if (!this.el)
            return;
        const url = '/time-stream';
        this.main.protocol(url);
        this.ws = new WebSocket(url);
        this.ws.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                if (data.type === 'timeUpdate') {
                    const date = new Date(data.timestamp);
                    const formatted = date.toLocaleString();
                    this.el.textContent = formatted;
                }
            }
            catch (err) {
                console.error(err);
            }
        };
        this.ws.onerror = (err) => {
            console.error('Time WS error', err);
        };
        this.ws.onclose = () => {
            console.log('Time stream disconnected, reconnecting...');
            setTimeout(this.connect, 3000);
        };
    }
}
