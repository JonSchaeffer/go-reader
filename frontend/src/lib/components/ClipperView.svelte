<script>
	import { Bookmark, Copy, Check, ExternalLink } from '@lucide/svelte';

	// Default to the same API base the app uses
	const DEFAULT_API = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api';

	let apiUrl = DEFAULT_API;
	let copied = false;

	// Minified bookmarklet JS with the API URL injected
	$: bookmarkletCode = `javascript:(function(){var API='${apiUrl}';var url=location.href;var d=document.createElement('div');d.style.cssText='position:fixed;top:16px;right:16px;z-index:99999;background:#1e293b;color:#f1f5f9;padding:12px 16px;border-radius:8px;font:14px/1.4 system-ui;box-shadow:0 4px 12px rgba(0,0,0,.5);border:1px solid #334155;transition:border-color .2s';d.textContent='Saving to YARR\u2026';document.body.appendChild(d);fetch(API+'/library',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:url})}).then(function(r){if(r.ok){d.textContent='\u2713 Saved to YARR!';d.style.borderColor='#22c55e';}else if(r.status===409){d.textContent='\u2022 Already in your library';d.style.borderColor='#f59e0b';}else{d.textContent='\u2717 Failed to save ('+r.status+')';d.style.borderColor='#ef4444';}setTimeout(function(){d.remove();},2500);}).catch(function(e){d.textContent='\u2717 Error: '+e.message;d.style.borderColor='#ef4444';setTimeout(function(){d.remove();},3000);});})();`;

	async function copyCode() {
		await navigator.clipboard.writeText(bookmarkletCode);
		copied = true;
		setTimeout(() => copied = false, 2000);
	}
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="flex items-center gap-3 border-b border-white/10 px-4 py-3">
		<Bookmark size={15} class="text-surface-400" />
		<span class="text-sm font-semibold text-surface-100">Web Clipper</span>
		<span class="rounded bg-blue-500/20 px-1.5 py-0.5 text-[10px] font-medium text-blue-400">Beta</span>
	</div>

	<div class="flex-1 overflow-y-auto px-6 py-6 space-y-8 max-w-xl">

		<!-- Step 1: Configure API URL -->
		<section>
			<h3 class="mb-1 text-sm font-semibold text-surface-100">1. Your backend URL</h3>
			<p class="mb-3 text-xs text-surface-400">
				Make sure this matches where your YARR backend is running. For local dev it's already set correctly.
			</p>
			<input
				bind:value={apiUrl}
				class="w-full rounded-md border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-surface-100 placeholder-surface-500 focus:border-blue-400 focus:outline-none"
				placeholder="http://localhost:8080/api"
			/>
		</section>

		<!-- Step 2: Drag to bookmarks -->
		<section>
			<h3 class="mb-1 text-sm font-semibold text-surface-100">2. Add to your bookmarks bar</h3>
			<p class="mb-3 text-xs text-surface-400">
				Drag the button below to your browser's bookmarks bar. Make sure the bookmarks bar is visible
				(<span class="text-surface-300">View → Show Bookmarks Bar</span> or <kbd class="rounded bg-slate-700 px-1 py-0.5 text-[10px]">⌘⇧B</kbd>).
			</p>
			<div class="flex items-center gap-3">
				<!-- svelte-ignore a11y-missing-attribute -->
				<a
					href={bookmarkletCode}
					class="flex items-center gap-2 rounded-lg border-2 border-dashed border-slate-500 bg-slate-800 px-4 py-3 text-sm font-medium text-surface-100 hover:border-blue-400 hover:text-blue-400 transition-colors cursor-grab select-none"
					on:click|preventDefault
					draggable="true"
				>
					<Bookmark size={15} />
					Save to YARR
				</a>
				<span class="text-xs text-surface-500">← drag this to your bookmarks bar</span>
			</div>
		</section>

		<!-- Step 3: How to use -->
		<section>
			<h3 class="mb-1 text-sm font-semibold text-surface-100">3. Clip any page</h3>
			<p class="text-xs text-surface-400 leading-relaxed">
				Navigate to any article or webpage you want to save, then click
				<span class="rounded bg-slate-700 px-1.5 py-0.5 text-surface-200">Save to YARR</span>
				in your bookmarks bar. A small notification will confirm the save. The page will be fetched with full
				content extraction and added to your Library.
			</p>
		</section>

		<!-- Advanced: copy raw code -->
		<section class="border-t border-white/10 pt-6">
			<h3 class="mb-1 text-sm font-semibold text-surface-100">Advanced — copy raw code</h3>
			<p class="mb-3 text-xs text-surface-400">
				If drag-and-drop doesn't work, create a new bookmark manually and paste this as the URL.
			</p>
			<div class="relative">
				<pre class="overflow-x-auto rounded-lg border border-slate-700 bg-slate-800 p-3 pr-10 text-[11px] text-surface-400 whitespace-pre-wrap break-all leading-relaxed">{bookmarkletCode}</pre>
				<button
					on:click={copyCode}
					class="absolute right-2 top-2 rounded p-1.5 text-surface-500 hover:bg-slate-700 hover:text-surface-100 transition-colors"
					title="Copy to clipboard"
				>
					{#if copied}
						<Check size={13} class="text-green-400" />
					{:else}
						<Copy size={13} />
					{/if}
				</button>
			</div>
		</section>

	</div>
</div>
