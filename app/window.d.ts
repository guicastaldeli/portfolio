declare global {
    interface Window {
        APP_ENV: 'dev',
        SERVER_URL: string,
        SERVER_WS: string,
        API_URL: string
    }
}

export {}