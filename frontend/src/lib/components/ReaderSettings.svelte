<script>
	import { createEventDispatcher } from 'svelte';
	import { Settings } from '@lucide/svelte';
	import { readerSettings } from '$lib/stores';

	let open = false;
	const dispatch = createEventDispatcher();

	const FONT_SIZES  = [['sm','A',  'text-xs'], ['md','A', 'text-sm'], ['lg','A', 'text-base'], ['xl','A', 'text-lg']];
	const FONT_LABELS = [['sm','S'], ['md','M'], ['lg','L'], ['xl','XL']];
	const FAMILIES    = [['sans','Sans'], ['serif','Serif'], ['mono','Mono']];
	const WIDTHS      = [['narrow','Narrow'], ['medium','Medium'], ['wide','Wide']];
	const THEMES      = [
		['dark',  'Dark',  'bg-slate-800 border-slate-600 text-slate-100'],
		['light', 'Light', 'bg-white border-slate-300 text-slate-900'],
		['sepia', 'Sepia', 'bg-[#f5edd8] border-[#c8b89a] text-[#433422]'],
	];

	function set(key, value) {
		readerSettings.update((s) => ({ ...s, [key]: value }));
	}

	function handleClickOutside(e) {
		if (!e.target.closest('.reader-settings-panel') && !e.target.closest('.reader-settings-btn')) {
			open = false;
		}
	}
</script>

<svelte:window on:mousedown={handleClickOutside} />

<div class="relative">
	<button
		class="reader-settings-btn flex items-center gap-1 rounded px-2 py-1 text-xs text-surface-400 transition-colors hover:bg-slate-700 hover:text-surface-100 {open ? 'bg-slate-700 text-surface-100' : ''}"
		on:click={() => open = !open}
		title="Reader settings"
	>
		<Settings size={13} />
		<span>Settings</span>
	</button>

	{#if open}
		<div class="reader-settings-panel absolute right-0 top-full z-[300] mt-1 w-64 rounded-xl border border-slate-600 bg-slate-800 p-4 shadow-2xl space-y-4">

			<!-- Font size -->
			<div>
				<p class="mb-2 text-[10px] font-semibold uppercase tracking-wide text-surface-500">Font Size</p>
				<div class="flex gap-1">
					{#each FONT_LABELS as [key, label]}
						<button
							on:click={() => set('fontSize', key)}
							class="flex-1 rounded py-1 text-xs font-medium transition-colors
								{$readerSettings.fontSize === key
								? 'bg-blue-600 text-white'
								: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
						>{label}</button>
					{/each}
				</div>
			</div>

			<!-- Font family -->
			<div>
				<p class="mb-2 text-[10px] font-semibold uppercase tracking-wide text-surface-500">Font</p>
				<div class="flex gap-1">
					{#each FAMILIES as [key, label]}
						<button
							on:click={() => set('fontFamily', key)}
							class="flex-1 rounded py-1 text-xs transition-colors
								{$readerSettings.fontFamily === key
								? 'bg-blue-600 text-white font-medium'
								: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
							style={key === 'serif' ? "font-family: Georgia, serif" : key === 'mono' ? "font-family: monospace" : ""}
						>{label}</button>
					{/each}
				</div>
			</div>

			<!-- Line width -->
			<div>
				<p class="mb-2 text-[10px] font-semibold uppercase tracking-wide text-surface-500">Width</p>
				<div class="flex gap-1">
					{#each WIDTHS as [key, label]}
						<button
							on:click={() => set('lineWidth', key)}
							class="flex-1 rounded py-1 text-xs transition-colors
								{$readerSettings.lineWidth === key
								? 'bg-blue-600 text-white font-medium'
								: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
						>{label}</button>
					{/each}
				</div>
			</div>

			<!-- Theme -->
			<div>
				<p class="mb-2 text-[10px] font-semibold uppercase tracking-wide text-surface-500">Theme</p>
				<div class="flex gap-1">
					{#each THEMES as [key, label, cls]}
						<button
							on:click={() => set('theme', key)}
							class="flex-1 rounded border py-1 text-xs font-medium transition-all {cls}
								{$readerSettings.theme === key ? 'ring-2 ring-blue-400' : 'opacity-70 hover:opacity-100'}"
						>{label}</button>
					{/each}
				</div>
			</div>

		</div>
	{/if}
</div>
