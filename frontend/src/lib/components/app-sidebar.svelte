<script lang="ts">
	import SunIcon from '@lucide/svelte/icons/sun';
	import MoonIcon from '@lucide/svelte/icons/moon';
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import Settings from '@lucide/svelte/icons/settings';
	import Wand from '@lucide/svelte/icons/wand';
	import FilePlus from '@lucide/svelte/icons/file-plus';
	import FolderPlus from '@lucide/svelte/icons/folder-plus';
	import Braces from '@lucide/svelte/icons/braces';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from './ui/button';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import NavFiles from './nav-files.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { page } from '$app/stores';
	import { Kbd } from '$lib/components/ui/kbd/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	// New utilities
	import { hasValidExtension, getExtensionValidationError } from '$lib/utils/fileExtensions';
	import { handleError, handleSuccess } from '$lib/utils/errorHandler';
	import {
		CreateFile,
		CreateDir,
		OpenFile,
		GetFileContent,
		GetCurrentFilesState
	} from '$lib/wailsjs/go/main/App';
	import { fileStore } from '$lib/stores/fileStore.svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	// State for create new file/folder dialog
	let dialogOpen = $state(false);
	let createInputValue = $state('');
	let createType = $state<'file' | 'folder'>('file');
	let createError = $state('');

	function handleNewFileClick() {
		createType = 'file';
		dialogOpen = true;
		createInputValue = '';
		createError = '';
	}

	function handleNewFolderClick() {
		createType = 'folder';
		dialogOpen = true;
		createInputValue = '';
		createError = '';
	}

	async function handleCreateDone() {
		if (!createInputValue.trim()) {
			dialogOpen = false;
			return;
		}

		try {
			if (createType === 'file') {
				// Validate file extension
				const fileName = createInputValue.trim();
				if (!hasValidExtension(fileName)) {
					createError = getExtensionValidationError();
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
				handleSuccess(`File "${fileName}" created successfully`);
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
				handleSuccess(`Folder "${createInputValue.trim()}" created successfully`);
			}

			// Reset state and close dialog
			dialogOpen = false;
			createInputValue = '';
			createError = '';
		} catch (error) {
			createError = `Failed to create ${createType}: ${error}`;
			handleError(error, `Failed to create ${createType}`);
		}
	}

	// Detect OS for keyboard shortcuts
	const isMac =
		typeof navigator !== 'undefined' &&
		(navigator.platform.toUpperCase().indexOf('MAC') >= 0 ||
			navigator.userAgent.toUpperCase().indexOf('MAC') >= 0);

	// Handle keyboard shortcuts
	function handleKeydown(event: KeyboardEvent) {
		// Cmd+N (Mac) or Ctrl+N (Windows/Linux) for new file
		if ((event.metaKey || event.ctrlKey) && event.key === 'n' && !event.shiftKey) {
			event.preventDefault();
			if (!dialogOpen) {
				handleNewFileClick();
			}
		}
		// Cmd+Shift+N (Mac) or Ctrl+Shift+N (Windows/Linux) for new folder
		if ((event.metaKey || event.ctrlKey) && event.key === 'n' && event.shiftKey) {
			event.preventDefault();
			if (!dialogOpen) {
				handleNewFolderClick();
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

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
									Format
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
									Variables
								{/snippet}
								{#snippet child({ props })}
									<a href="/variables" {...props}>
										<Braces />
										<span>Variables</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
						<!-- <Sidebar.MenuItem>
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
						</Sidebar.MenuItem> -->
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
				<ButtonGroup.Root class="flex w-full">
					<Button class="flex-1 flex-col gap-1 p-2" variant="outline" onclick={handleNewFileClick}>
						<FilePlus class="h-4 w-4" />
						<Kbd class="text-[8px]">{isMac ? '⌘N' : 'Ctrl+N'}</Kbd>
					</Button>
					<Button
						class="flex-1 flex-col gap-1 p-2"
						variant="outline"
						onclick={handleNewFolderClick}
					>
						<FolderPlus class="h-4 w-4" />
						<Kbd class="text-[8px]">{isMac ? '⌘⇧N' : 'Ctrl+⇧N'}</Kbd>
					</Button>
				</ButtonGroup.Root>
			</Sidebar.Footer>
		</Sidebar.Root>
	{/if}
</Sidebar.Root>

<!-- Create File/Folder Dialog -->
<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>
				{createType === 'file' ? 'Create New File' : 'Create New Folder'}
			</Dialog.Title>
			<Dialog.Description>
				{createType === 'file'
					? 'Enter a name for your new file. Supported extensions: .hurl, .md, .markdown'
					: 'Enter a name for your new folder.'}
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">
					{createType === 'file' ? 'File name' : 'Folder name'}
				</Label>
				<Input
					id="name"
					bind:value={createInputValue}
					placeholder={createType === 'file' ? 'example.hurl' : 'my-folder'}
					autocapitalize="off"
					autocomplete="off"
					autocorrect="off"
					spellcheck="false"
					onkeydown={(e) => {
						if (e.key === 'Enter') {
							handleCreateDone();
						}
					}}
				/>
				{#if createError}
					<p class="text-sm text-destructive">{createError}</p>
				{/if}
			</div>
		</div>

		<Dialog.Footer>
			<Button
				type="button"
				variant="outline"
				onclick={() => {
					dialogOpen = false;
					createError = '';
				}}
			>
				Cancel
			</Button>
			<Button type="button" onclick={handleCreateDone}>Create</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
