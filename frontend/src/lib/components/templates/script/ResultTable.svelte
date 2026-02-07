<script lang="ts">
	import { scriptStore } from '$lib/stores/script.svelte';
	import { VList } from 'virtua/svelte';

	const { state: status, data, error } = $derived(scriptStore.result);
	const cols = $derived(data?.columns ?? []);
	const rows = $derived(data?.data ?? []);
	const grid = $derived(`repeat(${cols.length}, minmax(150px, 1fr))`);
</script>

<div class="flex flex-1 flex-col overflow-hidden rounded-xl bg-surface">
	<div class="flex h-11 items-center justify-between border-b border-muted/15 px-4">
		<div class="flex space-x-6 text-xs font-semibold text-muted">
			{#if status === 'running'}
				<span class="flex items-center">
					<i class="mr-1.5 size-4 animate-spin icon-[regular--spinner]"></i>
					Executing query...
				</span>
			{:else if status === 'completed'}
				<span class="flex items-center text-emerald-600 dark:text-emerald-400">
					<i class="mr-1.5 size-4 icon-[regular--check]"></i>
					Query executed successfully
				</span>
			{:else if status === 'error'}
				<span class="flex items-center text-rose-600 dark:text-rose-400">
					<i class="mr-1.5 size-4 icon-[regular--circle-exclamation]"></i>
					Error executing query
				</span>
			{/if}
		</div>
		{#if status === 'completed'}
			<div class="flex space-x-6 text-xs font-medium text-muted">
				<span class="flex items-center">
					<i class="mr-1.5 size-4 icon-[regular--clock]"></i>
					{data?.executionTime}
				</span>
				<span class="flex items-center">
					<i class="mr-1.5 size-4 icon-[regular--table-cells-rows]"></i>
					{data?.rowsAffected} rows affected
				</span>
				<span class="flex items-center">
					<i class="mr-1.5 size-4 icon-[regular--record-vinyl]"></i>
					{data?.data?.length ?? 0} record{(data?.data?.length ?? 0) !== 1 ? 's' : ''}
				</span>
			</div>
		{/if}
	</div>

	{#if status === 'completed'}
		<div
			class="flex flex-1 flex-col scrollable overflow-x-auto overflow-y-hidden bg-surface text-sm whitespace-nowrap"
		>
			{#if rows.length > 0}
				<div class="flex flex-1 flex-col" style="min-width: {cols.length * 150}px;">
					<div
						class="shrink-0 border-b border-muted/15 bg-surface"
						style="display:grid;grid-template-columns:{grid};"
					>
						{#each cols as col}
							<div
								class="border-r border-muted/15 px-5 py-3 font-semibold text-main last:border-r-0"
							>
								{col}
							</div>
						{/each}
					</div>

					<VList data={rows} class="h-[40vh]!" getKey={(_, i) => i}>
						{#snippet children(row)}
							<div
								class="group border-b border-border/60 transition-colors hover:bg-primary/5"
								style="display:grid;grid-template-columns:{grid};"
							>
								{#each row as cell}
									<div class="h-10 max-w-75 truncate px-4 py-2.5 text-muted group-hover:text-main">
										{cell === null ? 'NULL' : String(cell)}
									</div>
								{/each}
							</div>
						{/snippet}
					</VList>
				</div>
			{:else}
				<div class="px-4 py-12 text-center text-muted italic">No data available</div>
			{/if}
		</div>
	{/if}

	{#if status === 'error'}
		<div class="scrollable flex-1 overflow-auto bg-surface p-4">
			<div class="w-full rounded-lg border border-rose-500/20 bg-rose-500/10 p-4">
				<p class="text-sm text-rose-600 dark:text-rose-400">{error}</p>
			</div>
		</div>
	{/if}
</div>
