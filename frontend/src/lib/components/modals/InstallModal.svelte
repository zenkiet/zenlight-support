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
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button
			class="absolute inset-0 bg-black/50 backdrop-blur-sm dark:bg-black/60"
			onclick={close}
			aria-label="Close"
		></button>

		<div
			class="relative flex w-full max-w-md flex-col overflow-hidden rounded-2xl border border-muted/15 bg-surface shadow-2xl dark:border-muted/10"
			role="none"
			transition:scale={{ start: 0.96, duration: 250, easing: cubicOut }}
		>
			{#if status === 'idle'}
				<!-- Header -->
				<div class="flex items-center justify-between border-b border-muted/10 px-5 py-4">
					<div class="flex items-center gap-3">
						<div>
							<h3 class="text-base font-semibold text-main">Install {targetType}</h3>
							<p class="text-xs text-muted/60">
								{subtitlePrefix}
								<span class="font-medium text-blue-500 dark:text-blue-400">{targetName}</span>
							</p>
						</div>
					</div>
					<button
						onclick={close}
						aria-label="Close"
						class="absolute top-3 right-3 flex h-6 w-6 cursor-pointer items-center justify-center rounded-md text-muted/40 transition-colors hover:bg-muted/8 hover:text-muted"
					>
						<span class="text-sm icon-[regular--x]"></span>
					</button>
				</div>

				<!-- Content -->
				<div class="scrollable flex max-h-[60vh] flex-col overflow-y-auto p-4">
					<input
						type="file"
						multiple
						bind:this={fileInput}
						class="hidden"
						onchange={handleFileSelect}
					/>

					<!-- Drop zone -->
					<button
						class={`flex w-full cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed transition-all duration-200 ${files.length > 0 ? 'py-4' : 'py-10'} ${isDragging ? 'border-blue-400 bg-blue-50 dark:bg-blue-500/8' : 'border-muted/12 hover:border-muted/25 hover:bg-muted/3'}`}
						ondragover={(e) => {
							e.preventDefault();
							isDragging = true;
						}}
						ondragleave={() => (isDragging = false)}
						ondrop={handleDrop}
						onclick={() => fileInput?.click()}
					>
						<span
							class={`transition-colors icon-[regular--cloud-arrow-up] ${files.length > 0 ? 'mb-1.5 text-xl' : 'mb-2 text-2xl'} ${isDragging ? 'text-blue-400' : 'text-muted/35'}`}
						></span>
						<p class="text-sm text-muted/50">
							<span class="font-medium text-blue-500 dark:text-blue-400">Browse</span> or drop files
						</p>
						{#if files.length === 0}
							<p class="mt-1 text-xs text-muted/30" transition:fade>{fileHint}</p>
						{/if}
					</button>

					<!-- File list -->
					{#if files.length > 0}
						<div class="mt-3 flex flex-col gap-1.5">
							<div class="flex items-center justify-between px-0.5">
								<span class="text-[10px] font-bold tracking-widest text-muted/35 uppercase">
									Files ({files.length})
								</span>
								<span class="font-mono text-[10px] text-muted/30">{formatBytes(totalSize)}</span>
							</div>
							{#each files as file, i (file.name + file.size)}
								<div
									class="group flex items-center gap-2.5 rounded-lg border border-muted/10 p-2 transition-colors hover:bg-muted/4"
									transition:slide={{ duration: 150, axis: 'y' }}
								>
									<div
										class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-emerald-50 dark:bg-emerald-500/8"
									>
										<i class=" text-emerald-500 icon-[regular--file-code] dark:text-emerald-400"
										></i>
									</div>
									<div class="flex min-w-0 flex-1 flex-col">
										<span class="truncate text-sm font-medium text-main" title={file.name}
											>{file.name}</span
										>
										<span class="text-[11px] text-muted/40">{formatBytes(file.size)}</span>
									</div>
									<button
										onclick={() => removeFile(i)}
										aria-label="Remove file"
										class="flex h-6 w-6 cursor-pointer items-center justify-center rounded text-muted/25 opacity-0 transition-all group-hover:opacity-100 hover:text-rose-500"
									>
										<span class="text-xs icon-[regular--x]"></span>
									</button>
								</div>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Actions -->
				<div class="flex items-center justify-end gap-2 border-t border-muted/10 px-4 py-3">
					<button
						onclick={close}
						class="cursor-pointer rounded-lg px-3.5 py-2 text-xs font-medium text-muted/50 transition-colors hover:bg-muted/6 hover:text-muted"
					>
						Cancel
					</button>
					<button
						onclick={submit}
						disabled={files.length === 0}
						class="flex cursor-pointer items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-xs font-semibold text-white transition-all hover:bg-blue-500 active:scale-[0.97] disabled:cursor-not-allowed disabled:opacity-40"
					>
						<span class="text-xs icon-[regular--download]"></span>
						Install
					</button>
				</div>
			{:else}
				<!-- Status states -->
				<div
					class="flex min-h-70 flex-col items-center justify-center p-8 text-center"
					in:fade={{ duration: 200 }}
				>
					{#if status === 'installing'}
						<div class="relative mb-5">
							<div
								class="h-14 w-14 animate-spin rounded-full border-[3px] border-muted/10 border-t-blue-500"
							></div>
							<div class="absolute inset-0 flex items-center justify-center">
								<span class="animate-pulse text-lg text-blue-400 icon-[regular--gear]"></span>
							</div>
						</div>
						<h3 class="text-base font-semibold text-main">Installing...</h3>
						<p class="mt-1.5 text-xs text-muted/45">Updating files for {targetName}</p>
					{:else if status === 'success'}
						<div
							class="mb-5 flex h-16 w-16 items-center justify-center rounded-full bg-emerald-100 ring-1 ring-emerald-200 dark:bg-emerald-500/10 dark:ring-emerald-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<span
								class="size-7 text-emerald-600 icon-[regular--circle-check] dark:text-emerald-400"
							></span>
						</div>
						<h3 class="text-lg font-bold text-main" in:slide={{ delay: 100 }}>Installed!</h3>
						<p class="mt-1.5 text-xs text-muted/45" in:slide={{ delay: 200 }}>
							{targetType} updated successfully
						</p>
					{:else if status === 'error'}
						<div
							class="mb-5 flex h-16 w-16 items-center justify-center rounded-full bg-rose-100 ring-1 ring-rose-200 dark:bg-rose-500/10 dark:ring-rose-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<span
								class="size-7 text-rose-600 icon-[regular--diamond-exclamation] dark:text-rose-400"
							></span>
						</div>
						<h3
							class="text-lg font-bold text-rose-600 dark:text-rose-400"
							in:slide={{ delay: 100 }}
						>
							Failed
						</h3>
						<p class="mt-1.5 max-w-[80%] text-xs text-muted/50" in:slide={{ delay: 200 }}>
							{errMsg}
						</p>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}
