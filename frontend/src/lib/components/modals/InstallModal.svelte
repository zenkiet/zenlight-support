<script lang="ts">
	import { fade, scale, slide } from 'svelte/transition';
	import { cubicOut, elasticOut } from 'svelte/easing';
	import { formatBytes } from '$lib/helpers';

	/** Props */
	let {
		isOpen = false,
		targetName,
		targetType = 'Service',
		subtitlePrefix = 'Upload files for',
		fileHint = 'Supported: .exe, .dll, config files...',
		onClose,
		onInstall
	} = $props<{
		isOpen: boolean;
		targetName: string;
		targetType?: string;
		subtitlePrefix?: string;
		fileHint?: string;
		onClose: () => void;
		onInstall: (files: File[]) => Promise<void>;
	}>();

	/** State */
	let isDragging = $state(false);
	let files = $state<File[]>([]);
	let fileInput = $state<HTMLInputElement>();

	let status = $state<'idle' | 'installing' | 'success' | 'error'>('idle');
	let errMsg = $state<string | null>(null);
	const totalSize = $derived(files.reduce((acc, f) => acc + f.size, 0));

	/** Functions */
	function addFiles(newFiles: FileList | null) {
		if (!newFiles) return;
		const unique = Array.from(newFiles).filter(
			(c) => !files.some((f) => f.name === c.name && f.size === c.size)
		);
		files = [...files, ...unique];
	}

	function removeFile(index: number) {
		files = files.filter((_, i) => i !== index);
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		if (status === 'idle') addFiles(e.dataTransfer?.files || null);
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		addFiles(target.files);
		target.value = '';
	}

	async function submit() {
		status = 'installing';
		errMsg = null;

		try {
			await onInstall(files);
			status = 'success';

			setTimeout(() => {
				close();
			}, 2000);
		} catch (err: unknown) {
			status = 'error';
			console.error('Installation error:', err);
			errMsg = err as string;

			setTimeout(() => {
				close();
			}, 2000);
		}
	}

	function close() {
		isOpen = false;
		onClose();

		setTimeout(() => {
			reset();
		}, 300);
	}

	function reset() {
		files = [];
		isDragging = false;
		status = 'idle';
		errMsg = '';
	}
</script>

{#if isOpen}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
		<div
			class="relative flex min-h-75 w-full max-w-lg flex-col overflow-hidden rounded-xl border border-muted/15 bg-page ring-1 ring-muted/10"
			role="none"
			transition:scale={{ start: 0.96, duration: 250, easing: cubicOut }}
			onclick={(e) => e.stopPropagation()}
		>
			{#if status === 'idle'}
				<div class="flex shrink-0 items-center justify-between border-b border-muted/10 px-5 py-4">
					<div>
						<h3 class="font-semibold">Install {targetType}</h3>
						<p class="text-xs text-muted/60">
							{subtitlePrefix} <span class="text-blue-400">{targetName}</span>
						</p>
					</div>

					<button
						onclick={close}
						title="Close"
						class="cursor-pointer text-muted/60 transition-colors hover:text-muted"
					>
						<i class="icon-[regular--x]"> </i>
					</button>
				</div>

				<div class="scrollable flex max-h-[60vh] flex-col overflow-y-auto p-5">
					<input
						type="file"
						multiple
						bind:this={fileInput}
						class="hidden"
						onchange={handleFileSelect}
					/>
					<button
						class={`relative flex w-full flex-col items-center justify-center rounded-lg border-2 border-dashed transition-all duration-300 ${files.length > 0 ? 'py-6' : 'py-12'} ${isDragging ? 'border-blue-500 bg-blue-500/10' : 'border-muted/10 bg-muted/2 hover:border-muted/20 hover:bg-muted/5'}`}
						ondragover={(e) => {
							e.preventDefault();
							isDragging = true;
						}}
						ondragleave={() => (isDragging = false)}
						ondrop={handleDrop}
						onclick={() => fileInput?.click()}
					>
						<div class="pointer-events-none flex flex-col items-center">
							<i
								class={`transition-colors icon-[regular--cloud-arrow-up] ${files.length > 0 ? 'mb-2 size-6' : 'mb-3 size-8'} ${isDragging ? 'text-blue-400' : 'text-muted/50'}`}
							></i>
							<p class="text-sm text-muted/60">
								<span class="font-medium text-blue-400">Click to upload</span> or drag files
							</p>
							{#if files.length === 0}
								<p class="mt-1 text-xs text-muted/30" transition:fade>
									{fileHint}
								</p>
							{/if}
						</div>
					</button>

					{#if files.length > 0}
						<div class="mt-4 flex flex-col gap-2">
							<div class="flex items-center justify-between px-1">
								<span class="text-xs font-bold tracking-wider text-muted/30 uppercase"
									>Queue ({files.length})</span
								>
								<span class="font-mono text-xs text-muted/40">Total: {formatBytes(totalSize)}</span>
							</div>
							<div class="flex flex-col gap-2">
								{#each files as file, i (file.name + file.size)}
									<div
										class="group flex items-center gap-3 rounded-lg border border-muted/30 bg-muted/10 p-2 pr-3 transition-colors hover:border-muted/40 hover:bg-muted/20"
										transition:slide={{ duration: 200, axis: 'y' }}
									>
										<div
											class="flex h-10 w-10 shrink-0 items-center justify-center rounded bg-page/20"
										>
											<i class="size-5 text-emerald-400 icon-[regular--file-code]"></i>
										</div>
										<div class="flex min-w-0 flex-1 flex-col">
											<span class="truncate text-sm font-medium" title={file.name}>{file.name}</span
											>
											<span class="text-xs text-muted/40">{formatBytes(file.size)}</span>
										</div>

										<button
											onclick={() => removeFile(i)}
											title="Remove File"
											class="cursor-pointer text-muted/30 opacity-0 transition-all hover:text-rose-400"
										>
											<i class="icon-[regular--x]"></i>
										</button>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<div
					class="mt-auto flex shrink-0 items-center justify-end gap-3 border-t border-muted/10 bg-page/20 px-5 py-3"
				>
					<button
						onclick={close}
						class="rounded-lg px-4 py-2 text-xs font-medium text-muted/60 transition-colors hover:bg-muted/10 hover:text-muted"
						>Cancel</button
					>
					<button
						onclick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							submit();
						}}
						disabled={files.length === 0}
						class="flex cursor-pointer items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-xs font-semibold text-white transition-all hover:bg-blue-500 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
					>
						<i class="icon-[regular--download]"></i>
						<span>Install</span>
					</button>
				</div>
			{:else}
				<div
					class="absolute inset-0 flex flex-col items-center justify-center p-8 text-center"
					in:fade={{ duration: 300 }}
				>
					{#if status === 'installing'}
						<div class="relative mb-6">
							<div
								class="h-16 w-16 animate-spin rounded-full border-4 border-muted/15 border-t-blue-500"
							></div>
							<div class="absolute inset-0 flex items-center justify-center">
								<i class="animate-pulse text-2xl text-blue-400 icon-[regular--gear]"></i>
							</div>
						</div>
						<h3 class="text-lg font-medium">Installing {targetType}...</h3>
						<p class="mt-2 text-sm text-muted/40">Please wait while we update your files.</p>
					{:else if status === 'success'}
						<div
							class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-emerald-500/10 ring-1 ring-emerald-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<i class="size-8 text-emerald-400 icon-[regular--circle-check]"></i>
						</div>
						<h3 class="text-xl font-bold" in:slide={{ delay: 100 }}>Success!</h3>
						<p class="mt-2 text-sm text-muted/40" in:slide={{ delay: 200 }}>
							{targetType} has been updated successfully.
						</p>
					{:else if status === 'error'}
						<div
							class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-rose-500/10 ring-1 ring-rose-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<i class="size-8 text-rose-400 icon-[regular--diamond-exclamation]"></i>
						</div>
						<h3 class="text-xl font-bold text-rose-500" in:slide={{ delay: 100 }}>
							Installation Failed
						</h3>
						<p class="mt-2 max-w-[80%] text-sm text-rose-500/60" in:slide={{ delay: 200 }}>
							{errMsg}
						</p>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}
