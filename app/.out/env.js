class EnvConfig {
    config = null;
    configPromise;
    constructor() {
        this.configPromise = this.loadConfig();
    }
    async loadConfig() {
        let data = await this.fetchEnvFile('/.env/.env.prod');
        if (!data || Object.keys(data).length === 0) {
            console.log('Production env not found, falling back to dev Env');
            data = await this.fetchEnvFile('/.env/.env.dev');
        }
        const config = {
            APP_ENV: data.APP_ENV,
            SERVER_URL: data.SERVER_URL,
            API_URL: data.API_URL,
            WEB_URL: data.WEB_URL
        };
        this.config = config;
        return config;
    }
    async fetchEnvFile(path) {
        try {
            const response = await fetch(path);
            if (!response.ok) {
                console.warn(`Failed to fetch ${path}: ${response.status}`);
                return {};
            }
            const text = await response.text();
            return this.parseEnvFile(text);
        }
        catch (error) {
            console.warn(`Error fetching ${path}:`, error);
            return {};
        }
    }
    parseEnvFile(content) {
        const env = {};
        const lines = content.split('\n');
        for (const line of lines) {
            const trimmed = line.trim();
            if (!trimmed || trimmed.startsWith('#'))
                continue;
            const match = trimmed.match(/^([^=]+)=(.*)$/);
            if (match) {
                const key = match[1].trim();
                let value = match[2].trim();
                if ((value.startsWith('"') && value.endsWith('"')) ||
                    (value.startsWith("'") && value.endsWith("'"))) {
                    value = value.slice(1, -1);
                }
                env[key] = value;
            }
        }
        return env;
    }
    async getConfig() {
        return this.configPromise;
    }
    getConfigSync() {
        if (!this.config) {
            throw new Error('Config not loaded yet');
        }
        return this.config;
    }
}
const envConfig = new EnvConfig();
export async function loadEnv() {
    return await envConfig.getConfig();
}
export function getEnv() {
    return envConfig.getConfigSync();
}
