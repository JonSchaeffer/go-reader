<script>
	import { onMount } from 'svelte';
	import { Highlighter, Trash2, StickyNote, Download, Copy, Check } from '@lucide/svelte';
	import { HighlightService } from '$lib/services/highlightService';
	import { selectedArticle, visibleArticles, articleModalOpen } from '$lib/stores';
	import { libraryApi, articleApi } from '$lib/api.js';
	import { exportMarkdown, exportCSV, copyToClipboard } from '$lib/utils/exportHighlights';

	const COLOR_MAP = {
		yellow: '#fef08a',
		green:  '#bbf7d0',
		blue:   '#bfdbfe',
		pink:   '#fbcfe8',
	};

	let highlights = [];
	let loading = true;
	let copied = false;

	async function handleCopy() {
		await copyToClipboard(highlights);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}

	onMount(async () => {
		highlights = await HighlightService.getAll();
		loading = false;
	});

	function formatDate(dateString) {
		if (!dateString) return '';
		try {
			const date = new Date(dateString);
			const now = new Date();
			const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));
			if (diffDays === 0) return 'Today';
			if (diffDays === 1) return 'Yesterday';
			if (diffDays < 7) return `${diffDays} days ago`;
			return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
		} catch {
			return '';
		}
	}

	async function openSource(h) {
		let article;
		if (h.ItemType === 'library') {
			const item = await libraryApi.getById(h.ItemID);
			if (!item) return;
			article = {
				ID: item.ID,
				Title: item.Title,
				Description: item.Content,
				Author: item.Author,
				PublishDate: item.PublishedDate,
				Read: item.Read,
				Link: item.URL,
				Category: null,
				_type: 'library',
				_archived: item.Archived,
			};
		} else {
			const results = await articleApi.getById(h.ItemID);
			const item = Array.isArray(results) ? results[0] : results;
			if (!item) return;
			article = {
				ID: item.ID,
				Title: item.Title,
				Description: item.Description,
				Author: item.Author,
				PublishDate: item.PublishDate,
				Read: item.Read,
				Link: item.Link,
				Category: item.Category,
				_type: 'article',
			};
		}
		// Open in the modal without navigating away from highlights
		visibleArticles.set([article]);
		selectedArticle.set(article);
		articleModalOpen.set(true);
	}

	async function deleteHighlight(h) {
		await HighlightService.delete(h.ID);
		highlights = highlights.filter((x) => x.ID !== h.ID);
	}
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="flex items-center gap-3 border-b border-white/10 px-4 py-3">
		<Highlighter size={15} class="text-surface-400" />
		<span class="text-sm font-semibold text-surface-100">Highlights</span>
		{#if !loading}
			<span class="text-xs text-surface-500">{highlights.length} total</span>
		{/if}

		{#if !loading && highlights.length > 0}
			<div class="ml-auto flex items-center gap-1">
				<button
					on:click={handleCopy}
					title="Copy all to clipboard"
					class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100"
				>
					{#if copied}
						<Check size={13} class="text-green-400" />
						<span class="text-green-400">Copied</span>
					{:else}
						<Copy size={13} />
						<span>Copy</span>
					{/if}
				</button>
				<button
					on:click={() => exportMarkdown(highlights)}
					title="Export as Markdown"
					class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100"
				>
					<Download size={13} />
					<span>Markdown</span>
				</button>
				<button
					on:click={() => exportCSV(highlights)}
					title="Export as CSV"
					class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100"
				>
					<Download size={13} />
					<span>CSV</span>
				</button>
			</div>
		{/if}
	</div>

	<!-- List -->
	<div class="flex-1 overflow-y-auto">
		{#if loading}
			<div class="p-4 text-sm text-surface-400">Loading...</div>
		{:else if highlights.length === 0}
			<div class="flex flex-col items-center justify-center gap-3 py-24 text-center">
				<Highlighter size={28} class="text-surface-600" />
				<p class="text-sm text-surface-500">No highlights yet. Select text in any article to highlight it.</p>
			</div>
		{:else}
			<div class="divide-y divide-white/10">
				{#each highlights as h (h.ID)}
					<div class="group px-4 py-4 hover:bg-slate-800/50 transition-colors">
						<!-- Highlight text -->
						<div
							class="mb-2 rounded px-2 py-1.5 text-sm text-slate-900 leading-relaxed"
							style="background: {COLOR_MAP[h.Color] ?? COLOR_MAP.yellow}"
						>
							{h.SelectedText}
						</div>

						<!-- Note -->
						{#if h.Note}
							<div class="mb-2 flex items-start gap-1.5 text-xs text-surface-400">
								<StickyNote size={12} class="mt-0.5 flex-shrink-0" />
								<span class="italic">{h.Note}</span>
							</div>
						{/if}

						<!-- Source + actions -->
						<div class="flex items-center justify-between gap-2">
							<button
								on:click={() => openSource(h)}
								class="min-w-0 flex-1 text-left text-xs text-surface-500 hover:text-blue-400 transition-colors truncate"
								title={h.ItemTitle || 'Open article'}
							>
								{h.ItemTitle || 'Untitled'}
							</button>

							<div class="flex flex-shrink-0 items-center gap-2">
								<span class="text-xs text-surface-600">{formatDate(h.CreatedAt)}</span>
								<button
									on:click={() => deleteHighlight(h)}
									class="opacity-0 group-hover:opacity-100 rounded p-1 text-surface-600 hover:text-red-400 transition-all"
									title="Delete highlight"
								>
									<Trash2 size={12} />
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
