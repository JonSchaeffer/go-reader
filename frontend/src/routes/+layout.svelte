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
	import { selectedArticle, currentView, articleModalOpen, openMode } from '$lib/stores';
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
