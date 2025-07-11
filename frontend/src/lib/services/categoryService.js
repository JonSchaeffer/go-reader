import { categoryApi } from '../api.js';
import { categories } from '../stores.js';

/**
 * Service for managing categories with state management
 */
export class CategoryService {
	/**
	 * Load all categories
	 */
	static async loadCategories() {

		try {
			const response = await categoryApi.getAll();
			console.log('Categories API response:', response);
			
			const categoryList = Array.isArray(response) ? response : [];
			categories.set(categoryList);
			return categoryList;
		} catch (error) {
			throw error;
		}
	}

	/**
	 * Create a new category
	 */
	static async createCategory(name, color = '#3b82f6') {

		try {
			const newCategory = await categoryApi.create(name, color);
			console.log('Category created:', newCategory);
			
			// Add to local state
			categories.update(currentCategories => [...currentCategories, newCategory]);
			
			return newCategory;
		} catch (error) {
			console.error('Failed to create category:', error);
			throw error;
		}
	}

	/**
	 * Update an existing category
	 */
	static async updateCategory(id, name, color) {

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
			throw error;
		}
	}

	/**
	 * Delete a category
	 */
	static async deleteCategory(id, categoryName = 'this category') {

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
			throw error;
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
