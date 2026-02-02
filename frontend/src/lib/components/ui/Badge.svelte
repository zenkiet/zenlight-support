<script lang="ts">
  let { type, count, label } = $props< {
    type: 'success' | 'danger';
    count?: number;
    label: string;
  }>();

  const styles: Record<string, string> = {
    success: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400 border-emerald-200 dark:border-emerald-500/20',
    danger: 'bg-rose-100 text-rose-700 dark:bg-rose-500/10 dark:text-rose-400 border-rose-200 dark:border-rose-500/20'
  };

  const dotColor: Record<string, string> = {
    success: 'bg-emerald-500 shadow-glow-green',
    danger: 'bg-rose-500 shadow-glow-red'
  };
</script>

<span class={`inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium border ${styles[type]} transition-colors`}>
  <span class={`w-1.5 h-1.5 rounded-full ${dotColor[type]} ${type === 'success' ? 'animate-pulse' : ''}`}></span>
  {label}{#if count !== undefined}: {count}{/if}
</span>