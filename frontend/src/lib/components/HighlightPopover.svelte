<script>
	import { createEventDispatcher } from 'svelte';
	import { Trash2 } from '@lucide/svelte';

	export let x = 0;
	export let y = 0;
	export let existingHighlight = null; // if editing an existing one

	const dispatch = createEventDispatcher();

	const COLORS = [
		{ key: 'yellow', bg: '#fef08a' },
		{ key: 'green',  bg: '#bbf7d0' },
		{ key: 'blue',   bg: '#bfdbfe' },
		{ key: 'pink',   bg: '#fbcfe8' },
	];

	let selectedColor = existingHighlight?.Color ?? 'yellow';
	let note = existingHighlight?.Note ?? '';

	function pick(color) {
		selectedColor = color;
	}

	function save() {
		dispatch('save', { color: selectedColor, note });
	}

	function remove() {
		dispatch('delete');
	}

	function handleKeydown(e) {
		if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) save();
		if (e.key === 'Escape') dispatch('cancel');
	}
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="highlight-popover fixed z-[200] flex flex-col gap-2 rounded-xl border border-slate-600 bg-slate-800 p-3 shadow-2xl"
	style="left:{x}px; top:{y}px; min-width:200px"
	on:keydown={handleKeydown}
>
	<!-- Color swatches -->
	<div class="flex items-center gap-1.5">
		{#each COLORS as c}
			<button
				on:click={() => pick(c.key)}
				class="h-6 w-6 rounded-full border-2 transition-transform hover:scale-110 {selectedColor === c.key ? 'border-white' : 'border-transparent'}"
				style="background:{c.bg}"
				title={c.key}
			></button>
		{/each}
		{#if existingHighlight}
			<button
				on:click={remove}
				class="ml-auto rounded p-1 text-surface-500 hover:text-red-400 transition-colors"
				title="Remove highlight"
			>
				<Trash2 size={13} />
			</button>
		{/if}
	</div>

	<!-- Note input -->
	<textarea
		bind:value={note}
		placeholder="Add a note… (optional)"
		rows="2"
		class="w-full resize-none rounded-md border border-slate-600 bg-slate-700 px-2 py-1.5 text-xs text-surface-100 placeholder-surface-500 focus:border-blue-400 focus:outline-none"
	></textarea>

	<!-- Actions -->
	<div class="flex justify-end gap-2">
		<button
			on:click={() => dispatch('cancel')}
			class="rounded px-2 py-1 text-xs text-surface-400 hover:text-surface-100 transition-colors"
		>
			Cancel
		</button>
		<button
			on:click={save}
			class="rounded bg-blue-600 px-2.5 py-1 text-xs font-medium text-white hover:bg-blue-500 transition-colors"
		>
			{existingHighlight ? 'Update' : 'Highlight'}
		</button>
	</div>
</div>
