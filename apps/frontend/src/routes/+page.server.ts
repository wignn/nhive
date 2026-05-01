import type { PageServerLoad } from './$types';
import { gatewayFetch } from '$lib/server/gateway';

export const load: PageServerLoad = async () => {
	const [recentRes, topRes] = await Promise.allSettled([
		gatewayFetch(`/api/v1/novels?page=1&page_size=18&sort=updated`),
		gatewayFetch(`/api/v1/novels?page=1&page_size=10&sort=views`)
	]);

	const recent = recentRes.status === 'fulfilled' ? recentRes.value : { novels: [], cover_base_url: '' };
	const top = topRes.status === 'fulfilled' ? topRes.value : { novels: [], cover_base_url: '' };

	return {
		recent: recent.novels || [],
		trending: (top.novels || []).slice(0, 8),
		coverBase: recent.cover_base_url || top.cover_base_url
	};
};
