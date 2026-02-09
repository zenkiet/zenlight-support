<script lang="ts">
	import { slide, fade } from 'svelte/transition';
	import SettingsGeneral from '$lib/components/settings/SettingsGeneral.svelte';
	import SettingsAppearance from '$lib/components/settings/SettingsAppearance.svelte';
	import SettingsBackup from '$lib/components/settings/SettingsBackup.svelte';
	import SettingsAbout from '$lib/components/settings/SettingsAbout.svelte';

	type SettingsTab = 'general' | 'appearance' | 'backup' | 'about';

	const tabs: { id: SettingsTab; label: string }[] = [
		{ id: 'general', label: 'General' },
		{ id: 'appearance', label: 'Appearance' },
		{ id: 'backup', label: 'Backup' },
		{ id: 'about', label: 'About' }
	];

	let activeTab = $state<SettingsTab>('general');
</script>

<div
	class="scrollable mx-auto flex h-full w-full max-w-4xl flex-col"
	in:slide={{ axis: 'y', duration: 400 }}
>
	<header class="mb-8 flex flex-col gap-6">
		<div class="flex flex-col gap-2">
			<h1 class="text-3xl font-bold tracking-tight">Settings</h1>
			<p class="text-sm text-muted">Manage application preferences and system behavior</p>
		</div>

		<div class="flex gap-1 rounded-xl border border-muted/10 bg-surface/50 p-1">
			{#each tabs as tab}
				<button
					onclick={() => (activeTab = tab.id)}
					class={`relative flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all duration-200 outline-none
						${
							activeTab === tab.id
								? 'bg-blue-200/60 text-blue-700 dark:bg-blue-500/10 dark:text-blue-300'
								: 'text-muted/70 hover:bg-muted/5 hover:text-muted'
						}`}
				>
					<span>{tab.label}</span>
				</button>
			{/each}
		</div>
	</header>

	<div class="flex-1">
		{#key activeTab}
			<div in:fade={{ duration: 200 }}>
				{#if activeTab === 'general'}
					<SettingsGeneral />
				{:else if activeTab === 'appearance'}
					<SettingsAppearance />
				{:else if activeTab === 'backup'}
					<SettingsBackup />
				{:else if activeTab === 'about'}
					<SettingsAbout />
				{/if}
			</div>
		{/key}
	</div>
</div>
