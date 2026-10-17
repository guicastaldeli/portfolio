export class GetClientCount {
    main;
    ws = null;
    el = null;
    container = null;
    constructor(main) {
        this.main = main;
        this.container = document.querySelector('#client-count');
        this.el = document.querySelector('#c-clients');
    }
    /**
     * Connect
     */
    async connect() {
        if (!this.el)
            return;
        const url = '/count';
        this.main.protocol(url);
        this.ws = new WebSocket(url);
        this.ws.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                if (data.type === 'clientsUpdate') {
                    const content = data.count;
                    this.el.textContent = `${content}`;
                    this.container?.setAttribute('data-count', data.count.toString());
                }
            }
            catch (err) {
                console.error(err);
            }
        };
        this.ws.onerror = (err) => {
            console.error('Client WS error', err);
        };
        this.ws.onclose = () => {
            console.log('Client count disconnected, reconnecting...');
            setTimeout(this.connect, 3000);
        };
    }
}
