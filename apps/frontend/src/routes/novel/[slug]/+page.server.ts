import type { PageServerLoad, Actions } from './$types';
import { gatewayFetch } from '$lib/server/gateway';
import { error, fail } from '@sveltejs/kit';

const normalize = (id?: string) => (id || '').replace(/-/g, '').toLowerCase();

export const load: PageServerLoad = async ({ params, locals }) => {
	const slug = params.slug;

	const tasks: Promise<any>[] = [
		gatewayFetch(`/api/v1/novels/${slug}`),
		gatewayFetch(`/api/v1/novels/${slug}/chapters`)
	];

	if (locals.token) {
		tasks.push(gatewayFetch(`/api/v1/library`, { token: locals.token }).catch(() => null));
	}

	try {
		const [n, c, lib] = await Promise.all(tasks);

		const novelData = n.novel || n;
		let libEntry = null;

		if (lib && novelData?.id) {
			const list: any[] = Array.isArray(lib) ? lib : lib?.entries ?? [];
			libEntry = list.find((e) => normalize(e.novel_id) === normalize(novelData.id)) ?? null;
			if (libEntry && locals.token) {
				try {
					const p = await gatewayFetch(`/api/v1/progress/${novelData.id}`, { token: locals.token });
					libEntry.progress = p?.chapter_number || 0;
				} catch (e) {}
			}
		}

		return {
			novel: novelData,
			chapters: c.chapters || [],
			coverBase: n.cover_base_url,
			libEntry
		};
	} catch (e: any) {
		throw error(404, 'Novel not found');
	}
};

export const actions: Actions = {
	toggleLibrary: async ({ request, locals }) => {
		if (!locals.token) return fail(401, { error: 'Unauthorized' });

		const data = await request.formData();
		const novelId = data.get('novelId')?.toString();
		const action = data.get('action')?.toString(); // 'add' or 'remove'

		if (!novelId) return fail(400, { error: 'Missing novel ID' });

		try {
			if (action === 'remove') {
				await gatewayFetch(`/api/v1/library/${novelId}`, {
					method: 'DELETE',
					token: locals.token
				});
				return { success: true, action: 'removed' };
			} else {
				await gatewayFetch(`/api/v1/library/${novelId}`, {
					method: 'POST',
					token: locals.token
				});
				return { success: true, action: 'added' };
			}
		} catch (e: any) {
			return fail(400, { error: e.message || 'Library update failed' });
		}
	}
};
