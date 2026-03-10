<script>
	import { tick } from 'svelte';
	import { X, ChevronLeft, ChevronRight, ExternalLink, Circle, CheckCircle, Bookmark } from '@lucide/svelte';
	import { selectedArticle, visibleArticles, articleModalOpen, readerSettings } from '$lib/stores';
	import { ArticleService } from '$lib/services/articleService';
	import { LibraryService } from '$lib/services/libraryService';
	import { HighlightService } from '$lib/services/highlightService';
	import { readingTime, saveProgress, restoreProgress } from '$lib/utils/readingTime';
	import HighlightPopover from './HighlightPopover.svelte';
	import ReaderSettings from './ReaderSettings.svelte';

	let saving = false;
	let highlights = [];
	let popover = null;
	let contentEl = null;

	const FONT_SIZE   = { sm: '13px', md: '15px', lg: '17px', xl: '19px' };
	const FONT_FAMILY = {
		sans:  "system-ui, -apple-system, sans-serif",
		serif: "Georgia, 'Times New Roman', serif",
		mono:  "'JetBrains Mono', Consolas, monospace",
	};
	const LINE_WIDTH  = { narrow: '52ch', medium: '65ch', wide: '85ch' };
	const THEME_STYLE = {
		dark:  { bg: '',        text: '' },
		light: { bg: '#ffffff', text: '#0f172a' },
		sepia: { bg: '#f5edd8', text: '#433422' },
	};

	$: rs = $readerSettings;
	$: proseStyle = [
		`font-size: ${FONT_SIZE[rs.fontSize] ?? FONT_SIZE.md}`,
		`font-family: ${FONT_FAMILY[rs.fontFamily] ?? FONT_FAMILY.sans}`,
		`max-width: ${LINE_WIDTH[rs.lineWidth] ?? LINE_WIDTH.medium}`,
	].join(';');
	$: themeStyle = rs.theme !== 'dark'
		? `background: ${THEME_STYLE[rs.theme].bg}; color: ${THEME_STYLE[rs.theme].text}; padding: 2rem; border-radius: 8px;`
		: '';

	$: currentIndex = $visibleArticles.findIndex((a) => a.ID === $selectedArticle?.ID);
	$: hasPrev = currentIndex > 0;
	$: hasNext = currentIndex < $visibleArticles.length - 1;
	$: timeToRead = readingTime($selectedArticle?.Description);

	// Load highlights + restore scroll when modal opens or article changes
	$: if ($articleModalOpen && $selectedArticle) onArticleChange($selectedArticle);

	async function onArticleChange(article) {
		highlights = [];
		if (!article) return;
		const itemType = article._type === 'library' ? 'library' : 'article';
		highlights = await HighlightService.getForItem(itemType, article.ID);
		restoreProgress(itemType, article.ID, contentEl, tick);
	}

	$: renderedHtml = $selectedArticle?.Description
		? HighlightService.applyToHtml($selectedArticle.Description, highlights)
		: '';

	function prev() { if (hasPrev) selectedArticle.set($visibleArticles[currentIndex - 1]); }
	function next() { if (hasNext) selectedArticle.set($visibleArticles[currentIndex + 1]); }

	function close() { popover = null; articleModalOpen.set(false); }

	function handleKeydown(e) {
		if (!$articleModalOpen) return;
		if (e.key === 'Escape') { if (popover) { popover = null; } else { close(); } return; }
		if (popover) return;
		if (e.key === 'ArrowLeft') prev();
		if (e.key === 'ArrowRight') next();
	}

	function formatDate(dateString) {
		if (!dateString) return '';
		try {
			return new Date(dateString).toLocaleDateString('en-US', {
				year: 'numeric', month: 'long', day: 'numeric',
				hour: '2-digit', minute: '2-digit'
			});
		} catch { return dateString; }
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
		try { await LibraryService.saveArticle($selectedArticle.ID); }
		finally { saving = false; }
	}

	// Debounced scroll save
	let scrollTimer;
	function onScroll() {
		popover = null;
		clearTimeout(scrollTimer);
		if (!$selectedArticle) return;
		scrollTimer = setTimeout(() => {
			const itemType = $selectedArticle._type === 'library' ? 'library' : 'article';
			saveProgress(itemType, $selectedArticle.ID, contentEl);
		}, 1000);
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
		popover = { x: Math.min(rect.left, window.innerWidth - 220), y: rect.bottom + 6, selectedText: text, textOffset, existing: null };
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
		await onArticleChange($selectedArticle);
	}

	async function handleDelete() {
		if (!popover?.existing) return;
		await HighlightService.delete(popover.existing.ID);
		popover = null;
		await onArticleChange($selectedArticle);
	}
</script>

<svelte:window
	on:keydown={handleKeydown}
	on:mousedown={(e) => {
		if ($articleModalOpen && !e.target.closest('.highlight-popover') && !e.target.closest('.reader-settings-panel') && !e.target.closest('.reader-settings-btn'))
			popover = null;
	}}
/>

{#if $articleModalOpen && $selectedArticle}
	<!-- Backdrop -->
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 px-4" on:click|self={close}>
		<!-- Modal card -->
		<!-- svelte-ignore a11y-no-static-element-interactions -->
		<div class="relative flex max-h-[90vh] w-full max-w-3xl flex-col rounded-xl bg-slate-800 shadow-2xl" on:mouseup={handleMouseup}>
			<!-- Toolbar -->
			<div class="flex flex-shrink-0 items-center justify-between border-b border-slate-700 px-5 py-3">
				<!-- Prev / counter / Next -->
				<div class="flex items-center gap-1">
					<button on:click={prev} disabled={!hasPrev}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Previous article (←)"><ChevronLeft size={16} /></button>
					<span class="min-w-[4rem] text-center text-xs text-surface-500">
						{currentIndex + 1} / {$visibleArticles.length}
					</span>
					<button on:click={next} disabled={!hasNext}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Next article (→)"><ChevronRight size={16} /></button>
				</div>

				<!-- Right actions -->
				<div class="flex items-center gap-1">
					<button on:click={toggleRead}
						class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
						title={$selectedArticle.Read ? 'Mark as unread' : 'Mark as read'}>
						{#if $selectedArticle.Read}
							<Circle size={13} /><span>Unread</span>
						{:else}
							<CheckCircle size={13} /><span>Read</span>
						{/if}
					</button>
					{#if $selectedArticle._type !== 'library'}
						<button on:click={saveToLibrary} disabled={saving}
							class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-40"
							title="Save to library">
							<Bookmark size={13} /><span>Save</span>
						</button>
					{/if}
					{#if $selectedArticle.Link}
						<a href={$selectedArticle.Link} target="_blank" rel="noopener noreferrer"
							class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
							title="Open original">
							<ExternalLink size={13} /><span>Original</span>
						</a>
					{/if}
					<ReaderSettings />
					<button on:click={close}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
						title="Close (Esc)"><X size={16} /></button>
				</div>
			</div>

			<!-- Scrollable content -->
			<div class="overflow-y-auto px-8 py-6" bind:this={contentEl} on:scroll={onScroll}>
				<div style={themeStyle}>
					<h1 class="mb-3 text-xl font-semibold leading-snug text-surface-100">
						{$selectedArticle.Title}
					</h1>
					<div class="mb-6 flex flex-wrap items-center gap-1.5 text-xs text-surface-500">
						{#if $selectedArticle.Author}
							<span>{$selectedArticle.Author}</span><span>·</span>
						{/if}
						{#if $selectedArticle.PublishDate}
							<span>{formatDate($selectedArticle.PublishDate)}</span>
						{/if}
						{#if timeToRead}
							<span>·</span><span>{timeToRead}</span>
						{/if}
						{#if $selectedArticle.Category && $selectedArticle.Category !== 'Uncategorized'}
							<span>·</span>
							<span class="rounded bg-slate-700 px-1.5 py-0.5">{$selectedArticle.Category}</span>
						{/if}
					</div>

					{#if renderedHtml}
						<div class="prose prose-invert prose-sm max-w-none" style={proseStyle}>
							{@html renderedHtml}
						</div>
					{:else if $selectedArticle.Description}
						<div class="prose prose-invert prose-sm max-w-none" style={proseStyle}>
							{@html $selectedArticle.Description}
						</div>
					{:else}
						<p class="text-sm text-surface-400">No content available.</p>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

{#if popover && $articleModalOpen}
	<HighlightPopover
		x={popover.x}
		y={popover.y}
		existingHighlight={popover.existing}
		on:save={handleSave}
		on:delete={handleDelete}
		on:cancel={() => popover = null}
	/>
{/if}
