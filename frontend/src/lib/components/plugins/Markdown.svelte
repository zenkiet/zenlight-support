<script>
  import { marked } from 'marked';
  let { content = '' } = $props();

  marked.setOptions({
    breaks: true,
    gfm: true
  });

  let html = $derived(marked(content));
</script>

<div class="markdown-content">
  {@html html}
</div>