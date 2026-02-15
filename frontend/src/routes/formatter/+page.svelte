<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as NativeSelect from '$lib/components/ui/native-select/index.js';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import { goto } from '$app/navigation';
	import { Kbd } from '$lib/components/ui/kbd/index.js';
	import Wand from '@lucide/svelte/icons/wand';


	let value = $state('');
	let format = $state('json');
	let language = $derived(format === 'json' ? 'json' : 'html');

	function onchange(newContent: string) {
		value = newContent;
	}

	function formatContent() {
		if (!value.trim()) return;

		try {
			if (format === 'json') {
				value = JSON.stringify(JSON.parse(value), null, 2);
			} else {
				value = formatHtml(value);
			}
		} catch {
			// invalid input, leave as-is
		}
	}

	function formatHtml(html: string): string {
		let formatted = '';
		let indent = 0;
		const tab = '  ';

		// Split on tags while keeping them
		const tokens = html.replace(/>\s*</g, '>\n<').split('\n');

		for (const token of tokens) {
			const trimmed = token.trim();
			if (!trimmed) continue;

			// Closing tag — decrease indent before printing
			if (/^<\//.test(trimmed)) {
				indent = Math.max(0, indent - 1);
			}

			formatted += tab.repeat(indent) + trimmed + '\n';

			// Opening tag (not self-closing, not closing) — increase indent
			if (/^<[^/!][^>]*[^/]>$/.test(trimmed) && !/^<(br|hr|img|input|meta|link)\b/i.test(trimmed)) {
				indent++;
			}
		}

		return formatted.trimEnd();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			goto('/');
		}
		if (event.metaKey && event.key === 'f') {
			event.preventDefault();
			formatContent();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<Card.Root class="h-full rounded-none">
	<Card.Header>
		<Card.Title>Formatter</Card.Title>
	</Card.Header>
	<Card.Content class="flex-1">
		<MonacoEditor {value} {language} {onchange} disableFind />
	</Card.Content>
	<Card.Footer class="flex gap-2">
		<Button onclick={formatContent} class="gap-2"><Wand /> <Kbd>Cmd+F</Kbd></Button>
		<NativeSelect.Root bind:value={format}>
			<NativeSelect.Option value="json">JSON</NativeSelect.Option>
			<NativeSelect.Option value="html">HTML</NativeSelect.Option>
		</NativeSelect.Root>
<Button href="/" variant="outline" class="gap-2">Close <Kbd>ESC</Kbd></Button>
	</Card.Footer>
</Card.Root>
