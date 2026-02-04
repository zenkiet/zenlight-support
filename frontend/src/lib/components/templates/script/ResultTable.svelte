<script>
	import { scriptStore } from '$lib/stores/script.svelte';

	const status = $derived(scriptStore.result.state);
	const executionTime = $derived(scriptStore.result.data?.executionTime);
	const rowsAffected = $derived(scriptStore.result.data?.rowsAffected);

	const data = $derived(scriptStore.result.data);
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
					{executionTime}
				</span>
				<span class="flex items-center">
					<i class="mr-1.5 size-4 icon-[regular--table-cells-rows]"></i>
					{rowsAffected} rows affected
				</span>
				<span class="flex items-center">
					<i class="mr-1.5 size-4 icon-[regular--table]"></i>
					{data?.columns?.length ?? 0} columns
				</span>
			</div>
		{/if}
	</div>

	{#if status === 'completed'}
		<div class="flex-1 overflow-hidden bg-surface">
			<div class="scrollable max-h-100 overflow-auto">
				<table class="min-w-full text-left text-sm whitespace-nowrap">
					<thead class="sticky top-0 z-10 border-b border-muted/15 bg-surface">
						<tr>
							{#each data?.columns as col}
								<th
									class="border-r border-muted/15 px-5 py-3 font-semibold text-main last:border-r-0"
									>{col}</th
								>
							{/each}
						</tr>
					</thead>
					<tbody class="divide-y divide-border/60">
						{#if data && data.data && data.data.length > 0}
							{#each data.data as row}
								<tr class="group transition-colors hover:bg-primary/5">
									{#each row as cell}
										<td
											class="text-muted-foreground h-10 max-w-75 truncate px-4 py-2.5 align-middle whitespace-nowrap group-hover:text-main"
											title={String(cell)}
										>
											{#if typeof cell === 'number' && cell > 1000}
												<span class="font-mono text-emerald-600 dark:text-emerald-400"
													>${cell.toLocaleString()}</span
												>
											{:else if String(cell).toUpperCase() === 'COMPLETED'}
												<span class="font-medium text-blue-600 dark:text-blue-400">{cell}</span>
											{:else}
												{cell === null ? 'NULL' : String(cell)}
											{/if}
										</td>
									{/each}
								</tr>
							{/each}
						{:else}
							<tr>
								<td
									colspan={data?.columns?.length}
									class="px-4 py-12 text-center text-muted italic"
								>
									No data available
								</td>
							</tr>
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	{/if}

	{#if status === 'error'}
		<div class="scrollable flex-1 overflow-auto bg-surface p-4">
			<div class="w-full rounded-lg border border-rose-500/20 bg-rose-500/10 p-4">
				<p class="text-sm text-rose-600 dark:text-rose-400">
					{scriptStore.result.error}
				</p>
			</div>
		</div>
	{/if}
</div>
