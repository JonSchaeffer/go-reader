<script>
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Feed from '$lib/components/Feed.svelte';
</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="bg-primary-900 grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Grid Column -->
	<div class="grid grid-cols-[auto_1fr_auto]">
		<!-- Sidebar (Left) -->
		<Sidebar />

		<!-- Main -->
		<Feed />
		<!-- Sidebar (Right) -->
		<aside class="bg-surface-900 text-surface-100 p-4">(sidebar)</aside>
	</div>
	<!-- Footer -->
	<footer class="bg-surface-900 text-surface-100 p-4">(footer)</footer>
</div>
