<script lang="ts">
	import ArchiveXIcon from '@lucide/svelte/icons/archive-x';
	import FileIcon from '@lucide/svelte/icons/file';
	import InboxIcon from '@lucide/svelte/icons/inbox';
	import SendIcon from '@lucide/svelte/icons/send';
	import SunIcon from '@lucide/svelte/icons/sun';
	import MoonIcon from '@lucide/svelte/icons/moon';
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import * as InputGroup from "$lib/components/ui/input-group/index.js";
	import * as ButtonGroup from "$lib/components/ui/button-group/index.js";
	import Wand from '@lucide/svelte/icons/wand';
	import FilePlus from '@lucide/svelte/icons/file-plus';
	import FolderPlus from '@lucide/svelte/icons/folder-plus';
	import ChartPieIcon from '@lucide/svelte/icons/chart-pie';
	import CommandIcon from '@lucide/svelte/icons/command';
	import FrameIcon from '@lucide/svelte/icons/frame';
	  import { Input } from "$lib/components/ui/input/index.js";
	import GalleryVerticalEndIcon from '@lucide/svelte/icons/gallery-vertical-end';
	import MapIcon from '@lucide/svelte/icons/map';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import NavFiles from './nav-files.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import Settings from '@lucide/svelte/icons/settings';
	import { page } from '$app/stores';
	import { Button } from './ui/button';
	import { CreateFile, CreateDir, OpenFile, GetFileContent, GetCurrentFilesState } from '$lib/wailsjs/go/main/App';
	import { fileStore } from '$lib/stores/fileStore.svelte';
	// This is sample data
	const data = {
		teams: [
			{
				name: 'Acme Inc',
				logo: GalleryVerticalEndIcon,
				plan: 'Enterprise'
			},
			{
				name: 'Acme Corp.',
				logo: AudioWaveformIcon,
				plan: 'Startup'
			},
			{
				name: 'Evil Corp.',
				logo: CommandIcon,
				plan: 'Free'
			}
		],
		projects: [
			{
				name: 'Design Engineering',
				url: '#',
				icon: FrameIcon
			},
			{
				name: 'Sales & Marketing',
				url: '#',
				icon: ChartPieIcon
			},
			{
				name: 'Travel',
				url: '#',
				icon: MapIcon
			}
		],
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		navMain: [
			{
				title: 'Inbox',
				url: '#',
				icon: InboxIcon,
				isActive: true
			},
			{
				title: 'Drafts',
				url: '#',
				icon: FileIcon,
				isActive: false
			},
			{
				title: 'Sent',
				url: '#',
				icon: SendIcon,
				isActive: false
			},
			{
				title: 'Junk',
				url: '#',
				icon: ArchiveXIcon,
				isActive: false
			}
		],
		mails: [
			{
				name: 'William Smith',
				email: 'williamsmith@example.com',
				subject: 'Meeting Tomorrow',
				date: '09:34 AM',
				teaser:
					'Hi team, just a reminder about our meeting tomorrow at 10 AM.\nPlease come prepared with your project updates.'
			},
			{
				name: 'Alice Smith',
				email: 'alicesmith@example.com',
				subject: 'Re: Project Update',
				date: 'Yesterday',
				teaser:
					"Thanks for the update. The progress looks great so far.\nLet's schedule a call to discuss the next steps."
			},
			{
				name: 'Bob Johnson',
				email: 'bobjohnson@example.com',
				subject: 'Weekend Plans',
				date: '2 days ago',
				teaser:
					"Hey everyone! I'm thinking of organizing a team outing this weekend.\nWould you be interested in a hiking trip or a beach day?"
			},
			{
				name: 'Emily Davis',
				email: 'emilydavis@example.com',
				subject: 'Re: Question about Budget',
				date: '2 days ago',
				teaser:
					"I've reviewed the budget numbers you sent over.\nCan we set up a quick call to discuss some potential adjustments?"
			},
			{
				name: 'Michael Wilson',
				email: 'michaelwilson@example.com',
				subject: 'Important Announcement',
				date: '1 week ago',
				teaser:
					"Please join us for an all-hands meeting this Friday at 3 PM.\nWe have some exciting news to share about the company's future."
			},
			{
				name: 'Sarah Brown',
				email: 'sarahbrown@example.com',
				subject: 'Re: Feedback on Proposal',
				date: '1 week ago',
				teaser:
					"Thank you for sending over the proposal. I've reviewed it and have some thoughts.\nCould we schedule a meeting to discuss my feedback in detail?"
			},
			{
				name: 'David Lee',
				email: 'davidlee@example.com',
				subject: 'New Project Idea',
				date: '1 week ago',
				teaser:
					"I've been brainstorming and came up with an interesting project concept.\nDo you have time this week to discuss its potential impact and feasibility?"
			},
			{
				name: 'Olivia Wilson',
				email: 'oliviawilson@example.com',
				subject: 'Vacation Plans',
				date: '1 week ago',
				teaser:
					"Just a heads up that I'll be taking a two-week vacation next month.\nI'll make sure all my projects are up to date before I leave."
			},
			{
				name: 'James Martin',
				email: 'jamesmartin@example.com',
				subject: 'Re: Conference Registration',
				date: '1 week ago',
				teaser:
					"I've completed the registration for the upcoming tech conference.\nLet me know if you need any additional information from my end."
			},
			{
				name: 'Sophia White',
				email: 'sophiawhite@example.com',
				subject: 'Team Dinner',
				date: '1 week ago',
				teaser:
					"To celebrate our recent project success, I'd like to organize a team dinner.\nAre you available next Friday evening? Please let me know your preferences."
			}
		]
	};

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	let activeItem = $state(data.navMain[0]);
	let mails = $state(data.mails);
	const sidebar = useSidebar();

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
				const hasValidExtension = validExtensions.some(ext => fileName.toLowerCase().endsWith(ext));

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
							<a href="##" {...props}>
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
						{#each data.navMain as item (item.title)}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton
									tooltipContentProps={{
										hidden: false
									}}
									onclick={() => {
										activeItem = item;
										const mail = data.mails.sort(() => Math.random() - 0.5);
										mails = mail.slice(0, Math.max(5, Math.floor(Math.random() * 10) + 1));
										sidebar.setOpen(true);
									}}
									isActive={activeItem.title === item.title}
									class="px-2.5 md:px-2"
								>
									{#snippet tooltipContent()}
										{item.title}
									{/snippet}
									<item.icon />
									<span>{item.title}</span>
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
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
		<Sidebar.Root collapsible="none" class="hidden flex-1 md:flex overflow-x-hidden">
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
							<div class="text-xs text-destructive px-2">
								{createError}
							</div>
						{/if}
						<Button class="w-full" variant="outline" onclick={handleCancel}>
							Cancel
						</Button>
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
