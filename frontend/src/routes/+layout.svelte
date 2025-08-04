<script>
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Feed from '$lib/components/Feed.svelte';
	
	let showRightSidebar = false;
</script>

<svelte:head>
	<title>YARR</title>
	<meta name="description" content="Yet Another RSS Reader" />
</svelte:head>

<div class="bg-primary-900 grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Grid Column -->
	<div class="grid grid-cols-[auto_1fr_auto]">
		<!-- Sidebar (Left) -->
		<Sidebar />

		<!-- Main -->
		<Feed bind:showRightSidebar />
		<!-- Sidebar (Right) -->
		{#if showRightSidebar}
			<aside class="text-surface-100 bg-slate-800 p-4 border-l border-slate-700">(sidebar)</aside>
		{/if}
	</div>
	<!-- Footer -->
	<footer class="text-surface-100 bg-slate-800 p-4">(footer)</footer>
</div>
