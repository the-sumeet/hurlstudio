<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { fileStore } from '$lib/stores/fileStore.svelte';
	import { page } from '$app/stores';
	import * as Card from '$lib/components/ui/card/index.js';
	let { children } = $props();

	let showSidebar = $derived($page.url.pathname === '/');
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#if showSidebar}
	<Sidebar.Provider style="--sidebar-width: 300px;">
		<AppSidebar />
		<Sidebar.Inset class="overflow-x-hidden">
			<header
				class="sticky top-0 flex w-full shrink-0 items-center gap-2 border-b bg-background p-4"
			>
				<Sidebar.Trigger class="-ml-1" />
				<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
				<div class="min-w-0 flex-1">
					{#if fileStore.currentFile}
						<h1 class="truncate text-sm font-medium" title={fileStore.currentFile.name}>
							{fileStore.currentFile.name}
						</h1>
					{:else}
						<h1 class="text-sm font-medium text-muted-foreground">No file selected</h1>
					{/if}
				</div>
			</header>
			<div class="flex flex-1 flex-col gap-4">
				{@render children()}
			</div>
		</Sidebar.Inset>
	</Sidebar.Provider>
{:else}
	<div class="h-screen w-full">
		<!-- <header class="sticky top-0 flex w-full shrink-0 items-center gap-2 border-b bg-background p-4">
			<div class="min-w-0 flex-1">
				<h1 class="text-sm font-medium">Hurl Studio</h1>
			</div>
		</header> -->

		<div class="flex h-full flex-1 flex-col gap-4 overflow-hidden">
			{@render children()}
		</div>
	</div>
{/if}
