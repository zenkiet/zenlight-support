<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	// Props
	let {
		content,
		children,
		position = 'top',
		delay = 200
	} = $props<{
		content: string;
		children: import('svelte').Snippet;
		position?: 'top' | 'bottom';
		delay?: number;
	}>();

	// State
	let isHovered = $state(false);
	let timeoutId: ReturnType<typeof setTimeout>;

	function handleMouseEnter() {
		clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			isHovered = true;
		}, delay);
	}

	function handleMouseLeave() {
		clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			isHovered = false;
		}, 100);
	}
</script>

<div
	class="relative inline-grid"
	role="tooltip"
	onmouseenter={handleMouseEnter}
	onmouseleave={handleMouseLeave}
>
	{@render children()}

	{#if isHovered && content}
		<div
			class="pointer-events-none absolute left-1/2 z-50 w-max -translate-x-1/2 px-3 py-1.5"
			style={position === 'top'
				? 'bottom: 100%; padding-bottom: 8px;'
				: 'top: 100%; padding-top: 8px;'}
			in:fly={{ y: position === 'top' ? 4 : -4, duration: 200, easing: cubicOut }}
			out:fade={{ duration: 150 }}
		>
			<div
				class="relative flex items-center justify-center rounded-lg border border-muted/10 bg-surface font-mono text-xs font-medium"
			>
				<span class="px-2 py-1">{content}</span>

				<div
					class="absolute left-1/2 h-2 w-2 -translate-x-1/2 rotate-45 border-muted/10 bg-surface"
					class:border-r={position === 'top'}
					class:border-b={position === 'top'}
					class:border-l={position === 'bottom'}
					class:border-t={position === 'bottom'}
					class:bottom-[-5px]={position === 'top'}
					class:top-[-5px]={position === 'bottom'}
				></div>
			</div>
		</div>
	{/if}
</div>
