<script lang="ts">
	import { systemStore } from "$lib/stores/system.svelte";

</script>

<button
	class="flex h-10 w-10 items-center justify-center rounded-full bg-surface text-primary transition-all duration-300 hover:text-primary"
	onclick={() => systemStore.toggleTheme()}
	aria-label="Toggle Theme"
>
	{#if systemStore.scheme === 'dark'}
		<i class="text-amber-300 icon-[regular--sun-bright]"></i>
	{:else if systemStore.scheme === 'light'}
		<i class="text-slate-900/80 icon-[regular--moon-cloud]"></i>
	{:else}
		<i class="icon-[regular--desktop]"></i>
	{/if}
</button>
