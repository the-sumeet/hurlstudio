<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as NativeSelect from '$lib/components/ui/native-select/index.js';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import { goto } from '$app/navigation';
	import { Kbd } from '$lib/components/ui/kbd/index.js';
	import { LoadEnvVariables, SaveEnvVariables } from '$lib/wailsjs/go/main/App';
	import { onMount } from 'svelte';

	let value = $state('');
	let selectedEnvironment = $state('development');
	let isValid = $state(true);
	let isSaving = $state(false);
	let isLoading = $state(true);

	// Load environment variables on mount
	onMount(async () => {
		try {
			const envJson = await LoadEnvVariables();
			value = envJson;
			validateJson(envJson);
		} catch (error) {
			console.error('Failed to load environment variables:', error);
		} finally {
			isLoading = false;
		}
	});

	function onchange(newContent: string) {
		value = newContent;
		validateJson(newContent);
	}

	function validateJson(content: string): boolean {
		try {
			const parsed = JSON.parse(content);

			// Validate required structure
			if (!parsed.version || typeof parsed.version !== 'string') {
				isValid = false;
				return false;
			}

			if (!parsed.activeEnvironment || typeof parsed.activeEnvironment !== 'string') {
				isValid = false;
				return false;
			}

			if (!parsed.global || typeof parsed.global !== 'object' || Array.isArray(parsed.global)) {
				isValid = false;
				return false;
			}

			if (
				!parsed.environments ||
				typeof parsed.environments !== 'object' ||
				Array.isArray(parsed.environments)
			) {
				isValid = false;
				return false;
			}

			// Validate that all values in global and environments are strings
			for (const value of Object.values(parsed.global)) {
				if (typeof value !== 'string') {
					isValid = false;
					return false;
				}
			}

			for (const env of Object.values(parsed.environments)) {
				if (typeof env !== 'object' || Array.isArray(env)) {
					isValid = false;
					return false;
				}
				for (const value of Object.values(env as object)) {
					if (typeof value !== 'string') {
						isValid = false;
						return false;
					}
				}
			}

			isValid = true;
			return true;
		} catch (error) {
			isValid = false;
			return false;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			goto('/');
		}
	}

	async function handleSave() {
		if (!isValid) return;

		isSaving = true;
		try {
			await SaveEnvVariables(value);
			console.log('Environment variables saved successfully');
		} catch (error) {
			console.error('Failed to save environment variables:', error);
			alert(`Failed to save: ${error}`);
		} finally {
			isSaving = false;
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<Card.Root class="h-full rounded-none">
	<Card.Header>
		<Card.Title>Environment Variables</Card.Title>
		<Card.Action>
			{#if !isValid}
				<span class="text-xs text-destructive">Invalid JSON structure</span>
			{/if}
		</Card.Action>
	</Card.Header>
	<Card.Content class="flex-1">
		{#if isLoading}
			<div class="flex h-full items-center justify-center">
				<p class="text-sm text-muted-foreground">Loading...</p>
			</div>
		{:else}
			<MonacoEditor {value} {onchange} language="json" />
		{/if}
	</Card.Content>
	<Card.Footer class="flex gap-2">
		<Button variant="default" onclick={handleSave} disabled={!isValid || isSaving}>
			{isSaving ? 'Saving...' : 'Save'}
		</Button>
		<Button href="/" variant="outline" class="gap-2">Close <Kbd>ESC</Kbd></Button>
	</Card.Footer>
</Card.Root>
