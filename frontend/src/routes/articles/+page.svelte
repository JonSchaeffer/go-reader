<script>
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { afterNavigate } from '$app/navigation';
	import { ArticleService } from '$lib/services/articleService.js';
	import { FeedService } from '$lib/services/feedService.js';
	import { CategoryService } from '$lib/services/categoryService.js';
	import { articles, feeds, categories, loading, errors } from '$lib/stores.js';
	import { articleApi } from '$lib/api.js';

	let feedId = null;
	let categoryId = null;
	let searchTerm = '';
	let readFilter = 'all'; // 'all', 'unread', 'read'
	let refreshInterval = null;
	let isRefreshing = false;
	let hasAttemptedLoad = false;
	
	// Pagination state
	let currentPage = 0;
	let pageSize = 50;
	let totalArticles = 0;
	let hasMore = true;
	let isLoadingMore = false;
	
	// Article counts for totals
	let totalUnreadCount = 0;
	let totalReadCount = 0;
	
	// Scroll tracking for auto-load
	let lastScrollY = 0;

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

	// Function to format dates properly with timestamp
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
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch (error) {
			console.warn('Date parsing error:', error, dateString);
			return 'Invalid date';
		}
	}

	
	// Auto-load on scroll to bottom
	function handleScroll() {
		if (typeof window === 'undefined') return;
		
		// Auto-load more articles
		if (isLoadingMore || !hasMore || feedId) {
			return;
		}
		
		const scrollPosition = window.scrollY + window.innerHeight;
		const documentHeight = document.documentElement.scrollHeight;
		const distanceFromBottom = documentHeight - scrollPosition;
		
		// Load more when user is within 300px of bottom (increased threshold)
		if (distanceFromBottom <= 300) {
			loadMoreArticles();
		}
	}


	// Alternative approach: Intersection Observer
	let observer;
	let sentinelElement;
	
	function setupIntersectionObserver() {
		if (typeof window === 'undefined') return;
		
		observer = new IntersectionObserver((entries) => {
			entries.forEach(entry => {
				if (entry.isIntersecting && !isLoadingMore && hasMore && !feedId) {
					loadMoreArticles();
				}
			});
		}, {
			rootMargin: '200px'
		});
		
		// Wait for the sentinel element to be available
		setTimeout(() => {
			if (sentinelElement) {
				observer.observe(sentinelElement);
			}
		}, 100);
	}

	// Load article counts for totals - using sample-based estimation
	async function loadArticleCounts() {
		// We'll use the sample-based estimation from the initial load
		// This is more efficient and avoids the 500 error
		console.log('Using sample-based article counts');
	}

	async function loadData() {
		console.log('Loading articles data...');
		hasAttemptedLoad = true;
		try {
			// Check if we have filters from URL params
			feedId = $page.url.searchParams.get('feed');
			categoryId = $page.url.searchParams.get('category');
			console.log('Feed ID from URL:', feedId);
			console.log('Category ID from URL:', categoryId);

			// Load feeds and categories data
			await Promise.all([
				FeedService.loadFeeds(),
				CategoryService.loadCategories()
			]);

			if (feedId) {
				await ArticleService.loadArticlesByFeed(feedId);
			} else {
				const result = await ArticleService.loadAllArticles(0, pageSize);
				console.log('Initial load result:', result);
				totalArticles = result.total;
				hasMore = result.hasMore;
				currentPage = 0;
				console.log('Initial state - totalArticles:', totalArticles, 'hasMore:', hasMore, 'currentPage:', currentPage);
				
				// Calculate initial counts based on the first page sample
				if (result.articles && result.articles.length > 0) {
					const sampleUnreadCount = result.articles.filter(a => !a.Read).length;
					const sampleReadCount = result.articles.filter(a => a.Read).length;
					const sampleTotal = result.articles.length;
					
					// Estimate total counts based on the sample
					const unreadRatio = sampleUnreadCount / sampleTotal;
					const readRatio = sampleReadCount / sampleTotal;
					
					totalUnreadCount = Math.floor(totalArticles * unreadRatio);
					totalReadCount = Math.floor(totalArticles * readRatio);
					
					console.log('Estimated counts from sample - Unread:', totalUnreadCount, 'Read:', totalReadCount);
				}
			}
			console.log('Articles loaded successfully, count:', $articles.length);
		} catch (error) {
			console.error('Failed to load articles:', error);
		}
	}

	// Function to get feed name from RSS ID
	function getFeedName(rssId) {
		if (!rssId || !$feeds) return `Feed ${rssId}`;

		const feed = $feeds.find((f) => f.ID === rssId);
		if (feed && feed.Title) {
			// Truncate feed name if too long
			const title = decodeHtml(feed.Title);
			return title.length > 20 ? title.substring(0, 17) + '...' : title;
		}
		return `Feed ${rssId}`;
	}

	async function toggleArticleReadStatus(article, event) {
		// Prevent the row click from navigating to the article
		event.stopPropagation();
		event.preventDefault();

		try {
			const newStatus = await ArticleService.toggleReadStatus(article.ID, article.Read);

			// Update the article in the local articles array
			articles.update((currentArticles) =>
				currentArticles.map((a) => (a.ID === article.ID ? { ...a, Read: newStatus } : a))
			);
			
			// Update total counts
			if (newStatus) {
				// Article was marked as read
				totalUnreadCount--;
				totalReadCount++;
			} else {
				// Article was marked as unread
				totalUnreadCount++;
				totalReadCount--;
			}
		} catch (error) {
			console.error('Failed to toggle read status:', error);
		}
	}

	function openExternalLink(url, event) {
		// Prevent the row click from navigating to the article
		event.stopPropagation();
		event.preventDefault();

		if (url) {
			window.open(url, '_blank', 'noopener,noreferrer');
		}
	}

	function filterByFeed(rssId, event) {
		// Prevent the row click from navigating to the article
		event.stopPropagation();
		event.preventDefault();

		// Navigate to filtered view for this feed
		window.location.href = `/articles?feed=${rssId}`;
	}

	async function markAllAsRead() {
		if (totalUnreadCount === 0) {
			return; // No unread articles
		}

		if (!confirm(`Mark all ${totalUnreadCount} unread articles as read?`)) {
			return;
		}

		try {
			// Update all unread articles to read status
			const unreadArticles = $articles.filter(article => !article.Read);
			const updatePromises = unreadArticles.map(article => 
				ArticleService.toggleReadStatus(article.ID, article.Read)
			);
			
			await Promise.all(updatePromises);
			
			// Update total counts
			totalReadCount += totalUnreadCount;
			totalUnreadCount = 0;
		} catch (error) {
			console.error('Failed to mark all articles as read:', error);
		}
	}

	async function refreshArticles() {
		isRefreshing = true;
		try {
			if (feedId) {
				await ArticleService.loadArticlesByFeedSilently(feedId);
			} else {
				// Reset pagination and load first page
				currentPage = 0;
				const result = await ArticleService.loadAllArticles(0, pageSize);
				totalArticles = result.total;
				hasMore = result.hasMore;
			}
		} catch (error) {
			console.error('Failed to refresh articles:', error);
		} finally {
			isRefreshing = false;
		}
	}

	async function loadMoreArticles() {
		if (isLoadingMore || !hasMore || feedId) {
			return; // Skip if already loading, no more, or feed-specific
		}
		
		isLoadingMore = true;
		try {
			const nextPage = currentPage + 1;
			const offset = nextPage * pageSize;
			
			const result = await ArticleService.loadAllArticles(offset, pageSize);
			
			if (result.articles.length > 0) {
				currentPage = nextPage;
				hasMore = result.hasMore;
			} else {
				hasMore = false;
			}
		} catch (error) {
			console.error('Failed to load more articles:', error);
		} finally {
			isLoadingMore = false;
		}
	}

	async function refreshArticlesSilently() {
		try {
			if (feedId) {
				await ArticleService.loadArticlesByFeedSilently(feedId);
			} else {
				await ArticleService.refreshArticlesSilently();
			}
		} catch (error) {
			console.error('Failed to refresh articles silently:', error);
		}
	}

	function startAutoRefresh() {
		refreshInterval = setInterval(refreshArticlesSilently, 60000); // Refresh every minute silently
	}

	function stopAutoRefresh() {
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	onMount(() => {
		console.log('Articles page mounted');
		loadData();
		startAutoRefresh();
		
		// Add scroll listener for auto-loading (browser only)
		if (typeof window !== 'undefined') {
			window.addEventListener('scroll', handleScroll, { passive: true });
		}
		
		// Setup Intersection Observer as backup
		setupIntersectionObserver();
	});

	onDestroy(() => {
		stopAutoRefresh();
		
		// Remove event listeners (browser only)
		if (typeof window !== 'undefined') {
			window.removeEventListener('scroll', handleScroll);
		}
		
		// Cleanup Intersection Observer
		if (observer) {
			observer.disconnect();
		}
	});

	afterNavigate((navigation) => {
		console.log('Navigated to articles page');
		// Check if URL parameters have changed
		const currentFeedId = $page.url.searchParams.get('feed');
		const currentCategoryId = $page.url.searchParams.get('category');
		
		if (currentFeedId !== feedId || currentCategoryId !== categoryId || (!$articles.length && !hasAttemptedLoad)) {
			loadData();
		}
	});

	// React to URL parameter changes for real-time filtering
	$: {
		const urlFeedId = $page.url.searchParams.get('feed');
		const urlCategoryId = $page.url.searchParams.get('category');
		
		// Update local variables when URL changes
		if (urlFeedId !== feedId) {
			feedId = urlFeedId;
		}
		if (urlCategoryId !== categoryId) {
			categoryId = urlCategoryId;
		}
	}

	// Filter articles based on search term, read status, and category
	$: filteredArticles = $articles.filter((article) => {
		// Search filter
		const matchesSearch = !searchTerm || 
			decodeHtml(article.Title)?.toLowerCase().includes(searchTerm.toLowerCase()) ||
			decodeHtml(article.Description)?.toLowerCase().includes(searchTerm.toLowerCase());
		
		// Read status filter
		const matchesReadFilter = 
			readFilter === 'all' || 
			(readFilter === 'unread' && !article.Read) ||
			(readFilter === 'read' && article.Read);
		
		// Category filter - find the feed that owns this article and check its category
		const matchesCategory = !categoryId || (() => {
			const feed = $feeds.find(f => f.ID === article.RssID);
			if (!feed) return false;
			
			// Handle "uncategorized" filter (null categoryId in URL becomes string 'null')
			if (categoryId === 'null' || categoryId === '') {
				return !feed.CategoryID;
			}
			
			return feed.CategoryID === parseInt(categoryId);
		})();
		
		return matchesSearch && matchesReadFilter && matchesCategory;
	});
</script>

<svelte:head>
	<title>Articles - RSS Reader</title>
</svelte:head>

<!-- Articles header removed for mobile - info now in filter bar -->

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
					<button class="btn-ghost search-clear" on:click={() => (searchTerm = '')}>✕</button>
				{/if}
			</div>
		</div>

		<!-- Filter Buttons -->
		<div class="filter-bar">
			{#if feedId}
				<div class="feed-info">
					<span class="feed-name">📡 {getFeedName(parseInt(feedId))}</span>
					<a href="/articles" class="clear-filter">← All Articles</a>
				</div>
			{/if}
			
			<div class="filter-group">
				<span class="filter-label">Show:</span>
				<button 
					class="filter-btn {readFilter === 'all' ? 'active' : ''}"
					on:click={() => readFilter = 'all'}
				>
					All ({totalArticles})
				</button>
				<button 
					class="filter-btn {readFilter === 'unread' ? 'active' : ''}"
					on:click={() => readFilter = 'unread'}
				>
					📕 Unread ({totalUnreadCount})
				</button>
				<button 
					class="filter-btn {readFilter === 'read' ? 'active' : ''}"
					on:click={() => readFilter = 'read'}
				>
					📖 Read ({totalReadCount})
				</button>
				
				<!-- Separator -->
				<div class="filter-separator"></div>
				
				<!-- Bulk Actions -->
				<button 
					class="bulk-action-btn"
					on:click={markAllAsRead}
					disabled={totalUnreadCount === 0 || $loading.articles}
					title="Mark all unread articles as read"
				>
					{#if $loading.articles}
						⏳ Updating...
					{:else}
						✓ Mark All Read
					{/if}
				</button>
				
				<!-- Refresh Button -->
				<button 
					class="refresh-btn"
					on:click={refreshArticles}
					disabled={isRefreshing}
					title="Refresh articles (Auto-refresh every minute)"
				>
					{#if isRefreshing}
						⏳ Refreshing...
					{:else}
						🔄 Refresh
					{/if}
				</button>
			</div>
		</div>
	{/if}

	<!-- Error State -->
	{#if $errors.articles}
		<div class="error-state">
			<div class="error-icon">❌</div>
			<h3>Failed to Load Articles</h3>
			<p>{$errors.articles}</p>
			<button class="btn btn-primary" on:click={() => ArticleService.loadAllArticles()}
				>Retry</button
			>
		</div>

		<!-- Empty State - only show if we've attempted to load and there are truly no articles -->
	{:else if $articles.length === 0 && hasAttemptedLoad}
		<div class="empty-state">
			<div class="empty-icon">📄</div>
			<h3>No Articles Yet</h3>
			<p>Add some RSS feeds to start seeing articles here.</p>
			<a href="/feeds" class="btn btn-primary"> ➕ Add RSS Feeds </a>
		</div>

		<!-- Articles List -->
	{:else if filteredArticles.length === 0}
		<div class="no-results">
			<p>No articles match your search for "{searchTerm}"</p>
			<button class="btn btn-secondary" on:click={() => (searchTerm = '')}>Clear Search</button>
		</div>
	{:else}
		<div class="articles-list">
			{#each filteredArticles as article (article.ID || article.id || article.GUID)}
				<article
					class="article-row {article.Read ? 'read' : 'unread'}"
					on:click={() => (window.location.href = `/articles/${article.ID}`)}
				>
					<!-- Read/Unread Status Icon (clickable) -->
					<button
						class="status-icon"
						on:click={(e) => toggleArticleReadStatus(article, e)}
						title={article.Read ? 'Mark as unread' : 'Mark as read'}
					>
						{#if article.Read}
							<span class="read-icon">📖</span>
						{:else}
							<span class="unread-icon">📕</span>
						{/if}
					</button>

					<!-- Feed Source -->
					<button 
						class="feed-source clickable"
						on:click={(e) => filterByFeed(article.RssID, e)}
						title="Filter articles from {getFeedName(article.RssID)}"
					>
						{getFeedName(article.RssID)}
					</button>

					<!-- Article Title -->
					<div class="article-title">
						<span class="title-text">
							{decodeHtml(article.Title) || 'Untitled Article'}
						</span>
					</div>

					<!-- Date with timestamp -->
					<div class="article-date">
						{formatDate(article.PublishDate)}
					</div>

					<!-- External Link Icon -->
					<button
						class="external-link"
						on:click={(e) => openExternalLink(article.Link, e)}
						title="Open original article"
					>
						🔗
					</button>
				</article>
			{/each}
		</div>
		
		<!-- Auto-loading indicator -->
		{#if isLoadingMore}
			<div class="loading-indicator">
				<div class="loading-spinner">⏳</div>
				<p>Loading more articles...</p>
			</div>
		{/if}
		
		<!-- Sentinel element for Intersection Observer -->
		<div class="sentinel" bind:this={sentinelElement}></div>
	{/if}
</div>

<style>
	/* Feed Info in Filter Bar */
	.feed-info {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem;
		background: var(--bg-tertiary);
		border-radius: var(--radius);
		margin-bottom: 0.5rem;
	}

	.feed-name {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
	}

	.clear-filter {
		color: var(--primary);
		text-decoration: none;
		font-size: 0.8125rem;
	}

	.clear-filter:hover {
		text-decoration: underline;
	}

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

	/* Desktop search adjustments */
	@media (min-width: 769px) {
		.search-bar {
			max-width: 100%;
		}
	}

	/* Filter Bar */
	.filter-bar {
		margin-bottom: 1.5rem;
	}


	.filter-group {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.filter-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.filter-buttons {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.action-row {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.filter-separator {
		width: 1px;
		height: 2rem;
		background: var(--border);
		margin: 0 0.5rem;
	}

	.filter-label {
		font-size: 0.875rem;
		color: var(--text-secondary);
		font-weight: 500;
		margin-right: 0.5rem;
	}

	.filter-btn {
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-secondary);
		color: var(--text-secondary);
		cursor: pointer;
		transition: all 0.15s ease;
		white-space: nowrap;
	}

	.filter-btn:hover {
		background: var(--bg-tertiary);
		border-color: var(--primary);
	}

	.filter-btn.active {
		background: var(--primary);
		color: white;
		border-color: var(--primary);
	}

	.filter-separator {
		width: 1px;
		height: 2rem;
		background: var(--border);
		margin: 0 0.5rem;
	}

	.bulk-action-btn {
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		border: 1px solid var(--success);
		border-radius: var(--radius);
		background: var(--bg-secondary);
		color: var(--success);
		cursor: pointer;
		transition: all 0.15s ease;
		white-space: nowrap;
		font-weight: 500;
	}

	.bulk-action-btn:hover:not(:disabled) {
		background: var(--success);
		color: white;
	}

	.bulk-action-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		border-color: var(--border);
		color: var(--text-tertiary);
	}

	.refresh-btn {
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		border: 1px solid var(--primary);
		border-radius: var(--radius);
		background: var(--bg-secondary);
		color: var(--primary);
		cursor: pointer;
		transition: all 0.15s ease;
		white-space: nowrap;
		font-weight: 500;
	}

	.refresh-btn:hover:not(:disabled) {
		background: var(--primary);
		color: white;
	}

	.refresh-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		border-color: var(--border);
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
		gap: 0;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		overflow: hidden;
		background: var(--bg-secondary);
	}

	.article-row {
		display: flex;
		align-items: center;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid var(--border-light);
		transition: all 0.15s ease;
		gap: 1rem;
		min-height: 3rem;
		cursor: pointer;
	}

	.article-row:last-child {
		border-bottom: none;
	}

	.article-row:hover {
		background: var(--bg-tertiary);
	}

	.article-row.read {
		opacity: 0.75;
	}

	/* Status Icon */
	.status-icon {
		font-size: 1rem;
		width: 1.5rem;
		flex-shrink: 0;
		text-align: center;
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		transition: transform 0.15s ease;
		border-radius: 4px;
	}

	.status-icon:hover {
		transform: scale(1.1);
		background: var(--bg-primary);
	}

	.read-icon {
		opacity: 0.6;
	}

	.unread-icon {
		/* Keep normal opacity for unread items */
	}

	/* Feed Source */
	.feed-source {
		font-size: 0.75rem;
		color: var(--text-tertiary);
		background: var(--bg-primary);
		padding: 0.25rem 0.5rem;
		border-radius: 12px;
		white-space: nowrap;
		flex-shrink: 0;
		min-width: 4rem;
		text-align: center;
		font-weight: 500;
		border: 1px solid transparent;
		cursor: default;
	}

	.feed-source.clickable {
		cursor: pointer;
		transition: all 0.15s ease;
		border-color: var(--border-light);
	}

	.feed-source.clickable:hover {
		background: var(--primary-light);
		color: var(--primary);
		border-color: var(--primary);
		transform: translateY(-1px);
	}

	/* Article Title */
	.article-title {
		flex: 1;
		min-width: 0; /* Allow flex item to shrink */
	}

	.title-text {
		color: var(--text-primary);
		font-size: 0.875rem;
		font-weight: 500;
		line-height: 1.4;
		display: block;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* External Link */
	.external-link {
		background: none;
		border: none;
		font-size: 0.875rem;
		cursor: pointer;
		padding: 0.25rem;
		transition: all 0.15s ease;
		border-radius: 4px;
		color: var(--text-tertiary);
		flex-shrink: 0;
		width: 1.5rem;
		height: 1.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.external-link:hover {
		background: var(--bg-primary);
		color: var(--primary);
		transform: scale(1.1);
	}

	/* Article Date */
	.article-date {
		font-size: 0.75rem;
		color: var(--text-tertiary);
		white-space: nowrap;
		flex-shrink: 0;
		min-width: 7rem;
		text-align: right;
	}

	/* Loading Indicator */
	.loading-indicator {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 2rem 0;
		margin-top: 1rem;
		color: var(--text-secondary);
	}

	.loading-spinner {
		font-size: 1.5rem;
		margin-bottom: 0.5rem;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.loading-indicator p {
		font-size: 0.875rem;
		margin: 0;
	}

	/* Sentinel element for Intersection Observer */
	.sentinel {
		height: 1px;
		width: 100%;
		background: transparent;
	}

	/* Desktop styles - normal layout */
	@media (min-width: 769px) {
		.search-bar {
			position: static;
			background: transparent;
			padding: 0;
			margin-bottom: 1.5rem;
			border: none;
		}

		.filter-bar {
			position: static;
			background: transparent;
			padding: 0;
			margin-bottom: 1.5rem;
			border: none;
		}

		.filter-group {
			flex-direction: row;
			align-items: center;
			gap: 0.5rem;
			flex-wrap: wrap;
			padding: 0;
		}

		.filter-separator {
			display: block;
		}

		.content-body {
			padding-top: 0;
		}
	}

	/* Mobile styles - optimized layout */
	@media (max-width: 768px) {
		.filter-bar {
			position: fixed;
			top: 60px; /* Account for main app header */
			left: 0;
			right: 0;
			margin: 0;
			z-index: 8; /* Below main header but above content */
			background: var(--bg-primary);
			border-bottom: 1px solid var(--border-light);
			backdrop-filter: blur(10px);
			-webkit-backdrop-filter: blur(10px);
		}

		.filter-group {
			flex-direction: column;
			padding: 0.75rem 1rem;
			gap: 0.75rem;
		}

		.filter-group > * {
			width: 100%;
		}

		.filter-group > .filter-label {
			margin-bottom: 0.25rem;
		}

		.filter-group > .filter-btn {
			margin-bottom: 0.25rem;
		}

		.filter-group > .bulk-action-btn,
		.filter-group > .refresh-btn {
			margin-top: 0.25rem;
		}

		.filter-separator {
			display: none;
		}

		.filter-btn {
			font-size: 0.75rem;
			padding: 0.375rem 0.5rem;
		}

		.bulk-action-btn,
		.refresh-btn {
			font-size: 0.75rem;
			padding: 0.375rem 0.75rem;
		}

		/* Add top padding to content body to account for fixed filter bar */
		.content-body {
			padding-top: 160px; /* Main header + filter bar */
		}

		.search-bar {
			margin-bottom: 1rem;
		}

		.article-row {
			padding: 0.5rem 0.75rem;
			gap: 0.75rem;
		}

		.feed-source {
			display: none; /* Hide feed source on mobile to save space */
		}

		.title-text {
			font-size: 0.8125rem;
		}

		.external-link {
			font-size: 0.75rem;
			width: 1.25rem;
			height: 1.25rem;
		}

		.article-date {
			font-size: 0.6875rem;
			min-width: 5rem;
		}

		.status-icon {
			font-size: 0.875rem;
			width: 1.25rem;
		}

		.loading-indicator {
			padding: 1.5rem 0;
		}
	}
</style>

