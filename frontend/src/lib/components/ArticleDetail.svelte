<script>
	import { ExternalLink, X, Circle, CheckCircle, Expand, Bookmark } from '@lucide/svelte';
	import { selectedArticle, articleModalOpen } from '$lib/stores';
	import { ArticleService } from '$lib/services/articleService';
	import { LibraryService } from '$lib/services/libraryService';
	import { HighlightService } from '$lib/services/highlightService';
	import HighlightPopover from './HighlightPopover.svelte';

	let saving = false;
	let highlights = [];
	let popover = null; // { x, y, selectedText, textOffset, existing? }
	let contentEl = null;

	// Reload highlights whenever the article changes
	$: if ($selectedArticle) loadHighlights($selectedArticle);

	async function loadHighlights(article) {
		highlights = [];
		if (!article) return;
		const itemType = article._type === 'library' ? 'library' : 'article';
		highlights = await HighlightService.getForItem(itemType, article.ID);
	}

	$: renderedHtml = $selectedArticle?.Description
		? HighlightService.applyToHtml($selectedArticle.Description, highlights)
		: '';

	function formatDate(dateString) {
		if (!dateString) return '';
		try {
			return new Date(dateString).toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'long',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch {
			return dateString;
		}
	}

	async function toggleRead() {
		if (!$selectedArticle) return;
		if ($selectedArticle._type === 'library') {
			await LibraryService.toggleRead($selectedArticle.ID, $selectedArticle.Read, $selectedArticle._archived);
		} else {
			await ArticleService.toggleReadStatus($selectedArticle.ID, $selectedArticle.Read);
		}
		selectedArticle.update((a) => ({ ...a, Read: !a.Read }));
	}

	async function saveToLibrary() {
		if (!$selectedArticle || saving) return;
		saving = true;
		try {
			await LibraryService.saveArticle($selectedArticle.ID);
		} finally {
			saving = false;
		}
	}

	function close() {
		popover = null;
		selectedArticle.set(null);
	}

	function getTextOffset(containerEl, range) {
		const walker = document.createTreeWalker(containerEl, NodeFilter.SHOW_TEXT, null);
		let offset = 0, node;
		while ((node = walker.nextNode())) {
			if (node === range.startContainer) return offset + range.startOffset;
			offset += node.textContent.length;
		}
		return -1;
	}

	function handleMouseup(e) {
		if (e.target.closest('.highlight-popover')) return;

		const markEl = e.target.closest('mark[data-highlight-id]');
		if (markEl) {
			const hid = parseInt(markEl.getAttribute('data-highlight-id'));
			const existing = highlights.find((h) => h.ID === hid);
			if (existing) {
				const rect = markEl.getBoundingClientRect();
				popover = { x: Math.min(rect.left, window.innerWidth - 220), y: rect.bottom + 6, selectedText: existing.SelectedText, textOffset: existing.TextOffset ?? -1, existing };
				return;
			}
		}

		const sel = window.getSelection();
		if (!sel || sel.isCollapsed) { popover = null; return; }
		const text = sel.toString().trim();
		if (text.length < 2) { popover = null; return; }

		const range = sel.getRangeAt(0);
		const rect = range.getBoundingClientRect();
		const proseEl = contentEl?.querySelector('.prose');
		const textOffset = proseEl ? getTextOffset(proseEl, range) : -1;
		popover = {
			x: Math.min(rect.left, window.innerWidth - 220),
			y: rect.bottom + 6,
			selectedText: text,
			textOffset,
			existing: null
		};
	}

	async function handleSave({ detail }) {
		if (!popover || !$selectedArticle) return;
		const itemType = $selectedArticle._type === 'library' ? 'library' : 'article';
		if (popover.existing) {
			await HighlightService.update(popover.existing.ID, detail.color, detail.note);
		} else {
			await HighlightService.create(itemType, $selectedArticle.ID, popover.selectedText, detail.color, detail.note, popover.textOffset ?? -1);
		}
		popover = null;
		window.getSelection()?.removeAllRanges();
		await loadHighlights($selectedArticle);
	}

	async function handleDelete() {
		if (!popover?.existing) return;
		await HighlightService.delete(popover.existing.ID);
		popover = null;
		await loadHighlights($selectedArticle);
	}
</script>

<svelte:window on:mousedown={(e) => { if (!e.target.closest('.highlight-popover')) popover = null; }} />

<aside class="flex w-[420px] flex-shrink-0 flex-col overflow-hidden border-l border-slate-700 bg-slate-800">
	{#if $selectedArticle}
		<!-- Toolbar -->
		<div class="flex items-center justify-between border-b border-slate-700 px-3 py-2">
			<div class="flex items-center gap-2">
				<button
					on:click={() => articleModalOpen.set(true)}
					class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
					title="Open in reading view"
				>
					<Expand size={13} />
					<span>Expand</span>
				</button>
				<button
					on:click={toggleRead}
					class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
					title={$selectedArticle.Read ? 'Mark as unread' : 'Mark as read'}
				>
					{#if $selectedArticle.Read}
						<Circle size={13} />
						<span>Unread</span>
					{:else}
						<CheckCircle size={13} />
						<span>Read</span>
					{/if}
				</button>
				{#if $selectedArticle._type !== 'library'}
					<button
						on:click={saveToLibrary}
						disabled={saving}
						class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors disabled:opacity-40"
						title="Save to library"
					>
						<Bookmark size={13} />
						<span>Save</span>
					</button>
				{/if}
				{#if $selectedArticle.Link}
					<a
						href={$selectedArticle.Link}
						target="_blank"
						rel="noopener noreferrer"
						class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
						title="Open original"
					>
						<ExternalLink size={13} />
						<span>Original</span>
					</a>
				{/if}
			</div>
			<button
				on:click={close}
				class="rounded p-1 text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
			>
				<X size={16} />
			</button>
		</div>

		<!-- Content -->
		<!-- svelte-ignore a11y-no-static-element-interactions -->
		<div class="flex-1 overflow-y-auto p-5" bind:this={contentEl} on:mouseup={handleMouseup} on:scroll={() => popover = null}>
			<!-- Title -->
			<h2 class="text-surface-100 mb-3 text-lg font-semibold leading-snug">
				{$selectedArticle.Title}
			</h2>

			<!-- Meta -->
			<div class="text-surface-500 mb-4 flex flex-wrap items-center gap-1.5 text-xs">
				{#if $selectedArticle.Author}
					<span>{$selectedArticle.Author}</span>
					<span>·</span>
				{/if}
				{#if $selectedArticle.PublishDate}
					<span>{formatDate($selectedArticle.PublishDate)}</span>
				{/if}
				{#if $selectedArticle.Category && $selectedArticle.Category !== 'Uncategorized'}
					<span>·</span>
					<span class="rounded bg-slate-700 px-1.5 py-0.5">{$selectedArticle.Category}</span>
				{/if}
			</div>

			<!-- Body -->
			{#if renderedHtml}
				<div class="prose prose-invert prose-sm max-w-none text-surface-200">
					{@html renderedHtml}
				</div>
			{:else if $selectedArticle.Description}
				<div class="prose prose-invert prose-sm max-w-none text-surface-200">
					{@html $selectedArticle.Description}
				</div>
			{:else}
				<p class="text-surface-400 text-sm">No content available.</p>
			{/if}
		</div>
	{/if}
</aside>

{#if popover}
	<HighlightPopover
		x={popover.x}
		y={popover.y}
		existingHighlight={popover.existing}
		on:save={handleSave}
		on:delete={handleDelete}
		on:cancel={() => popover = null}
	/>
{/if}
