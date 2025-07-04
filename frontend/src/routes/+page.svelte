<script>
	import { onMount } from 'svelte';
	import { FeedService } from '$lib/services/feedService.js';
	import { ArticleService } from '$lib/services/articleService.js';
	import { feeds, articles, loading, errors } from '$lib/stores.js';

	let testResults = {
		apiConnection: 'Not tested',
		feedsLoad: 'Not tested',
		articlesLoad: 'Not tested'
	};

	onMount(async () => {
		await testApiConnection();
	});

	async function testApiConnection() {
		console.log('Testing API connection...');
		
		try {
			// Test basic API connectivity
			const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api';
			const response = await fetch(`${apiBase}/rss`);
			if (response.ok) {
				testResults.apiConnection = '✅ Connected';
			} else {
				testResults.apiConnection = `❌ HTTP ${response.status}`;
			}
		} catch (error) {
			testResults.apiConnection = `❌ ${error.message}`;
		}

		// Test feeds loading
		try {
			await FeedService.loadFeeds();
			testResults.feedsLoad = '✅ Feeds loaded';
		} catch (error) {
			testResults.feedsLoad = `❌ ${error.message}`;
		}

		// Test articles loading
		try {
			await ArticleService.loadAllArticles();
			testResults.articlesLoad = '✅ Articles loaded';
		} catch (error) {
			testResults.articlesLoad = `❌ ${error.message}`;
		}

		// Trigger reactivity
		testResults = { ...testResults };
	}

	async function addTestFeed() {
		try {
			await FeedService.addFeed('https://feeds.feedburner.com/oreilly');
			testResults.feedsLoad = '✅ Test feed added';
		} catch (error) {
			testResults.feedsLoad = `❌ Failed to add feed: ${error.message}`;
		}
		testResults = { ...testResults };
	}
</script>

<svelte:head>
	<title>RSS Reader - API Test</title>
</svelte:head>

<main style="padding: 2rem; font-family: system-ui;">
	<h1>RSS Reader - API Integration Test</h1>
	
	<div style="margin: 2rem 0;">
		<h2>API Connection Status</h2>
		<ul>
			<li><strong>Backend Connection:</strong> {testResults.apiConnection}</li>
			<li><strong>Feeds API:</strong> {testResults.feedsLoad}</li>
			<li><strong>Articles API:</strong> {testResults.articlesLoad}</li>
		</ul>
	</div>

	<div style="margin: 2rem 0;">
		<h2>Current Data</h2>
		<p><strong>Feeds Count:</strong> {$feeds.length}</p>
		<p><strong>Articles Count:</strong> {$articles.length}</p>
		
		{#if $loading.feeds}
			<p>Loading feeds...</p>
		{/if}
		
		{#if $loading.articles}
			<p>Loading articles...</p>
		{/if}
		
		{#if $errors.feeds}
			<p style="color: red;"><strong>Feed Error:</strong> {$errors.feeds}</p>
		{/if}
		
		{#if $errors.articles}
			<p style="color: red;"><strong>Article Error:</strong> {$errors.articles}</p>
		{/if}
	</div>

	<div style="margin: 2rem 0;">
		<h2>Test Actions</h2>
		<button on:click={testApiConnection} style="margin-right: 1rem;">
			Retest Connection
		</button>
		<button on:click={addTestFeed}>
			Add Test Feed
		</button>
	</div>

	{#if $feeds.length > 0}
		<div style="margin: 2rem 0;">
			<h2>Current Feeds</h2>
			<ul>
				{#each $feeds as feed}
					<li>
						<strong>{feed.title || 'Untitled'}</strong> - {feed.url}
						<small>(ID: {feed.id})</small>
					</li>
				{/each}
			</ul>
		</div>
	{/if}

	{#if $articles.length > 0}
		<div style="margin: 2rem 0;">
			<h2>Recent Articles ({$articles.slice(0, 5).length} of {$articles.length})</h2>
			<ul>
				{#each $articles.slice(0, 5) as article}
					<li>
						<strong>{article.Title || 'Untitled'}</strong>
						{#if article.Read}
							<span style="color: gray;">(Read)</span>
						{:else}
							<span style="color: blue;">(Unread)</span>
						{/if}
						<br>
						<small>{ArticleService.formatDate(article.PublishDate)}</small>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</main>

<style>
	button {
		padding: 0.5rem 1rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		background: white;
		cursor: pointer;
	}
	
	button:hover {
		background: #f5f5f5;
	}
	
	ul {
		margin: 1rem 0;
		padding-left: 2rem;
	}
	
	li {
		margin: 0.5rem 0;
	}
</style>
