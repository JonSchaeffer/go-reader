<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { FeedService } from '$lib/services/feedService.js';
	import { ArticleService } from '$lib/services/articleService.js';
	import { feeds, articles, loading, errors } from '$lib/stores.js';

	let stats = {
		totalFeeds: 0,
		totalArticles: 0,
		unreadArticles: 0,
		apiStatus: 'checking'
	};

	onMount(async () => {
		await loadDashboardData();
	});

	async function loadDashboardData() {
		try {
			// Load feeds and articles
			await Promise.all([
				FeedService.loadFeeds(),
				ArticleService.loadAllArticles()
			]);
			
			// Update stats
			stats.totalFeeds = $feeds.length;
			stats.totalArticles = $articles.length;
			stats.unreadArticles = $articles.filter(article => !article.Read).length;
			stats.apiStatus = 'connected';
		} catch (error) {
			console.error('Failed to load dashboard data:', error);
			stats.apiStatus = 'error';
		}
	}

	function getQuickActions() {
		return [
			{
				title: 'View All Feeds',
				description: `Browse your ${stats.totalFeeds} RSS feeds`,
				action: () => goto('/feeds'),
				icon: '📡',
				primary: true
			},
			{
				title: 'Read Articles',
				description: `${stats.unreadArticles} unread of ${stats.totalArticles} total`,
				action: () => goto('/articles'),
				icon: '📄',
				primary: stats.unreadArticles > 0
			},
			{
				title: 'Add New Feed',
				description: 'Subscribe to a new RSS feed',
				action: () => goto('/feeds?action=add'),
				icon: '➕',
				primary: false
			}
		];
	}

	$: quickActions = getQuickActions();
</script>

<svelte:head>
	<title>Dashboard - RSS Reader</title>
</svelte:head>

<div class="content-header">
	<h1 style="font-size: 1.875rem; font-weight: 700; color: var(--text-primary);">
		Welcome to RSS Reader
	</h1>
	<p style="color: var(--text-secondary); margin-top: 0.5rem;">
		Stay up to date with your favorite feeds and articles
	</p>
</div>

<div class="content-body">
	<!-- Status Cards -->
	<div class="stats-grid">
		<div class="stat-card">
			<div class="stat-icon">📡</div>
			<div class="stat-content">
				<div class="stat-value">{stats.totalFeeds}</div>
				<div class="stat-label">RSS Feeds</div>
			</div>
		</div>
		
		<div class="stat-card">
			<div class="stat-icon">📄</div>
			<div class="stat-content">
				<div class="stat-value">{stats.totalArticles}</div>
				<div class="stat-label">Total Articles</div>
			</div>
		</div>
		
		<div class="stat-card">
			<div class="stat-icon">🔵</div>
			<div class="stat-content">
				<div class="stat-value">{stats.unreadArticles}</div>
				<div class="stat-label">Unread Articles</div>
			</div>
		</div>
		
		<div class="stat-card">
			<div class="stat-icon">
				{#if stats.apiStatus === 'connected'}
					✅
				{:else if stats.apiStatus === 'error'}
					❌
				{:else}
					⏳
				{/if}
			</div>
			<div class="stat-content">
				<div class="stat-value" style="font-size: 0.875rem;">
					{#if stats.apiStatus === 'connected'}
						Connected
					{:else if stats.apiStatus === 'error'}
						Error
					{:else}
						Checking
					{/if}
				</div>
				<div class="stat-label">API Status</div>
			</div>
		</div>
	</div>

	<!-- Quick Actions -->
	<div style="margin-top: 2rem;">
		<h2 style="font-size: 1.25rem; font-weight: 600; margin-bottom: 1rem; color: var(--text-primary);">
			Quick Actions
		</h2>
		
		<div class="actions-grid">
			{#each quickActions as action}
				<button 
					class="action-card {action.primary ? 'primary' : ''}"
					on:click={action.action}
				>
					<div class="action-icon">{action.icon}</div>
					<div class="action-content">
						<div class="action-title">{action.title}</div>
						<div class="action-description">{action.description}</div>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Loading States -->
	{#if $loading.feeds || $loading.articles}
		<div style="margin-top: 2rem; padding: 1rem; background: var(--bg-secondary); border-radius: var(--radius); text-align: center;">
			<p style="color: var(--text-secondary);">Loading data...</p>
		</div>
	{/if}

	<!-- Error States -->
	{#if $errors.feeds || $errors.articles}
		<div style="margin-top: 2rem; padding: 1rem; background: #fef2f2; border: 1px solid #fecaca; border-radius: var(--radius);">
			<h3 style="color: #dc2626; margin-bottom: 0.5rem;">Connection Issue</h3>
			<p style="color: #7f1d1d; margin-bottom: 1rem;">
				{$errors.feeds || $errors.articles}
			</p>
			<button class="btn btn-secondary" on:click={loadDashboardData}>
				Retry Connection
			</button>
		</div>
	{/if}
</div>

<style>
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
		margin-bottom: 2rem;
	}

	.stat-card {
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.stat-icon {
		font-size: 2rem;
		width: 3rem;
		height: 3rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--primary-light);
		border-radius: var(--radius);
	}

	.stat-content {
		flex: 1;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-primary);
		line-height: 1;
	}

	.stat-label {
		font-size: 0.875rem;
		color: var(--text-secondary);
		margin-top: 0.25rem;
	}

	.actions-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 1rem;
	}

	.action-card {
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		cursor: pointer;
		transition: all 0.15s ease;
		text-align: left;
	}

	.action-card:hover {
		background: var(--bg-tertiary);
		border-color: var(--primary);
		transform: translateY(-1px);
	}

	.action-card.primary {
		border-color: var(--primary);
		background: var(--primary-light);
	}

	.action-icon {
		font-size: 1.5rem;
		width: 2.5rem;
		height: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--bg-primary);
		border-radius: var(--radius);
	}

	.action-content {
		flex: 1;
	}

	.action-title {
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-primary);
		margin-bottom: 0.25rem;
	}

	.action-description {
		font-size: 0.875rem;
		color: var(--text-secondary);
	}
</style>
