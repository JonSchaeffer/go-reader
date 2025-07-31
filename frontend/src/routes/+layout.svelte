<script>
	import '../app.css';
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import { Settings } from '@lucide/svelte';
	let currentSection = 'articles'; // Track active section
	let categoriesExpanded = false; // Track if categories are expanded

	// Dummy categories with feeds for now
	let categoriesWithFeeds = {
		Technology: ['TechCrunch', 'Ars Technica', 'Hacker News'],
		News: ['BBC News', 'Reuters', 'Associated Press'],
		Gaming: ['IGN', 'GameSpot', 'Polygon'],
		Science: ['NASA News', 'Scientific American', 'Nature']
	};
	let expandedCategories = {}; // Track which categories are expanded
</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="bg-primary-900 grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Grid Column -->
	<div class="grid grid-cols-[auto_1fr_auto]">
		<!-- Sidebar (Left) -->
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
							{#each Object.entries(categoriesWithFeeds) as [category, feeds]}
								<div>
									<!-- Category Button -->
									<button
										class="hover:bg-surface-800 flex w-full items-center rounded-sm p-1 text-sm transition-colors"
										on:click={() => {
											currentSection = `category-${category.toLowerCase()}`;
											expandedCategories[category] = !expandedCategories[category];
										}}
									>
										<span class="flex-1 text-left">{category}</span>
										<span
											class="text-xs transition-transform {expandedCategories[category]
												? 'rotate-90'
												: ''}">&gt;</span
										>
									</button>

									<!-- Feeds under this category -->
									{#if expandedCategories[category]}
										<div class="ml-3 space-y-0">
											{#each feeds as feed}
												<button
													class="hover:bg-surface-500 text-surface-300 flex w-full items-center rounded-sm p-1 text-xs transition-colors"
													on:click={() =>
														(currentSection = `feed-${feed.toLowerCase().replace(/\s+/g, '-')}`)}
												>
													<span>{feed}</span>
												</button>
											{/each}
										</div>
									{/if}
								</div>
							{/each}
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
		<!-- Main -->
		<main class="bg-surface-900 space-y-4 p-4">
			<slot />
			<!-- Placeholder content for testing -->
			<p class="bg-surface-800 text-surface-50 h-[312px] p-4">Paragraph 1</p>
			<p class="bg-surface-800 text-surface-50 h-[312px] p-4">Paragraph 2</p>
		</main>
		<!-- Sidebar (Right) -->
		<aside class="bg-surface-900 text-surface-100 p-4">(sidebar)</aside>
	</div>
	<!-- Footer -->
	<footer class="bg-surface-900 text-surface-100 p-4">(footer)</footer>
</div>
