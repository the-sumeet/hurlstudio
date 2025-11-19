<script lang="ts">
	import SunIcon from '@lucide/svelte/icons/sun';
	import MoonIcon from '@lucide/svelte/icons/moon';
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import Settings from '@lucide/svelte/icons/settings';
	import Wand from '@lucide/svelte/icons/wand';
	import FilePlus from '@lucide/svelte/icons/file-plus';
	import FolderPlus from '@lucide/svelte/icons/folder-plus';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from './ui/button';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import NavFiles from './nav-files.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { page } from '$app/stores';
	import {
		CreateFile,
		CreateDir,
		OpenFile,
		GetFileContent,
		GetCurrentFilesState
	} from '$lib/wailsjs/go/main/App';
	import { fileStore } from '$lib/stores/fileStore.svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	// State for create new file/folder
	let showCreateInput = $state(false);
	let createInputValue = $state('');
	let createType = $state<'file' | 'folder'>('file'); // Track if creating file or folder
	let createError = $state('');

	function handleNewFileClick() {
		createType = 'file';
		showCreateInput = true;
		createInputValue = '';
		createError = '';
	}

	function handleNewFolderClick() {
		createType = 'folder';
		showCreateInput = true;
		createInputValue = '';
		createError = '';
	}

	function handleCancel() {
		showCreateInput = false;
		createInputValue = '';
		createError = '';
	}

	async function handleCreateDone() {
		if (!createInputValue.trim()) {
			showCreateInput = false;
			return;
		}

		try {
			if (createType === 'file') {
				// Validate file extension
				const fileName = createInputValue.trim();
				const validExtensions = ['.md', '.markdown', '.hurl'];
				const hasValidExtension = validExtensions.some((ext) =>
					fileName.toLowerCase().endsWith(ext)
				);

				if (!hasValidExtension) {
					createError = 'Invalid file extension. File must end with .md, .markdown, or .hurl';
					return;
				}

				await CreateFile(fileName);

				// Get current directory and open the newly created file
				const currentState = await GetCurrentFilesState();
				if (currentState.currentDir) {
					const filePath = `${currentState.currentDir.path}/${fileName}`;
					const state = await OpenFile(filePath);

					if (state.currentFile) {
						const content = await GetFileContent(state.currentFile.path);
						fileStore.setCurrentFile(state.currentFile);
						fileStore.setContent(content);
					}
				}

				// Trigger a refresh of the file list
				window.dispatchEvent(new CustomEvent('refresh-files'));
			} else {
				// Create folder - no validation needed
				await CreateDir(createInputValue.trim());

				// Get current directory and open the newly created folder
				const currentState = await GetCurrentFilesState();
				if (currentState.currentDir) {
					const folderPath = `${currentState.currentDir.path}/${createInputValue.trim()}`;
					await OpenFile(folderPath);
				}

				// Trigger a refresh of the file list
				window.dispatchEvent(new CustomEvent('refresh-files'));
			}

			// Reset state
			showCreateInput = false;
			createInputValue = '';
			createError = '';
		} catch (error) {
			console.error(`Failed to create ${createType}:`, error);
			createError = `Failed to create ${createType}: ${error}`;
		}
	}
</script>

<Sidebar.Root
	bind:ref
	collapsible="icon"
	class="overflow-hidden [&>[data-sidebar=sidebar]]:flex-row"
	{...restProps}
	id="main-sidebar-root"
>
	<!-- This is the first sidebar -->
	<!-- We disable collapsible and adjust width to icon. -->
	<!-- This will make the sidebar appear as icons. -->
	<Sidebar.Root collapsible="none" class="!w-[calc(var(--sidebar-width-icon)_+_1px)] border-r">
		<Sidebar.Header>
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg" class="md:h-8 md:p-0">
						{#snippet child({ props })}
							<a href="/" {...props}>
								<div
									class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
								>
									<ArrowLeftRightIcon class="size-4" />
								</div>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-medium">Hurl Studio</span>
								</div>
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		</Sidebar.Header>
		<Sidebar.Content>
			<Sidebar.Group>
				<Sidebar.GroupContent class="px-1.5 md:px-0">
					<Sidebar.Menu>
						<Sidebar.MenuItem>
							<Sidebar.MenuButton
								tooltipContentProps={{
									hidden: false
								}}
								class="px-2.5 md:px-2"
							>
								{#snippet tooltipContent()}
									Settings
								{/snippet}
								{#snippet child({ props })}
									<a href="/formatter" {...props}>
										<Wand />
										<span>Formatter</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
						<Sidebar.MenuItem>
							<Sidebar.MenuButton
								tooltipContentProps={{
									hidden: false
								}}
								class="px-2.5 md:px-2"
							>
								{#snippet tooltipContent()}
									Settings
								{/snippet}
								{#snippet child({ props })}
									<a href="/settings" {...props}>
										<Settings />
										<span>Settings</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
						<Sidebar.MenuItem>
							<Sidebar.MenuButton
								tooltipContentProps={{
									hidden: false
								}}
								onclick={() => themeStore.toggle()}
								class="px-2.5 md:px-2"
							>
								{#snippet tooltipContent()}
									Toggle Theme
								{/snippet}
								{#if themeStore.current === 'dark'}
									<SunIcon />
									<span>Light</span>
								{:else}
									<MoonIcon />
									<span>Dark</span>
								{/if}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Sidebar.Group>
		</Sidebar.Content>
		<!-- <Sidebar.Footer>
			<NavUser user={data.user} />
		</Sidebar.Footer> -->
	</Sidebar.Root>

	{#if $page.url.pathname === '/'}
		<Sidebar.Root collapsible="none" class="hidden flex-1 overflow-x-hidden md:flex">
			<Sidebar.Content>
				<NavFiles />
			</Sidebar.Content>
			<!-- <Sidebar.Footer>
				<NavUser user={data.user} />
			</Sidebar.Footer> -->
			<Sidebar.Rail />
			<Sidebar.Footer class="p-2">
				{#if showCreateInput}
					<div class="flex flex-col gap-2">
						<InputGroup.Root>
							<InputGroup.Input
								placeholder={createType === 'file' ? 'Enter file name...' : 'Enter folder name...'}
								bind:value={createInputValue}
								autocapitalize="off"
								autocomplete="off"
								autocorrect="off"
								spellcheck="false"
							/>
							<InputGroup.Addon align="inline-end">
								<InputGroup.Button onclick={handleCreateDone}>Done</InputGroup.Button>
							</InputGroup.Addon>
						</InputGroup.Root>
						{#if createError}
							<div class="px-2 text-xs text-destructive">
								{createError}
							</div>
						{/if}
						<Button class="w-full" variant="outline" onclick={handleCancel}>Cancel</Button>
					</div>
				{:else}
					<ButtonGroup.Root class="flex w-full">
						<Button class="flex-1" size="icon" variant="outline" onclick={handleNewFileClick}>
							<FilePlus />
						</Button>
						<Button class="flex-1" size="icon" variant="outline" onclick={handleNewFolderClick}>
							<FolderPlus />
						</Button>
					</ButtonGroup.Root>
				{/if}
			</Sidebar.Footer>
		</Sidebar.Root>
	{/if}
</Sidebar.Root>
