<script lang="ts">
	let {
		title = 'Coming Soon',
		subtitle = 'This feature is being prepared for the next release.',
		note = '',
	} = $props<{
		title?: string;
		subtitle?: string;
		note?: string;
	}>();
</script>

<section class="flex h-full w-full items-center justify-center">
	<div class="w-full max-w-xl p-8 text-center">
		<div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-surface-highlight">
			<i class="icon-[regular--rocket-launch] size-20 text-primary"></i>
		</div>

    {#if note}
		<span
			class="mx-auto mb-3 inline-flex items-center rounded-full border border-muted/20 bg-muted/10 px-3 py-1 text-[11px] font-semibold tracking-wider text-muted uppercase"
		>
			{note}
		</span>
    {/if}

		<h1 class="text-2xl font-bold tracking-tight text-main">{title}</h1>
		<p class="mt-2 text-sm text-muted">{subtitle}</p>
	</div>
</section>
