<script>
	import { Settings, Rss, Tag, Bookmark, Highlighter, Search } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { CategoryService } from '$lib/services/categoryService';
	import { selectedFeedId, selectedArticle, currentView } from '$lib/stores';

	let categoriesWithFeeds = [];
	let expandedCategories = {};
	let loading = true;

	onMount(async () => {
		try {
			categoriesWithFeeds = await CategoryService.loadCategoriesWithFeeds();
		} catch (error) {
			console.error('Failed to load categories:', error);
			categoriesWithFeeds = [];
		} finally {
			loading = false;
		}
	});

	function selectFeed(feedId) {
		selectedFeedId.set(feedId);
		selectedArticle.set(null);
		currentView.set('feed');
	}

	function selectAll() {
		selectedFeedId.set(null);
		selectedArticle.set(null);
		currentView.set('feed');
	}
</script>

<aside class="text-surface-100 flex w-56 flex-col bg-slate-800 border-r border-slate-700">
	<!-- Header -->
	<header class="flex items-center justify-between bg-slate-800 p-4 border-b border-slate-700">
		<h2 class="text-surface-100 font-semibold text-lg">YARR</h2>
		<button
			on:click={() => currentView.set('categories-management')}
			class="text-surface-400 hover:text-surface-100 transition-colors"
			title="Manage categories"
		>
			<Settings size={16} />
		</button>
	</header>

	<!-- Navigation -->
	<nav class="flex-1 overflow-y-auto p-2 space-y-1">
		<!-- Search -->
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'search'
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={() => { currentView.set('search'); selectedArticle.set(null); }}
		>
			<Search size={14} />
			<span>Search</span>
		</button>

		<!-- All Articles -->
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'feed' && $selectedFeedId === null
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={selectAll}
		>
			<Rss size={14} />
			<span>All Articles</span>
		</button>

		<!-- Feeds by category -->
		{#if loading}
			<div class="text-surface-400 px-2 py-1 text-xs">Loading...</div>
		{:else}
			{#each categoriesWithFeeds as category}
				<div>
					<!-- Category header -->
					<button
						class="flex w-full items-center gap-1 rounded-md px-2 py-1 text-xs font-medium uppercase tracking-wide text-surface-400 hover:text-surface-200 transition-colors"
						on:click={() =>
							(expandedCategories[category.id] =
								expandedCategories[category.id] !== false ? false : true)}
					>
						<span
							class="inline-block h-2 w-2 rounded-full flex-shrink-0"
							style="background-color: {category.color}"
						></span>
						<span class="flex-1 text-left">{category.name}</span>
						<span class="text-xs transition-transform {expandedCategories[category.id] ? 'rotate-90' : ''}">&rsaquo;</span>
					</button>

					<!-- Feeds under category -->
					{#if expandedCategories[category.id] !== false}
						<div class="ml-3 mt-0.5 space-y-0.5">
							{#if category.feeds && category.feeds.length > 0}
								{#each category.feeds as feed}
									<button
										class="flex w-full items-center gap-1.5 rounded-md px-2 py-1 text-sm transition-colors
											{$selectedFeedId === feed.ID
											? 'bg-slate-700 text-surface-100'
											: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
										on:click={() => selectFeed(feed.ID)}
									>
										<span class="flex-1 truncate text-left">{feed.Title || feed.URL}</span>
										{#if feed.UnreadCount > 0}
											<span class="rounded-full bg-blue-500/20 px-1.5 py-0.5 text-[10px] font-medium text-blue-400">
												{feed.UnreadCount}
											</span>
										{/if}
									</button>
								{/each}
							{:else}
								<span class="text-surface-500 px-2 text-xs">No feeds</span>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		{/if}
	</nav>

	<!-- Bottom nav -->
	<div class="border-t border-slate-700 p-2 space-y-1">
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'library'
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={() => { currentView.set('library'); selectedArticle.set(null); }}
		>
			<Bookmark size={14} />
			<span>Library</span>
		</button>
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'highlights'
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={() => { currentView.set('highlights'); selectedArticle.set(null); }}
		>
			<Highlighter size={14} />
			<span>Highlights</span>
		</button>
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'feeds-management'
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={() => currentView.set('feeds-management')}
		>
			<Rss size={14} />
			<span>Manage Feeds</span>
		</button>
		<button
			class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors
				{$currentView === 'categories-management'
				? 'bg-slate-700 text-surface-100'
				: 'text-surface-300 hover:bg-slate-700 hover:text-surface-100'}"
			on:click={() => currentView.set('categories-management')}
		>
			<Tag size={14} />
			<span>Categories</span>
		</button>
	</div>
</aside>
