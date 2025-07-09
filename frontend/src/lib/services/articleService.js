import { articleApi } from '../api.js';
import { articles, setLoading, setError } from '../stores.js';

/**
 * Service for managing articles with state management
 */
export class ArticleService {
	/**
	 * Load all articles
	 */
	static async loadAllArticles() {
		setLoading('articles', true);
		setError('articles', null);

		try {
			const response = await articleApi.getAll();
			console.log('Articles API response:', response);
			
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to load articles:', error);
			setError('articles', 'Failed to load articles. Please check if the backend is running.');
			throw error;
		} finally {
			setLoading('articles', false);
		}
	}

	/**
	 * Load articles for a specific RSS feed
	 */
	static async loadArticlesByFeed(feedId, limit = 100) {
		setLoading('articles', true);
		setError('articles', null);

		try {
			const response = await articleApi.getByRssId(feedId, limit);
			console.log('Articles by feed response:', response);
			
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to load articles by feed:', error);
			setError('articles', 'Failed to load articles for this feed.');
			throw error;
		} finally {
			setLoading('articles', false);
		}
	}

	/**
	 * Search articles
	 */
	static async searchArticles(query, limit = 20) {
		if (!query || !query.trim()) {
			return this.loadAllArticles();
		}

		setLoading('articles', true);
		setError('articles', null);

		try {
			const response = await articleApi.search(query.trim(), limit);
			console.log('Search articles response:', response);
			
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to search articles:', error);
			setError('articles', 'Search failed. Please try again.');
			throw error;
		} finally {
			setLoading('articles', false);
		}
	}

	/**
	 * Get single article by ID
	 */
	static async getArticleById(articleId) {
		try {
			const response = await articleApi.getById(articleId);
			console.log('Article by ID response:', response);
			
			// Handle response that might be an array with one item
			if (Array.isArray(response) && response.length > 0) {
				return response[0];
			} else if (response && !Array.isArray(response)) {
				return response;
			} else {
				throw new Error('Article not found');
			}
		} catch (error) {
			console.error('Failed to get article by ID:', error);
			throw error;
		}
	}

	/**
	 * Toggle article read status
	 */
	static async toggleReadStatus(articleId, currentReadStatus) {
		try {
			const newReadStatus = !currentReadStatus;
			await articleApi.updateReadStatus(articleId, newReadStatus);
			
			// Update local state
			articles.update(currentArticles => 
				currentArticles.map(article => 
					article.ID === articleId 
						? { ...article, Read: newReadStatus }
						: article
				)
			);
			
			return newReadStatus;
		} catch (error) {
			console.error('Failed to update read status:', error);
			setError('articles', 'Failed to update article status');
			throw error;
		}
	}

	/**
	 * Delete an article
	 */
	static async deleteArticle(articleId, articleTitle = 'this article') {
		try {
			await articleApi.delete(articleId);
			console.log('Article deleted successfully');
			
			// Remove from local state
			articles.update(currentArticles => 
				currentArticles.filter(article => article.ID !== articleId)
			);
			
			return true;
		} catch (error) {
			console.error('Failed to delete article:', error);
			setError('articles', `Failed to delete ${articleTitle}`);
			throw error;
		}
	}

	/**
	 * Format article date for display
	 */
	static formatDate(dateString) {
		if (!dateString) return 'No date';
		
		try {
			return new Date(dateString).toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'short',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch {
			return dateString;
		}
	}

	/**
	 * Truncate article description
	 */
	static truncateDescription(description, maxLength = 200) {
		if (!description) return '';
		return description.length > maxLength 
			? description.substring(0, maxLength) + '...'
			: description;
	}

	/**
	 * Refresh articles - reload current view
	 */
	static async refreshArticles() {
		setLoading('articles', true);
		setError('articles', null);

		try {
			const response = await articleApi.getAll();
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to refresh articles:', error);
			setError('articles', 'Failed to refresh articles');
			throw error;
		} finally {
			setLoading('articles', false);
		}
	}

	/**
	 * Refresh articles silently - reload without showing loading state
	 */
	static async refreshArticlesSilently() {
		try {
			const response = await articleApi.getAll();
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to refresh articles silently:', error);
			// Don't set error state for silent refresh to avoid disrupting user experience
			throw error;
		}
	}

	/**
	 * Load articles by feed silently - reload without showing loading state
	 */
	static async loadArticlesByFeedSilently(feedId, limit = 100) {
		try {
			const response = await articleApi.getByRssId(feedId, limit);
			const articleList = Array.isArray(response) ? response : [];
			articles.set(articleList);
			return articleList;
		} catch (error) {
			console.error('Failed to load articles by feed silently:', error);
			// Don't set error state for silent refresh to avoid disrupting user experience
			throw error;
		}
	}
}