<script lang="ts">
	import { page } from '$app/state';
	import Footer from '../Footer.svelte';

	let { isOpen = false, onClose } = $props<{
		isOpen?: boolean;
		onClose?: () => void;
	}>();

	const navGroups = [
		{
			title: 'Features',
			items: [
				{ name: 'Services', icon: 'icon-[regular--server]', href: '/' },
				{ name: 'Projects', icon: 'icon-[regular--diagram-project]', href: '/statics' },
				{ name: 'Run Script', icon: 'icon-[regular--code]', href: '/script' }
			]
		},
		{
			title: 'Configuration',
			items: [
				{ name: 'Settings', icon: 'icon-[regular--gear-complex]', href: '/settings' }
				// { name: 'Logs', icon: 'icon-[regular--rectangle-history]', href: '/logs' }
			]
		}
	];

	const currentPath = $derived(page.url.pathname);
	function isItemActive(href: string): boolean {
		if (href === '/') return currentPath === '/';
		return currentPath.startsWith(href);
	}
</script>

<aside
	class={`
        fixed inset-y-0 left-0 z-50 flex h-full w-64 shrink-0 flex-col overflow-y-auto
        border-r border-border bg-surface pt-5 shadow-2xl transition-transform duration-300 ease-in-out md:shadow-none
        ${isOpen ? 'translate-x-0' : '-translate-x-full'}
        md:static md:translate-x-0
    `}
>
	<div class="flex h-full flex-col px-4 pb-4">
		<div class="mb-6 flex items-center justify-between gap-3">
			<h1 class="text-lg font-bold tracking-tight text-main">ZenB Tool</h1>
		</div>

		{#each navGroups as group}
			<nav class="mb-6 space-y-1">
				<p class="mb-2 px-3 text-xs font-semibold tracking-wider text-muted/70 uppercase">
					{group.title}
				</p>
				{#each group.items as item}
					{@render navItem(item)}
				{/each}
			</nav>
		{/each}
	</div>

	<div class="mt-auto">
		<Footer />
	</div>
</aside>

{#snippet navItem(item: { name: string; icon: string; href: string })}
	{@const active = isItemActive(item.href)}
	<a
		href={item.href}
		onclick={() => window.innerWidth < 768 && onClose?.()}
		class={`group relative flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-all duration-200
        ${
					active
						? 'bg-blue-50 text-primary dark:bg-blue-500/10'
						: 'hover:bg-surface-highlight text-muted hover:text-main'
				}`}
	>
		<div
			class={`absolute top-0 left-0 h-full w-1 rounded-l-md bg-blue-500 transition-opacity ${
				active ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'
			}`}
		></div>

		<i class="{item.icon} size-4"></i>
		{item.name}
	</a>
{/snippet}
