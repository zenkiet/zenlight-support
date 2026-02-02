<script>
	import Tooltip from "./Tooltip.svelte";
</script>

<Tooltip content="Comming soon">
  <button class="group relative flex flex-col items-center justify-center p-4 rounded-xl border-2 border-dashed border-border hover:border-primary hover:bg-blue-50 dark:hover:bg-blue-500/5 transition-all duration-300 cursor-pointer disabled:cursor-not-allowed disabled:opacity-30" disabled>
    <div class="w-16 h-16 rounded-full bg-surface-highlight group-hover:bg-blue-100 dark:group-hover:bg-blue-500/20 flex items-center justify-center mb-4 transition-colors">
      <i class="icon-[regular--plus] size-7 text-muted group-hover:text-primary transition-colors"></i>
    </div>
    <h3 class="text-sm font-bold text-main mb-1">Add Service</h3>
    <p class="text-xs text-muted text-center">Configure a new service watcher</p>
</button>
</Tooltip>

