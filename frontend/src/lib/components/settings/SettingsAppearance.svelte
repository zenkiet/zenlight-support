<script lang="ts">
	import { systemStore, type ThemeScheme } from '$lib/stores/system.svelte';

	const themes: { value: ThemeScheme; label: string; icon: string; desc: string }[] = [
		{
			value: 'light',
			label: 'Light',
			icon: 'icon-[regular--sun]',
			desc: 'Clean and bright interface'
		},
		{
			value: 'dark',
			label: 'Dark',
			icon: 'icon-[regular--moon]',
			desc: 'Easy on the eyes in low light'
		},
		{
			value: 'system',
			label: 'System',
			icon: 'icon-[regular--display]',
			desc: 'Follows your OS preference'
		}
	];
</script>

<div class="flex flex-col gap-6">
	<section class="flex flex-col gap-3">
		<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Color Scheme</h2>
		<div class="grid grid-cols-3 gap-3">
			{#each themes as theme}
				<button
					onclick={() => systemStore.setScheme(theme.value)}
					class={`group relative flex cursor-pointer flex-col items-center gap-3 rounded-xl border p-5 transition-all duration-200 outline-none
						${
							systemStore.scheme === theme.value
								? 'border-blue-500/50 bg-blue-50/50 ring-1 ring-blue-500/20 dark:border-blue-400/40 dark:bg-blue-500/8'
								: 'border-muted/10 bg-surface/50 hover:border-muted/20 hover:bg-muted/4'
						}`}
				>
					{#if systemStore.scheme === theme.value}
						<div
							class="absolute top-2.5 right-2.5 flex h-5 w-5 items-center justify-center rounded-full bg-blue-500"
						>
							<span class="text-xs text-white icon-[regular--check]"></span>
						</div>
					{/if}

					<div
						class={`flex h-12 w-12 items-center justify-center rounded-xl transition-colors
							${
								systemStore.scheme === theme.value
									? 'bg-blue-100 text-blue-600 dark:bg-blue-500/15 dark:text-blue-400'
									: 'bg-muted/8 text-muted/60 group-hover:text-muted'
							}`}
					>
						<span class={`${theme.icon} text-xl`}></span>
					</div>

					<div class="flex flex-col items-center gap-1">
						<span
							class={`text-sm font-medium ${systemStore.scheme === theme.value ? 'text-blue-700 dark:text-blue-300' : ''}`}
						>
							{theme.label}
						</span>
						<span class="text-center text-[11px] leading-tight text-muted/60">{theme.desc}</span>
					</div>
				</button>
			{/each}
		</div>
	</section>
</div>
