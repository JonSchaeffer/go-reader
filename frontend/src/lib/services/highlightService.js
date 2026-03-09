import { highlightApi } from '../api.js';

export class HighlightService {
	static async getForItem(itemType, itemId) {
		const result = await highlightApi.getByItem(itemType, itemId);
		return Array.isArray(result) ? result : [];
	}

	static async create(itemType, itemId, selectedText, color = 'yellow', note = '', textOffset = -1) {
		return highlightApi.create({ itemType, itemId, selectedText, color, note, textOffset });
	}

	static async update(id, color, note) {
		return highlightApi.update(id, color, note);
	}

	static async delete(id) {
		return highlightApi.delete(id);
	}

	/**
	 * Apply highlights to an HTML string using DOM-based text-node walking.
	 * Handles HTML entities and selections that cross inline tags (e.g. <strong>).
	 */
	static applyToHtml(html, highlights) {
		if (!highlights || highlights.length === 0) return html;

		const colorMap = {
			yellow: '#fef08a',
			green:  '#bbf7d0',
			blue:   '#bfdbfe',
			pink:   '#fbcfe8',
		};

		const container = document.createElement('div');
		container.innerHTML = html;

		// Longest first so sub-strings don't get wrapped before their parent phrase
		const sorted = [...highlights].sort((a, b) => b.SelectedText.length - a.SelectedText.length);

		for (const h of sorted) {
			const bg = colorMap[h.Color] ?? colorMap.yellow;
			wrapTextInNode(container, h.SelectedText, h.ID, bg, h.Note || '', h.TextOffset ?? -1);
		}

		return container.innerHTML;
	}
}

/**
 * Walk all text nodes under `root`, find `searchText`, and wrap matched
 * portions with a <mark> element. Handles matches that span multiple text
 * nodes (e.g. text split by an inline <em> or <strong> tag).
 */
function wrapTextInNode(root, searchText, highlightId, bgColor, note, textOffset = -1) {
	// Collect all text nodes in document order
	const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, null);
	const nodes = [];
	let n;
	while ((n = walker.nextNode())) nodes.push(n);

	// Build a single concatenated string with offsets per node
	let combined = '';
	const ranges = nodes.map((node) => {
		const start = combined.length;
		combined += node.textContent;
		return { node, start, end: combined.length };
	});

	// Use the stored offset when available; fall back to first occurrence
	const matchStart = (textOffset >= 0 && combined.slice(textOffset, textOffset + searchText.length) === searchText)
		? textOffset
		: combined.indexOf(searchText);
	if (matchStart === -1) return;
	const matchEnd = matchStart + searchText.length;

	// Nodes that overlap with [matchStart, matchEnd)
	const overlapping = ranges.filter((r) => r.start < matchEnd && r.end > matchStart);
	if (!overlapping.length) return;

	// Process each overlapping node — they're already in DOM order
	for (const r of overlapping) {
		const localStart = Math.max(0, matchStart - r.start);
		const localEnd   = Math.min(r.node.textContent.length, matchEnd - r.start);
		const matched    = r.node.textContent.slice(localStart, localEnd);
		if (!matched) continue;

		const before = r.node.textContent.slice(0, localStart);
		const after  = r.node.textContent.slice(localEnd);
		const parent = r.node.parentNode;

		const mark = document.createElement('mark');
		mark.style.cssText = `background:${bgColor};border-radius:2px;padding:0 1px`;
		mark.setAttribute('data-highlight-id', String(highlightId));
		if (note) mark.setAttribute('data-note', note);
		mark.textContent = matched;

		if (before) parent.insertBefore(document.createTextNode(before), r.node);
		parent.insertBefore(mark, r.node);
		if (after)  parent.insertBefore(document.createTextNode(after), r.node);
		parent.removeChild(r.node);
	}
}
