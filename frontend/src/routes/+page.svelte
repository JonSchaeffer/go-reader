<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		// Redirect to articles page immediately
		goto('/articles');
	});
</script>

<svelte:head>
	<title>RSS Reader</title>
</svelte:head>

<!-- Loading message while redirecting -->
<div style="display: flex; align-items: center; justify-content: center; height: 100vh; background: var(--bg-primary);">
	<div style="text-align: center; color: var(--text-secondary);">
		<div style="font-size: 2rem; margin-bottom: 1rem;">📄</div>
		<p>Loading RSS Reader...</p>
	</div>
</div>
