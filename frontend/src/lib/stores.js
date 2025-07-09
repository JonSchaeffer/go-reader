import { writable } from 'svelte/store';

/**
 * RSS Feeds store
 */
export const feeds = writable([]);

/**
 * Articles store
 */
export const articles = writable([]);

/**
 * Categories store
 */
export const categories = writable([]);

/**
 * Selected feed store
 */
export const selectedFeed = writable(null);

/**
 * Loading states
 */
export const loading = writable({
	feeds: false,
	articles: false,
	categories: false,
	adding: false,
	deleting: false
});

/**
 * Error states
 */
export const errors = writable({
	feeds: null,
	articles: null,
	categories: null,
	general: null
});

/**
 * Search query
 */
export const searchQuery = writable('');

/**
 * UI preferences
 */
export const preferences = writable({
	theme: 'light',
	articlesPerPage: 20,
	showReadArticles: true
});

/**
 * Helper functions for managing loading states
 */
export const setLoading = (key, value) => {
	loading.update(state => ({ ...state, [key]: value }));
};

export const setError = (key, value) => {
	errors.update(state => ({ ...state, [key]: value }));
};

export const clearErrors = () => {
	errors.set({ feeds: null, articles: null, categories: null, general: null });
};