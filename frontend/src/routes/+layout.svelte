<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { fileStore } from '$lib/stores/fileStore.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider style="--sidebar-width: 300px;">
	<AppSidebar />
	<Sidebar.Inset class="overflow-x-hidden">
		<header class="bg-background sticky top-0 flex shrink-0 items-center gap-2 border-b p-4 w-full">
			<Sidebar.Trigger class="-ml-1" />
			<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
			<div class="flex-1 min-w-0">
				{#if fileStore.currentFile}
					<h1 class="text-sm font-medium truncate" title={fileStore.currentFile.name}>
						{fileStore.currentFile.name}
					</h1>
				{:else}
					<h1 class="text-sm font-medium text-muted-foreground">
						No file selected
					</h1>
				{/if}
			</div>
		</header>
		<div class="flex flex-1 flex-col gap-4">
			{@render children()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>


