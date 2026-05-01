import type { PageServerLoad } from './$types';
import { gatewayFetch } from '$lib/server/gateway';

export const load: PageServerLoad = async ({ url }) => {
	const genre = url.searchParams.get('genre') || '';
	const sort = url.searchParams.get('sort') || 'updated';

	const params = new URLSearchParams();
	if (genre) params.set('genre', genre);
	params.set('sort', sort);
	params.set('page_size', '30');

	try {
		const res = await gatewayFetch(`/api/v1/novels?${params.toString().toLowerCase()}`);
		return {
			novels: res.novels || [],
			total: res.total || 0,
			coverBase: res.cover_base_url,
			genre,
			sort
		};
	} catch (e) {
		return {
			novels: [],
			total: 0,
			coverBase: '',
			genre,
			sort
		};
	}
};
