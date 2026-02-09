<script lang="ts">
  import { bridge } from '$lib/services/wails-bridge';
  import { resourceStore } from '$lib/stores/resources.svelte';

  class Task<T> {
    status = $state<'idle' | 'loading' | 'success' | 'error'>('idle');
    message = $state('');
    data = $state<T | null>(null);

    async execute(fn: () => Promise<T>, successMsg?: (res: T) => string) {
      this.status = 'loading';
      this.message = '';
      try {
        const result = await fn();
        this.status = 'success';
        this.data = result;
        if (successMsg) this.message = successMsg(result);
      } catch (e: any) {
        this.status = 'error';
        this.message = e.message || 'Operation failed';
      }
    }
  }

  const backupTask = new Task<string>();
  const restoreTask = new Task<void>();

  const handleBackup = () => backupTask.execute(
    async () => {
      const res = await bridge.config.export();
      if (!res.success) throw new Error(res.error);
      return res.data;
    },
    (data) => `Saved to ${data}`
  );

  const handleRestore = () => restoreTask.execute(
    async () => {
      const [res] = await Promise.all([bridge.config.import(), delay(500)]);
      if (!res.success) throw new Error(res.error);

      for (let i = 5; i > 0; i--) {
        restoreTask.message = `Restarting in ${i}...`;
        await delay(1000);
      }
      bridge.quit();
    },
    () => 'Config restored. Restart recommended.'
  );

  // 3. Derived State cho UI (Clean & Centralized)
  const getStatusConfig = (status: string) => ({
    color: status === 'error' ? 'text-rose-600 dark:text-rose-400' : 'text-emerald-600 dark:text-emerald-400',
    icon: status === 'error' ? 'icon-[regular--circle-exclamation]' : 'icon-[regular--circle-check]'
  });


  function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
</script>

{#snippet statusMessage(task: Task<unknown>)}
  {#if task.message}
    {@const config = getStatusConfig(task.status)}
    <div class="flex items-center gap-2 px-1 text-xs {config.color}">
      <i class="size-3.5 shrink-0 {config.icon}"></i>
      {task.message}
    </div>
  {/if}
{/snippet}

{#snippet actionButton(task: Task<unknown>, icon: string, label: string, onclick: () => void)}
  <button
    {onclick}
    disabled={task.status === 'loading'}
    class="flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border border-muted/10 bg-surface/50 px-4 py-3 text-sm font-medium transition-colors hover:bg-muted/8 disabled:opacity-50"
  >
    {#if task.status === 'loading'}
      <i class="size-4 animate-spin icon-[regular--spinner]"></i>
    {:else}
      <i class="size-4 text-muted/60 {icon}"></i>
    {/if}
    {label}
  </button>
{/snippet}

<div class="flex flex-col gap-6">
  <section class="flex flex-col gap-3">
    <h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Current Configuration</h2>
    <div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
      <div class="flex items-center gap-3 p-4">
        <i class="size-5 text-muted/50 icon-[regular--database]"></i>
        <span class="text-sm">
          <strong>{resourceStore.services.length}</strong> <span class="text-muted/60">services</span>
          <span class="px-1.5 text-muted/20">·</span>
          <strong>{resourceStore.directories.length}</strong> <span class="text-muted/60">directories</span>
        </span>
      </div>
    </div>
  </section>

  <section class="flex flex-col space-y-3">
    <div class="flex space-x-3">
      {@render actionButton(backupTask, 'icon-[regular--arrow-down-to-bracket]', 'Backup', handleBackup)}
      {@render actionButton(restoreTask, 'icon-[regular--arrow-up-from-bracket]', 'Restore', handleRestore)}
    </div>

    {@render statusMessage(backupTask)}
    {@render statusMessage(restoreTask)}

    <p class="px-1 text-[11px] leading-relaxed text-muted">
      Current config is automatically backed up before any restore.
    </p>
  </section>
</div>