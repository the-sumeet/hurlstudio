<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import FolderIcon from '@lucide/svelte/icons/folder';
	import FileIcon from '@lucide/svelte/icons/file';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import { onMount } from 'svelte';
	import {
		ListFiles,
		OpenFile,
		GoUp,
		CreateFile,
		CreateDir,
		DeleteFile,
		RenameFile,
		GetCurrentFilesState,
		GetFileContent,
		SaveLastOpenedState,
		LoadLastOpenedState
	} from '$lib/wailsjs/go/main/App';
	import FolderUp from '@lucide/svelte/icons/folder-up';
	import { fileStore } from '$lib/stores/fileStore.svelte';

	type FileEntry = {
		name: string;
		path: string;
		isDirectory: boolean;
		size: number;
		modTime: string;
	};

	type CurrentFilesState = {
		currentDir?: FileEntry;
		currentFile?: FileEntry;
	};

	let files = $state<FileEntry[]>([]);
	let currentState = $state<CurrentFilesState>({});
	const sidebar = useSidebar();

	async function loadFiles() {
		try {
			const state = await GetCurrentFilesState();
			currentState = state;
			if (state.currentDir) {
				const fileList = await ListFiles(state.currentDir.path);
				files = fileList || [];
			}
		} catch (error) {
			console.error('Failed to load files:', error);
		}
	}

	async function handleOpenFile(file: FileEntry) {
		try {
			await OpenFile(file.path);
			await loadFiles();

			// If it's a file (not a directory), load its content
			if (!file.isDirectory) {
				const content = await GetFileContent(file.path);
				fileStore.setCurrentFile(file);
				fileStore.setContent(content);
			}

			// Save the current state
			await SaveLastOpenedState();
		} catch (error) {
			console.error('Failed to open file:', error);
		}
	}

	async function handleGoUp() {
		try {
			await GoUp();
			await loadFiles();

			// Save the current state
			await SaveLastOpenedState();
		} catch (error) {
			console.error('Failed to go up:', error);
		}
	}

	async function handleDeleteFile(filePath: string) {
		try {
			await DeleteFile(filePath);
			await loadFiles();
		} catch (error) {
			console.error('Failed to delete file:', error);
		}
	}

	async function handleCreateFile() {
		const fileName = prompt('Enter file name:');
		if (fileName) {
			try {
				await CreateFile(fileName);
				await loadFiles();
			} catch (error) {
				console.error('Failed to create file:', error);
				alert('Failed to create file: ' + error);
			}
		}
	}

	async function handleCreateDir() {
		const dirName = prompt('Enter directory name:');
		if (dirName) {
			try {
				await CreateDir(dirName);
				await loadFiles();
			} catch (error) {
				console.error('Failed to create directory:', error);
				alert('Failed to create directory: ' + error);
			}
		}
	}

	async function handleRenameFile(filePath: string, currentName: string) {
		const newName = prompt('Enter new name:', currentName);
		if (newName && newName !== currentName) {
			try {
				await RenameFile(filePath, newName);
				await loadFiles();
			} catch (error) {
				console.error('Failed to rename file:', error);
				alert('Failed to rename: ' + error);
			}
		}
	}

	onMount(async () => {
		try {
			// Try to load last opened state
			const lastState = await LoadLastOpenedState();

			if (lastState.currentDir || lastState.currentFile) {
				// State was loaded successfully
				currentState = lastState;

				// Load files from the last directory
				if (lastState.currentDir) {
					const fileList = await ListFiles(lastState.currentDir.path);
					files = fileList || [];
				}

				// Restore the last opened file
				if (lastState.currentFile) {
					const content = await GetFileContent(lastState.currentFile.path);
					fileStore.setCurrentFile(lastState.currentFile);
					fileStore.setContent(content);
				}
			} else {
				// No saved state, load default
				await loadFiles();
			}
		} catch (error) {
			console.error('Failed to load last state:', error);
			// Fall back to normal loading
			await loadFiles();
		}

		// Listen for refresh events
		const handleRefresh = () => {
			loadFiles();
		};
		window.addEventListener('refresh-files', handleRefresh);

		// Cleanup on unmount
		return () => {
			window.removeEventListener('refresh-files', handleRefresh);
		};
	});
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<!-- <Sidebar.GroupLabel>
		<div class="flex items-center justify-between w-full">
			<span>Files</span>
			<div class="flex gap-1">
				<button
					onclick={handleCreateFile}
					class="hover:bg-sidebar-accent rounded p-1"
					title="New File"
				>
					<FileIcon class="h-3 w-3" />
				</button>
				<button
					onclick={handleCreateDir}
					class="hover:bg-sidebar-accent rounded p-1"
					title="New Folder"
				>
					<FolderIcon class="h-3 w-3" />
				</button>
			</div>
		</div>
	</Sidebar.GroupLabel> -->

	{#if currentState.currentDir}
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton onclick={handleGoUp}>
					<FolderUp />
					<span>..</span>
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	{/if}

	<Sidebar.Menu>
		{#each files || [] as file (file.path)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					onclick={() => handleOpenFile(file)}
					isActive={currentState.currentFile?.path === file.path}
				>
					{#if file.isDirectory}
						<FolderIcon />
					{:else}
						<FileIcon />
					{/if}
					<span class="truncate">{file.name}</span>
				</Sidebar.MenuButton>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuAction showOnHover {...props}>
								<EllipsisIcon />
								<span class="sr-only">More</span>
							</Sidebar.MenuAction>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-48 rounded-lg"
						side={sidebar.isMobile ? 'bottom' : 'right'}
						align={sidebar.isMobile ? 'end' : 'start'}
					>
						<DropdownMenu.Item onclick={() => handleRenameFile(file.path, file.name)}>
							<PencilIcon class="text-muted-foreground" />
							<span>Rename</span>
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item onclick={() => handleDeleteFile(file.path)}>
							<Trash2Icon class="text-muted-foreground" />
							<span>Delete</span>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
