<script lang="ts">
	import * as Resizable from "$lib/components/ui/resizable/index.js";
	import MonacoEditor from "$lib/components/MonacoEditor.svelte";

	let hurlContent = `# Example Hurl request
GET https://api.github.com/users/github

HTTP 200
[Asserts]
header "Content-Type" contains "application/json"`;
</script>

<Resizable.PaneGroup direction="horizontal" class="h-screen">
	<Resizable.Pane defaultSize={50}>
		<MonacoEditor bind:value={hurlContent} language="plaintext" theme="vs-dark" />
	</Resizable.Pane>
	<Resizable.Handle />
	<Resizable.Pane defaultSize={50}>
		<div class="p-4">
			<h2 class="text-lg font-bold mb-2">Response</h2>
			<pre class="text-sm">Response will appear here...</pre>
		</div>
	</Resizable.Pane>
</Resizable.PaneGroup>
