<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { afterNavigate } from '$app/navigation';
	import { ArticleService } from '$lib/services/articleService.js';
	import { articles, loading, errors } from '$lib/stores.js';

	let feedId = null;
	let searchTerm = '';

	// Function to decode HTML entities
	function decodeHtml(html) {
		if (!html || typeof window === 'undefined') return html;
		
		const temp = document.createElement('div');
		temp.innerHTML = html;
		return temp.textContent || temp.innerText || html;
	}

	// Function to safely truncate HTML content
	function truncateHtml(html, maxLength = 300) {
		if (!html) return '';
		
		// If the HTML is short enough, return as-is
		if (html.length <= maxLength) {
			return html;
		}
		
		// Truncate and add ellipsis
		const truncated = html.substring(0, maxLength);
		const lastSpace = truncated.lastIndexOf(' ');
		const cutPoint = lastSpace > maxLength * 0.8 ? lastSpace : maxLength;
		
		return html.substring(0, cutPoint) + '...';
	}

	// Function to format dates properly
	function formatDate(dateString) {
		if (!dateString) return 'No date';
		
		try {
			// Handle different date formats
			let date;
			
			// Try parsing as-is first
			date = new Date(dateString);
			
			// If invalid, try parsing different formats
			if (isNaN(date.getTime()) && typeof dateString === 'string') {
				// Handle Go time format (RFC3339)
				if (dateString.includes('T') && dateString.includes('Z')) {
					date = new Date(dateString);
				}
				// Handle other common formats
				else if (dateString.includes('-')) {
					date = new Date(dateString.replace(' ', 'T'));
				}
			}
			
			// Check if date is valid
			if (isNaN(date.getTime())) {
				console.warn('Invalid date format:', dateString);
				return 'Invalid date';
			}
			
			return date.toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'short',
				day: 'numeric'
			});
		} catch (error) {
			console.warn('Date parsing error:', error, dateString);
			return 'Invalid date';
		}
	}

	async function loadData() {
		console.log('Loading articles data...');
		try {
			// Check if we have a feed filter from URL params
			feedId = $page.url.searchParams.get('feed');
			console.log('Feed ID from URL:', feedId);
			
			if (feedId) {
				await ArticleService.loadArticlesByFeed(feedId);
			} else {
				await ArticleService.loadAllArticles();
			}
			console.log('Articles loaded successfully, count:', $articles.length);
		} catch (error) {
			console.error('Failed to load articles:', error);
		}
	}

	onMount(() => {
		console.log('Articles page mounted');
		loadData();
	});

	afterNavigate(() => {
		console.log('Navigated to articles page');
		loadData();
	});




	// Filter articles based on search term
	$: filteredArticles = $articles.filter(article => 
		!searchTerm || 
		decodeHtml(article.Title)?.toLowerCase().includes(searchTerm.toLowerCase()) ||
		decodeHtml(article.Description)?.toLowerCase().includes(searchTerm.toLowerCase())
	);
</script>

<svelte:head>
	<title>Articles - RSS Reader</title>
</svelte:head>

<div class="content-header">
	<div>
		<h1 style="font-size: 1.875rem; font-weight: 700; color: var(--text-primary);">
			Articles {feedId ? '(Filtered)' : ''}
		</h1>
		<p style="color: var(--text-secondary); margin-top: 0.5rem;">
			{#if feedId}
				Showing articles from selected feed ({$articles.length} articles)
			{:else}
				Browse all articles from your RSS feeds ({$articles.length} articles)
			{/if}
		</p>
	</div>
</div>

<div class="content-body">
	<!-- Search Bar -->
	{#if $articles.length > 0}
		<div class="search-bar">
			<div class="search-input-wrapper">
				<span class="search-icon">🔍</span>
				<input 
					type="text" 
					bind:value={searchTerm}
					placeholder="Search articles..."
					class="search-input"
				/>
				{#if searchTerm}
					<button class="btn-ghost search-clear" on:click={() => searchTerm = ''}>✕</button>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Loading State -->
	{#if $loading.articles}
		<div class="loading-state">
			<div class="loading-spinner">⏳</div>
			<p>Loading articles...</p>
		</div>

	<!-- Error State -->
	{:else if $errors.articles}
		<div class="error-state">
			<div class="error-icon">❌</div>
			<h3>Failed to Load Articles</h3>
			<p>{$errors.articles}</p>
			<button class="btn btn-primary" on:click={() => ArticleService.loadAllArticles()}>Retry</button>
		</div>

	<!-- Empty State -->
	{:else if $articles.length === 0}
		<div class="empty-state">
			<div class="empty-icon">📄</div>
			<h3>No Articles Yet</h3>
			<p>Add some RSS feeds to start seeing articles here.</p>
			<a href="/feeds" class="btn btn-primary">
				➕ Add RSS Feeds
			</a>
		</div>

	<!-- Articles List -->
	{:else}
		{#if filteredArticles.length === 0}
			<div class="no-results">
				<p>No articles match your search for "{searchTerm}"</p>
				<button class="btn btn-secondary" on:click={() => searchTerm = ''}>Clear Search</button>
			</div>
		{:else}
			<div class="articles-list">
				{#each filteredArticles as article (article.ID || article.id || article.GUID)}
					<article class="article-card {article.Read ? 'read' : 'unread'}">
						<div class="article-header">
							<h3 class="article-title">
								<a href="/articles/{article.ID}">
									{decodeHtml(article.Title) || 'Untitled Article'}
								</a>
							</h3>
							<div class="article-meta">
								<time class="article-date">
									{formatDate(article.PublishDate)}
								</time>
								{#if !article.Read}
									<span class="unread-badge">New</span>
								{/if}
							</div>
						</div>
						
						{#if article.Description}
							<div class="article-description">
								{@html truncateHtml(article.Description, 300)}
							</div>
						{/if}
						
						<div class="article-footer">
							<a href="/articles/{article.ID}" class="btn btn-primary">
								Read Article
							</a>
							<a href={article.Link} target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
								Original Source ↗
							</a>
						</div>
					</article>
				{/each}
			</div>
		{/if}
	{/if}
</div>

<style>
	/* Search Bar */
	.search-bar {
		margin-bottom: 1.5rem;
		max-width: 400px;
	}

	.search-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.search-icon {
		position: absolute;
		left: 0.75rem;
		color: var(--text-tertiary);
		pointer-events: none;
	}

	.search-input {
		width: 100%;
		padding: 0.75rem 0.75rem 0.75rem 2.5rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-secondary);
		color: var(--text-primary);
		font-size: 0.875rem;
		transition: border-color 0.15s ease;
	}

	.search-input:focus {
		outline: none;
		border-color: var(--primary);
	}

	.search-clear {
		position: absolute;
		right: 0.5rem;
		padding: 0.25rem;
		color: var(--text-tertiary);
	}

	/* State Styles */
	.loading-state,
	.error-state,
	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: var(--text-secondary);
	}

	.loading-spinner,
	.error-icon,
	.empty-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.error-state h3,
	.empty-state h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-primary);
		margin-bottom: 0.5rem;
	}

	.no-results {
		text-align: center;
		padding: 2rem;
		color: var(--text-secondary);
	}

	/* Articles List */
	.articles-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.article-card {
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 1.5rem;
		transition: all 0.15s ease;
	}

	.article-card:hover {
		border-color: var(--primary);
		box-shadow: var(--shadow);
	}

	.article-card.unread {
		border-left: 4px solid var(--primary);
	}

	.article-card.read {
		opacity: 0.7;
	}

	.article-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.article-title {
		flex: 1;
		margin: 0;
		font-size: 1.125rem;
		font-weight: 600;
		line-height: 1.4;
	}

	.article-title a {
		color: var(--text-primary);
		text-decoration: none;
		transition: color 0.15s ease;
	}

	.article-title a:hover {
		color: var(--primary);
	}

	.article-meta {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.article-date {
		font-size: 0.75rem;
		color: var(--text-tertiary);
	}

	.unread-badge {
		background: var(--primary);
		color: white;
		font-size: 0.75rem;
		padding: 0.125rem 0.5rem;
		border-radius: 10px;
		font-weight: 500;
	}

	.article-description {
		color: var(--text-secondary);
		line-height: 1.6;
		margin-bottom: 1rem;
		font-size: 0.875rem;
		overflow: hidden;
	}

	/* Constrain all content in article descriptions */
	.article-description :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--radius);
	}

	.article-description :global(iframe),
	.article-description :global(video),
	.article-description :global(embed),
	.article-description :global(object) {
		max-width: 100%;
		height: auto;
	}

	.article-description :global(table) {
		max-width: 100%;
		overflow-x: auto;
		display: block;
		white-space: nowrap;
	}

	.article-footer {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.article-header {
			flex-direction: column;
			align-items: stretch;
		}

		.article-meta {
			flex-direction: row;
			justify-content: space-between;
			align-items: center;
		}

		.article-footer {
			flex-direction: column;
		}

		.article-footer .btn {
			width: 100%;
		}
	}
</style>