import { writable } from 'svelte/store';

// Reader settings — persisted to localStorage
const READER_DEFAULTS = { fontSize: 'md', fontFamily: 'sans', lineWidth: 'medium', theme: 'dark' };
function createReaderSettings() {
	const stored = typeof localStorage !== 'undefined'
		? JSON.parse(localStorage.getItem('readerSettings') || 'null')
		: null;
	const store = writable({ ...READER_DEFAULTS, ...stored });
	store.subscribe((v) => {
		if (typeof localStorage !== 'undefined') localStorage.setItem('readerSettings', JSON.stringify(v));
	});
	return store;
}
export const readerSettings = createReaderSettings();

/**
 * RSS Feeds store
 */
export const feeds = writable([]);

/**
 * Currently selected feed ID (null = all feeds)
 */
export const selectedFeedId = writable(null);

/**
 * Library saved items
 */
export const savedItems = writable([]);

/**
 * Currently selected article (null = none)
 */
export const selectedArticle = writable(null);

/**
 * Current main view
 */
export const currentView = writable('feed'); // 'feed' | 'feeds-management' | 'categories-management'

/**
 * The current filtered/visible article list (written by Feed, read by modal for prev/next)
 */
export const visibleArticles = writable([]);

/**
 * Whether the full-screen reading modal is open
 */
export const articleModalOpen = writable(false);

/**
 * Default article open behaviour: 'sidebar' | 'modal'
 */
export const openMode = writable('sidebar');

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

export const clearErrors = () => {
	errors.set({ feeds: null, articles: null, categories: null, general: null });
};
