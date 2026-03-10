<script>
	import { Search, BookOpen, Bookmark, X } from '@lucide/svelte';
	import { searchApi, articleApi, libraryApi } from '$lib/api.js';
	import { selectedArticle, visibleArticles, articleModalOpen } from '$lib/stores';

	let query = '';
	let results = [];
	let loading = false;
	let searched = false;
	let debounceTimer;

	function onInput() {
		clearTimeout(debounceTimer);
		if (!query.trim()) { results = []; searched = false; return; }
		debounceTimer = setTimeout(runSearch, 300);
	}

	async function runSearch() {
		if (!query.trim()) return;
		loading = true;
		searched = true;
		try {
			const data = await searchApi.search(query.trim());
			results = Array.isArray(data) ? data : [];
		} finally {
			loading = false;
		}
	}

	function clear() {
		query = '';
		results = [];
		searched = false;
	}

	async function openResult(r) {
		let article;
		if (r.Type === 'library') {
			const item = await libraryApi.getById(r.ID);
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
			const res = await articleApi.getById(r.ID);
			const item = Array.isArray(res) ? res[0] : res;
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
		visibleArticles.set([article]);
		selectedArticle.set(article);
		articleModalOpen.set(true);
	}

	const articleCount = (r) => results.filter((x) => x.Type === 'article').length;
	const libraryCount = (r) => results.filter((x) => x.Type === 'library').length;
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header / search bar -->
	<div class="border-b border-white/10 px-4 py-3">
		<div class="flex items-center gap-2 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 focus-within:border-blue-400">
			<Search size={14} class="flex-shrink-0 text-surface-500" />
			<input
				bind:value={query}
				on:input={onInput}
				on:keydown={(e) => e.key === 'Enter' && runSearch()}
				placeholder="Search articles and library…"
				autofocus
				class="flex-1 bg-transparent text-sm text-surface-100 placeholder-surface-500 focus:outline-none"
			/>
			{#if query}
				<button on:click={clear} class="text-surface-500 hover:text-surface-100 transition-colors">
					<X size={13} />
				</button>
			{/if}
		</div>
	</div>

	<!-- Results -->
	<div class="flex-1 overflow-y-auto">
		{#if loading}
			<div class="p-4 text-sm text-surface-400">Searching…</div>
		{:else if searched && results.length === 0}
			<div class="flex flex-col items-center justify-center gap-3 py-24 text-center">
				<Search size={28} class="text-surface-600" />
				<p class="text-sm text-surface-500">No results for <span class="text-surface-300">"{query}"</span></p>
			</div>
		{:else if results.length > 0}
			<!-- Summary line -->
			<div class="border-b border-white/10 px-4 py-2 text-xs text-surface-500">
				{results.length} result{results.length === 1 ? '' : 's'}
				·
				{results.filter(r => r.Type === 'article').length} articles
				·
				{results.filter(r => r.Type === 'library').length} library
			</div>

			<div class="divide-y divide-white/10">
				{#each results as r (r.Type + r.ID)}
					<button
						on:click={() => openResult(r)}
						class="w-full px-4 py-3 text-left hover:bg-slate-800 transition-colors"
					>
						<div class="flex items-start gap-2">
							<!-- Type icon -->
							<div class="mt-0.5 flex-shrink-0 text-surface-500">
								{#if r.Type === 'library'}
									<Bookmark size={13} />
								{:else}
									<BookOpen size={13} />
								{/if}
							</div>

							<div class="min-w-0 flex-1">
								<p class="mb-0.5 line-clamp-1 text-sm font-medium text-surface-100 {r.Read ? 'opacity-60' : ''}">
									{r.Title || 'Untitled'}
								</p>
								{#if r.Excerpt}
									<p class="mb-1 line-clamp-2 text-xs text-surface-400">
										{r.Excerpt}
									</p>
								{/if}
								<div class="flex items-center gap-1.5 text-xs text-surface-600">
									<span class="capitalize">{r.Type === 'library' ? 'Library' : 'Article'}</span>
									{#if r.Author}
										<span>·</span>
										<span>{r.Author}</span>
									{/if}
								</div>
							</div>
						</div>
					</button>
				{/each}
			</div>
		{:else}
			<div class="flex flex-col items-center justify-center gap-3 py-24 text-center">
				<Search size={28} class="text-surface-600" />
				<p class="text-sm text-surface-500">Search across all your articles and library items.</p>
			</div>
		{/if}
	</div>
</div>
