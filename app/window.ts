class Window {
    public vars = {
        APP_ENV: '',
        PORT: 3000,
        SERVER_ADDR: '',
        SERVER_WS: '',
        SERVER_URL: '',
        API_URL: '',
        WEB_URL: ''
    }

    init() {
        if(this.vars.APP_ENV === 'prod') {
            this.vars.SERVER_ADDR=`0.0.0.0:${this.vars.PORT}`
            this.vars.SERVER_WS="https://portfolio-server-npwe.onrender.com"
            this.vars.SERVER_URL="https://portfolio-server-npwe.onrender.com"
            
            this.vars.API_URL="https://portfolio-server-npwe.onrender.com/api"
            this.vars.WEB_URL="https://portfolio-eight-zeta-19.vercel.app"
        } else {
            this.vars.APP_ENV="dev"
        
            this.vars.SERVER_WS="ws://localhost:3000/ws"
            this.vars.SERVER_URL="http://localhost:3000"
            this.vars.SERVER_ADDR="localhost:3000"
        
            this.vars.API_URL="http://localhost:3000/api"
            this.vars.WEB_URL="http://localhost:3000"
        }
    }
}

const window = new Window();
export default window;