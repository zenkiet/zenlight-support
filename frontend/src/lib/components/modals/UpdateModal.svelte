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
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
    <div
      class="glass-panel relative flex w-full max-w-lg flex-col overflow-hidden rounded-xl border border-white/10 bg-[#1a1b26]/95 shadow-2xl ring-1 ring-white/5 transition-all duration-300"
      style="min-height: 320px;"
      transition:scale={{ start: 0.96, duration: 250, easing: cubicOut }}
    >
      {#if status === 'idle'}
        <div in:fade={{ duration: 200 }}>
          <div class="flex items-start justify-between border-b border-white/5 px-6 py-5">
            <div class="flex flex-col gap-1">
              <h3 class="text-base font-semibold text-white">Update Available</h3>
              <p class="text-xs text-white/40">A new version is ready to install.</p>
            </div>
            <div
              class="rounded border border-white/10 bg-white/5 px-2.5 py-1 font-mono text-xs font-medium text-white/80 shadow-inner"
            >
              {systemStore.info.latestVersion}
            </div>
          </div>

          <div class="px-6 py-5">
            <label
              for="changelog"
              class="mb-2 block text-[10px] font-bold tracking-wider text-white/30 uppercase"
            >
              Release Notes
            </label>
            <div
              class="scrollable h-48 overflow-y-auto rounded-lg border border-white/5 bg-black/20 p-3 shadow-inner"
            >
              <Markdown content={systemStore.info.releaseNotes} />
            </div>
          </div>

          <div
            class="flex items-center justify-end gap-3 rounded-b-lg border-t border-white/5 bg-black/20 px-6 py-4"
          >
            <button
              onclick={close}
              class="rounded-lg px-4 py-2 text-xs font-medium text-white/40 transition-colors hover:bg-white/5 hover:text-white"
            >
              Remind me later
            </button>

            <button
              onclick={handleUpdate}
              class="flex cursor-pointer items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-xs font-semibold text-white transition-all hover:bg-blue-500 active:scale-95"
            >
              <span class="icon-[regular--download]"></span>
              <span>Download & Install</span>
            </button>
          </div>
        </div>
      {:else}
        <div class="absolute inset-0 flex flex-col items-center justify-center p-8 text-center">
          <!-- Updating -->
          {#if status === 'updating'}
            <div class="relative mb-6" in:scale={{ duration: 300 }}>
              <div
                class="h-16 w-16 animate-spin rounded-full border-4 border-white/10 border-t-blue-500"
              ></div>
              <div class="absolute inset-0 flex items-center justify-center">
                <i class="size-6 animate-pulse text-blue-400 icon-[regular--rotate]"></i>
              </div>
            </div>
            <h3 class="text-lg font-medium text-white" in:slide={{ delay: 100, axis: 'y' }}>
              Updating Application...
            </h3>
            <p class="mt-2 text-sm text-white/40" in:slide={{ delay: 200, axis: 'y' }}>
              The app will restart automatically.
            </p>

            <!-- Success -->
          {:else if status === 'success'}
            <div
              class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-emerald-500/10 ring-1 ring-emerald-500/20"
              in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
            >
              <span class="size-8 text-emerald-400 icon-[regular--circle-check]"></span>
            </div>
            <h3 class="text-xl font-bold text-white" in:slide={{ delay: 100 }}>Update Complete!</h3>
            <p class="mt-2 text-sm text-white/40" in:slide={{ delay: 200 }}>
              Restarting application...
            </p>

            <!-- Error -->
          {:else if status === 'error'}
            <div
              class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-rose-500/10 ring-1 ring-rose-500/20"
              in:scale={{ start: 0.5, duration: 400, easing: elasticOut }}
            >
              <i class="size-8 text-rose-400 icon-[regular--circle-exclamation]"></i>
            </div>
            <h3 class="text-xl font-bold text-rose-500" in:slide={{ delay: 100 }}>Update Failed</h3>
            <p class="mt-2 max-w-[80%] text-sm text-white/60" in:slide={{ delay: 200 }}>
              {errorMsg}
            </p>

            <button
              onclick={() => (status = 'idle')}
              class="mt-6 rounded-lg bg-white/5 px-4 py-2 text-xs font-medium text-white hover:bg-white/10"
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

