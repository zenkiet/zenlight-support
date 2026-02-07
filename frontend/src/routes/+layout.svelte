<script lang="ts">
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import UpdateModal from '$lib/components/modals/UpdateModal.svelte';
	import { resourceStore } from '$lib/stores/resources.svelte';
	import { fade } from 'svelte/transition';
	import './layout.css';
	import { onMount } from 'svelte';
	import { systemStore } from '$lib/stores/system.svelte';

	let { children } = $props();

	let isSidebarOpen = $state(false);

	onMount(() => {
		resourceStore.init();
		systemStore.init();
	});
</script>

<div
	class="bg-surface-base flex h-screen w-screen flex-col overflow-hidden selection:bg-blue-500/30"
>
	<header
		class="flex h-14 shrink-0 items-center justify-between border-b border-muted/15 bg-surface px-4 md:hidden"
	>
		<div class="flex items-center gap-2">
			<button
				onclick={() => (isSidebarOpen = true)}
				class="hover:bg-surface-highlight inline-flex cursor-pointer rounded-md p-2 text-muted transition-colors hover:text-main"
				aria-label="Open menu"
			>
				<i class="size-5 icon-[regular--list]"></i>
			</button>
			<span class="font-bold text-main">ZenLight Support</span>
		</div>
	</header>

	<main class="relative flex flex-1 overflow-hidden">
		{#if resourceStore.initialized}
			<div class="flex h-full w-full overflow-hidden" in:fade={{ duration: 400 }}>
				{#if isSidebarOpen}
					<button
						class="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm transition-opacity md:hidden"
						in:fade={{ duration: 200 }}
						out:fade={{ duration: 200 }}
						onclick={() => (isSidebarOpen = false)}
						aria-label="Close sidebar overlay"
					></button>
				{/if}

				<Sidebar isOpen={isSidebarOpen} onClose={() => (isSidebarOpen = false)} />

				<main
					class="scrollable bg-surface-base relative w-full flex-1 overflow-y-auto scroll-smooth px-4 py-4 md:px-6"
				>
					{@render children()}
				</main>
			</div>
		{:else}
			<div class="flex h-full w-full items-center justify-center">
				<div class="text-center">
					<div
						class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-b-2 border-blue-500"
					></div>
					<p class="text-gray-600">Loading services...</p>
				</div>
			</div>
		{/if}
	</main>

	<UpdateModal />
</div>
