<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as NativeSelect from '$lib/components/ui/native-select/index.js';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import { goto } from '$app/navigation';
	import { Kbd } from '$lib/components/ui/kbd/index.js';

	let value = '';

	function onchange(newContent: string) {
		value = newContent;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			goto('/');
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<Card.Root class="h-full rounded-none">
	<Card.Header>
		<Card.Title>Formatter</Card.Title>
		<Card.Action>
			<NativeSelect.Root>
				<NativeSelect.Option value="json">JSON</NativeSelect.Option>
				<NativeSelect.Option value="html">HTML</NativeSelect.Option>
			</NativeSelect.Root>
		</Card.Action>
	</Card.Header>
	<Card.Content class="flex-1">
		<MonacoEditor {onchange} />
	</Card.Content>
	<Card.Footer class="flex gap-2">
		<Button href="/" variant="outline" class="gap-2">Close <Kbd>ESC</Kbd></Button>
	</Card.Footer>
</Card.Root>
