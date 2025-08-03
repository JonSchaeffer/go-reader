<script>
	import { onMount } from 'svelte';
	import { ArticleService } from '$lib/services/articleService';

	let feed = [];
	let loading = true;
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

<div class="flex-1 space-y-2 bg-slate-900 p-4">
	<div class="space-y-1">
		{#each feed.articles as article}
			<div
				class="card preset-filled-surface-100-600 cursor-pointer transition-colors hover:bg-white/10"
			>
				<article class="flex items-start p-3">
					<!-- Unread notification dot -->
					{#if !article.Read}
						<div class="mt-2 mr-3 flex-shrink-0">
							<div class="h-2 w-2 rounded-full bg-blue-400"></div>
						</div>
					{/if}

					<!-- Content -->
					<div class="min-w-0 flex-1">
						<!-- Title -->
						<header>
							<h5 class="mb-1 line-clamp-2 text-base leading-tight font-medium">
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
