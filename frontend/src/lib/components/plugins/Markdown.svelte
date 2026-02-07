<script lang="ts">
  import { marked } from 'marked';
  import { systemStore } from '$lib/stores/system.svelte';
  import { bridge } from '$lib/services/wails-bridge';

  let { content = '' } = $props();

  marked.setOptions({
    breaks: true,
    gfm: true
  });

  let html = $derived(marked(content));

  function handleClick(e: MouseEvent) {
    const target = (e.target as HTMLElement).closest('a');
    if (target?.href) {
      e.preventDefault();
      bridge.openURL(target.href);
    }
  }
</script>

<div class="markdown-content" class:dark={systemStore.isDark} onclick={handleClick} role="none">
  {@html html}
</div>

<style>
  .markdown-content {
    font-size: 13px;
    line-height: 1.7;
    color: #475569;
  }

  .markdown-content.dark {
    color: #94a3b8;
  }

  .markdown-content :global(h1),
  .markdown-content :global(h2),
  .markdown-content :global(h3) {
    color: #0f172a;
    margin-top: 1rem;
    margin-bottom: 0.5rem;
  }

  .markdown-content.dark :global(h1),
  .markdown-content.dark :global(h2),
  .markdown-content.dark :global(h3) {
    color: #f1f5f9;
  }

  .markdown-content :global(h1) {
    font-size: 1.25rem;
    font-weight: 700;
  }

  .markdown-content :global(h2) {
    font-size: 1.1rem;
    font-weight: 600;
  }

  .markdown-content :global(h3) {
    font-size: 0.95rem;
    font-weight: 600;
  }

  .markdown-content :global(p) {
    margin-bottom: 0.625rem;
  }

  .markdown-content :global(ul),
  .markdown-content :global(ol) {
    padding-left: 1.25rem;
    margin-bottom: 0.625rem;
  }

  .markdown-content :global(ul) {
    list-style-type: disc;
  }

  .markdown-content :global(ol) {
    list-style-type: decimal;
  }

  .markdown-content :global(li) {
    margin-bottom: 0.25rem;
  }

  .markdown-content :global(li::marker) {
    opacity: 0.4;
  }

  .markdown-content :global(strong) {
    font-weight: 600;
    color: #0f172a;
  }

  .markdown-content.dark :global(strong) {
    color: #f1f5f9;
  }

  .markdown-content :global(code) {
    font-family: 'JetBrains Mono', ui-monospace, monospace;
    font-size: 0.8em;
    padding: 0.15em 0.4em;
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(0, 0, 0, 0.08);
    color: #0f172a;
  }

  .markdown-content.dark :global(code) {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.08);
    color: #f1f5f9;
  }

  .markdown-content :global(pre) {
    margin-bottom: 0.75rem;
    padding: 0.75rem 1rem;
    border-radius: 8px;
    background: rgba(0, 0, 0, 0.04);
    border: 1px solid rgba(0, 0, 0, 0.06);
    overflow-x: auto;
  }

  .markdown-content.dark :global(pre) {
    background: rgba(0, 0, 0, 0.25);
    border-color: rgba(255, 255, 255, 0.05);
  }

  .markdown-content :global(pre code) {
    padding: 0;
    background: none;
    border: none;
    font-size: 0.8em;
  }

  .markdown-content :global(a) {
    color: #3b82f6;
    text-decoration: none;
    transition: color 0.15s;
  }

  .markdown-content :global(a:hover) {
    color: #2563eb;
    text-decoration: underline;
  }

  .markdown-content.dark :global(a) {
    color: #60a5fa;
  }

  .markdown-content.dark :global(a:hover) {
    color: #93bbfd;
  }

  .markdown-content :global(blockquote) {
    margin: 0.75rem 0;
    padding: 0.5rem 1rem;
    border-left: 3px solid rgba(59, 130, 246, 0.3);
    background: rgba(59, 130, 246, 0.04);
    border-radius: 0 6px 6px 0;
  }

  .markdown-content.dark :global(blockquote) {
    border-left-color: rgba(96, 165, 250, 0.4);
    background: rgba(96, 165, 250, 0.05);
  }

  .markdown-content :global(hr) {
    border: none;
    height: 1px;
    background: rgba(0, 0, 0, 0.1);
    margin: 1rem 0;
  }

  .markdown-content.dark :global(hr) {
    background: rgba(255, 255, 255, 0.08);
  }

  .markdown-content :global(table) {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 0.75rem;
    font-size: 0.85em;
  }

  .markdown-content :global(th),
  .markdown-content :global(td) {
    padding: 0.4rem 0.625rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
    text-align: left;
  }

  .markdown-content.dark :global(th),
  .markdown-content.dark :global(td) {
    border-color: rgba(255, 255, 255, 0.08);
  }

  .markdown-content :global(th) {
    font-weight: 600;
    color: #0f172a;
    background: rgba(0, 0, 0, 0.03);
  }

  .markdown-content.dark :global(th) {
    color: #f1f5f9;
    background: rgba(255, 255, 255, 0.04);
  }

  .markdown-content :global(> :first-child) {
    margin-top: 0;
  }

  .markdown-content :global(> :last-child) {
    margin-bottom: 0;
  }
</style>