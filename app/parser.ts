/**
 * .env parser.
 *
 * Handles:
 *   - KEY=value
 *   - KEY="value with spaces"
 *   - KEY='value'
 *   - # comments
 *   - blank lines
 *   - trailing whitespace
 *   - \n and \r\n line endings
 *
 */
export function ParseEnv(text: string): Record<string, string> {
    const out: Record<string, string> = {};
    const t = /\r?\n/;
    for(const rawLine of text.split(t)) {
        const line = rawLine.trim();
        if(line === '' || line.startsWith('#')) continue;

        const eq = line.indexOf('=');
        if(eq === -1) continue;

        let value = line.slice(eq + 1).trim();
        if((value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.startsWith("'"))) {
            value = value.slice(1, -1);
        }
        
        const key = line.slice(0, eq).trim();
        out[key] = value;
    }

    return out;
}