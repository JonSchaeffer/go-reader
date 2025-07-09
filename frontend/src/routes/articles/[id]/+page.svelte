<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArticleService } from '$lib/services/articleService.js';
	import { loading, errors } from '$lib/stores.js';

	let article = null;
	let articleId = null;

	// Function to decode HTML entities
	function decodeHtml(html) {
		if (!html || typeof window === 'undefined') return html;
		
		const temp = document.createElement('div');
		temp.innerHTML = html;
		return temp.textContent || temp.innerText || html;
	}

	onMount(async () => {
		articleId = $page.params.id;
		await loadArticle();
	});

	async function loadArticle() {
		if (!articleId) return;
		
		try {
			article = await ArticleService.getArticleById(articleId);
			
			// Auto-mark as read when article is opened
			if (article && !article.Read) {
				await markAsRead();
			}
		} catch (error) {
			console.error('Failed to load article:', error);
		}
	}

	async function markAsRead() {
		if (!article || article.Read) return;
		
		try {
			const newStatus = await ArticleService.toggleReadStatus(article.ID, article.Read);
			article.Read = newStatus;
		} catch (error) {
			console.error('Failed to mark as read:', error);
		}
	}

	async function toggleReadStatus() {
		if (!article) return;
		
		try {
			const newStatus = await ArticleService.toggleReadStatus(article.ID, article.Read);
			article.Read = newStatus;
		} catch (error) {
			console.error('Failed to toggle read status:', error);
		}
	}

	function formatDate(dateString) {
		if (!dateString) return 'No date';
		
		try {
			// Handle different date formats
			let date;
			
			// Try parsing as-is first
			date = new Date(dateString);
			
			// If invalid, try parsing ISO format with Z
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
				return 'Invalid date';
			}
			
			return date.toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'long',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch (error) {
			return 'Invalid date';
		}
	}

	function goBack() {
		if (window.history.length > 1) {
			window.history.back();
		} else {
			goto('/articles');
		}
	}
</script>

<svelte:head>
	<title>{article ? decodeHtml(article.Title) : 'Loading Article'} - RSS Reader</title>
</svelte:head>

<div class="article-page">
	<!-- Header with navigation -->
	<header class="article-header">
		<button class="back-button" on:click={goBack}>
			← Back to Articles
		</button>
		
		{#if article}
			<div class="article-actions">
				<button 
					class="btn btn-secondary"
					on:click={toggleReadStatus}
					disabled={$loading.articles}
				>
					{article.Read ? '📕 Mark Unread' : '📖 Mark Read'}
				</button>
				
				<a 
					href={article.Link} 
					target="_blank" 
					rel="noopener noreferrer" 
					class="btn btn-secondary"
				>
					🔗 View Original
				</a>
			</div>
		{/if}
	</header>

	<!-- Main Scrollable Content -->
	<div class="main-content">
		<!-- Loading State -->
		{#if $loading.articles && !article}
			<div class="loading-state">
				<div class="loading-spinner">⏳</div>
				<p>Loading article...</p>
			</div>

		<!-- Error State -->
		{:else if $errors.articles && !article}
			<div class="error-state">
				<div class="error-icon">❌</div>
				<h2>Failed to Load Article</h2>
				<p>{$errors.articles}</p>
				<button class="btn btn-primary" on:click={loadArticle}>Retry</button>
			</div>

		<!-- Article Content -->
		{:else if article}
			<article class="article-content">
				<!-- Article Title -->
				<header class="content-header">
					<h1 class="article-title">{decodeHtml(article.Title)}</h1>
					
					<div class="article-meta">
						<time class="article-date">
							{formatDate(article.PublishDate)}
						</time>
						
						{#if !article.Read}
							<span class="unread-badge">New</span>
						{/if}
					</div>
					
					{#if article.Link}
						<div class="article-link">
							<a href={article.Link} target="_blank" rel="noopener noreferrer">
								🔗 {article.Link}
							</a>
						</div>
					{/if}
				</header>

				<!-- Article Body -->
				<div class="article-body">
					{#if article.Description}
						{@html article.Description}
					{:else}
						<p class="no-content">No content available for this article.</p>
						<a href={article.Link} target="_blank" rel="noopener noreferrer" class="btn btn-primary">
							Read on Original Site →
						</a>
					{/if}
				</div>
			</article>

		<!-- Article Not Found -->
		{:else}
			<div class="error-state">
				<div class="error-icon">📄</div>
				<h2>Article Not Found</h2>
				<p>The article you're looking for doesn't exist or has been removed.</p>
				<button class="btn btn-primary" on:click={() => goto('/articles')}>
					← Back to Articles
				</button>
			</div>
		{/if}
	</div>
</div>

<style>
	.article-page {
		height: 100vh;
		background: var(--bg-primary);
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	/* Header */
	.article-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 2rem;
		border-bottom: 1px solid var(--border-light);
		background: var(--bg-secondary);
		flex-shrink: 0;
		z-index: 10;
	}

	.back-button {
		background: none;
		border: none;
		color: var(--primary);
		font-size: 0.875rem;
		cursor: pointer;
		padding: 0.5rem 0;
		transition: color 0.15s ease;
	}

	.back-button:hover {
		color: var(--primary-hover);
		text-decoration: underline;
	}

	.article-actions {
		display: flex;
		gap: 0.75rem;
	}

	/* Main Content Container */
	.main-content {
		flex: 1;
		overflow-y: auto;
		overflow-x: hidden;
	}

	/* Content */
	.article-content {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.content-header {
		margin-bottom: 2rem;
		padding-bottom: 1.5rem;
		border-bottom: 1px solid var(--border-light);
	}

	.article-title {
		font-size: 2.25rem;
		font-weight: 700;
		line-height: 1.2;
		color: var(--text-primary);
		margin: 0 0 1rem 0;
	}

	.article-meta {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.article-date {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.unread-badge {
		background: var(--primary);
		color: white;
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 10px;
		font-weight: 500;
	}

	.article-link {
		margin-top: 1rem;
	}

	.article-link a {
		color: var(--text-tertiary);
		font-size: 0.875rem;
		text-decoration: none;
		transition: color 0.15s ease;
	}

	.article-link a:hover {
		color: var(--primary);
		text-decoration: underline;
	}

	/* Article Body */
	.article-body {
		line-height: 1.7;
		color: var(--text-primary);
		font-size: 1.125rem;
	}

	.no-content {
		text-align: center;
		color: var(--text-secondary);
		margin: 2rem 0;
		font-style: italic;
	}

	/* Style content within article body */
	.article-body :global(p) {
		margin-bottom: 1.5rem;
	}

	.article-body :global(h1),
	.article-body :global(h2),
	.article-body :global(h3),
	.article-body :global(h4),
	.article-body :global(h5),
	.article-body :global(h6) {
		margin: 2rem 0 1rem 0;
		color: var(--text-primary);
		font-weight: 600;
	}

	.article-body :global(h1) { font-size: 1.875rem; }
	.article-body :global(h2) { font-size: 1.5rem; }
	.article-body :global(h3) { font-size: 1.25rem; }

	.article-body :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--radius);
		margin: 1.5rem 0;
		box-shadow: var(--shadow);
	}

	.article-body :global(blockquote) {
		border-left: 4px solid var(--primary);
		padding-left: 1.5rem;
		margin: 1.5rem 0;
		font-style: italic;
		color: var(--text-secondary);
	}

	.article-body :global(code) {
		background: var(--bg-tertiary);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-size: 0.875em;
	}

	.article-body :global(pre) {
		background: var(--bg-tertiary);
		padding: 1rem;
		border-radius: var(--radius);
		overflow-x: auto;
		margin: 1.5rem 0;
	}

	.article-body :global(pre code) {
		background: none;
		padding: 0;
	}

	.article-body :global(ul),
	.article-body :global(ol) {
		margin: 1.5rem 0;
		padding-left: 2rem;
	}

	.article-body :global(li) {
		margin-bottom: 0.5rem;
	}

	.article-body :global(a) {
		color: var(--primary);
		text-decoration: underline;
		transition: color 0.15s ease;
	}

	.article-body :global(a:hover) {
		color: var(--primary-hover);
	}

	.article-body :global(table) {
		width: 100%;
		border-collapse: collapse;
		margin: 1.5rem 0;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		overflow: hidden;
	}

	.article-body :global(th),
	.article-body :global(td) {
		padding: 0.75rem;
		text-align: left;
		border-bottom: 1px solid var(--border-light);
	}

	.article-body :global(th) {
		background: var(--bg-secondary);
		font-weight: 600;
	}

	/* State Styles */
	.loading-state,
	.error-state {
		text-align: center;
		padding: 4rem 2rem;
		max-width: 600px;
		margin: 0 auto;
	}

	.loading-spinner,
	.error-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.error-state h2 {
		color: var(--text-primary);
		margin-bottom: 1rem;
	}

	.error-state p {
		color: var(--text-secondary);
		margin-bottom: 2rem;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.article-header {
			flex-direction: column;
			gap: 1rem;
			align-items: stretch;
		}

		.article-actions {
			justify-content: center;
		}

		.article-content {
			padding: 1rem;
		}

		.article-title {
			font-size: 1.75rem;
		}

		.article-body {
			font-size: 1rem;
		}

		.article-meta {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}
	}
</style>