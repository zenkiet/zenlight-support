<script lang="ts">
	import { theme } from '$lib/stores/theme.svelte';
</script>

<button
	class="flex h-10 w-10 items-center justify-center rounded-full bg-surface text-primary transition-all duration-300 hover:text-primary"
	onclick={() => theme.toggle()}
	aria-label="Toggle Theme"
>
	{#if theme.isDark}
		<i class="text-amber-300 icon-[regular--sun-bright]"></i>
	{:else}
		<i class="text-slate-900/80 icon-[regular--moon-cloud]"></i>
	{/if}
</button>
