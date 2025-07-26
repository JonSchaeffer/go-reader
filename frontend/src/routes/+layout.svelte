<script>
	import '../app.css';
</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="bg-color-primary-900 grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Grid Column -->
	<div class="grid grid-cols-1 md:grid-cols-[auto_1fr_auto]">
		<!-- Sidebar (Left) -->
		<aside class="bg-surface-740 text-surface-100 p-4">RSS Reader</aside>
		<!-- Main -->
		<main class="bg-surface-50 space-y-4 p-4">
			<p class="bg-surface-700 text-surface-750 h-[512px] p-4">Paragraph 1</p>
			<p class="bg-surface-700 text-surface-750 h-[512px] p-4">Paragraph 2</p>
			<p class="bg-surface-700 text-surface-750 h-[512px] p-4">Paragraph 3</p>
			<slot />
		</main>
		<!-- Sidebar (Right) -->
		<aside class="bg-surface-900 text-surface-900 p-4">(sidebar)</aside>
	</div>
	<!-- Footer -->
	<footer class="bg-surface-900 text-surface-100 p-4">(footer)</footer>
</div>
