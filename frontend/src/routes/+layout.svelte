<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import { Separator } from "$lib/components/ui/separator/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { fileStore } from '$lib/stores/fileStore.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider style="--sidebar-width: 300px;">
	<AppSidebar />
	<Sidebar.Inset class="overflow-x-hidden">
		<header class="bg-background sticky top-0 flex shrink-0 items-center gap-2 border-b p-4 w-full">
			<Sidebar.Trigger class="-ml-1" />
			<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
			<div class="flex-1 min-w-0">
				{#if fileStore.currentFile}
					<h1 class="text-sm font-medium truncate" title={fileStore.currentFile.name}>
						fileStore.currentFile.name
						Lorem ipsum dolor sit amet consectetur adipisicing elit. Enim reiciendis, aliquam dolor ratione reprehenderit numquam veniam earum, recusandae rem sit error maiores obcaecati, totam dolorum cum exercitationem quos blanditiis veritatis placeat. Itaque, soluta! Animi quo fugit consequuntur corporis autem sunt? Repudiandae unde ratione dolore commodi culpa qui dolorem veniam tempore! Dolorem quae numquam nisi voluptas reiciendis esse alias, accusantium perspiciatis quasi error magnam atque, deleniti impedit quas? Provident rem animi quo. Consequuntur, sed ad eligendi excepturi quasi porro harum quibusdam obcaecati quo similique blanditiis corrupti! Repudiandae sapiente alias quasi. Consequatur maiores alias at eveniet pariatur soluta unde? Exercitationem natus corrupti aut doloremque obcaecati molestias at iure sunt incidunt. Repellat similique recusandae ullam fugit natus id aperiam sapiente tempore asperiores. Earum officiis officia quisquam veniam quidem. Soluta, omnis maxime illum dolor doloremque adipisci nesciunt nostrum quidem dolorum excepturi aspernatur vitae obcaecati reiciendis accusantium nihil. Libero optio vero sapiente consectetur molestias sequi saepe placeat perspiciatis culpa alias nesciunt cupiditate commodi iure, hic neque amet deleniti quas voluptate harum? Voluptate eaque recusandae, autem, et ea tempore maxime laudantium perferendis est sapiente voluptates aperiam assumenda fugit necessitatibus. Esse numquam cum, maxime illum ut ratione autem itaque, sed soluta sit dicta sunt architecto nesciunt possimus voluptates consectetur vel. Ipsa voluptate nesciunt blanditiis nihil repellendus non perspiciatis labore, quos est. Quibusdam aut aliquid, quod esse impedit enim, similique ducimus deleniti quo ea corporis ad facilis necessitatibus alias eos labore maxime. Consectetur veritatis porro sint praesentium ea doloribus sequi reprehenderit omnis! Eveniet veritatis corporis, totam nemo ducimus nisi accusamus incidunt aperiam voluptates fugit veniam laboriosam, quos odio, omnis nam molestias perspiciatis molestiae autem est asperiores ut? Quod voluptatibus blanditiis omnis rem, voluptate iure nostrum ratione perspiciatis, ipsa esse nobis distinctio alias a ut dolore delectus. Laboriosam similique itaque architecto delectus exercitationem eveniet quos, excepturi perspiciatis sint sapiente?
					</h1>
				{:else}
					<h1 class="text-sm font-medium text-muted-foreground">
						No file selected
					</h1>
				{/if}
			</div>
		</header>
		<div class="flex flex-1 flex-col gap-4">
			{@render children()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>


