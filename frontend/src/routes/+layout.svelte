<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from 'svelte-sonner';
	import { themeStore } from '$lib/stores/themeStore.svelte';

	let { children } = $props();
</script>

<Toaster theme={themeStore.current} />

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex h-screen flex-col overflow-hidden">
	{@render children()}
</div>
