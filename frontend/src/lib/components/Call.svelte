<script lang="ts">
	let { call }: { call: App.Call } = $props();
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import RequestDetails from '$lib/components/request-details.svelte';
	import ResponseDetails from '$lib/components/response-details.svelte';
	import { getHttpStatusColorClass } from '$lib/utils/httpStatus';

	// Get badge variant based on HTTP method
	function getMethodVariant(method: string): 'default' | 'secondary' | 'destructive' | 'outline' {
		switch (method.toUpperCase()) {
			case 'GET':
				return 'default'; // Blue
			case 'POST':
				return 'secondary'; // Green
			case 'PUT':
			case 'PATCH':
				return 'outline'; // Yellow
			case 'DELETE':
				return 'destructive'; // Red
			default:
				return 'outline';
		}
	}
</script>

<Card.Root class="flex h-full flex-col">
	<Card.Header class="shrink-0">
		<Card.Title>
			<div class="flex flex-col gap-2">
				<Badge variant={getMethodVariant(call.request.method)}>{call.request.method}</Badge>
				<span class="font-mono text-sm break-all">{call.request.url}</span>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.Content class="flex flex-1 flex-col overflow-hidden">
		<!-- Request Response -->
		<Card.Root class="flex flex-1 flex-col overflow-hidden">
			<Card.Content class="flex flex-1 flex-col overflow-hidden">
				<Tabs.Root value="response" class="flex h-full min-h-0 w-full flex-col">
					<Tabs.List class="shrink-0">
						<Tabs.Trigger value="response">Response</Tabs.Trigger>
						<Tabs.Trigger value="request">Request</Tabs.Trigger>
					</Tabs.List>
					<!-- Response -->
					<Tabs.Content value="response" class="flex min-h-0 flex-1 flex-col overflow-hidden">
						<Card.Root class="flex flex-1 flex-col overflow-hidden">
							<Card.Content class="flex-1">
								<ResponseDetails response={call.response} />
							</Card.Content>
						</Card.Root>
					</Tabs.Content>
					<!-- Request -->
					<Tabs.Content value="request" class="flex min-h-0 flex-1 flex-col overflow-hidden">
						<div class="flex flex-1 flex-col overflow-hidden">
							<RequestDetails request={call.request} />
						</div>
					</Tabs.Content>
				</Tabs.Root>
			</Card.Content>
			<Card.Footer class="flex-shrink-0">
				<div class="flex">
					<div class="flex gap-1">
						<Badge
							variant="outline"
							class={`border ${getHttpStatusColorClass(call.response.status, 'border-')}`}
							>Status {call.response.status}</Badge
						>
					</div>
				</div>
			</Card.Footer>
		</Card.Root>
	</Card.Content>
</Card.Root>
