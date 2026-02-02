<script>
	import Tooltip from './Tooltip.svelte';
</script>

<Tooltip content="Comming soon">
	<button
		class="group relative flex flex-col items-center justify-center rounded-xl border-2 border-dashed border-border p-4 transition-all duration-300 cursor-pointer hover:border-primary hover:bg-blue-50 dark:hover:bg-blue-500/5 disabled:cursor-not-allowed disabled:opacity-30"
		disabled
	>
		<div
			class="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-surface-highlight transition-colors group-hover:bg-blue-100 dark:group-hover:bg-blue-500/20"
		>
			<i class="size-7 icon-[regular--plus] text-muted transition-colors group-hover:text-primary"></i>
		</div>
		<h3 class="mb-1 text-sm font-bold text-main">Add Project</h3>
		<p class="text-center text-xs text-muted">Register a static project folder</p>
	</button>
</Tooltip>
