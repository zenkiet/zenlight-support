<script>
	import ResultTable from '$lib/components/templates/script/ResultTable.svelte';
	import SqlEditor from '$lib/components/templates/script/SqlEditor.svelte';
	import { scriptStore } from '$lib/stores/script.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { format } from 'sql-formatter';

	let data = $state('');

	const server = $derived(scriptStore.config?.server || 'Not Connected');
	const database = $derived(scriptStore.config?.database || 'Not Connected');

	onMount(() => {
		scriptStore.init();
	});

	onDestroy(() => {
		console.log('Resetting script store on unmount');
		scriptStore.reset();
	});

	function formatScript() {
		if (!data || !data.trim()) return;
		try {
			data = format(data, {
				language: 'transactsql',
				keywordCase: 'upper',
				tabWidth: 2,
				linesBetweenQueries: 2
			});
		} catch (error) {
			console.error('Format error:', error);
		}
	}
</script>

<header class="mb-6 flex shrink-0 flex-col gap-5 pt-2">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-4">
			<div>
				<div class="flex items-center gap-4">
					<h1 class="text-2xl font-bold tracking-tight">Run SQL Script</h1>
				</div>
				<span class="text-sm text-muted">
					Execute SQL scripts against your connected databases (only use SQL Server)
				</span>
			</div>
		</div>
	</div>
</header>

<div
	class="flex h-fit w-full overflow-hidden rounded-lg bg-page font-sans text-main selection:bg-primary/30"
>
	<main class="relative flex min-w-0 flex-1 flex-col">
		<header
			class="z-10 flex h-12 items-center justify-between rounded-t-lg border-b border-muted/10 bg-surface px-4 transition-colors duration-300"
		>
			<div class="flex w-full items-center space-x-6">
				<div class="flex items-center text-sm">
					<div class="mr-2 flex items-center gap-1.5 text-muted">
						<i class="size-4 text-muted icon-[regular--database]"></i>
						<span>Server: <span class="font-medium text-main"> {server} </span></span>
					</div>
				</div>

				<div class="flex items-center text-sm">
					<span class="mr-2 text-muted">Database:</span>
					<span class="flex items-center font-medium text-main">
						{database}
					</span>
				</div>
			</div>
		</header>

		<div class="flex h-10 items-center space-x-3 bg-surface px-4 transition-colors duration-300">
			<button
				class="cursor-pointer rounded text-muted hover:text-primary disabled:cursor-not-allowed disabled:text-muted/30"
				aria-label="Run Script"
				onclick={() => scriptStore.execute(data)}
				disabled={scriptStore.loading || !data.trim()}
			>
				{#if scriptStore.loading}
					<i class="animate-spin icon-[regular--spinner]"></i>
				{:else}
					<i class="icon-[regular--play]"></i>
				{/if}
			</button>

			<div class="mx-2 h-6 w-px bg-border"></div>

			<button
				class="cursor-pointer text-muted hover:text-main disabled:cursor-not-allowed disabled:text-muted/30"
				aria-label="Format Script"
				onclick={formatScript}
				disabled={scriptStore.loading || !data.trim()}
			>
				<i class="icon-[regular--broom]"></i>
			</button>

			<div class="grow"></div>
		</div>

		<div class="flex flex-1 flex-col overflow-hidden">
			<div class="flex h-100 flex-col">
				<SqlEditor value={data} onchange={(newValue) => (data = newValue)} />
			</div>

			<div class="flex min-h-0 flex-1 flex-col rounded-b-lg bg-surface">
				<ResultTable />
			</div>
		</div>
	</main>
</div>
