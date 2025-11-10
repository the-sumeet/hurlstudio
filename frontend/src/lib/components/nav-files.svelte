<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
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

	// State for delete confirmation dialog
	let deleteDialogOpen = $state(false);
	let fileToDelete = $state<{ path: string; name: string; isDirectory: boolean } | null>(null);

	// State for rename dialog
	let renameDialogOpen = $state(false);
	let fileToRename = $state<{ path: string; name: string; isDirectory: boolean } | null>(null);
	let newFileName = $state('');
	let renameError = $state('');

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

	function handleDeleteClick(filePath: string, fileName: string, isDirectory: boolean) {
		fileToDelete = { path: filePath, name: fileName, isDirectory };
		deleteDialogOpen = true;
	}

	async function confirmDelete() {
		if (!fileToDelete) return;

		try {
			// Check if the file being deleted is currently selected
			const isCurrentFile = fileStore.currentFile?.path === fileToDelete.path;

			await DeleteFile(fileToDelete.path);

			// If the deleted file was selected, clear the editor
			if (isCurrentFile) {
				fileStore.setCurrentFile(null);
				fileStore.setContent('');
			}

			await loadFiles();
			deleteDialogOpen = false;
			fileToDelete = null;
		} catch (error) {
			console.error('Failed to delete file:', error);
			alert(`Failed to delete: ${error}`);
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

	function handleRenameClick(filePath: string, fileName: string, isDirectory: boolean) {
		fileToRename = { path: filePath, name: fileName, isDirectory };
		newFileName = fileName;
		renameError = '';
		renameDialogOpen = true;
	}

	async function confirmRename() {
		if (!fileToRename || !newFileName.trim() || newFileName === fileToRename.name) {
			renameDialogOpen = false;
			return;
		}

		// Validate file extension if it's a file (not a directory)
		if (!fileToRename.isDirectory) {
			const trimmedName = newFileName.trim();
			const validExtensions = ['.md', '.markdown', '.hurl'];
			const hasValidExtension = validExtensions.some(ext => trimmedName.toLowerCase().endsWith(ext));

			if (!hasValidExtension) {
				renameError = 'Invalid file extension. File must end with .md, .markdown, or .hurl';
				return;
			}
		}

		try {
			const isCurrentFile = fileStore.currentFile?.path === fileToRename.path;

			await RenameFile(fileToRename.path, newFileName.trim());

			// If the renamed file was the currently selected file, update the fileStore
			if (isCurrentFile && currentState.currentDir) {
				const newPath = `${currentState.currentDir.path}/${newFileName.trim()}`;
				const content = await GetFileContent(newPath);

				// Update fileStore with new path
				const state = await GetCurrentFilesState();
				if (state.currentFile) {
					fileStore.setCurrentFile(state.currentFile);
					fileStore.setContent(content);
				}
			}

			await loadFiles();
			renameDialogOpen = false;
			fileToRename = null;
			newFileName = '';
			renameError = '';
		} catch (error) {
			console.error('Failed to rename:', error);
			renameError = `Failed to rename: ${error}`;
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
				<Tooltip.Root delayDuration={500}>
					<Tooltip.Trigger asChild>
						{#snippet child({ props })}
							<Sidebar.MenuButton
								{...props}
								onclick={() => handleOpenFile(file)}
								isActive={currentState.currentFile?.path === file.path}
							>
								{#if file.isDirectory}
									<FolderIcon />
								{:else}
									<FileIcon />
								{/if}
								<span>{file.name}</span>
							</Sidebar.MenuButton>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content side="right">
						<p class="text-xs">{file.path}</p>
					</Tooltip.Content>
				</Tooltip.Root>
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
						<DropdownMenu.Item onclick={() => handleRenameClick(file.path, file.name, file.isDirectory)}>
							<PencilIcon class="text-muted-foreground" />
							<span>Rename</span>
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item onclick={() => handleDeleteClick(file.path, file.name, file.isDirectory)}>
							<Trash2Icon class="text-muted-foreground" />
							<span>Delete</span>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>

<!-- Delete Confirmation Dialog -->
<AlertDialog.Root bind:open={deleteDialogOpen}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
			<AlertDialog.Description>
				{#if fileToDelete}
					This action cannot be undone. This will permanently delete the {fileToDelete.isDirectory ? 'folder' : 'file'} "<strong>{fileToDelete.name}</strong>"{fileToDelete.isDirectory ? ' and all its contents' : ''}.
				{/if}
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={confirmDelete}>Delete</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<!-- Rename Dialog -->
<Dialog.Root bind:open={renameDialogOpen}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Rename {fileToRename?.isDirectory ? 'Folder' : 'File'}</Dialog.Title>
			<Dialog.Description>
				Enter a new name for "{fileToRename?.name}".
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="newname" class="text-right">Name</Label>
				<Input
					id="newname"
					bind:value={newFileName}
					class="col-span-3"
					autocapitalize="off"
					autocomplete="off"
					autocorrect="off"
					spellcheck="false"
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button type="submit" onclick={confirmRename}>Rename</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
