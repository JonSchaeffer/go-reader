import { tagApi } from '$lib/api';
import { allTags } from '$lib/stores';

export const TagService = {
	async loadAll() {
		const tags = await tagApi.getAll();
		allTags.set(tags ?? []);
		return tags ?? [];
	},

	async getForItem(itemType, itemId) {
		return (await tagApi.getForItem(itemType, itemId)) ?? [];
	},

	/** Create tag (upsert) and return it */
	async create(name) {
		const tag = await tagApi.create(name.trim());
		allTags.update((tags) => {
			if (!tags.find((t) => t.ID === tag.ID)) return [...tags, tag];
			return tags;
		});
		return tag;
	},

	async addToItem(itemType, itemId, tagId) {
		await tagApi.addToItem(itemType, itemId, tagId);
	},

	async removeFromItem(itemType, itemId, tagId) {
		await tagApi.removeFromItem(itemType, itemId, tagId);
	},

	async deleteTag(id) {
		await tagApi.delete(id);
		allTags.update((tags) => tags.filter((t) => t.ID !== id));
	}
};
