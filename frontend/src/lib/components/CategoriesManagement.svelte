<script>
	import { onMount } from 'svelte';
	import { Trash2, Plus, Check, X, Tag } from '@lucide/svelte';
	import { CategoryService } from '$lib/services/categoryService';
	import { categories, feeds } from '$lib/stores';
	import { FeedService } from '$lib/services/feedService';

	const colorPresets = [
		'#ef4444', '#f97316', '#eab308', '#22c55e',
		'#06b6d4', '#3b82f6', '#8b5cf6', '#ec4899',
		'#64748b', '#78716c'
	];

	let newName = '';
	let newColor = colorPresets[5]; // blue default
	let editingId = null;
	let editName = '';
	let editColor = '';

	onMount(async () => {
		await Promise.all([CategoryService.loadCategories(), FeedService.loadFeeds()]);
	});

	async function createCategory() {
		if (!newName.trim()) return;
		await CategoryService.createCategory(newName.trim(), newColor);
		newName = '';
		newColor = colorPresets[5];
	}

	function startEdit(cat) {
		editingId = cat.ID;
		editName = cat.Name;
		editColor = cat.Color || colorPresets[5];
	}

	async function saveEdit() {
		await CategoryService.updateCategory(editingId, editName, editColor);
		editingId = null;
	}

	async function deleteCategory(id) {
		if (!confirm('Delete this category? Feeds will become uncategorized.')) return;
		await CategoryService.deleteCategory(id);
	}

	function feedCountForCategory(catId) {
		return $feeds.filter((f) => f.CategoryID === catId).length;
	}
</script>

<div class="flex flex-1 flex-col overflow-hidden bg-slate-900">
	<!-- Header -->
	<div class="border-b border-white/10 px-8 py-5">
		<h2 class="text-base font-semibold text-surface-100">Categories</h2>
		<p class="mt-0.5 text-xs text-surface-500">Organise your feeds into groups</p>
	</div>

	<div class="flex-1 overflow-y-auto">
		<div class="mx-auto max-w-xl px-8 py-6 space-y-8">

			<!-- Add new category -->
			<div>
				<p class="mb-3 text-xs font-medium uppercase tracking-wide text-surface-500">New Category</p>
				<form on:submit|preventDefault={createCategory} class="space-y-3">
					<!-- Name input -->
					<input
						bind:value={newName}
						placeholder="Category name"
						class="w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-surface-100 placeholder-surface-500 focus:border-blue-400 focus:outline-none"
					/>
					<!-- Color swatches -->
					<div class="flex items-center gap-2 flex-wrap">
						{#each colorPresets as color}
							<button
								type="button"
								on:click={() => (newColor = color)}
								class="h-6 w-6 rounded-full transition-transform hover:scale-110 {newColor === color ? 'ring-2 ring-white ring-offset-2 ring-offset-slate-900' : ''}"
								style="background-color: {color}"
							></button>
						{/each}
						<!-- Custom color fallback -->
						<label class="relative h-6 w-6 cursor-pointer overflow-hidden rounded-full border border-dashed border-slate-500 hover:border-slate-300" title="Custom colour">
							<input type="color" bind:value={newColor} class="absolute inset-0 h-full w-full cursor-pointer opacity-0" />
							<span class="flex h-full w-full items-center justify-center text-[10px] text-surface-500">+</span>
						</label>
					</div>
					<button
						type="submit"
						class="flex items-center gap-1.5 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-500 transition-colors"
					>
						<Plus size={13} />
						Add Category
					</button>
				</form>
			</div>

			<!-- Existing categories -->
			{#if $categories.length > 0}
				<div>
					<p class="mb-3 text-xs font-medium uppercase tracking-wide text-surface-500">Your Categories</p>
					<ul class="divide-y divide-white/5 rounded-xl border border-white/5 overflow-hidden">
						{#each $categories as cat (cat.ID)}
							<li class="bg-slate-800/50">
								{#if editingId === cat.ID}
									<!-- Edit mode -->
									<div class="px-4 py-3 space-y-3">
										<input
											bind:value={editName}
											class="w-full rounded-md border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-surface-100 focus:border-blue-400 focus:outline-none"
										/>
										<div class="flex items-center gap-2 flex-wrap">
											{#each colorPresets as color}
												<button
													type="button"
													on:click={() => (editColor = color)}
													class="h-6 w-6 rounded-full transition-transform hover:scale-110 {editColor === color ? 'ring-2 ring-white ring-offset-2 ring-offset-slate-800' : ''}"
													style="background-color: {color}"
												></button>
											{/each}
											<label class="relative h-6 w-6 cursor-pointer overflow-hidden rounded-full border border-dashed border-slate-500 hover:border-slate-300" title="Custom colour">
												<input type="color" bind:value={editColor} class="absolute inset-0 h-full w-full cursor-pointer opacity-0" />
												<span class="flex h-full w-full items-center justify-center text-[10px] text-surface-500">+</span>
											</label>
										</div>
										<div class="flex gap-2">
											<button
												on:click={saveEdit}
												class="flex items-center gap-1 rounded-md bg-blue-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-blue-500 transition-colors"
											>
												<Check size={12} /> Save
											</button>
											<button
												on:click={() => (editingId = null)}
												class="flex items-center gap-1 rounded-md px-3 py-1.5 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
											>
												<X size={12} /> Cancel
											</button>
										</div>
									</div>
								{:else}
									<!-- View mode -->
									{@const count = feedCountForCategory(cat.ID)}
									<div class="group flex items-center gap-3 px-4 py-3">
										<div
											class="h-3 w-3 flex-shrink-0 rounded-full"
											style="background-color: {cat.Color}"
										></div>
										<span class="flex-1 text-sm text-surface-100">{cat.Name}</span>
										{#if count > 0}
											<span class="text-xs text-surface-500">{count} feed{count === 1 ? '' : 's'}</span>
										{/if}
										<div class="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
											<button
												on:click={() => startEdit(cat)}
												class="rounded px-2 py-1 text-xs text-surface-400 hover:bg-slate-700 hover:text-surface-100 transition-colors"
											>
												Edit
											</button>
											<button
												on:click={() => deleteCategory(cat.ID)}
												class="rounded p-1 text-surface-600 hover:bg-slate-700 hover:text-red-400 transition-colors"
											>
												<Trash2 size={13} />
											</button>
										</div>
									</div>
								{/if}
							</li>
						{/each}
					</ul>
				</div>
			{:else}
				<div class="flex flex-col items-center justify-center gap-3 py-12 text-center">
					<Tag size={28} class="text-surface-600" />
					<p class="text-sm text-surface-500">No categories yet.</p>
				</div>
			{/if}

		</div>
	</div>
</div>
