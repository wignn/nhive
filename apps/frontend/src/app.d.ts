// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				id: string;
				username: string;
				email: string;
				role: string;
			};
			token?: string;
		}
		interface PageData {
			user?: {
				id: string;
				username: string;
				email: string;
				role: string;
			};
		}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
