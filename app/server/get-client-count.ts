import { Main } from "./main.js";

export class GetClientCount {
    private main: Main;
    private el: HTMLSpanElement | null = null;
    public ws: WebSocket | null = null; 
    
    constructor(main: Main) {
        this.main = main;
        this.el = document.querySelector('.main #client-count');
    }
    
    /**
     * Connect
     */
    public async connect(): Promise<void> {
        const url = '/count';
    
        this.main.protocol(url);
    
        this.ws = new WebSocket(url);
        this.ws.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                if(data.type === 'clientsUpdate') {
                    this.el!.textContent = `Clients: ${data.count}`;
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