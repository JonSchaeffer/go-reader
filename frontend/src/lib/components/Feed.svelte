<script>
	import { onMount } from 'svelte';
	import { ArticleService } from '$lib/services/articleService';

	let feed = [];
	let loading = true;
	const feedData = [
		{
			id: 1,
			title: 'Fair Access to Banking',
			description: 'Comments',
			source: 'Hacker News',
			author: 'Hacker News',
			timeAgo: '6 mins',
			timestamp: '10:11 pm',
			category: 'FEED',
			unread: true
		},
		{
			id: 2,
			title: 'Federal judge delays expiration of TPS for Hondura...',
			description: 'A judge stopped the Trump administration from ending Tem...',
			source: 'News : NPR',
			author: 'Sergio Martinez-Beltrán',
			timeAgo: '2 mins',
			timestamp: '9:57 pm',
			category: 'FEED',
			unread: true
		},
		{
			id: 3,
			title: 'THPS3+4 - now do Tony Hawk 6',
			description: "Tony Hawk's Pro Skater 3+4 is a great re-imaging of the ori...",
			source: 'Jake Baldino',
			author: 'Jake Baldino',
			timeAgo: '2 mins',
			timestamp: '9:57 pm',
			category: 'FEED',
			unread: true
		}
	];
	onMount(async () => {
		try {
			feed = await ArticleService.loadAllArticles();
			loading = false;
		} catch (error) {
			console.error('Failed to load feed:', error);
			feed = [];
			loading = false;
		}
	});
</script>

<div class="bg-surface-900 flex-1 space-y-2 p-4">
	<div class="space-y-3">
		{#each feed.articles as article}
			<div
				class="card preset-filled-surface-100-900 hover:preset-tonal-primary cursor-pointer border-l-4 {article.unread
					? 'border-primary-500'
					: 'border-transparent'}"
			>
				<article class="flex items-start gap-3 p-4">
					<!-- Source Icon/Thumbnail -->
					<!-- <div class="flex-shrink-0"> -->
					<!-- 	<div class="w-12 h-12 rounded-lg flex items-center justify-center text-white font-bold text-lg" style="background-color: {article.sourceColor}"> -->
					<!-- 		{article.sourceIcon} -->
					<!-- 	</div> -->
					<!-- </div> -->

					<!-- Content -->
					<div class="min-w-0 flex-1">
						<!-- Title -->
						<header>
							<h3 class="mb-1 line-clamp-2 text-base leading-tight font-medium">
								{article.Title}
							</h3>
						</header>

						<!-- Description -->
						{#if article.description}
							<p class="text-surface-600-300 mb-2 line-clamp-2 text-sm">
								{article.description}
							</p>
						{/if}

						<!-- Source Info -->
						<footer class="text-surface-500-400 flex items-center gap-2 text-xs">
							<span class="text-surface-700-200">{article.source}</span>
							<span>•</span>
							<span>author</span>
							<span>•</span>
							<span>Timeago</span>
						</footer>
					</div>

					<!-- Right Side Info -->
					<div class="flex-shrink-0 text-right">
						<small class="text-surface-500-400 text-xs">
							category • {article.PublishDate}
						</small>
					</div>
				</article>
			</div>
		{/each}
	</div>
</div>
