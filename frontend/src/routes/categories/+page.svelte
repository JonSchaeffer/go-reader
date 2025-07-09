<script>
	import { onMount } from 'svelte';
	import { CategoryService } from '$lib/services/categoryService.js';
	import { categories, loading, errors } from '$lib/stores.js';

	let showAddCategory = false;
	let editingCategory = null;
	let newCategoryName = '';
	let newCategoryColor = '#3b82f6';
	let searchTerm = '';

	// Default color options
	const colorOptions = [
		'#3b82f6', // Blue
		'#10b981', // Green  
		'#f59e0b', // Yellow
		'#ef4444', // Red
		'#8b5cf6', // Purple
		'#06b6d4', // Cyan
		'#f97316', // Orange
		'#84cc16', // Lime
		'#ec4899', // Pink
		'#64748b'  // Gray
	];

	onMount(async () => {
		console.log('Categories page mounted');
		await loadData();
	});

	async function loadData() {
		try {
			await CategoryService.loadCategories();
		} catch (error) {
			console.error('Failed to load categories:', error);
		}
	}

	async function handleAddCategory() {
		if (!newCategoryName.trim()) return;
		
		try {
			await CategoryService.createCategory(newCategoryName.trim(), newCategoryColor);
			newCategoryName = '';
			newCategoryColor = '#3b82f6';
			showAddCategory = false;
		} catch (error) {
			console.error('Failed to add category:', error);
		}
	}

	async function handleEditCategory() {
		if (!editingCategory || !newCategoryName.trim()) return;
		
		try {
			await CategoryService.updateCategory(editingCategory.ID, newCategoryName.trim(), newCategoryColor);
			editingCategory = null;
			newCategoryName = '';
			newCategoryColor = '#3b82f6';
		} catch (error) {
			console.error('Failed to update category:', error);
		}
	}

	async function handleDeleteCategory(category) {
		if (!confirm(`Are you sure you want to delete "${category.Name}"? This will remove the category from all feeds.`)) {
			return;
		}
		
		try {
			await CategoryService.deleteCategory(category.ID, category.Name);
		} catch (error) {
			console.error('Failed to delete category:', error);
		}
	}

	function startEdit(category) {
		editingCategory = category;
		newCategoryName = category.Name;
		newCategoryColor = category.Color;
		showAddCategory = true;
	}

	function cancelEdit() {
		editingCategory = null;
		newCategoryName = '';
		newCategoryColor = '#3b82f6';
		showAddCategory = false;
	}

	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			if (editingCategory) {
				handleEditCategory();
			} else {
				handleAddCategory();
			}
		}
	}

	// Filter categories based on search term
	$: filteredCategories = $categories.filter(category => 
		!searchTerm || 
		category.Name?.toLowerCase().includes(searchTerm.toLowerCase())
	);
</script>

<svelte:head>
	<title>Categories - RSS Reader</title>
</svelte:head>

<div class="content-header">
	<div>
		<h1 style="font-size: 1.875rem; font-weight: 700; color: var(--text-primary);">
			Categories
		</h1>
		<p style="color: var(--text-secondary); margin-top: 0.5rem;">
			Organize your RSS feeds into categories ({$categories.length} categories)
		</p>
	</div>
	
	<div class="header-actions">
		<button 
			class="btn btn-primary"
			on:click={() => { showAddCategory = true; editingCategory = null; }}
			disabled={$loading.adding}
		>
			{#if $loading.adding}
				⏳ Adding...
			{:else}
				📁 Add Category
			{/if}
		</button>
	</div>
</div>

<div class="content-body">
	<!-- Add/Edit Category Modal -->
	{#if showAddCategory}
		<div class="modal-overlay" on:click={cancelEdit} on:keydown={(e) => e.key === 'Escape' && cancelEdit()}>
			<div class="modal" on:click|stopPropagation>
				<div class="modal-header">
					<h3>{editingCategory ? 'Edit Category' : 'Add New Category'}</h3>
					<button class="btn-ghost" on:click={cancelEdit}>✕</button>
				</div>
				
				<div class="modal-body">
					<div class="form-group">
						<label for="categoryName">Category Name</label>
						<input 
							id="categoryName"
							type="text" 
							bind:value={newCategoryName}
							on:keypress={handleKeyPress}
							placeholder="e.g., Tech News, Personal Blogs"
							required
							class="form-input"
						/>
					</div>

					<div class="form-group">
						<label for="categoryColor">Color</label>
						<div class="color-picker">
							{#each colorOptions as color}
								<button 
									class="color-option {newCategoryColor === color ? 'selected' : ''}"
									style="background-color: {color}"
									on:click={() => newCategoryColor = color}
									title={color}
								></button>
							{/each}
						</div>
						<input 
							id="categoryColor"
							type="color" 
							bind:value={newCategoryColor}
							class="color-input"
						/>
					</div>
				</div>
				
				<div class="modal-footer">
					<button class="btn btn-secondary" on:click={cancelEdit}>Cancel</button>
					<button 
						class="btn btn-primary" 
						on:click={editingCategory ? handleEditCategory : handleAddCategory}
						disabled={!newCategoryName.trim() || $loading.adding}
					>
						{#if $loading.adding}
							⏳ {editingCategory ? 'Updating...' : 'Adding...'}
						{:else}
							{editingCategory ? 'Update Category' : 'Add Category'}
						{/if}
					</button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Search Bar -->
	{#if $categories.length > 0}
		<div class="search-bar">
			<div class="search-input-wrapper">
				<span class="search-icon">🔍</span>
				<input 
					type="text" 
					bind:value={searchTerm}
					placeholder="Search categories..."
					class="search-input"
				/>
				{#if searchTerm}
					<button class="btn-ghost search-clear" on:click={() => searchTerm = ''}>✕</button>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Loading State -->
	{#if $loading.categories}
		<div class="loading-state">
			<div class="loading-spinner">⏳</div>
			<p>Loading categories...</p>
		</div>

	<!-- Error State -->
	{:else if $errors.categories}
		<div class="error-state">
			<div class="error-icon">❌</div>
			<h3>Failed to Load Categories</h3>
			<p>{$errors.categories}</p>
			<button class="btn btn-primary" on:click={loadData}>Retry</button>
		</div>

	<!-- Empty State -->
	{:else if $categories.length === 0}
		<div class="empty-state">
			<div class="empty-icon">📁</div>
			<h3>No Categories Yet</h3>
			<p>Create categories to organize your RSS feeds into groups.</p>
			<button class="btn btn-primary" on:click={() => showAddCategory = true}>
				📁 Create Your First Category
			</button>
		</div>

	<!-- Categories List -->
	{:else}
		<div class="categories-grid">
			{#if filteredCategories.length === 0}
				<div class="no-results">
					<p>No categories match your search for "{searchTerm}"</p>
					<button class="btn btn-secondary" on:click={() => searchTerm = ''}>Clear Search</button>
				</div>
			{:else}
				{#each filteredCategories as category (category.ID)}
					<div class="category-card">
						<div class="category-header">
							<div class="category-icon" style="background-color: {category.Color}">
								📁
							</div>
							<div class="category-info">
								<h4 class="category-name">{category.Name}</h4>
								<p class="category-meta">
									Created {new Date(category.CreatedAt).toLocaleDateString()}
								</p>
							</div>
							<div class="category-actions">
								<button 
									class="btn-ghost category-action"
									on:click={() => startEdit(category)}
									title="Edit Category"
								>
									✏️
								</button>
								<button 
									class="btn-ghost category-action"
									on:click={() => handleDeleteCategory(category)}
									disabled={$loading.deleting}
									title="Delete Category"
								>
									🗑️
								</button>
							</div>
						</div>

						<div class="category-preview">
							<div class="color-preview" style="background-color: {category.Color}"></div>
							<span class="color-code">{category.Color}</span>
						</div>

						<div class="category-footer">
							<a href="/feeds?category={category.ID}" class="btn btn-secondary">
								View Feeds
							</a>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
</div>

<style>
	/* Header Actions */
	.content-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
	}

	.header-actions {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	/* Modal Styles */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
	}

	.modal {
		background: var(--bg-primary);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		width: 90%;
		max-width: 500px;
		max-height: 90vh;
		overflow: hidden;
		box-shadow: var(--shadow-lg);
	}

	.modal-header {
		padding: 1.5rem;
		border-bottom: 1px solid var(--border-light);
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.modal-header h3 {
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--text-primary);
		margin: 0;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.modal-footer {
		padding: 1rem 1.5rem;
		border-top: 1px solid var(--border-light);
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}

	/* Form Styles */
	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
		margin-bottom: 0.5rem;
	}

	.form-input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-primary);
		color: var(--text-primary);
		font-size: 0.875rem;
		transition: border-color 0.15s ease;
	}

	.form-input:focus {
		outline: none;
		border-color: var(--primary);
	}

	/* Color Picker */
	.color-picker {
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		gap: 0.5rem;
		margin-bottom: 0.75rem;
	}

	.color-option {
		width: 2.5rem;
		height: 2.5rem;
		border-radius: var(--radius);
		border: 2px solid transparent;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.color-option:hover {
		transform: scale(1.1);
		box-shadow: var(--shadow);
	}

	.color-option.selected {
		border-color: var(--text-primary);
		transform: scale(1.1);
		box-shadow: var(--shadow);
	}

	.color-input {
		width: 100%;
		height: 3rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-primary);
		cursor: pointer;
	}

	/* Search Bar */
	.search-bar {
		margin-bottom: 1.5rem;
		max-width: 400px;
	}

	.search-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.search-icon {
		position: absolute;
		left: 0.75rem;
		color: var(--text-tertiary);
		pointer-events: none;
	}

	.search-input {
		width: 100%;
		padding: 0.75rem 0.75rem 0.75rem 2.5rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg-secondary);
		color: var(--text-primary);
		font-size: 0.875rem;
		transition: border-color 0.15s ease;
	}

	.search-input:focus {
		outline: none;
		border-color: var(--primary);
	}

	.search-clear {
		position: absolute;
		right: 0.5rem;
		padding: 0.25rem;
		color: var(--text-tertiary);
	}

	/* State Styles */
	.loading-state,
	.error-state,
	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: var(--text-secondary);
	}

	.loading-spinner,
	.error-icon,
	.empty-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.error-state h3,
	.empty-state h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-primary);
		margin-bottom: 0.5rem;
	}

	.no-results {
		grid-column: 1 / -1;
		text-align: center;
		padding: 2rem;
		color: var(--text-secondary);
	}

	/* Categories Grid */
	.categories-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.category-card {
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 1.5rem;
		transition: all 0.15s ease;
	}

	.category-card:hover {
		border-color: var(--primary);
		box-shadow: var(--shadow);
	}

	.category-header {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.category-icon {
		font-size: 1.5rem;
		width: 2.5rem;
		height: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius);
		flex-shrink: 0;
		color: white;
	}

	.category-info {
		flex: 1;
		min-width: 0;
	}

	.category-name {
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-primary);
		margin: 0 0 0.25rem 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.category-meta {
		font-size: 0.75rem;
		color: var(--text-tertiary);
		margin: 0;
	}

	.category-actions {
		display: flex;
		gap: 0.25rem;
	}

	.category-action {
		padding: 0.375rem;
		font-size: 0.875rem;
	}

	.category-preview {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 1rem;
		padding: 0.75rem;
		background: var(--bg-primary);
		border-radius: var(--radius);
	}

	.color-preview {
		width: 1.5rem;
		height: 1.5rem;
		border-radius: var(--radius);
		flex-shrink: 0;
	}

	.color-code {
		font-family: monospace;
		font-size: 0.875rem;
		color: var(--text-secondary);
	}

	.category-footer {
		display: flex;
		gap: 0.75rem;
	}

	.category-footer .btn {
		flex: 1;
		justify-content: center;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.content-header {
			flex-direction: column;
			align-items: stretch;
		}

		.categories-grid {
			grid-template-columns: 1fr;
		}

		.color-picker {
			grid-template-columns: repeat(5, 1fr);
		}

		.modal {
			margin: 1rem;
			width: auto;
		}
	}
</style>