if(window.APP_ENV === 'prod') {
    window.PORT=3000
    window.SERVER_ADDR="0.0.0.0:${PORT}"
    window.SERVER_WS="https://portfolio-server-npwe.onrender.com"
    window.SERVER_URL="https://portfolio-server-npwe.onrender.com"
    
    window.API_URL="https://portfolio-server-npwe.onrender.com/api"
    window.WEB_URL="https://portfolio-eight-zeta-19.vercel.app"
} else {
    window.APP_ENV="dev"

    window.SERVER_WS="ws://localhost:3000/ws"
    window.SERVER_URL="http://localhost:3000"
    window.SERVER_ADDR="localhost:3000"

    window.API_URL="http://localhost:3000/api"
    window.WEB_URL="http://localhost:3000"
}