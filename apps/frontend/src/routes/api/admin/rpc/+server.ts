import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { gatewayFetch } from '$lib/server/gateway';

export const POST: RequestHandler = async ({ request, locals }) => {
	if (!locals.token || locals.user?.role !== 'admin') {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	const { action, payload } = await request.json();

	try {
		switch (action) {
			case 'listNovels':
				return json(await gatewayFetch('/api/v1/admin/novels', { token: locals.token }));
			case 'listChapters':
				return json(await gatewayFetch('/api/v1/admin/chapters', { token: locals.token }));
			case 'listUsers':
				return json(await gatewayFetch('/api/v1/admin/users', { token: locals.token }));
			case 'createNovel':
				return json(await gatewayFetch('/api/v1/admin/novels', { method: 'POST', body: JSON.stringify(payload), token: locals.token }));
			case 'updateNovel':
				return json(await gatewayFetch(`/api/v1/admin/novels/${payload.id}`, { method: 'PUT', body: JSON.stringify(payload.data), token: locals.token }));
			case 'deleteNovel':
				return json(await gatewayFetch(`/api/v1/admin/novels/${payload.id}`, { method: 'DELETE', token: locals.token }));
			case 'createChapter':
				return json(await gatewayFetch('/api/v1/admin/chapters', { method: 'POST', body: JSON.stringify(payload), token: locals.token }));
			case 'updateChapter':
				return json(await gatewayFetch(`/api/v1/admin/chapters/${payload.id}`, { method: 'PUT', body: JSON.stringify({ title: payload.title, content: payload.content }), token: locals.token }));
			case 'deleteChapter':
				return json(await gatewayFetch(`/api/v1/admin/chapters/${payload.id}`, { method: 'DELETE', token: locals.token }));
			case 'setUserRole':
				return json(await gatewayFetch(`/api/v1/admin/users/${payload.id}/role`, { method: 'PUT', body: JSON.stringify({ role: payload.role }), token: locals.token }));
			case 'listGenres':
				return json(await gatewayFetch('/api/v1/admin/genres', { token: locals.token }));
			case 'createGenre':
				return json(await gatewayFetch('/api/v1/admin/genres', { method: 'POST', body: JSON.stringify(payload), token: locals.token }));
			case 'deleteGenre':
				return json(await gatewayFetch(`/api/v1/admin/genres/${payload.id}`, { method: 'DELETE', token: locals.token }));
			default:
				return json({ error: 'Unknown action' }, { status: 400 });
		}
	} catch (e: any) {
		return json({ error: e.message || 'Action failed' }, { status: 500 });
	}
};
