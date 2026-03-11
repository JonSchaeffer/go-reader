<script>
	import { onMount, onDestroy, tick } from 'svelte';
	import { Maximize2, PanelRight, CheckCheck, RefreshCw, Search, X } from '@lucide/svelte';
	import { ArticleService } from '$lib/services/articleService';
	import { articleApi } from '$lib/api.js';
	import { articles, selectedFeedId, selectedArticle, visibleArticles, articleModalOpen, openMode } from '$lib/stores';

	let loading = true;
	let loadingMore = false;
	let markingRead = false;
	let filterMode = 'all';
	let firstLoad = true;
	let offset = 0;
	let hasMore = false;
	const PAGE_SIZE = 50;

	// Inline search
	let searchOpen = false;
	let searchQuery = '';
	let searchResults = [];
	let searching = false;
	let searchInputEl = null;
	let searchTimer = null;

	$: isSearching = searchOpen && searchQuery.trim().length > 0;

	async function openSearch() {
		searchOpen = true;
		await tick();
		searchInputEl?.focus();
	}

	function closeSearch() {
		searchOpen = false;
		searchQuery = '';
		searchResults = [];
	}

	function onSearchInput() {
		clearTimeout(searchTimer);
		const q = searchQuery.trim();
		if (!q) { searchResults = []; return; }
		searchTimer = setTimeout(async () => {
			searching = true;
			try {
				const res = await articleApi.search(q, 50);
				searchResults = res ?? [];
			} finally {
				searching = false;
			}
		}, 300);
	}

	function onSearchKeydown(e) {
		if (e.key === 'Escape') closeSearch();
	}

	$: currentArticles = $articles || [];

	$: filteredArticles = currentArticles.filter((article) => {
		if ($selectedArticle?.ID === article.ID) return true;
		if (filterMode === 'all') return true;
		if (filterMode === 'unread') return !article.Read;
		if (filterMode === 'read') return article.Read;
		return false;
	});

	$: visibleArticles.set(filteredArticles);

	function formatTime(dateString) {
		if (!dateString) return '';
		try {
			const date = new Date(dateString);
			const now = new Date();
			const diffInHours = Math.abs(now - date) / (1000 * 60 * 60);
			if (diffInHours < 24) {
				return date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true });
			}
			return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
		} catch {
			return dateString;
		}
	}

	function openArticle(article) {
		selectedArticle.set(article);
		if (!article.Read) ArticleService.toggleReadStatus(article.ID, article.Read);
		if ($openMode === 'modal') articleModalOpen.set(true);
	}

	function openInModal(article) {
		selectedArticle.set(article);
		if (!article.Read) ArticleService.toggleReadStatus(article.ID, article.Read);
		articleModalOpen.set(true);
	}

	function toggleOpenMode() {
		openMode.update((m) => (m === 'sidebar' ? 'modal' : 'sidebar'));
	}

	async function loadMore() {
		loadingMore = true;
		const nextOffset = offset + PAGE_SIZE;
		const result = await ArticleService.loadAllArticles(nextOffset, PAGE_SIZE);
		offset = nextOffset;
		hasMore = result?.hasMore ?? false;
		loadingMore = false;
	}

	async function markAllRead() {
		if (!$selectedFeedId) return;
		markingRead = true;
		try {
			await articleApi.markAllRead($selectedFeedId);
			// Update local state — mark all currently loaded articles for this feed as read
			articles.update((all) =>
				all.map((a) => (a.RssID === $selectedFeedId ? { ...a, Read: true } : a))
			);
		} finally {
			markingRead = false;
		}
	}

	// React to feed selection changes
	const unsubscribe = selectedFeedId.subscribe(async (feedId) => {
		if (firstLoad) return;
		filterMode = 'all';
		offset = 0;
		hasMore = false;
		loading = true;
		if (feedId !== null) {
			await ArticleService.loadArticlesByFeed(feedId);
		} else {
			const result = await ArticleService.loadAllArticles(0, PAGE_SIZE);
			hasMore = result?.hasMore ?? false;
		}
		loading = false;
	});

	onMount(async () => {
		const savedMode = localStorage.getItem('openMode');
		if (savedMode === 'sidebar' || savedMode === 'modal') openMode.set(savedMode);

		const result = await ArticleService.loadAllArticles(0, PAGE_SIZE);
		hasMore = result?.hasMore ?? false;
		loading = false;
		firstLoad = false;
	});

	onDestroy(unsubscribe);

	$: if (typeof window !== 'undefined') localStorage.setItem('openMode', $openMode);
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="flex items-center justify-between border-b border-white/10 px-4 py-3">
		<div class="flex items-baseline gap-4" class:invisible={isSearching}>
			<span class="text-sm font-semibold text-surface-100">
				{$selectedFeedId ? 'Feed' : 'All Articles'}
			</span>
			<div class="flex gap-3">
				{#each [['all', 'All'], ['unread', 'Unread'], ['read', 'Read']] as [mode, label]}
					<button
						class="border-b-2 pb-0.5 text-sm font-medium transition-colors {filterMode === mode
							? 'border-blue-400 text-blue-400'
							: 'border-transparent text-surface-400 hover:text-surface-100'}"
						on:click={() => (filterMode = mode)}
					>
						{label}
					</button>
				{/each}
			</div>
		</div>

		<div class="flex items-center gap-2">
			<!-- Inline search -->
			{#if searchOpen}
				<div class="flex items-center gap-1 rounded-md border border-slate-600 bg-slate-800 px-2 py-1">
					<Search size={13} class="flex-shrink-0 text-surface-500" />
					<input
						bind:this={searchInputEl}
						bind:value={searchQuery}
						on:input={onSearchInput}
						on:keydown={onSearchKeydown}
						placeholder="Search articles…"
						class="w-48 bg-transparent text-sm text-surface-100 placeholder-surface-500 outline-none"
					/>
					{#if searching}
						<RefreshCw size={12} class="flex-shrink-0 animate-spin text-surface-500" />
					{/if}
					<button on:click={closeSearch} class="flex-shrink-0 text-surface-500 hover:text-surface-300 transition-colors">
						<X size={13} />
					</button>
				</div>
			{:else}
				<span class="text-surface-500 text-xs">{filteredArticles.length} articles</span>
				<button
					on:click={openSearch}
					title="Search"
					class="rounded p-1 text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100"
				>
					<Search size={14} />
				</button>
			{/if}

			<!-- Mark all read (only when a specific feed is selected) -->
			{#if $selectedFeedId && !searchOpen}
				<button
					on:click={markAllRead}
					disabled={markingRead}
					title="Mark all as read"
					class="rounded p-1 text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100 disabled:opacity-40"
				>
					<CheckCheck size={14} />
				</button>
			{/if}

			<!-- Default open mode toggle -->
			{#if !searchOpen}
				<button
					on:click={toggleOpenMode}
					title="Default: open in {$openMode === 'sidebar' ? 'sidebar' : 'full screen'} — click to switch"
					class="rounded p-1 text-surface-500 transition-colors hover:bg-slate-800 hover:text-surface-100"
				>
					{#if $openMode === 'sidebar'}
						<PanelRight size={14} />
					{:else}
						<Maximize2 size={14} />
					{/if}
				</button>
			{/if}
		</div>
	</div>

	<!-- Article list -->
	<div class="flex-1 overflow-y-auto">
		{#if isSearching}
			{#if searching && searchResults.length === 0}
				<div class="text-surface-400 p-4 text-sm">Searching…</div>
			{:else if searchResults.length === 0}
				<div class="text-surface-400 p-4 text-sm">No results for "{searchQuery}"</div>
			{:else}
				<div class="border-b border-white/5 px-4 py-2 text-xs text-surface-500">{searchResults.length} result{searchResults.length === 1 ? '' : 's'}</div>
				{#each searchResults as article (article.ID)}
					<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
					<div class="group relative border-b border-white/10 last:border-b-0 transition-colors
						{$selectedArticle?.ID === article.ID ? 'bg-slate-700' : 'hover:bg-slate-800'}">
						<button class="w-full p-3 text-left" on:click={() => openArticle(article)}>
							<div class="flex items-start gap-2 pr-6">
								<div class="mt-1.5 h-2 w-2 flex-shrink-0">
									{#if !article.Read}
										<div class="h-2 w-2 rounded-full bg-blue-400"></div>
									{/if}
								</div>
								<div class="min-w-0 flex-1">
									<p class="text-surface-100 mb-0.5 line-clamp-2 text-sm font-medium leading-snug {article.Read ? 'opacity-60' : ''}">{article.Title}</p>
									{#if article.Description}
										<p class="text-surface-400 mb-1 line-clamp-2 text-xs">{article.Description?.replace(/<[^>]*>/g, '').slice(0, 120)}</p>
									{/if}
									<div class="text-surface-500 flex items-center gap-1.5 text-xs">
										{#if article.Author}<span>{article.Author}</span><span>·</span>{/if}
										<span>{formatTime(article.PublishDate)}</span>
									</div>
								</div>
							</div>
						</button>
						{#if $openMode === 'sidebar'}
							<button on:click={() => openInModal(article)} title="Open in full screen"
								class="absolute right-2 top-1/2 -translate-y-1/2 rounded p-1.5 opacity-0 group-hover:opacity-100 transition-opacity text-surface-500 hover:bg-slate-600 hover:text-surface-100">
								<Maximize2 size={13} />
							</button>
						{/if}
					</div>
				{/each}
			{/if}
		{:else if loading}
			<div class="text-surface-400 p-4 text-sm">Loading...</div>
		{:else if filteredArticles.length === 0}
			<div class="text-surface-400 p-4 text-sm">
				No {filterMode === 'all' ? '' : filterMode + ' '}articles.
			</div>
		{:else}
			{#each filteredArticles as article (article.ID)}
				<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
				<div
					class="group relative border-b border-white/10 last:border-b-0 transition-colors
						{$selectedArticle?.ID === article.ID ? 'bg-slate-700' : 'hover:bg-slate-800'}"
				>
					<button class="w-full p-3 text-left" on:click={() => openArticle(article)}>
						<div class="flex items-start gap-2 pr-6">
							<div class="mt-1.5 h-2 w-2 flex-shrink-0">
								{#if !article.Read}
									<div class="h-2 w-2 rounded-full bg-blue-400"></div>
								{/if}
							</div>
							<div class="min-w-0 flex-1">
								<p class="text-surface-100 mb-0.5 line-clamp-2 text-sm font-medium leading-snug {article.Read ? 'opacity-60' : ''}">
									{article.Title}
								</p>
								{#if article.Description}
									<p class="text-surface-400 mb-1 line-clamp-2 text-xs">
										{article.Description?.replace(/<[^>]*>/g, '').slice(0, 120)}
									</p>
								{/if}
								<div class="text-surface-500 flex items-center gap-1.5 text-xs">
									{#if article.Author}
										<span>{article.Author}</span>
										<span>·</span>
									{/if}
									<span>{formatTime(article.PublishDate)}</span>
									{#if article.Category && article.Category !== 'Uncategorized'}
										<span>·</span>
										<span>{article.Category}</span>
									{/if}
								</div>
							</div>
						</div>
					</button>

					{#if $openMode === 'sidebar'}
						<button
							on:click={() => openInModal(article)}
							title="Open in full screen"
							class="absolute right-2 top-1/2 -translate-y-1/2 rounded p-1.5
								opacity-0 group-hover:opacity-100 transition-opacity
								text-surface-500 hover:bg-slate-600 hover:text-surface-100"
						>
							<Maximize2 size={13} />
						</button>
					{/if}
				</div>
			{/each}

			<!-- Load more (only for all-articles view) -->
			{#if !$selectedFeedId && hasMore}
				<div class="flex justify-center p-4">
					<button
						on:click={loadMore}
						disabled={loadingMore}
						class="flex items-center gap-2 rounded-md px-4 py-2 text-sm text-surface-400 hover:bg-slate-800 hover:text-surface-100 disabled:opacity-40 transition-colors"
					>
						{#if loadingMore}
							<RefreshCw size={13} class="animate-spin" />
							Loading…
						{:else}
							Load more
						{/if}
					</button>
				</div>
			{/if}
		{/if}
	</div>
</div>
