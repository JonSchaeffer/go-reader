<script>
	import { onMount } from 'svelte';
	import { Plus, RefreshCw, Archive, Maximize2, PanelRight, Bookmark } from '@lucide/svelte';
	import { LibraryService } from '$lib/services/libraryService';
	import { savedItems, selectedArticle, articleModalOpen, openMode } from '$lib/stores';

	let loading = true;
	let loadingMore = false;
	let saving = false;
	let saveError = '';
	let saveUrl = '';
	let showSaveForm = false;
	let filterMode = 'unread'; // 'all' | 'unread' | 'archived'
	let hasMore = false;
	let offset = 0;
	const PAGE_SIZE = 50;

	$: filteredItems = ($savedItems ?? []).filter((item) => {
		if ($selectedArticle?._type === 'library' && $selectedArticle?.ID === item.ID) return true;
		if (filterMode === 'all') return !item.Archived;
		if (filterMode === 'unread') return !item.Read && !item.Archived;
		if (filterMode === 'archived') return item.Archived;
		return false;
	});

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

	function openItem(item) {
		// Normalise saved item into the same shape as an article so existing viewers work
		selectedArticle.set({
			ID: item.ID,
			Title: item.Title,
			Description: item.Content,
			Author: item.Author,
			PublishDate: item.PublishedDate,
			Read: item.Read,
			Link: item.URL,
			Category: null,
			_type: 'library',
			_archived: item.Archived
		});
		if ($openMode === 'modal') articleModalOpen.set(true);
	}

	function openItemInModal(item) {
		selectedArticle.set({
			ID: item.ID,
			Title: item.Title,
			Description: item.Content,
			Author: item.Author,
			PublishDate: item.PublishedDate,
			Read: item.Read,
			Link: item.URL,
			Category: null,
			_type: 'library',
			_archived: item.Archived
		});
		articleModalOpen.set(true);
	}

	async function saveFromUrl() {
		if (!saveUrl.trim()) return;
		saving = true;
		saveError = '';
		try {
			await LibraryService.saveURL(saveUrl.trim());
			saveUrl = '';
			showSaveForm = false;
			// Reload to show the new item
			const result = await LibraryService.loadItems(0, PAGE_SIZE, filterMode === 'archived');
			hasMore = result.hasMore;
			offset = 0;
		} catch (e) {
			saveError = e.message?.includes('already') ? 'Already in your library.' : 'Failed to save. Check the URL and try again.';
		} finally {
			saving = false;
		}
	}

	async function archiveItem(item) {
		await LibraryService.archive(item.ID, item.Read);
	}

	async function loadMore() {
		loadingMore = true;
		const nextOffset = offset + PAGE_SIZE;
		const result = await LibraryService.loadItems(nextOffset, PAGE_SIZE, filterMode === 'archived');
		offset = nextOffset;
		hasMore = result.hasMore;
		loadingMore = false;
	}

	async function switchFilter(mode) {
		filterMode = mode;
		loading = true;
		offset = 0;
		const result = await LibraryService.loadItems(0, PAGE_SIZE, mode === 'archived');
		hasMore = result.hasMore;
		loading = false;
	}

	onMount(async () => {
		const result = await LibraryService.loadItems(0, PAGE_SIZE, false);
		hasMore = result.hasMore;
		loading = false;
	});
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="flex items-baseline justify-between border-b border-white/10 px-4 py-3">
		<div class="flex items-baseline gap-4">
			<span class="text-sm font-semibold text-surface-100">Library</span>
			<div class="flex gap-3">
				{#each [['all', 'All'], ['unread', 'Unread'], ['archived', 'Archived']] as [mode, label]}
					<button
						class="border-b-2 pb-0.5 text-sm font-medium transition-colors {filterMode === mode
							? 'border-blue-400 text-blue-400'
							: 'border-transparent text-surface-400 hover:text-surface-100'}"
						on:click={() => switchFilter(mode)}
					>
						{label}
					</button>
				{/each}
			</div>
		</div>

		<div class="flex items-center gap-2">
			<span class="text-surface-500 text-xs">{filteredItems.length} items</span>
			<button
				on:click={() => { showSaveForm = !showSaveForm; saveError = ''; }}
				class="flex items-center gap-1 rounded p-1 text-sm transition-colors
					{showSaveForm ? 'text-surface-100' : 'text-surface-500 hover:bg-slate-800 hover:text-surface-100'}"
				title="Save a URL"
			>
				<Plus size={15} />
			</button>
		</div>
	</div>

	<!-- Save URL form -->
	{#if showSaveForm}
		<div class="border-b border-white/10 bg-slate-800/50 px-4 py-3">
			<form on:submit|preventDefault={saveFromUrl} class="flex gap-2">
				<input
					bind:value={saveUrl}
					placeholder="https://example.com/article"
					class="flex-1 rounded-md border border-slate-600 bg-slate-800 px-3 py-1.5 text-sm text-surface-100 placeholder-surface-500 focus:border-blue-400 focus:outline-none"
					disabled={saving}
					autofocus
				/>
				<button
					type="submit"
					disabled={saving}
					class="flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-500 disabled:opacity-50 transition-colors"
				>
					{#if saving}
						<RefreshCw size={12} class="animate-spin" />
					{:else}
						Save
					{/if}
				</button>
			</form>
			{#if saveError}
				<p class="mt-1.5 text-xs text-red-400">{saveError}</p>
			{/if}
		</div>
	{/if}

	<!-- Item list -->
	<div class="flex-1 overflow-y-auto">
		{#if loading}
			<div class="text-surface-400 p-4 text-sm">Loading...</div>
		{:else if filteredItems.length === 0}
			<div class="flex flex-col items-center justify-center gap-3 py-24 text-center">
				<Bookmark size={28} class="text-surface-600" />
				<p class="text-sm text-surface-500">
					{filterMode === 'archived' ? 'Nothing archived yet.' : 'Nothing saved yet. Add a URL or save articles from your feed.'}
				</p>
			</div>
		{:else}
			{#each filteredItems as item (item.ID)}
				<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
				<div
					class="group relative border-b border-white/10 last:border-b-0 transition-colors
						{$selectedArticle?.ID === item.ID && $selectedArticle?._type === 'library'
						? 'bg-slate-700'
						: 'hover:bg-slate-800'}"
				>
					<button class="w-full p-3 text-left" on:click={() => openItem(item)}>
						<div class="flex items-start gap-2 pr-14">
							<!-- Unread dot -->
							<div class="mt-1.5 h-2 w-2 flex-shrink-0">
								{#if !item.Read}
									<div class="h-2 w-2 rounded-full bg-blue-400"></div>
								{/if}
							</div>
							<div class="min-w-0 flex-1">
								<p class="text-surface-100 mb-0.5 line-clamp-2 text-sm font-medium leading-snug {item.Read ? 'opacity-60' : ''}">
									{item.Title || item.URL}
								</p>
								{#if item.Content}
									<p class="text-surface-400 mb-1 line-clamp-2 text-xs">
										{item.Content.replace(/<[^>]*>/g, '').slice(0, 120)}
									</p>
								{/if}
								<div class="text-surface-500 flex items-center gap-1.5 text-xs">
									{#if item.Author}
										<span>{item.Author}</span>
										<span>·</span>
									{/if}
									<span>{formatTime(item.SavedAt)}</span>
									<span>·</span>
									<span class="capitalize">{item.SourceType === 'rss_article' ? 'RSS' : 'URL'}</span>
								</div>
							</div>
						</div>
					</button>

					<!-- Row actions (on hover) -->
					<div class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
						{#if $openMode === 'sidebar'}
							<button
								on:click={() => openItemInModal(item)}
								title="Open in full screen"
								class="rounded p-1.5 text-surface-500 hover:bg-slate-600 hover:text-surface-100"
							>
								<Maximize2 size={13} />
							</button>
						{/if}
						{#if filterMode !== 'archived'}
							<button
								on:click={() => archiveItem(item)}
								title="Archive"
								class="rounded p-1.5 text-surface-500 hover:bg-slate-600 hover:text-surface-100"
							>
								<Archive size={13} />
							</button>
						{/if}
					</div>
				</div>
			{/each}

			{#if hasMore}
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
