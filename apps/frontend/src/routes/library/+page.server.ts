import type { PageServerLoad } from './$types';
import { gatewayFetch } from '$lib/server/gateway';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.token) throw redirect(303, '/login');

	try {
		const d = await gatewayFetch(`/api/v1/library`, { token: locals.token });
		const raw: any[] = Array.isArray(d) ? d : d?.entries ?? [];
		const normalized = raw.map((e) => ({ ...e, status: e.status ?? "reading" }));

		let enriched = normalized;
		const needs = normalized.some((e) => !e.novel_title || !e.novel_slug);
		if (needs) {
			try {
				const r: any = await gatewayFetch(`/api/v1/novels?page=1&page_size=200&sort=updated`);
				const map = new Map<string, any>();
				(r.novels || []).forEach((n: any) => map.set(n.id, n));
				enriched = normalized.map((e) => {
					const n = map.get(e.novel_id);
					return n ? { ...e, novel_title: n.title, novel_slug: n.slug, total: n.total_chapters } : e;
				});
			} catch {}
		}

		try {
			const progressPromises = enriched.map(e => 
				gatewayFetch(`/api/v1/progress/${e.novel_id}`, { token: locals.token })
					.then(p => ({ novel_id: e.novel_id, progress: p?.chapter_number || 0 }))
					.catch(() => ({ novel_id: e.novel_id, progress: 0 }))
			);
			const progressResults = await Promise.all(progressPromises);
			const progressMap = new Map(progressResults.map(p => [p.novel_id, p.progress]));
			
			enriched = enriched.map(e => ({
				...e,
				progress: progressMap.get(e.novel_id) || 0
			}));
		} catch {}

		return { entries: enriched };
	} catch (e) {
		return { entries: [] };
	}
};
