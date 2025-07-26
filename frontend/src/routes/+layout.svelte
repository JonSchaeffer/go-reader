<script>
	import '../app.css';
</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="h-screen flex flex-col">
	<!-- Header -->
	<header class="bg-surface-50-900-token border-b border-surface-300-600-token">
		<div class="flex items-center justify-between p-4">
			<a href="/" class="text-xl font-bold">RSS Reader</a>
			<button class="btn btn-sm variant-ghost-surface">🌙</button>
		</div>
	</header>

	<!-- Main Content -->
	<main class="flex-1 overflow-hidden">
		<slot />
	</main>
</div>
