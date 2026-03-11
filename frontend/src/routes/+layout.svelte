<script>
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Feed from '$lib/components/Feed.svelte';
	import ArticleDetail from '$lib/components/ArticleDetail.svelte';
	import ArticleModal from '$lib/components/ArticleModal.svelte';
	import LibraryFeed from '$lib/components/LibraryFeed.svelte';
	import HighlightsView from '$lib/components/HighlightsView.svelte';
	import ClipperView from '$lib/components/ClipperView.svelte';
	import SettingsView from '$lib/components/SettingsView.svelte';
	import ReviewView from '$lib/components/ReviewView.svelte';
	import { selectedArticle, currentView, articleModalOpen, openMode, reviewCount } from '$lib/stores';
	import { reviewApi } from '$lib/api';
	import { onMount } from 'svelte';

	onMount(async () => {
		const res = await reviewApi.getCount();
		if (res?.count) reviewCount.set(res.count);
	});
</script>

<svelte:head>
	<title>YARR</title>
	<meta name="description" content="Yet Another RSS Reader" />
</svelte:head>

<div class="flex h-screen overflow-hidden bg-slate-900 transition-[filter] duration-200 {$articleModalOpen ? 'blur-sm brightness-75' : ''}">
	<Sidebar />

	{#if $currentView === 'library'}
		<LibraryFeed />
	{:else if $currentView === 'highlights'}
		<HighlightsView />
	{:else if $currentView === 'review'}
		<ReviewView />
	{:else if $currentView === 'clipper'}
		<ClipperView />
	{:else if $currentView === 'settings'}
		<SettingsView />
	{:else}
		<Feed />
	{/if}

	{#if $selectedArticle && !$articleModalOpen}
		<ArticleDetail />
	{/if}
</div>

<ArticleModal />
