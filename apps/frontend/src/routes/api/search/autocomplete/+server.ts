import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { gatewayFetch } from '$lib/server/gateway';

export const GET: RequestHandler = async ({ url }) => {
	const q = url.searchParams.get('q');
	if (!q) return json({ suggestions: [] });
	try {
		const res = await gatewayFetch(`/api/v1/search/autocomplete?q=${encodeURIComponent(q)}`);
		return json(res);
	} catch (e) {
		return json({ suggestions: [] }, { status: 500 });
	}
};
