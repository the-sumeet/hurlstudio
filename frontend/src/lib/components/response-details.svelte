<script lang="ts">
	// Using global App.Response type from app.d.ts
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { GetResponseBody } from '$lib/wailsjs/go/main/App';
	import { fileStore } from '$lib/stores/fileStore.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import { Button } from './ui/button/index.js';
	import Wand from '@lucide/svelte/icons/wand';
	import Copy from '@lucide/svelte/icons/copy';
	import { getHttpStatusColorClass } from '$lib/utils/httpStatus';

	let { response }: { response: App.Response } = $props();

	let actualBody = $state<string>('');
	let isLoadingBody = $state(false);
	let formattedBody = $state<string>('');
	let isFormatted = $state(false);

	// Load actual body content if response.body is a file path
	$effect(() => {
		if (response.body && fileStore.currentFile) {
			loadBodyContent();
		}
	});

	// Reset formatted state when body changes
	$effect(() => {
		if (actualBody) {
			isFormatted = false;
			formattedBody = actualBody;
		}
	});

	async function loadBodyContent() {
		if (!response.body || !fileStore.currentFile) return;

		// Check if body looks like a file path (e.g., "store/response_1.txt")
		if (response.body.includes('/') || response.body.includes('\\')) {
			isLoadingBody = true;
			try {
				const bodyContent = await GetResponseBody(fileStore.currentFile.path, response.body);
				actualBody = bodyContent;
			} catch (error) {
				console.error('Failed to load response body:', error);
				actualBody = response.body; // Fallback to showing the path
			} finally {
				isLoadingBody = false;
			}
		} else {
			// Body is already the actual content
			actualBody = response.body;
		}
	}

	// Determine status color
	let statusColor = $derived(getHttpStatusColorClass(response.status));

	// Determine editor theme based on current theme
	let editorTheme = $derived(themeStore.current === 'dark' ? 'vs-dark' : 'vs');

	// Determine language for body based on Content-Type header
	let bodyLanguage = $derived.by(() => {
		const contentType = response.headers
			.find((h) => h.name.toLowerCase() === 'content-type')
			?.value.toLowerCase();

		if (!contentType) return 'plaintext';

		if (contentType.includes('json')) return 'json';
		if (contentType.includes('html')) return 'html';
		if (contentType.includes('xml')) return 'xml';
		if (contentType.includes('javascript') || contentType.includes('ecmascript'))
			return 'javascript';
		if (contentType.includes('css')) return 'css';
		if (contentType.includes('yaml')) return 'yaml';

		return 'plaintext';
	});

	// Format details as text for textarea
	let detailsText = $derived.by(() => {
		let text = `Status: ${response.status}\n`;
		text += `Version: ${response.http_version}\n`;

		if (response.headers.length > 0) {
			text += `\nHeaders:\n`;
			response.headers.forEach((header) => {
				text += `  ${header.name}: ${header.value}\n`;
			});
		}

		if (response.cookies.length > 0) {
			text += `\nCookies:\n`;
			response.cookies.forEach((cookie) => {
				text += `  ${cookie.name}: ${cookie.value}\n`;
			});
		}

		if (response.certificate) {
			text += `\nCertificate:\n`;
			text += `  Subject: ${response.certificate.subject}\n`;
			text += `  Issuer: ${response.certificate.issuer}\n`;
			text += `  Valid: ${response.certificate.start_date} - ${response.certificate.expire_date}\n`;
		}

		return text;
	});

	// Copy body content to clipboard
	async function handleCopy() {
		try {
			await navigator.clipboard.writeText(formattedBody);
		} catch (error) {
			console.error('Failed to copy to clipboard:', error);
		}
	}

	// Format body content based on content type
	function handleFormat() {
		if (!actualBody) return;

		try {
			if (bodyLanguage === 'json') {
				// Format JSON
				const parsed = JSON.parse(actualBody);
				formattedBody = JSON.stringify(parsed, null, 2);
				isFormatted = true;
			} else if (bodyLanguage === 'xml' || bodyLanguage === 'html') {
				// Basic XML/HTML formatting
				formattedBody = formatXml(actualBody);
				isFormatted = true;
			} else {
				// For other types, just use original
				formattedBody = actualBody;
			}
		} catch (error) {
			console.error('Failed to format content:', error);
			formattedBody = actualBody;
		}
	}

	// Simple XML/HTML formatter
	function formatXml(xml: string): string {
		let formatted = '';
		let indent = 0;
		const tab = '  ';

		xml.split(/>\s*</).forEach((node) => {
			if (node.match(/^\/\w/)) indent--; // Closing tag
			formatted += tab.repeat(indent) + '<' + node + '>\n';
			if (node.match(/^<?\w[^>]*[^\/]$/)) indent++; // Opening tag
		});

		return formatted.substring(1, formatted.length - 2);
	}
</script>

<Tabs.Root value="body" class="flex h-full min-h-0 w-full flex-col">
	<Tabs.List class="shrink-0">
		<Tabs.Trigger value="body">Body</Tabs.Trigger>
		<Tabs.Trigger value="details">Details</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="body" class="flex min-h-0 flex-1 flex-col overflow-hidden">
		{#if isLoadingBody}
			<p class="p-2 text-sm text-muted-foreground">Loading body...</p>
		{:else if actualBody}
			<div class="mb-2 flex w-full shrink-0 justify-end">
				<ButtonGroup.Root>
					<Button size="icon" variant="ghost" onclick={handleFormat} title="Format content">
						<Wand />
					</Button>
					<Button size="icon" variant="ghost" onclick={handleCopy} title="Copy to clipboard">
						<Copy />
					</Button>
				</ButtonGroup.Root>
			</div>
			<div class="flex-1 overflow-hidden">
				<MonacoEditor
					value={formattedBody}
					language={bodyLanguage}
					theme={editorTheme}
					readonly={true}
				/>
			</div>
		{:else}
			<p class="p-2 text-sm text-muted-foreground">No body content</p>
		{/if}
	</Tabs.Content>
	<Tabs.Content value="details" class="flex min-h-0 flex-1 flex-col overflow-hidden">
		<textarea
			readonly
			class="w-full flex-1 resize-none overflow-auto rounded border bg-muted p-2 font-mono text-xs"
			style="white-space: pre;"
			value={detailsText}
		></textarea>
	</Tabs.Content>
</Tabs.Root>
