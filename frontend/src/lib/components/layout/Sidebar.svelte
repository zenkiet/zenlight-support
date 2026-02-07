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
				{ name: 'Projects', icon: 'icon-[regular--diagram-project]', href: '/projects' },
				{ name: 'Run Script', icon: 'icon-[regular--code]', href: '/script' }
			]
		},
		{
			title: 'Configuration',
			items: [{ name: 'Settings', icon: 'icon-[regular--gear-complex]', href: '/settings' }]
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
		fixed inset-y-0 left-0 z-50 flex h-full w-60 shrink-0 flex-col
		border-r border-muted/10 bg-surface transition-transform duration-300 ease-in-out
		${isOpen ? 'translate-x-0' : '-translate-x-full'}
		md:static md:translate-x-0
	`}
>
	<!-- Brand -->
	<div class="flex h-14 items-center gap-3 border-b border-muted/10 px-5">
		<span class="font-quicksand font-bold tracking-tight text-main">ZenLight Support</span>
	</div>

	<!-- Navigation -->
	<nav class="scrollable flex-1 overflow-y-auto px-3 py-4">
		{#each navGroups as group, gi}
			{#if gi > 0}
				<div class="my-3 h-px bg-muted/8"></div>
			{/if}
			<p class="mb-2 px-2 text-[10px] font-semibold tracking-widest text-muted/50 uppercase">
				{group.title}
			</p>
			<ul class="flex flex-col gap-0.5">
				{#each group.items as item}
					{@render navItem(item)}
				{/each}
			</ul>
		{/each}
	</nav>

	<Footer />
</aside>

{#snippet navItem(item: { name: string; icon: string; href: string })}
	{@const active = isItemActive(item.href)}
	<li>
		<a
			href={item.href}
			onclick={() => window.innerWidth < 768 && onClose?.()}
			class={`group relative flex items-center gap-3 rounded-lg px-2.5 py-2 text-[13px] font-medium transition-all duration-150
				${
					active
						? 'bg-blue-500/8 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400'
						: 'text-muted/70 hover:bg-muted/6 hover:text-main'
				}`}
		>
			{#if active}
				<div
					class="absolute top-1/2 left-0 h-5 w-0.75 -translate-y-1/2 rounded-r-full bg-blue-500"
				></div>
			{/if}

			<span
				class={`${item.icon} size-4 transition-colors ${active ? 'text-blue-500 dark:text-blue-400' : 'text-muted/50 group-hover:text-muted'}`}
			></span>
			{item.name}
		</a>
	</li>
{/snippet}
