import { Main } from "./main";

export class GetTimeStream {
    private main: Main;
    private el: HTMLSpanElement | null = null;
    public ws: WebSocket | null = null; 

    constructor(main: Main) {
        this.main = main;
        this.el = document.querySelector('.main time-stream');
    }

    /**
     * Connect
     */
    public connect() {
        const url = '/time-stream';

        this.main.protocol(url);

        this.ws = new WebSocket(url);
        this.ws.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                if(data.type === 'timeUpdate') {
                    this.el!.textContent = `Time: ${data.formatted}`;
                }
            } catch(err) {
                console.error(err);
            }
        }
        this.ws.onerror = (err) => {
            console.error('Time WS error', err);
        }
        this.ws.onclose = () => {
            console.log('Time stream disconnected, reconnecting...');
            setTimeout(this.connect, 3000);
        }
    }
}