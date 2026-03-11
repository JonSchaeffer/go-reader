<script>
	import { onMount } from 'svelte';
	import { ChevronLeft, ChevronRight, Check, BookOpen } from '@lucide/svelte';
	import { reviewApi } from '$lib/api';
	import { reviewCount, selectedArticle, visibleArticles, articleModalOpen } from '$lib/stores';
	import { articleApi, libraryApi } from '$lib/api';

	const QUEUE_SIZE = 5;

	let queue = [];
	let index = 0;
	let loading = true;
	let completing = false;
	let done = false;

	$: current = queue[index] ?? null;
	$: progress = queue.length > 0 ? index + 1 : 0;

	const COLOR_BG = {
		yellow: 'bg-yellow-400/20 border-yellow-400/40 text-yellow-200',
		green:  'bg-green-400/20  border-green-400/40  text-green-200',
		blue:   'bg-blue-400/20   border-blue-400/40   text-blue-200',
		pink:   'bg-pink-400/20   border-pink-400/40   text-pink-200',
	};

	const COLOR_MARK = {
		yellow: '#fbbf24',
		green:  '#4ade80',
		blue:   '#60a5fa',
		pink:   '#f472b6',
	};

	onMount(async () => {
		await loadQueue();
	});

	async function loadQueue() {
		loading = true;
		try {
			queue = await reviewApi.getQueue(QUEUE_SIZE) ?? [];
			done = queue.length === 0;
			index = 0;
		} finally {
			loading = false;
		}
	}

	async function markAndAdvance() {
		if (!current || completing) return;
		completing = true;
		await reviewApi.complete(current.ID);
		// Update badge count
		reviewCount.update((n) => Math.max(0, n - 1));
		completing = false;

		if (index < queue.length - 1) {
			index++;
		} else {
			done = true;
		}
	}

	function prev() {
		if (index > 0) index--;
	}

	async function openSource() {
		if (!current) return;
		try {
			let article;
			if (current.ItemType === 'library') {
				article = await libraryApi.getById(current.ItemID);
				article = { ...article, _type: 'library', Description: article.Content, PublishDate: article.PublishedDate };
			} else {
				const res = await articleApi.getById(current.ItemID);
				article = Array.isArray(res) ? res[0] : res;
				article = { ...article, _type: 'article' };
			}
			visibleArticles.set([article]);
			selectedArticle.set(article);
			articleModalOpen.set(true);
		} catch (e) {
			console.error('Failed to open source', e);
		}
	}
</script>

<div class="flex flex-1 flex-col items-center justify-center bg-slate-900 px-6">
	{#if loading}
		<p class="text-surface-500 text-sm">Loading your review queue…</p>

	{:else if done}
		<!-- Done state -->
		<div class="flex flex-col items-center gap-4 text-center">
			<div class="flex h-16 w-16 items-center justify-center rounded-full bg-green-500/20">
				<Check size={32} class="text-green-400" />
			</div>
			<h2 class="text-xl font-semibold text-surface-100">You're all caught up!</h2>
			<p class="max-w-sm text-sm text-surface-500">
				You've reviewed all your highlights for today. Come back tomorrow for a fresh queue.
			</p>
			<button
				on:click={loadQueue}
				class="mt-2 rounded-lg border border-slate-600 px-4 py-2 text-sm text-surface-400 transition-colors hover:bg-slate-800 hover:text-surface-100"
			>
				Check again
			</button>
		</div>

	{:else if current}
		<!-- Progress bar -->
		<div class="mb-8 w-full max-w-xl">
			<div class="mb-2 flex items-center justify-between text-xs text-surface-500">
				<span>Daily Review</span>
				<span>{progress} / {queue.length}</span>
			</div>
			<div class="h-1 w-full rounded-full bg-slate-700">
				<div
					class="h-1 rounded-full bg-blue-500 transition-all duration-300"
					style="width: {(progress / queue.length) * 100}%"
				></div>
			</div>
		</div>

		<!-- Card -->
		<div class="w-full max-w-xl rounded-2xl border border-slate-700 bg-slate-800 shadow-2xl">
			<!-- Highlight text -->
			<div class="p-8">
				<div
					class="rounded-xl border-l-4 p-5 text-base leading-relaxed {COLOR_BG[current.Color] ?? COLOR_BG.yellow}"
					style="border-left-color: {COLOR_MARK[current.Color] ?? COLOR_MARK.yellow}"
				>
					"{current.SelectedText}"
				</div>

				<!-- Note -->
				{#if current.Note}
					<div class="mt-4 rounded-lg bg-slate-700/50 px-4 py-3 text-sm text-surface-300 italic">
						{current.Note}
					</div>
				{/if}
			</div>

			<!-- Source + actions -->
			<div class="flex items-center justify-between border-t border-slate-700 px-6 py-4">
				<!-- Source title -->
				<button
					on:click={openSource}
					class="flex items-center gap-2 text-left text-xs text-surface-500 transition-colors hover:text-blue-400"
					title="Open source article"
				>
					<BookOpen size={13} class="flex-shrink-0" />
					<span class="line-clamp-1 max-w-xs">{current.ItemTitle || 'Unknown source'}</span>
				</button>

				<!-- Nav + complete -->
				<div class="flex items-center gap-2">
					<button
						on:click={prev}
						disabled={index === 0}
						class="rounded p-1.5 text-surface-500 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Previous"
					><ChevronLeft size={16} /></button>

					<button
						on:click={markAndAdvance}
						disabled={completing}
						class="flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-500 disabled:opacity-50"
					>
						<Check size={14} />
						{index < queue.length - 1 ? 'Got it' : 'Finish'}
					</button>

					<button
						on:click={() => { if (index < queue.length - 1) index++; }}
						disabled={index >= queue.length - 1}
						class="rounded p-1.5 text-surface-500 transition-colors hover:bg-slate-700 hover:text-surface-100 disabled:opacity-30"
						title="Skip"
					><ChevronRight size={16} /></button>
				</div>
			</div>
		</div>

		<!-- Review count hint -->
		<p class="mt-6 text-xs text-surface-600">
			{#if current.ReviewCount === 0}
				First time reviewing this highlight
			{:else}
				Reviewed {current.ReviewCount} time{current.ReviewCount === 1 ? '' : 's'} before
			{/if}
		</p>
	{/if}
</div>
