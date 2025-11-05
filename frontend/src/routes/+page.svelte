<script lang="ts">
	import * as Resizable from '$lib/components/ui/resizable/index.js';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import { fileStore } from '$lib/stores/fileStore.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { SaveFile } from '$lib/wailsjs/go/main/App';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	let editorContent = $derived(fileStore.content);

	// Determine language based on file extension
	let language = $derived.by(() => {
		if (!fileStore.currentFile) return 'plaintext';
		const ext = fileStore.currentFile.name.split('.').pop()?.toLowerCase();
		if (ext === 'md' || ext === 'markdown') return 'markdown';
		if (ext === 'hurl') return 'plaintext'; // We can add custom Hurl syntax later
		return 'plaintext';
	});

	// Determine editor theme based on current theme
	let editorTheme = $derived(themeStore.current === 'dark' ? 'vs-dark' : 'vs');

	async function handleContentChange(newContent: string) {
		if (fileStore.currentFile) {
			try {
				await SaveFile(fileStore.currentFile.path, newContent);
				fileStore.setContent(newContent);
			} catch (error) {
				console.error('Failed to save file:', error);
			}
		}
	}
</script>

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
		<Resizable.PaneGroup direction="horizontal" class="h-screen">
			<Resizable.Pane defaultSize={50}>
				{#if fileStore.currentFile}
					<MonacoEditor
						value={editorContent}
						{language}
						theme={editorTheme}
						onchange={handleContentChange}
					/>
				{:else}
					<div class="flex h-full items-center justify-center text-muted-foreground">
						<p>Select a file to edit</p>
					</div>
				{/if}
			</Resizable.Pane>
			<Resizable.Handle />
			<Resizable.Pane defaultSize={50}>
				<div class="p-4">
					<h2 class="mb-2 text-lg font-bold">Response</h2>
					<pre class="text-sm">Response will appear here...</pre>
				</div>
			</Resizable.Pane>
		</Resizable.PaneGroup>
	</Sidebar.Inset>
</Sidebar.Provider>
