import { Main } from "./main.js";

export class GetClientCount {
    private main: Main;
    public ws: WebSocket | null = null;

    private el: HTMLSpanElement | null = null;
    private container: HTMLDivElement | null = null;

    constructor(main: Main) {
        this.main = main;
        this.container = document.querySelector('#client-count'); 
        this.el = document.querySelector('#c-clients');
    }
    
    /**
     * Connect
     */
    public async connect(): Promise<void> {
        if(!this.el) return;
        
        const url = '/count';
    
        this.main.protocol(url);
    
        this.ws = new WebSocket(url);
        this.ws.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                if(data.type === 'clientsUpdate') {
                    const content = data.count;
                    this.el!.textContent = `${content}`;
                    this.container?.setAttribute('data-count', data.count.toString());
                }
            } catch(err) {
                console.error(err);
            }
        }
        this.ws.onerror = (err) => {
            console.error('Client WS error', err);
        }
        this.ws.onclose = () => {
            console.log('Client count disconnected, reconnecting...');
            setTimeout(this.connect, 3000);
        }
    }
}