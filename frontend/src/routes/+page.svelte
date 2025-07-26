<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		// Redirect to articles page immediately
		goto('/articles');
	});
</script>

<div class="container mx-auto p-8 flex items-center justify-center min-h-screen">
	<div class="text-center">
		<div class="text-4xl mb-4">📄</div>
		<p>Loading RSS Reader...</p>
	</div>
</div>
