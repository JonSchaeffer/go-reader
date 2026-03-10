<script>
	import { onMount } from 'svelte';
	import { X, Tag } from '@lucide/svelte';
	import { allTags } from '$lib/stores';
	import { TagService } from '$lib/services/tagService';

	export let itemType;   // 'library' | 'article'
	export let itemId;

	let itemTags = [];      // tags applied to this item
	let inputVal = '';
	let inputEl = null;
	let focused = false;

	$: suggestions = inputVal.trim().length > 0
		? $allTags.filter(
				(t) =>
					t.Name.toLowerCase().includes(inputVal.toLowerCase()) &&
					!itemTags.find((it) => it.ID === t.ID)
			)
		: [];

	onMount(async () => {
		await TagService.loadAll();
		itemTags = await TagService.getForItem(itemType, itemId);
	});

	async function addTag(tag) {
		if (itemTags.find((t) => t.ID === tag.ID)) return;
		await TagService.addToItem(itemType, itemId, tag.ID);
		itemTags = [...itemTags, tag];
		inputVal = '';
	}

	async function createAndAdd() {
		const name = inputVal.trim();
		if (!name) return;
		const existing = $allTags.find((t) => t.Name.toLowerCase() === name.toLowerCase());
		const tag = existing ?? (await TagService.create(name));
		await addTag(tag);
	}

	async function removeTag(tag) {
		await TagService.removeFromItem(itemType, itemId, tag.ID);
		itemTags = itemTags.filter((t) => t.ID !== tag.ID);
	}

	function onKeydown(e) {
		if (e.key === 'Enter') { e.preventDefault(); createAndAdd(); }
		if (e.key === 'Backspace' && inputVal === '' && itemTags.length > 0) {
			removeTag(itemTags[itemTags.length - 1]);
		}
		if (e.key === 'Escape') { inputVal = ''; focused = false; inputEl?.blur(); }
	}
</script>

<div class="flex flex-wrap items-center gap-1">
	<!-- Existing tag pills -->
	{#each itemTags as tag (tag.ID)}
		<span class="flex items-center gap-0.5 rounded-full bg-blue-900/60 px-2 py-0.5 text-[11px] text-blue-300">
			{tag.Name}
			<button
				on:click={() => removeTag(tag)}
				class="ml-0.5 rounded-full p-0.5 hover:bg-blue-700/60 transition-colors"
				title="Remove tag"
			><X size={9} /></button>
		</span>
	{/each}

	<!-- Input -->
	<div class="relative">
		<div class="flex items-center gap-0.5 rounded px-1.5 py-0.5 text-xs text-surface-500 hover:text-surface-300 transition-colors {focused ? 'text-surface-300' : ''}">
			{#if itemTags.length === 0 && !focused}
				<Tag size={11} />
				<span class="text-[11px]">Add tag</span>
			{/if}
			<input
				bind:this={inputEl}
				bind:value={inputVal}
				on:focus={() => (focused = true)}
				on:blur={() => setTimeout(() => (focused = false), 150)}
				on:keydown={onKeydown}
				placeholder={focused || itemTags.length > 0 ? 'Add tag…' : ''}
				class="w-20 bg-transparent text-[11px] text-surface-200 placeholder-surface-600 outline-none"
			/>
		</div>

		<!-- Suggestions dropdown -->
		{#if focused && suggestions.length > 0}
			<div class="absolute left-0 top-full z-50 mt-1 w-40 rounded-lg border border-slate-600 bg-slate-800 py-1 shadow-xl">
				{#each suggestions as s (s.ID)}
					<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
					<div
						class="cursor-pointer px-3 py-1.5 text-xs text-surface-300 hover:bg-slate-700 hover:text-surface-100"
						on:mousedown|preventDefault={() => addTag(s)}
					>
						{s.Name}
					</div>
				{/each}
				{#if inputVal.trim() && !$allTags.find((t) => t.Name.toLowerCase() === inputVal.toLowerCase())}
					<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
					<div
						class="cursor-pointer px-3 py-1.5 text-xs text-blue-400 hover:bg-slate-700"
						on:mousedown|preventDefault={createAndAdd}
					>
						Create "{inputVal.trim()}"
					</div>
				{/if}
			</div>
		{/if}
		{#if focused && suggestions.length === 0 && inputVal.trim()}
			<div class="absolute left-0 top-full z-50 mt-1 w-40 rounded-lg border border-slate-600 bg-slate-800 py-1 shadow-xl">
				<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
				<div
					class="cursor-pointer px-3 py-1.5 text-xs text-blue-400 hover:bg-slate-700"
					on:mousedown|preventDefault={createAndAdd}
				>
					Create "{inputVal.trim()}"
				</div>
			</div>
		{/if}
	</div>
</div>
