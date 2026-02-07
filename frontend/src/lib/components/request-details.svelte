<script lang="ts">
	// Using global App.Request type from app.d.ts

	let { request }: { request: App.Request } = $props();

	// Format request details as text for textarea
	let requestText = $derived.by(() => {
		let text = `Method: ${request.method}\n`;
		text += `URL: ${request.url}\n`;

		if (request.headers.length > 0) {
			text += `\nHeaders:\n`;
			request.headers.forEach((header) => {
				text += `  ${header.name}: ${header.value}\n`;
			});
		}

		if (request.query_string.length > 0) {
			text += `\nQuery Parameters:\n`;
			request.query_string.forEach((param) => {
				text += `  ${param.name}: ${param.value}\n`;
			});
		}

		if (request.cookies.length > 0) {
			text += `\nCookies:\n`;
			request.cookies.forEach((cookie) => {
				text += `  ${cookie.name}: ${cookie.value}\n`;
			});
		}

		return text;
	});
</script>

<textarea
	readonly
	class="h-full w-full flex-1 resize-none overflow-auto rounded border bg-muted p-2 font-mono text-xs"
	style="white-space: pre;"
	value={requestText}
></textarea>
