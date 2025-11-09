<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
</script>

<Card.Root class="h-full rounded-none">
	<Card.Header>
		<Card.Title>Settings</Card.Title>
	</Card.Header>
	<Card.Content class="flex-1"></Card.Content>
	<Card.Footer class="flex gap-2">
		<Button href="/" variant="outline">Close</Button>
	</Card.Footer>
</Card.Root>
