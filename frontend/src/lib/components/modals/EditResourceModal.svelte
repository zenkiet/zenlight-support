<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import type { domain } from '../../../../wailsjs/go/models';

	let {
		isOpen = false,
		config,
		onClose,
		onSave
	} = $props<{
		isOpen: boolean;
		config: domain.ResourceConfig;
		onClose: () => void;
		onSave: (data: Partial<domain.ResourceConfig>) => void;
	}>();

	let name = $state('');
	let description = $state('');
	let path = $state('');
	let serviceName = $state('');
	let installable = $state(false);

	const isService = $derived(config?.type === 'service');

	$effect(() => {
		if (isOpen && config) {
			name = config.name ?? '';
			description = config.description ?? '';
			path = config.path ?? '';
			serviceName = config.serviceName ?? '';
			installable = config.installable ?? false;
		}
	});

	function handleSave() {
		onSave({
			name,
			description,
			path,
			serviceName: isService ? serviceName : undefined,
			installable
		});
		onClose();
	}
</script>

{#if isOpen}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
		<div
			class="w-full max-w-md rounded-2xl border border-muted/10 bg-surface shadow-2xl"
			transition:scale={{ duration: 200, start: 0.95, easing: cubicOut }}
		>
			<!-- Header -->
			<div class="border-b border-muted/10 px-6 py-5">
				<div class="flex flex-col gap-0.5">
					<h3 class="text-base font-semibold text-main">
						Edit {isService ? 'service' : 'project'}
					</h3>
					<p class="text-xs text-muted/60">
						Modify the settings for <span class="font-medium text-blue-500 dark:text-blue-400"
							>{config?.name}</span
						>.
					</p>
				</div>
			</div>

			<!-- Body -->
			<div class="flex flex-col gap-4 px-5 py-4">
				{@render field('Name', 'name')}
				{@render field('Description', 'description')}
				{@render field('Path', 'path')}

				{#if isService}
					{@render field('Service Name', 'serviceName')}
				{/if}

				<!-- Installable toggle -->
				<div class="flex items-center justify-between rounded-lg border border-muted/10 px-4 py-3">
					<div>
						<p class="text-sm font-medium text-main">Installable</p>
						<p class="text-xs text-muted">Allow file uploads for this resource</p>
					</div>
					<button
						type="button"
						class="relative h-6 w-11 shrink-0 cursor-pointer rounded-full transition-colors duration-200 {installable
							? 'bg-primary'
							: 'bg-muted/30'}"
						onclick={() => (installable = !installable)}
						aria-label="Toggle installable"
					>
						<span
							class="absolute top-0.5 left-0.5 size-5 rounded-full bg-white shadow-sm transition-transform duration-200 {installable
								? 'translate-x-5'
								: 'translate-x-0'}"
						></span>
					</button>
				</div>
			</div>

			<!-- Footer -->
			<div class="flex justify-end gap-2 border-t border-muted/10 px-5 py-3">
				<button
					class="btn secondary"
					onclick={onClose}
				>
					Cancel
				</button>
				<button
					class="btn primary"
					onclick={handleSave}
				>
					Save
				</button>
			</div>
		</div>
	</div>
{/if}

{#snippet field(label: string, key: 'name' | 'description' | 'path' | 'serviceName')}
	<div class="flex flex-col gap-1.5">
		<label for={key} class="text-xs font-medium text-muted">{label}</label>
		{#if key === 'name'}
			<input id={key} type="text" bind:value={name} class="input" placeholder="Enter name" />
		{:else if key === 'description'}
			<textarea
				id={key}
				bind:value={description}
				rows="2"
				class="input resize-none"
				placeholder="Enter description"
			></textarea>
		{:else if key === 'path'}
			<input id={key} type="text" bind:value={path} class="input" placeholder="Enter path" />
		{:else if key === 'serviceName'}
			<input
				id={key}
				type="text"
				bind:value={serviceName}
				class="input"
				placeholder="Enter service name"
			/>
		{/if}
	</div>
{/snippet}
