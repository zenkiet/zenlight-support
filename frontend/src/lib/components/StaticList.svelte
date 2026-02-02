<script lang="ts">
	import { resourceStore } from '$lib/stores/resources.svelte';
	import AddStatic from './AddStatic.svelte';
	import StaticCard from './StaticCard.svelte';
	import { fade } from 'svelte/transition';

	let { searchQuery } = $props<{
		searchQuery: string;
	}>();

	const filteredProjects = $derived.by(() =>
		resourceStore.directories.filter(
			(d) =>
				d.config?.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				d.config?.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
				d.config?.path.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<div class="flex flex-1 flex-col gap-4">
	{#if filteredProjects.length === 0 && searchQuery}
		<div class="flex flex-col items-center justify-center py-20 text-muted/50" in:fade>
			<i class="mb-4 text-4xl icon-[regular--magnifying-glass]"></i>
			<p class="text-sm font-medium">No projects found for "{searchQuery}"</p>
			<button
				class="mt-2 text-xs text-blue-400 hover:text-blue-300 hover:underline"
				onclick={() => (searchQuery = '')}
			>
				Clear filter
			</button>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
			{#each filteredProjects as project (project.config?.id)}
				<div in:fade={{ duration: 200 }}>
					<StaticCard {project} />
				</div>
			{/each}

			<AddStatic />
		</div>
	{/if}
</div>
