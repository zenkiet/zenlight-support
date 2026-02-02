<script lang="ts">
	import ServiceList from '$lib/components/ServiceList.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { serviceStore } from '$lib/stores/services.svelte';

	let searchQuery = $state('');
</script>

<header class="flex shrink-0 flex-col gap-5 pt-2 mb-6">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-4">
			<h1 class="text-2xl font-bold tracking-tight">Services Management</h1>
			<div class="flex items-center gap-2 border-l border-muted/20 pl-4">
				<Badge label="Running" type="success" count={serviceStore.runningCount} />
				<Badge label="Stopped" type="danger" count={serviceStore.stoppedCount} />
			</div>
		</div>
	</div>

	<div class="flex items-center justify-between backdrop-blur-md">
		<div class="group relative max-w-md flex-1">
			<i class="absolute top-1/2 left-3 -translate-y-1/2 text-lg text-muted/40 transition-colors icon-[regular--magnifying-glass] group-focus-within:text-muted/80"></i>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter services by name, PID or description..."
				class="h-10 w-full rounded-lg border border-transparent bg-transparent pr-4 pl-10 text-sm placeholder-muted/50 transition-all focus:border-white/10 focus:bg-muted/10 focus:ring-0 focus:outline-none"
			/>
		</div>
	</div>
</header>

<ServiceList {searchQuery} />