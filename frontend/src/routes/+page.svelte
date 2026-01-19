<script>
  import { EventsOn } from '../../wailsjs/runtime';
  import { GetServiceStatus, WatchLogs } from '../../wailsjs/go/app/App';
	import { WindowMinimise, WindowHide, Quit } from '../../wailsjs/runtime';

  let status = $state('Checking...');
  let isHealthy = $state(false);
  let logs = $state([]);

  $effect(() => {
    const interval = setInterval(async () => {
      try {
        const res = await GetServiceStatus();
        status = res.status;
        isHealthy = res.is_healthy;
      } catch (err) {
        console.error("Failed to get status:", err);
        status = "Connection Error";
      }
    }, 2000);

    // WatchLogs('logs/temp.log');

    const stopLogListener = EventsOn('new-log', (line) => {
        logs = [...logs.slice(-4), line];
    });

    return () => {
      clearInterval(interval);
      stopLogListener();
    };
  });

	const handleClose = () => WindowHide();
  const handleMin = () => WindowMinimise();
</script>

<main class="widget-container">
	<div
		class="flex h-full w-full cursor-default flex-col gap-2 overflow-hidden border border-white/20 bg-white/10 p-4"
	>
		<div class="flex items-center justify-between">
			<h3 class="text-sm font-semibold text-white">Blogic Report Service</h3>
			<div class="flex items-center gap-1.5">
				<div
					class="h-2 w-2 rounded-full {isHealthy ? 'bg-green-400' : 'animate-pulse bg-red-400'}"
				></div>
				<span class="font-mono text-xs text-white/80">{status}</span>
			</div>
		</div>

		<div
			class="group relative flex-1 overflow-hidden rounded-lg bg-black/20 p-2 font-mono text-[10px] text-white/70"
		>
			<div class="scrollbar-hide absolute inset-0 overflow-y-auto pl-2">
				{#if logs.length === 0}
					<div class="text-white/30 italic">Waiting for logs...</div>
				{/if}
				{#each logs as log}
					<div
						class="truncate border-b border-white/5 py-0.5 transition-colors last:border-0 hover:text-white"
					>
						{log}
					</div>
				{/each}
			</div>
		</div>
	</div>
</main>
