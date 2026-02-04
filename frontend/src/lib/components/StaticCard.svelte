<script lang="ts">
	import { formatDate, formatFullDateTime, formatRelativeTime, readFileAsBytes } from '$lib/helpers';
	import type { Resource } from '$lib/stores/resources.svelte';
	import InstallModal from './modals/InstallModal.svelte';

	let { project } = $props<{ project: Resource }>();

	let showInstallModal = $state(false);

	const dateModified = $derived(
		project.metrics.lastModified ? formatFullDateTime(project.metrics.lastModified) : '--'
	)
	const totalSize = $derived(
		project.metrics.totalSize ? `${(project.metrics.totalSize / 1024 / 1024).toFixed(1)} MB` : '--'
	);

	async function handleInstall(files: File[]) {
		if (!files || files.length === 0) return;

		const payload = await Promise.all(
			files.map(async (file) => ({
				name: file.name,
				data: await readFileAsBytes(file),
				size: file.size
			}))
		);

		await project.install(payload);
	}
</script>

<div
	class="group hover:border-surface-highlight hover:bg-surface-highlight relative flex flex-col gap-4 rounded-xl border border-surface bg-surface"
>
	<div class="flex items-start justify-between gap-3 px-4 pt-4">
		<div class="min-w-0 flex-1">
			<div class="flex flex-col">
				<h3 class="truncate font-semibold" title={project.config?.name}>
					<button
						class="transitio n-colors cursor-pointer duration-200
					hover:text-accent hover:underline"
						onclick={() => project.openExplorer()}
					>
						{project.config?.name}
					</button>
				</h3>
				<p
					class="line-clamp-2 min-h-[2.5em] text-xs leading-relaxed text-muted/80"
					title={project.config?.description}
				>
					{project.config?.description}
				</p>
			</div>
		</div>

		<div class="flex gap-2">
			<button
				title="Install project"
				onclick={() => (showInstallModal = true)}
				class="cursor-pointer text-muted/80 transition-colors hover:text-main"
			>
				<span class="icon-[regular--folder-arrow-up]"></span>
			</button>
		</div>
	</div>

	<div
		class="grid grid-cols-2 overflow-hidden rounded-b-lg border-t border-muted/10 text-center text-xs font-medium"
	>
		{@render metric('Date Modified', dateModified, 'text-blue-600 dark:text-blue-300')}
		{@render metric('Total Size', totalSize, 'text-emerald-600 dark:text-emerald-300')}
	</div>
</div>

<InstallModal
	isOpen={showInstallModal}
	targetName={project.config?.name}
	targetType="Project"
	subtitlePrefix="Replace folder contents for"
	fileHint="Supported: .zip, .html, asset files..."
	onClose={() => (showInstallModal = false)}
	onInstall={handleInstall}
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
