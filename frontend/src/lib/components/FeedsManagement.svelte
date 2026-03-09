<script>
	import { onMount } from 'svelte';
	import { Trash2, Plus, RefreshCw, Rss } from '@lucide/svelte';
	import { FeedService } from '$lib/services/feedService';
	import { CategoryService } from '$lib/services/categoryService';
	import { rssApi } from '$lib/api.js';
	import { feeds, categories } from '$lib/stores';

	let newFeedUrl = '';
	let adding = false;
	let addError = '';
	let showAddForm = false;
	let openPickerId = null;
	let refreshingId = null;

	onMount(async () => {
		await Promise.all([FeedService.loadFeeds(), CategoryService.loadCategories()]);
	});

	async function addFeed() {
		if (!newFeedUrl.trim()) return;
		adding = true;
		addError = '';
		try {
			await FeedService.addFeed(newFeedUrl.trim());
			newFeedUrl = '';
			showAddForm = false;
		} catch {
			addError = 'Failed to add feed. Check the URL and try again.';
		} finally {
			adding = false;
		}
	}

	async function deleteFeed(id) {
		if (!confirm('Delete this feed and all its articles?')) return;
		await FeedService.deleteFeed(id);
	}

	async function refreshFeed(id) {
		refreshingId = id;
		try {
			await rssApi.refresh(id);
			await FeedService.loadFeeds();
		} finally {
			refreshingId = null;
		}
	}

	async function assignCategory(feedId, categoryId) {
		openPickerId = null;
		await rssApi.assignCategory(feedId, categoryId || null);
		await FeedService.loadFeeds();
	}

	function getCategoryForFeed(feed) {
		return $categories.find((c) => c.ID === feed.CategoryID) ?? null;
	}

	function handleWindowClick(e) {
		if (!e.target.closest('[data-picker]')) openPickerId = null;
	}
</script>

<svelte:window on:click={handleWindowClick} />

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="flex items-center justify-between border-b border-white/10 px-8 py-5">
		<div>
			<h2 class="text-base font-semibold text-surface-100">Feeds</h2>
			<p class="mt-0.5 text-xs text-surface-500">{$feeds.length} feed{$feeds.length === 1 ? '' : 's'} subscribed</p>
		</div>
		<button
			on:click={() => { showAddForm = !showAddForm; addError = ''; }}
			class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors
				{showAddForm ? 'bg-slate-700 text-surface-100' : 'bg-blue-600 text-white hover:bg-blue-500'}"
		>
			<Plus size={14} />
			Add Feed
		</button>
	</div>

	<!-- Add form (expandable) -->
	{#if showAddForm}
		<div class="border-b border-white/10 bg-slate-800/50 px-8 py-4">
			<form on:submit|preventDefault={addFeed} class="flex gap-2">
				<input
					bind:value={newFeedUrl}
					placeholder="https://example.com/feed.xml"
					class="flex-1 rounded-md border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-surface-100 placeholder-surface-500 focus:border-blue-400 focus:outline-none"
					disabled={adding}
					autofocus
				/>
				<button
					type="submit"
					class="flex items-center gap-1.5 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-500 disabled:opacity-50 transition-colors"
					disabled={adding}
				>
					{#if adding}
						<RefreshCw size={13} class="animate-spin" />
						Adding…
					{:else}
						Subscribe
					{/if}
				</button>
			</form>
			{#if addError}
				<p class="mt-2 text-xs text-red-400">{addError}</p>
			{/if}
		</div>
	{/if}

	<!-- Feed list -->
	<div class="flex-1 overflow-y-auto">
		{#if $feeds.length === 0}
			<div class="flex flex-col items-center justify-center gap-3 py-24 text-center">
				<Rss size={32} class="text-surface-600" />
				<p class="text-sm text-surface-500">No feeds yet. Add one above to get started.</p>
			</div>
		{:else}
			<ul>
				{#each $feeds as feed (feed.ID)}
					{@const cat = getCategoryForFeed(feed)}
					<li class="group flex items-center gap-4 border-b border-white/5 px-8 py-4 last:border-b-0 hover:bg-slate-800/40">
						<!-- Icon -->
						<div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md bg-slate-700">
							<Rss size={14} class="text-surface-400" />
						</div>

						<!-- Text -->
						<div class="min-w-0 flex-1">
							<p class="truncate text-sm font-medium text-surface-100">
								{feed.Title || feed.URL}
							</p>
							<p class="mt-0.5 truncate text-xs text-surface-500">{feed.URL}</p>
						</div>

						<!-- Category picker -->
						<div class="relative flex-shrink-0" data-picker>
							<button
								data-picker
								on:click|stopPropagation={() =>
									(openPickerId = openPickerId === feed.ID ? null : feed.ID)}
								class="flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs transition-colors hover:bg-slate-700"
							>
								{#if cat}
									<span
										class="h-2 w-2 rounded-full flex-shrink-0"
										style="background-color: {cat.Color}"
									></span>
									<span class="text-surface-300">{cat.Name}</span>
								{:else}
									<span class="text-surface-500">Uncategorized</span>
								{/if}
							</button>

							{#if openPickerId === feed.ID}
								<div
									data-picker
									class="absolute right-0 top-full z-10 mt-1 min-w-[160px] overflow-hidden rounded-lg border border-slate-600 bg-slate-800 shadow-xl"
								>
									<button
										data-picker
										on:click|stopPropagation={() => assignCategory(feed.ID, null)}
										class="flex w-full items-center gap-2 px-3 py-2 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100"
									>
										<span class="h-2 w-2 rounded-full bg-slate-500"></span>
										Uncategorized
									</button>
									{#each $categories as c (c.ID)}
										<button
											data-picker
											on:click|stopPropagation={() => assignCategory(feed.ID, c.ID)}
											class="flex w-full items-center gap-2 px-3 py-2 text-xs text-surface-200 hover:bg-slate-700"
										>
											<span
												class="h-2 w-2 rounded-full flex-shrink-0"
												style="background-color: {c.Color}"
											></span>
											{c.Name}
										</button>
									{/each}
								</div>
							{/if}
						</div>

						<!-- Delete -->
						<button
							on:click={() => refreshFeed(feed.ID)}
							disabled={refreshingId === feed.ID}
							class="flex-shrink-0 rounded p-1.5 text-surface-600 opacity-0 transition-all group-hover:opacity-100 hover:bg-slate-700 hover:text-surface-100 disabled:opacity-40"
							title="Refresh feed"
						>
							<RefreshCw size={14} class={refreshingId === feed.ID ? 'animate-spin' : ''} />
						</button>
						<button
							on:click={() => deleteFeed(feed.ID)}
							class="flex-shrink-0 rounded p-1.5 text-surface-600 opacity-0 transition-all group-hover:opacity-100 hover:bg-slate-700 hover:text-red-400"
							title="Remove feed"
						>
							<Trash2 size={14} />
						</button>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</div>
