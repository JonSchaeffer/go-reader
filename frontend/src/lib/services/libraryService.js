import { libraryApi } from '../api.js';
import { savedItems } from '../stores.js';

export class LibraryService {
	static async loadItems(offset = 0, limit = 50, archived = false) {
		const response = await libraryApi.getAll(offset, limit, archived);
		const list = response?.items ?? [];

		if (offset === 0) {
			savedItems.set(list);
		} else {
			savedItems.update((current) => [...current, ...list]);
		}

		return {
			items: list,
			total: response?.total ?? 0,
			hasMore: response?.hasMore ?? false
		};
	}

	static async saveArticle(articleId) {
		return libraryApi.save({ articleId });
	}

	static async saveURL(url) {
		return libraryApi.save({ url });
	}

	static async toggleRead(id, currentRead, currentArchived) {
		await libraryApi.update(id, !currentRead, currentArchived);
		savedItems.update((items) =>
			items.map((item) => (item.ID === id ? { ...item, Read: !currentRead } : item))
		);
	}

	static async archive(id, currentRead) {
		await libraryApi.update(id, currentRead, true);
		// Remove from active list
		savedItems.update((items) => items.filter((item) => item.ID !== id));
	}

	static async unarchive(id, currentRead) {
		await libraryApi.update(id, currentRead, false);
		savedItems.update((items) => items.filter((item) => item.ID !== id));
	}

	static async deleteItem(id) {
		await libraryApi.delete(id);
		savedItems.update((items) => items.filter((item) => item.ID !== id));
	}

	static async search(query) {
		const results = await libraryApi.search(query);
		const list = Array.isArray(results) ? results : [];
		savedItems.set(list);
		return list;
	}
}
