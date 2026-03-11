<script>
	import { readerSettings, openMode } from '$lib/stores';
	import FeedsManagement from './FeedsManagement.svelte';
	import CategoriesManagement from './CategoriesManagement.svelte';

	let activeSection = 'reader'; // 'reader' | 'app' | 'feeds' | 'categories'

	const FONT_LABELS  = [['sm','S'], ['md','M'], ['lg','L'], ['xl','XL']];
	const FAMILIES     = [['sans','Sans'], ['serif','Serif'], ['mono','Mono']];
	const WIDTHS       = [['narrow','Narrow'], ['medium','Medium'], ['wide','Wide']];
	const THEMES       = [
		['dark',  'Dark',  'bg-slate-800 border-slate-600 text-slate-100'],
		['light', 'Light', 'bg-white border-slate-300 text-slate-900'],
		['sepia', 'Sepia', 'bg-[#f5edd8] border-[#c8b89a] text-[#433422]'],
	];

	function setReader(key, value) {
		readerSettings.update((s) => ({ ...s, [key]: value }));
	}

	const SECTIONS = [
		{ id: 'reader',     label: 'Reader' },
		{ id: 'app',        label: 'App' },
		{ id: 'feeds',      label: 'Manage Feeds' },
		{ id: 'categories', label: 'Categories' },
	];
</script>

<div class="flex flex-1 overflow-hidden bg-slate-900">
	<!-- Settings sidebar -->
	<nav class="w-44 flex-shrink-0 border-r border-slate-700 bg-slate-800/50 p-3 space-y-0.5">
		{#each SECTIONS as s}
			<button
				class="w-full rounded-md px-3 py-2 text-left text-sm transition-colors
					{activeSection === s.id
					? 'bg-slate-700 text-surface-100 font-medium'
					: 'text-surface-400 hover:bg-slate-700/50 hover:text-surface-200'}"
				on:click={() => (activeSection = s.id)}
			>{s.label}</button>
		{/each}
	</nav>

	<!-- Content area -->
	<div class="flex-1 overflow-y-auto">
		{#if activeSection === 'reader'}
			<div class="mx-auto max-w-lg px-8 py-8 space-y-8">
				<h2 class="text-lg font-semibold text-surface-100">Reader Settings</h2>

				<!-- Font size -->
				<div>
					<p class="mb-2 text-xs font-semibold uppercase tracking-wide text-surface-500">Font Size</p>
					<div class="flex gap-2">
						{#each FONT_LABELS as [key, label]}
							<button
								on:click={() => setReader('fontSize', key)}
								class="flex-1 rounded-lg py-2 text-sm font-medium transition-colors
									{$readerSettings.fontSize === key
									? 'bg-blue-600 text-white'
									: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
							>{label}</button>
						{/each}
					</div>
				</div>

				<!-- Font family -->
				<div>
					<p class="mb-2 text-xs font-semibold uppercase tracking-wide text-surface-500">Font</p>
					<div class="flex gap-2">
						{#each FAMILIES as [key, label]}
							<button
								on:click={() => setReader('fontFamily', key)}
								class="flex-1 rounded-lg py-2 text-sm transition-colors
									{$readerSettings.fontFamily === key
									? 'bg-blue-600 text-white font-medium'
									: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
								style={key === 'serif' ? 'font-family: Georgia, serif' : key === 'mono' ? 'font-family: monospace' : ''}
							>{label}</button>
						{/each}
					</div>
				</div>

				<!-- Line width -->
				<div>
					<p class="mb-2 text-xs font-semibold uppercase tracking-wide text-surface-500">Line Width</p>
					<div class="flex gap-2">
						{#each WIDTHS as [key, label]}
							<button
								on:click={() => setReader('lineWidth', key)}
								class="flex-1 rounded-lg py-2 text-sm transition-colors
									{$readerSettings.lineWidth === key
									? 'bg-blue-600 text-white font-medium'
									: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
							>{label}</button>
						{/each}
					</div>
				</div>

				<!-- Theme -->
				<div>
					<p class="mb-2 text-xs font-semibold uppercase tracking-wide text-surface-500">Theme</p>
					<div class="flex gap-2">
						{#each THEMES as [key, label, cls]}
							<button
								on:click={() => setReader('theme', key)}
								class="flex-1 rounded-lg border py-2 text-sm font-medium transition-all {cls}
									{$readerSettings.theme === key ? 'ring-2 ring-blue-400' : 'opacity-70 hover:opacity-100'}"
							>{label}</button>
						{/each}
					</div>
				</div>

				<p class="text-xs text-surface-600">Settings are saved automatically to your browser.</p>
			</div>

		{:else if activeSection === 'app'}
			<div class="mx-auto max-w-lg px-8 py-8 space-y-8">
				<h2 class="text-lg font-semibold text-surface-100">App Settings</h2>

				<!-- Default open mode -->
				<div>
					<p class="mb-1 text-sm font-medium text-surface-200">Default Article View</p>
					<p class="mb-3 text-xs text-surface-500">How articles open when you click them.</p>
					<div class="flex gap-2">
						<button
							on:click={() => openMode.set('sidebar')}
							class="flex-1 rounded-lg py-2 text-sm transition-colors
								{$openMode === 'sidebar'
								? 'bg-blue-600 text-white font-medium'
								: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
						>Sidebar Panel</button>
						<button
							on:click={() => openMode.set('modal')}
							class="flex-1 rounded-lg py-2 text-sm transition-colors
								{$openMode === 'modal'
								? 'bg-blue-600 text-white font-medium'
								: 'bg-slate-700 text-surface-300 hover:bg-slate-600'}"
						>Full Screen Modal</button>
					</div>
				</div>
			</div>

		{:else if activeSection === 'feeds'}
			<FeedsManagement embedded />

		{:else if activeSection === 'categories'}
			<CategoriesManagement embedded />
		{/if}
	</div>
</div>
