<script lang="ts">
	import { formatRelativeTime, readFileAsBytes } from '$lib/helpers';
	import type { Resource } from '$lib/stores/resources.svelte';
	import InstallModal from './modals/InstallModal.svelte';
	import EditResourceModal from './modals/EditResourceModal.svelte';

	let { service } = $props<{ service: Resource }>();

	let showInstallModal = $state(false);
	let showEditModal = $state(false);

	async function handleInstall(files: File[]) {
		if (!files || files.length === 0) return;

		const payload = files.map(async (file) => {
			const data = await readFileAsBytes(file);
			return {
				name: file.name,
				data,
				extension: file.name.split('.').pop() || ''
			};
		});
		await service.install(await Promise.all(payload));
	}

	async function handleSave(newConfig: typeof service.config) {
		await service.save(newConfig);
	}
</script>

<div
	class="group hover:border-surface-highlight hover:bg-surface-highlight relative flex flex-col gap-4 rounded-xl border border-surface bg-surface"
>
	<div class="flex items-start justify-between gap-3 px-4 pt-4">
		<div class="min-w-0 flex-1">
			<div class="mb-1 flex items-center gap-2">
				<span class="font-mono text-xs font-semibold text-muted"
					>{service.metrics.pid ? `#${service.metrics.pid}` : '--'}</span
				>
			</div>

			<div class="flex flex-col">
				<div class="flex items-center">
					<span
						class="mr-2 h-2 w-2 shrink-0 rounded-full {service.isRunning
							? 'bg-emerald-500'
							: service.loading
								? 'bg-muted/50'
								: 'bg-rose-500'}"
					></span>
					<h3 class="truncate font-semibold" title={service.config?.name}>
						<button
							class="cursor-pointer transition-colors duration-200 hover:text-accent hover:underline"
							onclick={() => service.openExplorer()}
						>
							{service.config?.name}
						</button>
					</h3>
				</div>

				<p
					class="line-clamp-2 min-h-[2.5em] text-xs leading-relaxed text-muted/80"
					title={service.config?.description}
				>
					{service.config?.description}
				</p>
			</div>
		</div>

		<div class="flex gap-2">
			<button
				title="edit service"
				onclick={() => (showEditModal = true)}
				class="cursor-pointer text-muted/80 transition-colors hover:text-main"
			>
				<span class="icon-[regular--pen-to-square]"></span>
			</button>

			{#if service.config?.installable}
				<button
					title="install service"
					onclick={() => (showInstallModal = true)}
					class="cursor-pointer text-muted/80 transition-colors hover:text-main"
				>
					<span class="icon-[regular--cloud-arrow-up]"></span>
				</button>
			{/if}

			<button
				onclick={service.isRunning ? () => service.stop() : () => service.start()}
				disabled={service.loading}
				class="cursor-pointer text-muted/80 transition-colors hover:text-main"
			>
				{#if service.loading}
					<span class="animate-spin icon-[regular--spinner]"></span>
				{:else if service.isRunning}
					<span class="icon-[regular--pause]"></span>
				{:else}
					<span class="icon-[regular--play]"></span>
				{/if}
			</button>
		</div>
	</div>

	<div
		class="grid grid-cols-3 overflow-hidden rounded-b-lg border-t border-muted/10 text-center text-xs font-medium"
	>
		{@render metric(
			'CPU',
			service.metrics.cpu ? `${service.metrics.cpu.toFixed(1)}%` : '--',
			'text-blue-600 dark:text-blue-300'
		)}
		{@render metric(
			'RAM',
			service.metrics.mem ? `${(service.metrics.mem / 1024 / 1024).toFixed(0)} MB` : '--',
			'text-purple-600 dark:text-purple-300'
		)}
		{@render metric(
			'UPTIME',
			service.metrics.createTime ? formatRelativeTime(service.metrics.createTime) : '--',
			'text-emerald-600 dark:text-emerald-300'
		)}
	</div>
</div>

<InstallModal
	isOpen={showInstallModal}
	targetName={service.config?.name}
	onClose={() => (showInstallModal = false)}
	onInstall={handleInstall}
/>

<EditResourceModal
	isOpen={showEditModal}
	config={service.config!}
	onClose={() => (showEditModal = false)}
	onSave={handleSave}
/>

{#snippet metric(label: string, value: string, colorClass: string)}
	<div
		class="flex flex-col items-center justify-center gap-1 border-r border-muted/20 bg-page/10 py-2 transition-colors last:border-0 hover:bg-muted/10"
	>
		<span class="text-[9px] font-bold tracking-wider text-muted uppercase">{label}</span>
		<span class={`font-mono text-xs font-medium ${value === '--' ? 'text-muted/80' : colorClass}`}>
			{value}
		</span>
	</div>
{/snippet}
