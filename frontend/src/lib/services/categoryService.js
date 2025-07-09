import { categoryApi } from '../api.js';
import { categories, setLoading, setError } from '../stores.js';

/**
 * Service for managing categories with state management
 */
export class CategoryService {
	/**
	 * Load all categories
	 */
	static async loadCategories() {
		setLoading('categories', true);
		setError('categories', null);

		try {
			const response = await categoryApi.getAll();
			console.log('Categories API response:', response);
			
			const categoryList = Array.isArray(response) ? response : [];
			categories.set(categoryList);
			return categoryList;
		} catch (error) {
			console.error('Failed to load categories:', error);
			setError('categories', 'Failed to load categories. Please check if the backend is running.');
			throw error;
		} finally {
			setLoading('categories', false);
		}
	}

	/**
	 * Create a new category
	 */
	static async createCategory(name, color = '#3b82f6') {
		setLoading('adding', true);
		setError('categories', null);

		try {
			const newCategory = await categoryApi.create(name, color);
			console.log('Category created:', newCategory);
			
			// Add to local state
			categories.update(currentCategories => [...currentCategories, newCategory]);
			
			return newCategory;
		} catch (error) {
			console.error('Failed to create category:', error);
			setError('categories', 'Failed to create category.');
			throw error;
		} finally {
			setLoading('adding', false);
		}
	}

	/**
	 * Update an existing category
	 */
	static async updateCategory(id, name, color) {
		setLoading('categories', true);
		setError('categories', null);

		try {
			await categoryApi.update(id, name, color);
			console.log('Category updated successfully');
			
			// Update local state
			categories.update(currentCategories =>
				currentCategories.map(category =>
					category.ID === id ? { ...category, Name: name, Color: color } : category
				)
			);
			
			return true;
		} catch (error) {
			console.error('Failed to update category:', error);
			setError('categories', 'Failed to update category.');
			throw error;
		} finally {
			setLoading('categories', false);
		}
	}

	/**
	 * Delete a category
	 */
	static async deleteCategory(id, categoryName = 'this category') {
		setLoading('deleting', true);
		setError('categories', null);

		try {
			await categoryApi.delete(id);
			console.log('Category deleted successfully');
			
			// Remove from local state
			categories.update(currentCategories =>
				currentCategories.filter(category => category.ID !== id)
			);
			
			return true;
		} catch (error) {
			console.error('Failed to delete category:', error);
			setError('categories', `Failed to delete ${categoryName}.`);
			throw error;
		} finally {
			setLoading('deleting', false);
		}
	}

	/**
	 * Get category name by ID
	 */
	static getCategoryName(categoryId, categoriesList) {
		if (!categoryId || !categoriesList) return 'Uncategorized';
		
		const category = categoriesList.find(cat => cat.ID === categoryId);
		return category ? category.Name : 'Unknown Category';
	}

	/**
	 * Get category color by ID
	 */
	static getCategoryColor(categoryId, categoriesList) {
		if (!categoryId || !categoriesList) return '#64748b'; // Default gray
		
		const category = categoriesList.find(cat => cat.ID === categoryId);
		return category ? category.Color : '#64748b';
	}
}