<script lang="ts">
	import { systemStore } from '$lib/stores/system.svelte';
	import { formatDate } from '$lib/helpers';
	import { bridge } from '$lib/services/wails-bridge';
	import logo from '../../../icons/logo.png';

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

	const links = [
		{
			label: 'License',
			url: 'https://github.com/zenkiet/zenlight-support/blob/main/LICENSE'
		},
		{
			label: 'Sponsor',
			url: 'https://buymeacoffee.com/zenkiet'
		},
		{
			label: 'GitHub',
			url: 'https://github.com/zenkiet/zenlight-support'
		}
	];
</script>

<div class="flex flex-col gap-6">
	<section class="flex flex-col gap-3">
		<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Application</h2>
		<div
			class="flex flex-col items-center gap-5 rounded-xl border border-muted/10 bg-surface/50 p-8 text-center"
		>
			<img src={logo} alt="Logo" class="h-20 w-20 bg-transparent object-cover" />

			<div class="flex flex-col gap-1.5">
				<h3 class="text-lg font-semibold font-quicksand">ZenLight Support</h3>
				<p class="text-xs text-muted/70">Build {formatDate(systemStore.info.build ?? '')}</p>
			</div>

			<div
				class="inline-flex items-center gap-2 rounded-full border border-muted/10 bg-muted/5 px-4 py-1.5 font-mono text-xs text-muted"
			>
				<span class="icon-[regular--tag] text-xs text-muted/50"></span>
				{systemStore.info?.currentVersion}
			</div>

			<button
				disabled={checking}
				onclick={handleCheckUpdate}
				class="flex w-full max-w-xs cursor-pointer items-center justify-center gap-2 rounded-lg border border-muted/10 bg-muted/10 py-2.5 text-sm font-medium transition-all hover:bg-muted/15 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-surface dark:hover:bg-surface-highlight"
			>
				{#if checking}
					<i class="animate-spin icon-[regular--spinner]"></i>
					<span>Checking...</span>
				{:else}
					<span class="icon-[regular--arrows-rotate] text-sm"></span>
					<span>Check for Updates</span>
				{/if}
			</button>
		</div>
	</section>

	<section class="flex flex-col gap-3">
		<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Links</h2>
		<div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
			{#each links as link, i}
				{#if i > 0}
					<div class="h-px w-full bg-muted/10"></div>
				{/if}
				<button
					onclick={() => bridge.openURL(link.url)}
					class="group flex w-full items-center justify-between p-4 text-left transition-colors hover:bg-muted/4"
				>
					<div class="flex items-center gap-3">
						<span class="text-sm font-medium">{link.label}</span>
					</div>
					<span
						class="icon-[regular--arrow-up-right-from-square] text-xs text-muted/40 transition-colors group-hover:text-muted"
					></span>
				</button>
			{/each}
		</div>
	</section>
</div>
