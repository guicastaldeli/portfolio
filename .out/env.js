import { ParseEnv } from "./parser.js";
let cache = null;
let loading = null;
let ENV_KEY = '';
const KEY_ADDR = '/app/.env/env.key';
const ENV_ADDR = '/app/.env/env.temp';
export function GetEnv(key) {
    const val = cache?.[key];
    return val;
}
// Derive Key
async function deriveKey() {
    const enc = new TextEncoder();
    const raw = await crypto.subtle.digest('SHA-256', enc.encode(ENV_KEY));
    return crypto.subtle.importKey('raw', raw, { name: 'AES-GCM' }, false, ['decrypt']);
}
// Decrypt
async function decrypt(b64) {
    const binary = atob(b64.trim());
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    const NONCE_LEN = 12;
    const nonce = bytes.slice(0, NONCE_LEN);
    const ciphertext = bytes.slice(NONCE_LEN);
    const key = await deriveKey();
    const plain = await crypto.subtle.decrypt({ name: 'AES-GCM', iv: nonce }, key, ciphertext);
    const val = new TextDecoder().decode(plain);
    return val;
}
/**
 * Load
 */
async function Load() {
    if (cache)
        return;
    if (loading)
        return loading;
    loading = (async () => {
        try {
            const keyRes = await fetch(KEY_ADDR, { cache: 'no-store' });
            if (!keyRes.ok)
                throw new Error(`HTTP ${keyRes.status} for ${KEY_ADDR}`);
            ENV_KEY = (await keyRes.text()).trim();
            const res = await fetch(ENV_ADDR, { cache: 'no-store' });
            if (!res.ok)
                throw new Error(`HTTP ${res.status} for ${ENV_ADDR}`);
            const text = await res.text();
            const plaintext = await decrypt(text);
            cache = ParseEnv(plaintext);
            //console.log(`[Env] Loaded ${Object.keys(cache).length} keys from ${ENV_ADDR}`);
        }
        catch (err) {
            console.error(`[Env] Failed to load ${ENV_ADDR}:`, err);
            cache = {};
        }
    })();
    return loading;
}
/**
 *
 * Env
 *
 */
export const Env = {
    Ready: Load,
    All() { return { ...(cache ?? {}) }; },
    Has(key) { return cache ? key in cache : false; },
    Filename() { return ENV_ADDR; }
};
