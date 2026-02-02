<script lang="ts">
	import { slide } from 'svelte/transition';
	import { BrowserOpenURL } from '../../../wailsjs/runtime';
	import { systemStore } from '$lib/stores/system.svelte';
	import logo from '../../icons/logo.png';
	import { formatDate } from '$lib/helpers';

	/** State */
	let checking = $state(false);

	async function handleCheckUpdate() {
		checking = true;
		try {
			const version = await systemStore.getVersion();
			if (version.available) {
				systemStore.info = version;
			}
		} catch (e: unknown) {
			console.error('Error checking for updates:', e);
		} finally {
			checking = false;
		}
	}
</script>

<div
	class="scrollable mx-auto flex h-full w-full max-w-4xl flex-col"
	in:slide={{ axis: 'y', duration: 400 }}
>
	<header class="mb-10 flex items-center justify-between">
		<div class="flex items-center gap-4">
			<div class="flex flex-col gap-3">
				<h1 class="text-3xl font-bold tracking-tight">Settings</h1>
				<p class="text-sm text-muted">Manage application preferences and system behavior</p>
			</div>
		</div>
	</header>

	<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
		<!-- Settings -->
		<div class="flex flex-col gap-6 lg:col-span-8">
			<section class="flex flex-col gap-4">
				<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">General</h2>
				<div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
					{@render toggleItem({
						icon: 'icon-[regular--rocket-launch]',
						label: 'Launch at Startup',
						desc: 'Automatically start service when you log in',
						bindValue: systemStore.setting.runOnStartup
					})}

					<div class="h-px w-full bg-muted/10"></div>

					{@render toggleItem({
						icon: 'icon-[regular--command]',
						label: 'Run in Background',
						desc: 'Keep application running when closed',
						bindValue: systemStore.setting.runBackground
					})}
				</div>
			</section>

			<section class="flex flex-col gap-4">
				<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Notifications</h2>
				<div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
					{@render toggleItem({
						icon: 'icon-[regular--bell]',
						label: 'System Notifications',
						desc: 'Get slack notified when a service status changes',
						bindValue: systemStore.setting.notifications
					})}
				</div>
			</section>
		</div>

		<!-- About -->
		<div class="flex flex-col gap-4 lg:col-span-4">
			<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">About</h2>

			<div
				class="flex flex-col items-center gap-3 rounded-xl border border-muted/10 bg-surface/50 p-6 text-center"
			>
				<img src={logo} alt="Logo" class="h-24 w-24 bg-transparent object-cover" />

				<div class="flex flex-col gap-1">
					<h3 class="font-semibold">Service Watcher</h3>
					<p class="text-xs text-muted">Build {formatDate(systemStore.info.build ?? '')}</p>
				</div>

				<div
					class="inline-flex items-center gap-2 rounded-full border border-muted/10 bg-muted/5 px-3 py-1 font-mono text-xs text-muted"
				>
					{systemStore.info?.currentVersion}
				</div>

				<button
					disabled={checking}
					onclick={handleCheckUpdate}
					class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg border border-muted/10 bg-muted/10 py-2 text-sm font-medium transition-all hover:bg-muted/10 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-surface"
				>
					{#if checking}
						<i class="animate-spin icon-[regular--spinner]"></i>
						<span>Checking...</span>
					{:else}
						<span>Check for Updates</span>
					{/if}
				</button>

				<div
					class="flex w-full justify-between border-t border-muted/10 pt-4 text-xs text-muted/50"
				>
					<a
						href="https://github.com/zenkiet/window-service-watcher/blob/main/LICENSE"
						target="_blank"
						onclick={() =>
							BrowserOpenURL('https://github.com/zenkiet/window-service-watcher/blob/main/LICENSE')}
						class="transition-colors hover:text-muted">License</a
					>
					<a
						href="https://buymeacoffee.com/zenkiet"
						target="_blank"
						onclick={() => BrowserOpenURL('https://buymeacoffee.com/zenkiet')}
						class="transition-colors hover:text-muted">Sponsor</a
					>
					<a
						href="https://github.com/zenkiet/window-service-watcher"
						target="_blank"
						onclick={() => BrowserOpenURL('https://github.com/zenkiet/window-service-watcher')}
						class="transition-colors hover:text-muted">GitHub</a
					>
				</div>
			</div>
		</div>

		<div class="lg:hidden"></div>
	</div>

</div>

{#snippet toggleItem({
	icon,
	label,
	desc,
	bindValue
}: {
	icon: string;
	label: string;
	desc: string;
	bindValue: boolean;
})}
	<button
		class="group flex w-full items-center justify-between p-4 text-left transition-colors outline-none hover:bg-muted/4"
		onclick={() => {
			if (label === 'Launch at Startup')
				systemStore.setting.runOnStartup = !systemStore.setting.runOnStartup;
			if (label === 'Run in Background')
				systemStore.setting.runBackground = !systemStore.setting.runBackground;
			if (label === 'System Notifications')
				systemStore.setting.notifications = !systemStore.setting.notifications;
		}}
	>
		<div class="flex items-center gap-4">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-surface/50 text-muted/80 ring-1 ring-muted/10 transition-colors group-hover:text-muted group-hover:ring-muted/20"
			>
				<span class={`${icon} text-lg`}></span>
			</div>
			<div class="flex flex-col gap-0.5">
				<span class="text-sm font-medium">{label}</span>
				<span class="text-xs text-muted/80">{desc}</span>
			</div>
		</div>

		<div
			class={`relative h-6 w-11 rounded-full border border-transparent transition-colors duration-200 ease-in-out focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none ${bindValue ? 'bg-blue-600' : 'bg-muted/20'}`}
		>
			<span
				aria-hidden="true"
				class={`pointer-events-none inline-block h-5.5 w-5.5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${bindValue ? 'translate-x-5' : 'translate-x-0.5'}`}
			></span>
		</div>
	</button>
{/snippet}
