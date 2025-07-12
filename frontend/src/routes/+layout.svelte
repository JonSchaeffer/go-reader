<script>
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { FeedService } from '$lib/services/feedService.js';
	import { CategoryService } from '$lib/services/categoryService.js';
	import { feeds, categories } from '$lib/stores.js';
	import '../app.css';

	// Theme management
	let theme = 'light';
	let sidebarOpen = false;
	let collapsedCategories = new Set();

	onMount(async () => {
		// Load theme from localStorage (theme should already be applied by app.html script)
		if (browser) {
			const savedTheme = localStorage.getItem('theme') || 'light';
			theme = savedTheme;
			// Ensure theme is properly applied even if app.html script failed
			document.documentElement.setAttribute('data-theme', savedTheme);
			
			// Load collapsed categories from localStorage
			const collapsed = localStorage.getItem('collapsedCategories');
			if (collapsed) {
				collapsedCategories = new Set(JSON.parse(collapsed));
			}

			// Re-apply theme when page becomes visible (handles navigation edge cases)
			document.addEventListener('visibilitychange', () => {
				if (!document.hidden) {
					const currentTheme = localStorage.getItem('theme') || 'light';
					if (currentTheme !== theme) {
						theme = currentTheme;
					}
					document.documentElement.setAttribute('data-theme', theme);
				}
			});
		}

		// Load feeds and categories data
		try {
			await Promise.all([
				FeedService.loadFeeds(),
				CategoryService.loadCategories()
			]);
		} catch (error) {
			console.error('Failed to load sidebar data:', error);
		}
	});

	// Re-apply theme whenever the component updates (e.g., after navigation)
	$: if (browser && theme) {
		document.documentElement.setAttribute('data-theme', theme);
	}

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		
		if (browser) {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('theme', theme);
		}
	}

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function toggleCategory(categoryId) {
		if (collapsedCategories.has(categoryId)) {
			collapsedCategories.delete(categoryId);
		} else {
			collapsedCategories.add(categoryId);
		}
		collapsedCategories = collapsedCategories; // Trigger reactivity
		
		// Save to localStorage
		if (browser) {
			localStorage.setItem('collapsedCategories', JSON.stringify([...collapsedCategories]));
		}
	}

	// Group feeds by category
	$: groupedFeeds = (() => {
		const grouped = {
			uncategorized: [],
			categories: {}
		};

		// Initialize categories
		$categories.forEach(category => {
			grouped.categories[category.ID] = {
				...category,
				feeds: []
			};
		});

		// Group feeds
		$feeds.forEach(feed => {
			if (feed.CategoryID && grouped.categories[feed.CategoryID]) {
				grouped.categories[feed.CategoryID].feeds.push(feed);
			} else {
				grouped.uncategorized.push(feed);
			}
		});

		return grouped;
	})();

	// Navigation items
	const navItems = [
		{
			section: 'Content',
			items: [
				{
					label: 'Feeds',
					href: '/feeds',
					icon: '📡',
					active: false
				},
				{
					label: 'Articles',
					href: '/articles',
					icon: '📄',
					active: false
				}
			]
		},
		{
			section: 'Organization',
			items: [
				{
					label: 'Categories',
					href: '/categories',
					icon: '📁',
					active: false
				},
				{
					label: 'Lists',
					href: '/lists',
					icon: '📋',
					disabled: true,
					badge: 'Soon'
				}
			]
		}
	];

	// Create computed nav items with active states
	$: navItemsWithActive = navItems.map(section => ({
		...section,
		items: section.items.map(item => ({
			...item,
			active: $page.url.pathname.startsWith(item.href) && item.href !== '/'
		}))
	}));

</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="app-container">
	<!-- Header -->
	<header class="header">
		<a href="/" class="header-title">
			RSS Reader
		</a>
		
		<div class="header-actions">
			<!-- Mobile menu toggle -->
			<button 
				class="btn-ghost md:hidden"
				on:click={toggleSidebar}
				aria-label="Toggle sidebar"
			>
				☰
			</button>
			
			<!-- Theme toggle -->
			<button 
				class="theme-toggle"
				on:click={toggleTheme}
				aria-label="Toggle theme"
				title="Toggle {theme === 'light' ? 'dark' : 'light'} mode"
			>
				{theme === 'light' ? '🌙' : '☀️'}
			</button>
		</div>
	</header>

	<!-- Sidebar -->
	<aside class="sidebar {sidebarOpen ? 'open' : ''}">
		<div class="sidebar-header">
			<h2 class="sidebar-title">Navigation</h2>
		</div>
		
		<nav class="sidebar-nav">
			{#each navItemsWithActive as section}
				<div class="nav-section">
					<h3 class="nav-section-title">{section.section}</h3>
					<ul class="nav-list">
						{#each section.items as item}
							<li class="nav-item">
								<a
									href={item.disabled ? '#' : item.href}
									class="nav-link {item.active ? 'active' : ''} {item.disabled ? 'disabled' : ''}"
									on:click={() => sidebarOpen = false}
								>
									<span class="nav-icon">{item.icon}</span>
									<span class="nav-text">{item.label}</span>
									{#if item.badge}
										<span class="nav-badge">{item.badge}</span>
									{/if}
								</a>
							</li>
						{/each}
					</ul>
				</div>
			{/each}

			<!-- Feeds by Category -->
			{#if $feeds.length > 0}
				<div class="nav-section">
					<h3 class="nav-section-title">
						<a href="/articles" class="feeds-title-link" on:click={() => sidebarOpen = false}>
							My Feeds
						</a>
					</h3>
					
					<!-- Categorized Feeds -->
					{#each Object.values(groupedFeeds.categories) as category}
						{#if category.feeds.length > 0}
							<div class="category-group">
								<div class="category-header">
									<a 
										href="/articles?category={category.ID}"
										class="category-link"
										on:click={() => sidebarOpen = false}
									>
										<span class="category-icon" style="background-color: {category.Color}">📁</span>
										<span class="category-name">{category.Name}</span>
									</a>
									<button 
										class="category-toggle {collapsedCategories.has(category.ID) ? 'collapsed' : ''}"
										on:click={() => toggleCategory(category.ID)}
										title="Toggle category"
									>
										{collapsedCategories.has(category.ID) ? '▶' : '▼'}
									</button>
								</div>
								
								{#if !collapsedCategories.has(category.ID)}
									<ul class="feed-list">
										{#each category.feeds as feed}
											<li class="feed-item">
												<a 
													href="/articles?feed={feed.ID}" 
													class="feed-link"
													on:click={() => sidebarOpen = false}
												>
													<span class="feed-icon">📡</span>
													<span class="feed-name">{feed.Title || 'Untitled Feed'}</span>
												</a>
											</li>
										{/each}
									</ul>
								{/if}
							</div>
						{/if}
					{/each}

					<!-- Uncategorized Feeds -->
					{#if groupedFeeds.uncategorized.length > 0}
						<div class="category-group">
							<div class="category-header">
								<a 
									href="/articles?category=null"
									class="category-link"
									on:click={() => sidebarOpen = false}
								>
									<span class="category-icon">📂</span>
									<span class="category-name">Uncategorized</span>
								</a>
								<button 
									class="category-toggle {collapsedCategories.has('uncategorized') ? 'collapsed' : ''}"
									on:click={() => toggleCategory('uncategorized')}
									title="Toggle category"
								>
									{collapsedCategories.has('uncategorized') ? '▶' : '▼'}
								</button>
							</div>
							
							{#if !collapsedCategories.has('uncategorized')}
								<ul class="feed-list">
									{#each groupedFeeds.uncategorized as feed}
										<li class="feed-item">
											<a 
												href="/articles?feed={feed.ID}" 
												class="feed-link"
												on:click={() => sidebarOpen = false}
											>
												<span class="feed-icon">📡</span>
												<span class="feed-name">{feed.Title || 'Untitled Feed'}</span>
											</a>
										</li>
									{/each}
								</ul>
							{/if}
						</div>
					{/if}
				</div>
			{/if}
		</nav>
	</aside>

	<!-- Main Content -->
	<main class="main-content">
		<slot />
	</main>
</div>

<!-- Mobile sidebar backdrop -->
{#if sidebarOpen}
	<button 
		class="fixed inset-0 bg-black bg-opacity-50 z-10 md:hidden backdrop-button"
		on:click={toggleSidebar}
		aria-label="Close sidebar"
	></button>
{/if}

<style>
	/* Additional responsive utilities */
	@media (min-width: 768px) {
		.md\:hidden {
			display: none;
		}
	}
	
	.fixed {
		position: fixed;
	}
	
	.inset-0 {
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
	}
	
	.bg-black {
		background-color: rgb(0 0 0);
	}
	
	.bg-opacity-50 {
		background-color: rgb(0 0 0 / 0.5);
	}
	
	.z-10 {
		z-index: 10;
	}

	/* Category Group Styles */
	.category-group {
		margin-bottom: 0.5rem;
	}

	.category-header {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.25rem 0.75rem;
		border-radius: var(--radius);
		transition: all 0.15s ease;
	}

	.category-link {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.25rem 0;
		color: var(--text-secondary);
		font-size: 0.875rem;
		font-weight: 500;
		text-decoration: none;
		flex-grow: 1;
		transition: color 0.15s ease;
	}

	.category-link:hover {
		color: var(--text-primary);
	}

	.category-header:hover {
		background: var(--bg-tertiary);
	}

	.category-icon {
		width: 1.25rem;
		height: 1.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
		font-size: 0.75rem;
		flex-shrink: 0;
	}

	.category-name {
		flex: 1;
		text-align: left;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.category-toggle {
		background: none;
		border: none;
		color: var(--text-tertiary);
		font-size: 0.75rem;
		padding: 0.25rem;
		cursor: pointer;
		transition: transform 0.15s ease, color 0.15s ease;
		flex-shrink: 0;
		border-radius: var(--radius-sm);
	}

	.category-toggle:hover {
		color: var(--text-primary);
		background: var(--bg-secondary);
	}

	.category-toggle.collapsed {
		transform: rotate(-90deg);
	}

	/* Feed List Styles */
	.feed-list {
		margin: 0.25rem 0 0 0;
		padding: 0;
		list-style: none;
	}

	.feed-item {
		margin: 0;
	}

	.feed-link {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.75rem 0.375rem 2rem;
		color: var(--text-tertiary);
		text-decoration: none;
		font-size: 0.8125rem;
		border-radius: var(--radius);
		transition: all 0.15s ease;
	}

	.feed-link:hover {
		background: var(--bg-tertiary);
		color: var(--text-primary);
	}

	.feed-icon {
		font-size: 0.75rem;
		flex-shrink: 0;
	}

	.feed-name {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* Active state for current feed */
	.feed-link.active {
		background: var(--primary-light);
		color: var(--primary);
	}

	/* Feeds title link */
	.feeds-title-link {
		color: inherit;
		text-decoration: none;
		transition: color 0.15s ease;
	}

	.feeds-title-link:hover {
		color: var(--text-primary);
	}
</style>