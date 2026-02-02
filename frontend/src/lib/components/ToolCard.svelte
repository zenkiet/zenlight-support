<script lang="ts">
  let { data } = $props<{ data: {
     id: string;
  name: string;
  description: string;
  isInstalled: boolean;
  } }>();
</script>

<div class="group flex flex-col relative bg-surface rounded-xl border border-l-4 border-l-blue-500 border-y-border border-r-border shadow-sm dark:shadow-lg transition-all duration-300 overflow-hidden">
  <div class="absolute top-0 right-0 p-4 opacity-10 pointer-events-none">
    <i class="icon-[regular--puzzle-piece] text-6xl text-blue-500" ></i>
  </div>

  <div class="px-5 py-2 flex-1 relative z-10">
    <div class="flex justify-between items-start mb-2">
      <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-blue-100 text-blue-800 dark:bg-blue-500/10 dark:text-blue-300 border border-transparent dark:border-blue-500/20">
        TOOL
      </span>
    </div>

    <h3 class="text-base font-bold text-main mb-2">{data.name}</h3>
    <p class="text-sm text-muted mb-6">{data.description}</p>

    <div class="flex gap-3">
      <button class="flex-1 flex items-center justify-center gap-2 px-3 bg-primary hover:bg-blue-600 text-white text-sm font-medium rounded-lg transition-colors shadow-sm shadow-blue-500/20">
        <i class="icon-[regular--arrow-down-to-line] size-5"></i>
        Install
      </button>
      <button class="px-2 py-1 bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 text-main text-sm font-medium rounded-lg transition-colors border border-border" title="Details">
        <i class="icon-[regular--folder] size-5"></i>
      </button>
    </div>
  </div>
</div>