export function formatTime(seconds) {
    if (seconds <= 0 || isNaN(seconds)) return '0s';

    const h=Math.floor(seconds / 3600);
    const m=Math.floor((seconds % 3600) / 60);
    const s=Math.floor(seconds % 60);

    const parts=[];
    if (h > 0) parts.push(`${h}h`);
    if (m > 0) parts.push(`${m}m`);
    if (s > 0 || parts.length === 0) parts.push(`${s}s`);

    return parts.join(' ');
}

export function formatMs(ms) {
    return formatTime(Math.max(0, Math.floor(ms / 1000)));
}
