<script>
	import { X, ChevronLeft, ChevronRight, ExternalLink, Circle, CheckCircle, Bookmark } from '@lucide/svelte';
	import { selectedArticle, visibleArticles, articleModalOpen } from '$lib/stores';
	import { ArticleService } from '$lib/services/articleService';
	import { LibraryService } from '$lib/services/libraryService';

	let saving = false;

	$: currentIndex = $visibleArticles.findIndex((a) => a.ID === $selectedArticle?.ID);
	$: hasPrev = currentIndex > 0;
	$: hasNext = currentIndex < $visibleArticles.length - 1;

	function prev() {
		if (hasPrev) selectedArticle.set($visibleArticles[currentIndex - 1]);
	}

	function next() {
		if (hasNext) selectedArticle.set($visibleArticles[currentIndex + 1]);
	}

	function close() {
		articleModalOpen.set(false);
	}

	function handleKeydown(e) {
		if (!$articleModalOpen) return;
		if (e.key === 'Escape') close();
		if (e.key === 'ArrowLeft') prev();
		if (e.key === 'ArrowRight') next();
	}

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
</script>

<svelte:window on:keydown={handleKeydown} />

{#if $articleModalOpen && $selectedArticle}
	<!-- Backdrop -->
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 px-4"
		on:click|self={close}
	>
		<!-- Modal card -->
		<div class="relative flex max-h-[90vh] w-full max-w-3xl flex-col rounded-xl bg-slate-800 shadow-2xl">
			<!-- Toolbar -->
			<div class="flex flex-shrink-0 items-center justify-between border-b border-slate-700 px-5 py-3">
				<!-- Prev / counter / Next -->
				<div class="flex items-center gap-1">
					<button
						on:click={prev}
						disabled={!hasPrev}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Previous article (←)"
					>
						<ChevronLeft size={16} />
					</button>
					<span class="min-w-[4rem] text-center text-xs text-surface-500">
						{currentIndex + 1} / {$visibleArticles.length}
					</span>
					<button
						on:click={next}
						disabled={!hasNext}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Next article (→)"
					>
						<ChevronRight size={16} />
					</button>
				</div>

				<!-- Right actions -->
				<div class="flex items-center gap-1">
					<button
						on:click={toggleRead}
						class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
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
							class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-40"
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
							class="flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
							title="Open original"
						>
							<ExternalLink size={13} />
							<span>Original</span>
						</a>
					{/if}
					<button
						on:click={close}
						class="rounded p-1.5 text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100"
						title="Close (Esc)"
					>
						<X size={16} />
					</button>
				</div>
			</div>

			<!-- Scrollable content -->
			<div class="overflow-y-auto px-8 py-6">
				<h1 class="mb-3 text-xl font-semibold leading-snug text-surface-100">
					{$selectedArticle.Title}
				</h1>
				<div class="mb-6 flex flex-wrap items-center gap-1.5 text-xs text-surface-500">
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

				{#if $selectedArticle.Description}
					<div class="prose prose-invert prose-sm max-w-none text-surface-200">
						{@html $selectedArticle.Description}
					</div>
				{:else}
					<p class="text-sm text-surface-400">No content available.</p>
				{/if}
			</div>
		</div>
	</div>
{/if}
