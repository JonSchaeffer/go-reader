/**
 * API client for the Go RSS Reader backend
 */

// Use environment variable or fallback to localhost
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api';

/**
 * Generic API request function with error handling
 */
async function apiRequest(endpoint, options = {}) {
	const url = `${API_BASE}${endpoint}`;
	
	const config = {
		headers: {
			'Content-Type': 'application/json',
			...options.headers
		},
		...options
	};

	try {
		const response = await fetch(url, config);
		
		if (!response.ok) {
			throw new Error(`HTTP ${response.status}: ${response.statusText}`);
		}

		// Handle empty responses
		const text = await response.text();
		if (!text) return null;
		
		try {
			return JSON.parse(text);
		} catch {
			return text;
		}
	} catch (error) {
		console.error(`API Error (${endpoint}):`, error);
		throw error;
	}
}

/**
 * RSS Feed API functions
 */
export const rssApi = {
	/**
	 * Get all RSS feeds
	 */
	async getAll() {
		return apiRequest('/rss');
	},

	/**
	 * Get specific RSS feed by ID
	 */
	async getById(id) {
		return apiRequest(`/rss?id=${id}`);
	},

	/**
	 * Create new RSS feed
	 */
	async create(url) {
		return apiRequest('/rss', {
			method: 'POST',
			body: JSON.stringify({ url })
		});
	},

	/**
	 * Update RSS feed
	 */
	async update(id, data) {
		const params = new URLSearchParams({ id: id.toString(), ...data });
		return apiRequest(`/rss?${params}`, {
			method: 'PUT'
		});
	},

	/**
	 * Delete RSS feed
	 */
	async delete(id) {
		return apiRequest(`/rss?id=${id}`, {
			method: 'DELETE'
		});
	},

	/**
	 * Get RSS feed statistics
	 */
	async getStats(id) {
		return apiRequest(`/rss/stats?id=${id}`);
	}
};

/**
 * Article API functions
 */
export const articleApi = {
	/**
	 * Get all articles
	 */
	async getAll() {
		return apiRequest('/articles');
	},

	/**
	 * Get specific article by ID
	 */
	async getById(id) {
		return apiRequest(`/articles/single?id=${id}`);
	},

	/**
	 * Get articles by RSS feed ID
	 */
	async getByRssId(rssId, limit = 100) {
		return apiRequest(`/articles/by-rss?rssid=${rssId}&limit=${limit}`);
	},

	/**
	 * Update article read status
	 */
	async updateReadStatus(id, read) {
		return apiRequest(`/articles/update?id=${id}&read=${read}`, {
			method: 'PUT'
		});
	},

	/**
	 * Search articles
	 */
	async search(query, limit = 20) {
		return apiRequest(`/articles/search?query=${encodeURIComponent(query)}&limit=${limit}`);
	},

	/**
	 * Delete article
	 */
	async delete(id) {
		return apiRequest(`/articles/delete?id=${id}`, {
			method: 'DELETE'
		});
	}
};

/**
 * API client instance for direct use
 */
export const api = {
	get: (endpoint) => apiRequest(endpoint),
	post: (endpoint, data) => apiRequest(endpoint, { method: 'POST', body: JSON.stringify(data) }),
	put: (endpoint, data) => apiRequest(endpoint, { method: 'PUT', body: JSON.stringify(data) }),
	delete: (endpoint) => apiRequest(endpoint, { method: 'DELETE' })
};

export default api;
