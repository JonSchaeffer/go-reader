<script>
	import { onMount } from 'svelte';
	import { ArticleService } from '$lib/services/articleService';
	import { articles } from '$lib/stores';
	import { PanelRightClose, PanelLeftClose } from '@lucide/svelte';
	import { AppBar } from '@skeletonlabs/skeleton-svelte';

	let feed = [];
	let loading = true;
	let filterMode = 'unread'; // 'unread', 'read'
	export let showRightSidebar = false;

	// Save filter mode to localStorage when it changes
	$: if (typeof window !== 'undefined') {
		localStorage.setItem('feedFilterMode', filterMode);
	}

	// Use articles from the store instead of local feed state
	$: currentArticles = $articles || [];

	// Filter articles based on read status
	$: filteredArticles = currentArticles.filter((article) => {
		if (filterMode === 'unread') return !article.Read;
		if (filterMode === 'read') return article.Read;
		return false;
	});

	// Format date to human readable time
	function formatTime(dateString) {
		if (!dateString) return '';

		try {
			const date = new Date(dateString);
			const now = new Date();
			const diffInHours = Math.abs(now - date) / (1000 * 60 * 60);

			// If within last 24 hours, show time only
			if (diffInHours < 24) {
				return date.toLocaleTimeString('en-US', {
					hour: 'numeric',
					minute: '2-digit',
					hour12: true
				});
			}

			// If older, show date
			return date.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric'
			});
		} catch (error) {
			return dateString;
		}
	}

	onMount(async () => {
		// Load saved filter mode from localStorage
		if (typeof window !== 'undefined') {
			const savedFilter = localStorage.getItem('feedFilterMode');
			if (savedFilter && (savedFilter === 'unread' || savedFilter === 'read')) {
				filterMode = savedFilter;
			}
		}

		// Load articles
		try {
			await ArticleService.loadAllArticles();
			loading = false;
		} catch (error) {
			console.error('Failed to load feed:', error);
			loading = false;
		}
	});
</script>

<div class="flex-1 space-y-1 bg-slate-900 p-3">
	<div class="spacey-1 bg-slate-900">
		<AppBar background="bg-slate-900">
			{#snippet headline()}
				<div class="flex w-full items-center justify-between">
					<!-- Title and filter buttons -->
					<div class="flex items-center gap-6">
						<h3 class="h3 text-surface-100">Feed</h3>

						<!-- Filter buttons -->
						<div class="flex gap-2">
							<button
								class="border-b-2 text-lg font-medium transition-colors {filterMode === 'unread'
									? 'text-primary-100 border-primary-100'
									: 'text-surface-300 hover:text-surface-100 border-transparent'}"
								on:click={() => (filterMode = 'unread')}
							>
								Unread
							</button>
							<button
								class="border-b-2 text-lg font-medium transition-colors {filterMode === 'read'
									? 'text-primary-100 border-primary-100'
									: 'text-surface-300 hover:text-surface-100 border-transparent'}"
								on:click={() => (filterMode = 'read')}
							>
								Read
							</button>
						</div>
					</div>

					<!-- Sidebar toggle -->
					<button
						class="btn btn-sm variant-ghost-surface"
						on:click={() => (showRightSidebar = !showRightSidebar)}
						title="Toggle sidebar"
					>
						{#if showRightSidebar}
							<PanelRightClose size="20" />
						{:else}
							<PanelLeftClose size="20" />
						{/if}
					</button>
				</div>
			{/snippet}
		</AppBar>
	</div>
	<div class="border-t border-white/20 pt-3">
		<div class="space-y-1">
			{#each filteredArticles as article}
				<div
					class="card preset-filled-surface-100-600 cursor-pointer border-b border-white/10 transition-colors last:border-b-0 hover:bg-white/10"
				>
					<article class="flex items-start p-3">
						<!-- Notification dot space (always present) -->
						<div class="mt-2 mr-3 flex h-2 w-2 flex-shrink-0 items-center justify-center">
							{#if !article.Read}
								<div class="h-2 w-2 rounded-full bg-blue-400"></div>
							{/if}
						</div>

						<!-- Content -->
						<div class="min-w-0 flex-1">
							<!-- Title -->
							<header>
								<h5 class="mb-1 line-clamp-2 text-sm leading-tight font-medium">
									{article.Title}
								</h5>
							</header>

							<!-- Description -->
							{#if article.description}
								<p class="text-surface-600-300 mb-2 line-clamp-2 text-sm">
									{article.description}
								</p>
							{/if}

							<!-- Source Info -->
							<footer class="text-surface-400-500 flex items-center gap-2 text-xs">
								<span class="text-surface-400-500">{article.Author}</span>
								<span>•</span>
								<span>{formatTime(article.PublishDate)}</span>
							</footer>
						</div>

						<!-- Right Side Info -->
						<div class="flex-shrink-0 text-right">
							<small class="text-surface-400-500 text-xs">
								{#if article.Category && article.Category !== 'Uncategorized'}
									{article.Category} •
								{/if}
								{formatTime(article.PublishDate)}
							</small>
						</div>
					</article>
				</div>
			{/each}
		</div>
	</div>
</div>
