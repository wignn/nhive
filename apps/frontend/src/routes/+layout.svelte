<script lang="ts">
	import '../app.css';
	import { Toaster } from 'svelte-sonner';
	import { page } from '$app/stores';
	
	import SiteHeader from '$lib/components/layout/SiteHeader.svelte';
	import SiteFooter from '$lib/components/layout/SiteFooter.svelte';

	let { children } = $props();

	let isAdmin = $derived($page.url.pathname.startsWith('/admin'));
</script>

<div class="flex min-h-screen flex-col">
	<SiteHeader />
	<div class="flex-1">
		{@render children()}
	</div>
	{#if !isAdmin}
		<SiteFooter />
	{/if}
</div>

<Toaster richColors position="top-center" theme="dark" />
