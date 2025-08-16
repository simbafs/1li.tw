import { BASE } from './api'

export function formatShortPath(shortPath: string) {
    return `${BASE}/r/${shortPath}`
}
