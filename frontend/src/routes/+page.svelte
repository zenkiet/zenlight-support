<script lang="ts">
	import ServiceList from '$lib/components/ServiceList.svelte';
	import { serviceStore } from '$lib/stores/services.svelte';
	import { fade } from 'svelte/transition';

	// Local State
	let searchQuery = $state('');
</script>

<div class="bg-glass flex h-screen w-full flex-col overflow-hidden text-neutral-200" in:fade={{ duration: 400 }}>
	<header class="flex shrink-0 flex-col gap-5 px-6 pt-6 pb-2">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-4">
				<h1 class="text-2xl font-bold tracking-tight text-white">Service Watcher</h1>

				<div class="flex items-center gap-2 border-l border-white/10 pl-4">
            <!-- <div class="flex items-center gap-1.5 rounded-full bg-white/5 px-2.5 py-1 text-xs font-medium text-white/50 ring-1 ring-inset ring-white/10">
                <span class="h-1.5 w-1.5 rounded-full bg-white/40"></span>
                Total: <span class="text-white">{serviceStore.totalCount}</span>
            </div> -->
            <div class="flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-2.5 py-1 text-xs font-medium text-emerald-400 ring-1 ring-inset ring-emerald-500/20">
                <span class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                Running: <span class="text-white">{serviceStore.runningCount}</span>
            </div>
             <!-- {#if stoppedServices > 0} -->
                <div class="flex items-center gap-1.5 rounded-full bg-rose-500/10 px-2.5 py-1 text-xs font-medium text-rose-400 ring-1 ring-inset ring-rose-500/20">
                    <span class="h-1.5 w-1.5 rounded-full bg-rose-500"></span>
                    Stopped: <span class="text-white">{serviceStore.stoppedCount}</span>
                </div>
            <!-- {/if} -->
        </div>
			</div>

			<a
				href="/settings"
				class="group flex h-9 w-9 items-center justify-center rounded-lg border border-white/5 bg-white/5 text-white/40 transition-all hover:bg-white/10 hover:text-white"
				title="Settings"
			>
				<i
					class="text-lg transition-transform duration-500 icon-[duotone--gear] group-hover:rotate-90"
				></i>
			</a>
		</div>

		<div class="flex items-center justify-between backdrop-blur-md">
			<div class="group relative max-w-md flex-1">
				<i
					class="absolute top-1/2 left-3 -translate-y-1/2 text-lg text-white/20 transition-colors icon-[duotone--magnifying-glass] group-focus-within:text-white/50"
				></i>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Filter services by name, PID or description..."
					class="h-10 w-full rounded-lg border border-transparent bg-transparent pr-4 pl-10 text-sm text-white placeholder-white/30 transition-all focus:border-white/10 focus:bg-white/5 focus:ring-0 focus:outline-none"
				/>
			</div>
		</div>
	</header>

	<main class="scrollable flex-1 overflow-y-auto scroll-smooth px-6 py-4">
		<ServiceList {searchQuery}  />

		<div class="h-10"></div>
	</main>
</div>
