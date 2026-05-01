import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/private';

export const POST: RequestHandler = async ({ request, locals }) => {
	if (!locals.token || locals.user?.role !== 'admin') {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	try {
		const fd = await request.formData();
		const image = fd.get('image') as File;
		if (!image) return json({ error: 'No image provided' }, { status: 400 });

		const gatewayUrl = env.NOVELHIVE_GATEWAY_URL || 'http://localhost:8080';
		const internalKey = env.NOVELHIVE_INTERNAL_KEY || 'dev-internal-key';

		// Forward form data to gateway
		const body = new FormData();
		body.append('image', image);

		const res = await fetch(`${gatewayUrl}/api/v1/admin/upload`, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${locals.token}`,
				'X-Internal-Key': internalKey
			},
			body: body
		});

		if (!res.ok) {
			const text = await res.text();
			throw new Error(text || 'Upload failed');
		}

		const data = await res.json();
		return json(data);
	} catch (e: any) {
		return json({ error: e.message || 'Upload failed' }, { status: 500 });
	}
};
