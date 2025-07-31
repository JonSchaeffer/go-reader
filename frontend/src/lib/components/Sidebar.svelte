<script>
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import { Settings } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { CategoryService } from '$lib/services/categoryService';
	
	let currentSection = 'articles'; // Track active section
	let categoriesExpanded = false; // Track if categories are expanded
	let categoriesWithFeeds = [];
	let expandedCategories = {}; // Track which categories are expanded
	let loading = true;

	onMount(async () => {
		try {
			categoriesWithFeeds = await CategoryService.loadCategoriesWithFeeds();
			loading = false;
		} catch (error) {
			console.error('Failed to load categories:', error);
			categoriesWithFeeds = [];
			loading = false;
		}
	});
</script>

<aside class="bg-surface-900 text-surface-100 flex w-64 flex-col">
	<AppBar>
		{#snippet headline()}
			<h2 class="h2 text-surface-100">RSS Reader</h2>
		{/snippet}
		{#snippet trail()}
			<Settings size={20} />
		{/snippet}
	</AppBar>

	<!-- Navigation Menu -->
	<nav class="flex-1 space-y-1 p-2">
		<!-- Articles Section -->
		<button
			class="flex w-full items-center gap-1 rounded-lg p-1 transition-colors {currentSection ===
			'articles'
				? 'bg-surface-800'
				: 'hover:bg-surface-800'}"
			on:click={() => (currentSection = 'articles')}
		>
			<span>Articles</span>
		</button>

		<!-- Categories Section -->
		<div>
			<button
				class="flex w-full items-center gap-1 rounded-lg p-1 transition-colors {currentSection ===
				'categories'
					? 'bg-surface-800'
					: 'hover:bg-surface-800'}"
				on:click={() => {
					currentSection = 'categories';
					categoriesExpanded = !categoriesExpanded;
				}}
			>
				<span class="flex-1 text-left">Categories</span>
				<span class="text-sm transition-transform {categoriesExpanded ? 'rotate-90' : ''}"
					>&gt;</span
				>
			</button>

			<!-- Collapsible Categories List -->
			{#if categoriesExpanded}
				<div class="mt-1 ml-3 space-y-0">
					{#if loading}
						<div class="text-surface-400 p-1 text-sm">Loading categories...</div>
					{:else if categoriesWithFeeds.length === 0}
						<div class="text-surface-400 p-1 text-sm">No categories found</div>
					{:else}
						{#each categoriesWithFeeds as category}
							<div>
								<!-- Category Button -->
								<button
									class="hover:bg-surface-800 flex w-full items-center rounded-sm p-1 text-sm transition-colors"
									style="color: {category.color}"
									on:click={() => {
										currentSection = `category-${category.id}`;
										expandedCategories[category.id] = !expandedCategories[category.id];
									}}
								>
									<span class="flex-1 text-left">{category.name}</span>
									<span
										class="text-xs transition-transform {expandedCategories[category.id]
											? 'rotate-90'
											: ''}">&gt;</span
									>
								</button>

								<!-- Feeds under this category -->
								{#if expandedCategories[category.id]}
									<div class="ml-3 space-y-0">
										{#if category.feeds && category.feeds.length > 0}
											{#each category.feeds as feed}
												<button
													class="hover:bg-surface-500 text-surface-300 flex w-full items-center rounded-sm p-1 text-xs transition-colors"
													on:click={() =>
														(currentSection = `feed-${feed.ID}`)}
												>
													<span>{feed.Title}</span>
												</button>
											{/each}
										{:else}
											<div class="text-surface-400 p-1 text-xs">No feeds in this category</div>
										{/if}
									</div>
								{/if}
							</div>
						{/each}
					{/if}
				</div>
			{/if}
		</div>

		<!-- Bookmarks Section -->
		<button
			class="flex w-full items-center gap-1 rounded-lg p-1 transition-colors {currentSection ===
			'bookmarks'
				? 'bg-surface-800'
				: 'hover:bg-surface-800'}"
			on:click={() => (currentSection = 'bookmarks')}
		>
			<span>Bookmarks</span>
		</button>

		<!-- Settings Section -->
		<button
			class="flex w-full items-center gap-1 rounded-lg p-1 transition-colors {currentSection ===
			'settings'
				? 'bg-surface-800'
				: 'hover:bg-surface-800'}"
			on:click={() => (currentSection = 'settings')}
		>
			<span>Settings</span>
		</button>
	</nav>
</aside>