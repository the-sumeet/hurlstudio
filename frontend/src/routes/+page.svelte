<script lang="ts">
	import * as Resizable from "$lib/components/ui/resizable/index.js";
	import MonacoEditor from "$lib/components/MonacoEditor.svelte";
	import { fileStore } from '$lib/stores/fileStore.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { SaveFile } from "$lib/wailsjs/go/main/App";

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

<Resizable.PaneGroup direction="horizontal" class="h-screen">
	<Resizable.Pane defaultSize={50}>
		{#if fileStore.currentFile}
			<div class="h-full flex flex-col">
				<div class="px-4 py-2 border-b">
					<h2 class="text-sm font-medium">{fileStore.currentFile.name}</h2>
				</div>
				<div class="flex-1">
					<MonacoEditor
						value={editorContent}
						{language}
						theme={editorTheme}
						onchange={handleContentChange}
					/>
				</div>
			</div>
		{:else}
			<div class="h-full flex items-center justify-center text-muted-foreground">
				<p>Select a file to edit</p>
			</div>
		{/if}
	</Resizable.Pane>
	<Resizable.Handle />
	<Resizable.Pane defaultSize={50}>
		<div class="p-4">
			<h2 class="text-lg font-bold mb-2">Response</h2>
			<pre class="text-sm">Response will appear here...</pre>
		</div>
	</Resizable.Pane>
</Resizable.PaneGroup>
