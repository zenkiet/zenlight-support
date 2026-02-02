<script lang="ts">
	import { resourceStore } from '$lib/stores/resources.svelte';
	import AddService from './AddService.svelte';
	import ServiceCard from './ServiceCard.svelte';
	import { fade } from 'svelte/transition';

	let { searchQuery } = $props<{
		searchQuery: string;
	}>();

	const filteredServices = $derived.by(() =>
		resourceStore.services.filter(
			(service) =>
				service.config?.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				service.config?.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
				service.metrics?.pid?.toString().includes(searchQuery)
		)
	);
</script>

<div class="flex flex-1 flex-col gap-4">
	{#if filteredServices.length === 0 && searchQuery}
		<div class="flex flex-col items-center justify-center py-20 text-muted/50" in:fade>
			<i class="mb-4 text-4xl icon-[regular--magnifying-glass]"></i>
			<p class="text-sm font-medium">No services found for "{searchQuery}"</p>
			<button
				class="mt-2 text-xs text-blue-400 hover:text-blue-300 hover:underline"
				onclick={() => (searchQuery = '')}
			>
				Clear filter
			</button>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
    {#each filteredServices as service (service.config?.id)}
        <div in:fade={{ duration: 200 }}>
            <ServiceCard {service} />
        </div>
    {/each}

		<AddService />
		</div>
	{/if}
</div>
