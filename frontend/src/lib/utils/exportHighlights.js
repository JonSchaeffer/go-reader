/**
 * Download a string as a file.
 */
function download(filename, content, mimeType) {
	const blob = new Blob([content], { type: mimeType });
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = filename;
	a.click();
	URL.revokeObjectURL(url);
}

function safeFilename(title) {
	return (title || 'highlights').replace(/[^a-z0-9]/gi, '-').replace(/-+/g, '-').toLowerCase().slice(0, 60);
}

/**
 * Export highlights as Markdown.
 * Groups by source article/item when exporting all.
 */
export function exportMarkdown(highlights, title = 'All Highlights') {
	const groups = new Map();
	for (const h of highlights) {
		const key = h.ItemTitle || h.ItemID || 'Unknown';
		if (!groups.has(key)) groups.set(key, []);
		groups.get(key).push(h);
	}

	const lines = [`# ${title}`, '', `*Exported ${new Date().toLocaleDateString()}*`, ''];

	for (const [source, hs] of groups) {
		if (groups.size > 1) lines.push(`## ${source}`, '');
		for (const h of hs) {
			lines.push(`> ${h.SelectedText}`);
			if (h.Note) lines.push(`> `, `> *Note: ${h.Note}*`);
			lines.push('');
		}
	}

	download(`${safeFilename(title)}.md`, lines.join('\n'), 'text/markdown');
}

/**
 * Export highlights as CSV (compatible with Notion, Obsidian, Anki, etc.)
 */
export function exportCSV(highlights, title = 'highlights') {
	const escape = (s) => `"${String(s ?? '').replace(/"/g, '""')}"`;

	const header = ['Text', 'Note', 'Color', 'Source', 'Date'].map(escape).join(',');
	const rows = highlights.map((h) =>
		[
			h.SelectedText,
			h.Note ?? '',
			h.Color,
			h.ItemTitle ?? '',
			new Date(h.CreatedAt).toLocaleDateString(),
		].map(escape).join(',')
	);

	download(`${safeFilename(title)}.csv`, [header, ...rows].join('\n'), 'text/csv');
}

/**
 * Copy all highlights as plain text to the clipboard.
 */
export async function copyToClipboard(highlights) {
	const lines = highlights.flatMap((h) => {
		const parts = [`"${h.SelectedText}"`];
		if (h.Note) parts.push(`Note: ${h.Note}`);
		if (h.ItemTitle) parts.push(`— ${h.ItemTitle}`);
		return [...parts, ''];
	});
	await navigator.clipboard.writeText(lines.join('\n'));
}
