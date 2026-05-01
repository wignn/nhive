import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const sessionStr = event.cookies.get('novelhive_session');
	
	if (sessionStr) {
		try {
			const sessionData = JSON.parse(Buffer.from(sessionStr, 'base64').toString('utf-8'));
			if (sessionData.user && sessionData.token) {
				event.locals.user = sessionData.user;
				event.locals.token = sessionData.token;
			}
		} catch (e) {
			// invalid session
			event.cookies.delete('novelhive_session', { path: '/' });
		}
	}

	const response = await resolve(event);
	return response;
};
