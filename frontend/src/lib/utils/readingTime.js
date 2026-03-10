/**
 * Estimate reading time from an HTML string.
 * Returns a string like "4 min read" or null if no content.
 */
export function readingTime(html) {
	if (!html) return null;
	const text = html.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim();
	const words = text.split(' ').filter((w) => w.length > 0).length;
	if (words < 50) return null;
	const minutes = Math.max(1, Math.ceil(words / 238));
	return `${minutes} min read`;
}

/**
 * Save reading progress (scroll %) to localStorage.
 */
export function saveProgress(itemType, itemId, scrollEl) {
	if (!scrollEl || !itemId) return;
	const max = scrollEl.scrollHeight - scrollEl.clientHeight;
	if (max <= 0) return;
	const pct = scrollEl.scrollTop / max;
	if (pct > 0.01) localStorage.setItem(`progress:${itemType}:${itemId}`, pct.toFixed(4));
}

/**
 * Restore saved scroll position. Call after content has rendered (pass tick).
 */
export async function restoreProgress(itemType, itemId, scrollEl, tick) {
	if (!scrollEl || !itemId) return;
	const saved = localStorage.getItem(`progress:${itemType}:${itemId}`);
	if (!saved) return;
	const pct = parseFloat(saved);
	if (isNaN(pct) || pct <= 0.01) return;
	await tick();
	scrollEl.scrollTop = pct * (scrollEl.scrollHeight - scrollEl.clientHeight);
}
