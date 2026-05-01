import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { gatewayFetch } from '$lib/server/gateway';

export const GET: RequestHandler = async ({ url }) => {
	const q = url.searchParams.get('q');
	if (!q) return json({ hits: [] });
	try {
		const res = await gatewayFetch(`/api/v1/search?q=${encodeURIComponent(q)}`);
		return json(res);
	} catch (e) {
		return json({ hits: [] }, { status: 500 });
	}
};
