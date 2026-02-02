<script lang="ts">
	import StaticList from '$lib/components/StaticList.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { resourceStore } from '$lib/stores/resources.svelte';

	let searchQuery = $state('');

	const totalCount = $derived(resourceStore.directories.length);
</script>

<header class="mb-6 flex shrink-0 flex-col gap-5 pt-2">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-4">
			<div>
				<div class="flex items-center gap-4">
					<h1 class="text-2xl font-bold tracking-tight">Project Management</h1>
					<div class="flex items-center gap-2 border-l border-muted/20 pl-4">
						<Badge label="Total" type="info" count={totalCount} />
					</div>
				</div>
				<span class="text-sm text-muted">
					Manage your static project folders registered in ZenB Tool
				</span>
			</div>
		</div>
	</div>

	<div class="flex items-center justify-between backdrop-blur-md">
		<div class="group relative max-w-md flex-1">
			<i
				class="absolute top-1/2 left-3 -translate-y-1/2 text-lg text-muted/40 transition-colors icon-[regular--magnifying-glass] group-focus-within:text-muted/80"
			></i>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter projects by name, path or description..."
				class="h-10 w-full rounded-lg border border-transparent bg-transparent pr-4 pl-10 text-sm placeholder-muted/50 transition-all focus:border-white/10 focus:bg-muted/10 focus:ring-0 focus:outline-none"
			/>
		</div>
	</div>
</header>

<StaticList {searchQuery} />
