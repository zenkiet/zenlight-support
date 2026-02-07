<script lang="ts">
	import { systemStore } from '$lib/stores/system.svelte';
	import { Quit } from '../../../../wailsjs/runtime/runtime';
	import { fade, scale, slide } from 'svelte/transition';
	import { cubicOut, elasticOut } from 'svelte/easing';
	import Markdown from '../plugins/Markdown.svelte';

	type UpdateStatus = 'idle' | 'updating' | 'success' | 'error';

	let isOpen = $state(false);
	let status = $state<UpdateStatus>('idle');
	let errorMsg = $state<string | null>(null);

	$effect(() => {
		if (systemStore.info.available) {
			isOpen = true;
			status = 'idle';
		}
	});

	async function handleUpdate() {
		status = 'updating';
		errorMsg = '';

		const startTime = Date.now();
		try {
			await systemStore.update();

			const elapsed = Date.now() - startTime;
			if (elapsed < 800) {
				await new Promise((r) => setTimeout(r, 800 - elapsed));
			}

			status = 'success';

			setTimeout(() => {
				Quit();
			}, 1500);
		} catch (err) {
			status = 'error';
			errorMsg = err as string;
		}
	}

	function close() {
		isOpen = false;
		status = 'idle';
		errorMsg = '';
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm dark:bg-black/60"
	>
		<!-- Backdrop click to close -->
		<button class="absolute inset-0" onclick={close} aria-label="Close"></button>

		<!-- Modal -->
		<div
			class="relative flex w-full max-w-lg flex-col overflow-hidden rounded-2xl border border-muted/15 bg-surface shadow-2xl dark:border-muted/10"
			style="min-height: 320px;"
			transition:scale={{ start: 0.96, duration: 250, easing: cubicOut }}
		>
			{#if status === 'idle'}
				<div in:fade={{ duration: 200 }}>
					<!-- Header -->
					<div class="flex items-start justify-between border-b border-muted/10 px-6 py-5">
						<div class="flex items-center gap-3">
							<div class="flex flex-col gap-0.5">
								<h3 class="text-base font-semibold text-main">Update Available</h3>
								<p class="text-xs text-muted/60">A new version is ready to install</p>
							</div>
						</div>
						<div
							class="rounded-full border border-muted/15 bg-muted/5 px-3 py-1 font-mono text-xs font-medium text-muted dark:border-muted/10"
						>
							{systemStore.info.latestVersion}
						</div>
					</div>

					<!-- Release Notes -->
					<div class="px-6 py-4">
						<p class="mb-2 text-[10px] font-bold tracking-widest text-muted/40 uppercase">
							Release Notes
						</p>
						<div
							class="scrollable h-48 overflow-y-auto rounded-xl border border-muted/10 bg-muted/4 p-4 dark:bg-muted/5"
						>
							<Markdown content={systemStore.info.releaseNotes} />
						</div>
					</div>

					<!-- Actions -->
					<div class="flex items-center justify-end gap-2.5 border-t border-muted/10 px-6 py-4">
						<button
							onclick={close}
							class="cursor-pointer rounded-lg px-4 py-2 text-xs font-medium text-muted/50 transition-colors hover:bg-muted/6 hover:text-muted"
						>
							Remind me later
						</button>

						<button
							onclick={handleUpdate}
							class="flex cursor-pointer items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-xs font-semibold text-white transition-all hover:bg-blue-500 active:scale-[0.97]"
						>
							<span class="icon-[regular--download]"></span>
							<span>Download & Install</span>
						</button>
					</div>
				</div>
			{:else}
				<div class="absolute inset-0 flex flex-col items-center justify-center p-8 text-center">
					{#if status === 'updating'}
						<div class="relative mb-6" in:scale={{ duration: 300 }}>
							<div
								class="h-16 w-16 animate-spin rounded-full border-[3px] border-muted/10 border-t-blue-500"
							></div>
							<div class="absolute inset-0 flex items-center justify-center">
								<span class="size-5 animate-pulse text-blue-400 icon-[regular--rotate]"></span>
							</div>
						</div>
						<h3 class="text-lg font-medium text-main" in:slide={{ delay: 100, axis: 'y' }}>
							Updating Application...
						</h3>
						<p class="mt-2 text-sm text-muted/50" in:slide={{ delay: 200, axis: 'y' }}>
							The app will restart automatically.
						</p>
					{:else if status === 'success'}
						<div
							class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-emerald-100 ring-1 ring-emerald-200 dark:bg-emerald-500/10 dark:ring-emerald-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<span
								class="size-8 text-emerald-600 icon-[regular--circle-check] dark:text-emerald-400"
							></span>
						</div>
						<h3 class="text-xl font-bold text-main" in:slide={{ delay: 100 }}>Update Complete!</h3>
						<p class="mt-2 text-sm text-muted/50" in:slide={{ delay: 200 }}>
							Restarting application...
						</p>
					{:else if status === 'error'}
						<div
							class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-rose-100 ring-1 ring-rose-200 dark:bg-rose-500/10 dark:ring-rose-500/20"
							in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
						>
							<span
								class="size-8 text-rose-600 icon-[regular--circle-exclamation] dark:text-rose-400"
							></span>
						</div>
						<h3
							class="text-xl font-bold text-rose-600 dark:text-rose-400"
							in:slide={{ delay: 100 }}
						>
							Update Failed
						</h3>
						<p class="mt-2 max-w-[80%] text-sm text-muted/60" in:slide={{ delay: 200 }}>
							{errorMsg}
						</p>

						<button
							onclick={() => (status = 'idle')}
							class="mt-6 cursor-pointer rounded-lg border border-muted/10 bg-muted/5 px-4 py-2 text-xs font-medium text-muted transition-colors hover:bg-muted/10"
							in:fade={{ delay: 500 }}
						>
							Try Again
						</button>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}
