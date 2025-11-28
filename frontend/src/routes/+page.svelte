<script lang="ts">
	import * as Resizable from '$lib/components/ui/resizable/index.js';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import Play from '@lucide/svelte/icons/play';
	import Copy from '@lucide/svelte/icons/copy';
	import { fileStore } from '$lib/stores/fileStore.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { SaveFile, RunHurl, GetExistingReport, RunHurlEntry } from '$lib/wailsjs/go/main/App';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as NativeSelect from '$lib/components/ui/native-select/index.js';
	import Call from '$lib/components/Call.svelte';
	import { marked } from 'marked';
	import { Kbd } from '$lib/components/ui/kbd/index.js';

	let editorContent = $derived(fileStore.content);
	let output = $state('');
	let isRunning = $state(false);
	let report = $state<App.HurlReport | null>(null);
	let selectedEntryIndex = $state(0);

	// Load existing report when file changes
	$effect(() => {
		if (fileStore.currentFile && fileStore.currentFile.name.endsWith('.hurl')) {
			loadExistingReport();
		} else {
			// Clear report if not a hurl file
			report = null;
			output = '';
		}
	});

	async function loadExistingReport() {
		if (!fileStore.currentFile) return;

		try {
			const existingReport = await GetExistingReport(fileStore.currentFile.path);
			if (existingReport) {
				// Try to parse as JSON report
				try {
					const parsed = JSON.parse(existingReport);
					// Hurl returns an array with a single report
					if (Array.isArray(parsed) && parsed.length > 0) {
						report = parsed[0];
						output = existingReport;
					}
				} catch (e) {
					// Not a valid JSON, ignore
					console.log('Existing report is not valid JSON');
				}
			} else {
				// No existing report found
				report = null;
				output = '';
			}
		} catch (error) {
			console.error('Failed to load existing report:', error);
		}
	}

	// Get selected entry
	let selectedEntry = $derived.by(() => {
		if (report && report.entries.length > 0) {
			return report.entries[selectedEntryIndex];
		}
		return null;
	});

	// Determine language based on file extension
	let language = $derived.by(() => {
		if (!fileStore.currentFile) return 'plaintext';
		const ext = fileStore.currentFile.name.split('.').pop()?.toLowerCase();
		if (ext === 'md' || ext === 'markdown') return 'markdown';
		if (ext === 'hurl') return 'plaintext'; // We can add custom Hurl syntax later
		return 'plaintext';
	});

	// Check if current file is markdown
	let isMarkdownFile = $derived.by(() => {
		if (!fileStore.currentFile) return false;
		const ext = fileStore.currentFile.name.split('.').pop()?.toLowerCase();
		return ext === 'md' || ext === 'markdown';
	});

	// Check if current file is a Hurl file
	let isHurlFile = $derived.by(() => {
		if (!fileStore.currentFile) return false;
		const ext = fileStore.currentFile.name.split('.').pop()?.toLowerCase();
		return ext === 'hurl';
	});

	// Determine editor theme based on current theme
	let editorTheme = $derived(themeStore.current === 'dark' ? 'vs-dark' : 'vs');

	// Render markdown content with resolved paths
	let renderedMarkdown = $derived.by(() => {
		if (isMarkdownFile && editorContent && fileStore.currentFile) {
			// Configure marked to use a custom renderer for images
			const renderer = new marked.Renderer();
			const baseDir = fileStore.currentFile.path.substring(
				0,
				fileStore.currentFile.path.lastIndexOf('/')
			);

			renderer.image = ({ href, title, text }) => {
				// Resolve relative paths to absolute paths
				let resolvedHref = href || '';
				if (
					typeof href === 'string' &&
					href &&
					!href.startsWith('http://') &&
					!href.startsWith('https://') &&
					!href.startsWith('data:')
				) {
					// Handle relative paths
					if (href.startsWith('./') || href.startsWith('../') || !href.startsWith('/')) {
						// Resolve relative to the markdown file's directory
						const parts = href.split('/');
						const dirParts = baseDir.split('/');

						for (const part of parts) {
							if (part === '..') {
								dirParts.pop();
							} else if (part !== '.' && part !== '') {
								dirParts.push(part);
							}
						}

						resolvedHref = dirParts.join('/');
					}
				}

				const titleAttr = title ? ` title="${title}"` : '';
				return `<img src="${resolvedHref}" alt="${text || ''}"${titleAttr}>`;
			};

			return marked(editorContent, { renderer });
		}
		return '';
	});

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

	async function handleRun() {
		if (!fileStore.currentFile) return;

		isRunning = true;
		output = 'Running...';
		report = null;
		selectedEntryIndex = 0;

		try {
			const result = await RunHurl(fileStore.currentFile.path);
			output = result || 'Success! (No output)';

			// Try to parse as JSON report
			try {
				const parsed = JSON.parse(result);
				// Hurl returns an array with a single report
				if (Array.isArray(parsed) && parsed.length > 0) {
					report = parsed[0];
				}
			} catch (e) {
				// Not a JSON response, keep as plain text
				console.log('Response is not JSON, displaying as plain text');
			}
		} catch (error) {
			output = `Error: ${error}`;
			console.error('Failed to run hurl:', error);
		} finally {
			isRunning = false;
		}
	}

	async function handleRunEntry(entryIndex: number) {
		if (!fileStore.currentFile) return;

		isRunning = true;
		output = `Running entry ${entryIndex}...`;
		report = null;
		selectedEntryIndex = 0;

		try {
			const result = await RunHurlEntry(fileStore.currentFile.path, entryIndex);
			output = result || 'Success! (No output)';

			// Try to parse as JSON report
			try {
				const parsed = JSON.parse(result);
				// Hurl returns an array with a single report
				if (Array.isArray(parsed) && parsed.length > 0) {
					report = parsed[0];
					// Set the selected entry to the one we just ran (entryIndex - 1 for 0-based)
					selectedEntryIndex = 0;
				}
			} catch (e) {
				// Not a JSON response, keep as plain text
				console.log('Response is not JSON, displaying as plain text');
			}
		} catch (error) {
			output = `Error: ${error}`;
			console.error('Failed to run hurl entry:', error);
		} finally {
			isRunning = false;
		}
	}

	// Copy curl command to clipboard
	async function copyCurlCommand() {
		if (!selectedEntry?.curl_cmd) return;

		try {
			await navigator.clipboard.writeText(selectedEntry.curl_cmd);
		} catch (error) {
			console.error('Failed to copy curl command:', error);
		}
	}

	// Handle keyboard shortcuts
	function handleKeydown(event: KeyboardEvent) {
		// Cmd+R (Mac) or Ctrl+R (Windows/Linux) to run
		if ((event.metaKey || event.ctrlKey) && event.key === 'r') {
			event.preventDefault(); // Prevent browser refresh
			if (fileStore.currentFile && !isMarkdownFile && !isRunning) {
				handleRun();
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<Sidebar.Provider style="--sidebar-width: 300px;">
	<AppSidebar />
	<Sidebar.Inset class="overflow-x-hidden">
		<header class="sticky top-0 flex w-full shrink-0 items-center gap-2 border-b bg-background p-4">
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
			{#if fileStore.currentFile && !isMarkdownFile}
				<div class="flex gap-2">
					<Button onclick={handleRun} disabled={isRunning} class="gap-2">
						<Play />
						{isRunning ? 'Running...' : 'Run'}
						{#if !isRunning}
							<Kbd>⌘R</Kbd>
						{/if}
					</Button>
				</div>
			{/if}
		</header>
		<Resizable.PaneGroup direction="horizontal" class="h-screen">
			<!-- Request -->
			<Resizable.Pane defaultSize={50}>
				{#if fileStore.currentFile}
					<MonacoEditor
						value={editorContent}
						{language}
						theme={editorTheme}
						onchange={handleContentChange}
						onRunEntry={isHurlFile ? handleRunEntry : undefined}
					/>
				{:else}
					<div class="flex h-full items-center justify-center text-muted-foreground">
						<p>Select a file to edit</p>
					</div>
				{/if}
			</Resizable.Pane>

			<Resizable.Handle withHandle />

			<!-- Response / Preview -->
			<Resizable.Pane defaultSize={50}>
				{#if isMarkdownFile}
					<!-- Markdown Preview -->
					<div class="h-full overflow-auto p-4">
						<div class="prose prose-sm max-w-none dark:prose-invert">
							{@html renderedMarkdown}
						</div>
					</div>
				{:else}
					<!-- Hurl Response -->
					<div class="flex h-full flex-col gap-4 p-2">
						<!-- Entries and Curl Copy -->
						{#if report && report.entries.length > 0}
							<div class="flex flex-shrink-0 gap-2">
								{#if report.entries.length > 1}
									<div class="flex-1">
										<NativeSelect.Root bind:value={selectedEntryIndex}>
											{#each report.entries as entry, index}
												<NativeSelect.Option value={index}>
													{entry.calls[0]?.request.method}
													{entry.calls[0]?.request.url}
												</NativeSelect.Option>
											{/each}
										</NativeSelect.Root>
									</div>
								{/if}
								{#if selectedEntry?.curl_cmd}
									<Button
										size="sm"
										variant="outline"
										onclick={copyCurlCommand}
										title="Copy curl command"
									>
										<Copy class="h-4 w-4" />
										Curl
									</Button>
								{/if}
							</div>
						{/if}

						<div class="flex min-h-0 flex-1 snap-y snap-mandatory flex-col gap-2 overflow-y-auto">
							{#if selectedEntry && selectedEntry.calls.length > 0}
								{#each selectedEntry.calls as call}
									<div class="h-full snap-start">
										<Call {call} />
									</div>
								{/each}
							{:else if output}
								<pre class="text-sm whitespace-pre-wrap">{output}</pre>
							{:else}
								<p class="text-sm text-muted-foreground">Response will appear here...</p>
							{/if}
						</div>
					</div>
				{/if}
			</Resizable.Pane>
		</Resizable.PaneGroup>
	</Sidebar.Inset>
</Sidebar.Provider>
