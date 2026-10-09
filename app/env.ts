import { ParseEnv } from "./parser.js";

let cache: Record<string, string> | null = null;
let loading: Promise<void> | null = null;

const ENV_FILENAME = '.env.dev';

async function Load(): Promise<void> {
    if(cache) return;
    if(loading) return loading;

    loading = (async () => {
        const url = `./env/${ENV_FILENAME}`;

        try {
            const res = await fetch(url);
            if(!res.ok) throw new Error(`HTTP ${res.status} for ${url}`);

            const text = await res.text();
            cache = ParseEnv(text);
            console.log(`[Env] Loaded ${Object.keys(cache).length} keys from ${url}`);
        } catch(err) {
            console.error(`[Env] Failed to load ${url}:`, err);
            cache = {};
        }
    })();

    return loading;
}

/**
 * 
 * Getter
 * 
 */
export function GetEnv(key: string): string | undefined {
    const val = cache?.[key];
    return val;
}

/**
 * 
 * Env
 * 
 */
export const Env = {
    Ready: Load,
    All(): Record<string, string> { return { ...(cache ?? {}) }; },
    Has(key: string): boolean { return cache ? key in cache : false; },
    Filename(): string { return ENV_FILENAME; }
};