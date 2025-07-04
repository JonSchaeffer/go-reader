<script>
	import { onMount } from 'svelte';
	import { FeedService } from '$lib/services/feedService.js';
	import { feeds, loading, errors } from '$lib/stores.js';

	let showAddFeed = false;
	let newFeedUrl = '';
	let searchTerm = '';

	onMount(async () => {
		console.log('Feeds page mounted, loading feeds...');
		try {
			await FeedService.loadFeeds();
			console.log('Feeds loaded successfully');
		} catch (error) {
			console.error('Failed to load feeds on mount:', error);
		}
	});

	async function handleAddFeed() {
		if (!newFeedUrl.trim()) return;
		
		try {
			await FeedService.addFeed(newFeedUrl.trim());
			newFeedUrl = '';
			showAddFeed = false;
		} catch (error) {
			console.error('Failed to add feed:', error);
		}
	}

	async function handleDeleteFeed(feedId) {
		if (!confirm('Are you sure you want to delete this feed? This will also delete all associated articles.')) {
			return;
		}
		
		try {
			await FeedService.deleteFeed(feedId);
		} catch (error) {
			console.error('Failed to delete feed:', error);
		}
	}

	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			handleAddFeed();
		}
	}

	// Reactive loading when page is accessed
	$: if (typeof window !== 'undefined') {
		if ($feeds.length === 0 && !$loading.feeds && !$errors.feeds) {
			console.log('Reactive feeds loading triggered');
			FeedService.loadFeeds();
		}
	}

	// Filter feeds based on search term
	$: filteredFeeds = $feeds.filter(feed => 
		!searchTerm || 
		feed.title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
		feed.url?.toLowerCase().includes(searchTerm.toLowerCase())
	);
</script>

<svelte:head>
	<title>RSS Feeds - RSS Reader</title>
</svelte:head>

<style>
	/* Header Actions */
	.content-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
	}

	.header-actions {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	/* Modal Styles */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
	}

	.modal {
		background: var(--bg-primary);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		width: 90%;
		max-width: 500px;
		max-height: 90vh;
		overflow: hidden;
		box-shadow: var(--shadow-lg);
	}

	.modal-header {
		padding: 1.5rem;
		border-bottom: 1px solid var(--border-light);
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.modal-header h3 {
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--text-primary);
		margin: 0;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.modal-footer {
		padding: 1rem 1.5rem;
		border-top: 1px solid var(--border-light);
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}

	/* Form Styles */
	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
		margin-bottom: 0.5rem;
	}

	.form-group input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-primary);
		color: var(--text-primary);
		font-size: 0.875rem;
		transition: border-color 0.15s ease;
	}

	.form-group input:focus {
		outline: none;
		border-color: var(--primary);
	}

	.form-group small {
		display: block;
		margin-top: 0.5rem;
		font-size: 0.75rem;
		color: var(--text-tertiary);
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
		grid-column: 1 / -1;
		text-align: center;
		padding: 2rem;
		color: var(--text-secondary);
	}

	/* Feeds Grid */
	.feeds-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
		gap: 1.5rem;
	}

	.feed-card {
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 1.5rem;
		transition: all 0.15s ease;
	}

	.feed-card:hover {
		border-color: var(--primary);
		box-shadow: var(--shadow);
	}

	.feed-header {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.feed-icon {
		font-size: 1.5rem;
		width: 2.5rem;
		height: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--primary-light);
		border-radius: var(--radius);
		flex-shrink: 0;
	}

	.feed-info {
		flex: 1;
		min-width: 0;
	}

	.feed-title {
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-primary);
		margin: 0 0 0.25rem 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.feed-url {
		font-size: 0.75rem;
		color: var(--text-tertiary);
		margin: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.feed-actions {
		display: flex;
		gap: 0.25rem;
	}

	.feed-action {
		padding: 0.375rem;
		font-size: 0.875rem;
	}

	.feed-meta {
		margin-bottom: 1rem;
	}

	.feed-stats {
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.stat-label {
		font-size: 0.75rem;
		color: var(--text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.025em;
	}

	.stat-value {
		font-size: 0.875rem;
		color: var(--text-secondary);
		font-weight: 500;
	}

	.stat-value.enabled {
		color: var(--success);
	}

	.feed-footer {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.feed-footer .btn {
		flex: 1;
		min-width: 120px;
		justify-content: center;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.content-header {
			flex-direction: column;
			align-items: stretch;
		}

		.feeds-grid {
			grid-template-columns: 1fr;
		}

		.feed-footer {
			flex-direction: column;
		}

		.feed-footer .btn {
			flex: none;
		}

		.modal {
			margin: 1rem;
			width: auto;
		}
	}
</style>

<div class="content-header">
	<div>
		<h1 style="font-size: 1.875rem; font-weight: 700; color: var(--text-primary);">
			RSS Feeds
		</h1>
		<p style="color: var(--text-secondary); margin-top: 0.5rem;">
			Manage your RSS feed subscriptions ({$feeds.length} feeds)
		</p>
	</div>
	
	<div class="header-actions">
		<button 
			class="btn btn-primary"
			on:click={() => showAddFeed = true}
			disabled={$loading.adding}
		>
			{#if $loading.adding}
				⏳ Adding...
			{:else}
				➕ Add Feed
			{/if}
		</button>
	</div>
</div>

<div class="content-body">
	<!-- Add Feed Modal -->
	{#if showAddFeed}
		<div class="modal-overlay" on:click={() => showAddFeed = false} on:keydown={(e) => e.key === 'Escape' && (showAddFeed = false)} role="dialog" aria-modal="true">
			<div class="modal" on:click|stopPropagation>
				<div class="modal-header">
					<h3>Add New RSS Feed</h3>
					<button class="btn-ghost" on:click={() => showAddFeed = false}>✕</button>
				</div>
				
				<div class="modal-body">
					<div class="form-group">
						<label for="feedUrl">RSS Feed URL</label>
						<input 
							id="feedUrl"
							type="url" 
							bind:value={newFeedUrl}
							on:keypress={handleKeyPress}
							placeholder="https://example.com/feed.xml"
							required
							/>
						<small>Enter the URL of an RSS or Atom feed</small>
					</div>
				</div>
				
				<div class="modal-footer">
					<button class="btn btn-secondary" on:click={() => showAddFeed = false}>Cancel</button>
					<button 
						class="btn btn-primary" 
						on:click={handleAddFeed}
						disabled={!newFeedUrl.trim() || $loading.adding}
					>
						{#if $loading.adding}
							⏳ Adding Feed...
						{:else}
							Add Feed
						{/if}
					</button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Search Bar -->
	{#if $feeds.length > 0}
		<div class="search-bar">
			<div class="search-input-wrapper">
				<span class="search-icon">🔍</span>
				<input 
					type="text" 
					bind:value={searchTerm}
					placeholder="Search feeds..."
					class="search-input"
				/>
				{#if searchTerm}
					<button class="btn-ghost search-clear" on:click={() => searchTerm = ''}>✕</button>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Loading State -->
	{#if $loading.feeds}
		<div class="loading-state">
			<div class="loading-spinner">⏳</div>
			<p>Loading your RSS feeds...</p>
		</div>

	<!-- Error State -->
	{:else if $errors.feeds}
		<div class="error-state">
			<div class="error-icon">❌</div>
			<h3>Failed to Load Feeds</h3>
			<p>{$errors.feeds}</p>
			<button class="btn btn-primary" on:click={() => FeedService.loadFeeds()}>Retry</button>
		</div>

	<!-- Empty State -->
	{:else if $feeds.length === 0}
		<div class="empty-state">
			<div class="empty-icon">📡</div>
			<h3>No RSS Feeds Yet</h3>
			<p>Get started by adding your first RSS feed subscription.</p>
			<button class="btn btn-primary" on:click={() => showAddFeed = true}>
				➕ Add Your First Feed
			</button>
		</div>

	<!-- Feeds List -->
	{:else}
		<div class="feeds-grid">
			{#if filteredFeeds.length === 0}
				<div class="no-results">
					<p>No feeds match your search for "{searchTerm}"</p>
					<button class="btn btn-secondary" on:click={() => searchTerm = ''}>Clear Search</button>
				</div>
			{:else}
				{#each filteredFeeds as feed (feed.id)}
					<div class="feed-card">
						<div class="feed-header">
							<div class="feed-icon">📡</div>
							<div class="feed-info">
								<h4 class="feed-title">{feed.title || 'Untitled Feed'}</h4>
								<p class="feed-url">{feed.url}</p>
							</div>
							<div class="feed-actions">
								<button 
									class="btn-ghost feed-action"
									on:click={() => handleDeleteFeed(feed.id)}
									disabled={$loading.deleting}
									title="Delete Feed"
								>
									🗑️
								</button>
							</div>
						</div>
						
						<div class="feed-meta">
							<div class="feed-stats">
								<span class="stat">
									<span class="stat-label">Added:</span>
									<span class="stat-value">{new Date(feed.created_at).toLocaleDateString()}</span>
								</span>
								{#if feed.fivefilters_url}
									<span class="stat">
										<span class="stat-label">Full-text:</span>
										<span class="stat-value enabled">✓ Enabled</span>
									</span>
								{/if}
							</div>
						</div>
						
						<div class="feed-footer">
							<a href="/articles?feed={feed.id}" class="btn btn-secondary">
								View Articles
							</a>
							<a href={feed.url} target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
								Original Feed ↗
							</a>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
</div>