<script lang="ts">
	let { call, asserts = [] }: { call: App.Call; asserts?: App.Assert[] } = $props();
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import RequestDetails from '$lib/components/request-details.svelte';
	import ResponseDetails from '$lib/components/response-details.svelte';
	import { getHttpStatusColorClass } from '$lib/utils/httpStatus';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import CircleCheck from '@lucide/svelte/icons/circle-check';
	import CircleX from '@lucide/svelte/icons/circle-x';

	let hasFailedAsserts = $derived(asserts.some((a) => !a.success));
	let defaultTab = $derived(hasFailedAsserts ? 'asserts' : 'response');

	function getMethodVariant(method: string): 'default' | 'secondary' | 'destructive' | 'outline' {
		switch (method.toUpperCase()) {
			case 'GET':
				return 'default';
			case 'POST':
				return 'secondary';
			case 'PUT':
			case 'PATCH':
				return 'outline';
			case 'DELETE':
				return 'destructive';
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
		<Card.Root class="flex flex-1 flex-col overflow-hidden">
			<Card.Content class="flex flex-1 flex-col overflow-hidden">
				<Tabs.Root value={defaultTab} class="flex h-full min-h-0 w-full flex-col">
					<Tabs.List class="shrink-0">
						<Tabs.Trigger value="response">Response</Tabs.Trigger>
						<Tabs.Trigger value="request">Request</Tabs.Trigger>
						{#if asserts.length > 0}
							<Tabs.Trigger value="asserts">
								Asserts
								{#if hasFailedAsserts}
									<span class="ml-1 rounded-full bg-destructive px-1.5 py-0.5 text-xs text-destructive-foreground">
										{asserts.filter((a) => !a.success).length}
									</span>
								{/if}
							</Tabs.Trigger>
						{/if}
					</Tabs.List>
					<Tabs.Content value="response" class="flex min-h-0 flex-1 flex-col overflow-hidden">
						<Card.Root class="flex flex-1 flex-col overflow-hidden">
							<Card.Content class="flex-1">
								<ResponseDetails response={call.response} />
							</Card.Content>
						</Card.Root>
					</Tabs.Content>
					<Tabs.Content value="request" class="flex min-h-0 flex-1 flex-col overflow-hidden">
						<div class="flex flex-1 flex-col overflow-hidden">
							<RequestDetails request={call.request} />
						</div>
					</Tabs.Content>
					{#if asserts.length > 0}
						<Tabs.Content value="asserts" class="flex min-h-0 flex-1 flex-col overflow-y-auto">
							<div class="flex flex-col gap-2 p-2">
								{#each asserts as assert}
									{#if assert.success}
										<Alert.Root variant="default" class="border-green-500/40 bg-green-500/5">
											<CircleCheck class="h-4 w-4 text-green-500" />
											<Alert.Title class="text-green-700 dark:text-green-400">Pass — Line {assert.line}</Alert.Title>
										</Alert.Root>
									{:else}
										<Alert.Root variant="destructive">
											<CircleX class="h-4 w-4" />
											<Alert.Title>Fail — Line {assert.line}</Alert.Title>
											{#if assert.message}
												<Alert.Description>
													<pre class="mt-1 overflow-x-auto whitespace-pre-wrap text-xs">{assert.message}</pre>
												</Alert.Description>
											{/if}
										</Alert.Root>
									{/if}
								{/each}
							</div>
						</Tabs.Content>
					{/if}
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
