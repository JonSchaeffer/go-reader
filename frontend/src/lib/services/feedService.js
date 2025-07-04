import { rssApi } from '../api.js';
import { feeds, setLoading, setError } from '../stores.js';

/**
 * Service for managing RSS feeds with state management
 */
export class FeedService {
	/**
	 * Load all RSS feeds
	 */
	static async loadFeeds() {
		setLoading('feeds', true);
		setError('feeds', null);

		try {
			const response = await rssApi.getAll();
			console.log('Feeds API response:', response);
			
			// Handle different response structures from Go API
			let feedList = [];
			if (Array.isArray(response)) {
				feedList = response;
			} else if (response && Array.isArray(response.entries)) {
				feedList = response.entries;
			}
			
			feeds.set(feedList);
			return feedList;
		} catch (error) {
			console.error('Failed to load feeds:', error);
			setError('feeds', 'Failed to load RSS feeds. Please check if the backend is running.');
			throw error;
		} finally {
			setLoading('feeds', false);
		}
	}

	/**
	 * Add a new RSS feed
	 */
	static async addFeed(url) {
		if (!url || !url.trim()) {
			throw new Error('Please enter a valid RSS URL');
		}

		setLoading('adding', true);
		setError('feeds', null);

		try {
			const response = await rssApi.create(url.trim());
			console.log('Add feed response:', response);
			
			// Reload feeds to get the updated list
			await this.loadFeeds();
			
			return response;
		} catch (error) {
			console.error('Failed to add feed:', error);
			const errorMessage = error.message.includes('HTTP') 
				? 'Failed to add RSS feed. Please check the URL and try again.'
				: error.message;
			setError('feeds', errorMessage);
			throw error;
		} finally {
			setLoading('adding', false);
		}
	}

	/**
	 * Delete an RSS feed
	 */
	static async deleteFeed(feedId, feedTitle = 'this feed') {
		setLoading('deleting', true);
		setError('feeds', null);

		try {
			await rssApi.delete(feedId);
			console.log('Feed deleted successfully');
			
			// Remove from local state
			feeds.update(currentFeeds => 
				currentFeeds.filter(feed => feed.id !== feedId)
			);
			
			return true;
		} catch (error) {
			console.error('Failed to delete feed:', error);
			setError('feeds', `Failed to delete ${feedTitle}`);
			throw error;
		} finally {
			setLoading('deleting', false);
		}
	}

	/**
	 * Get feed by ID
	 */
	static async getFeedById(feedId) {
		try {
			const response = await rssApi.getById(feedId);
			console.log('Feed by ID response:', response);
			return response;
		} catch (error) {
			console.error('Failed to get feed by ID:', error);
			throw error;
		}
	}

	/**
	 * Get feed statistics
	 */
	static async getFeedStats(feedId) {
		try {
			const response = await rssApi.getStats(feedId);
			console.log('Feed stats response:', response);
			return response;
		} catch (error) {
			console.error('Failed to get feed stats:', error);
			throw error;
		}
	}
}