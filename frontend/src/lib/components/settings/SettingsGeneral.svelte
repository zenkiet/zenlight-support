<script lang="ts">
	import { systemStore } from '$lib/stores/system.svelte';
</script>

<div class="flex flex-col gap-6">
	<section class="flex flex-col gap-3">
		<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Startup</h2>
		<div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
			{@render toggleItem({
				icon: 'icon-[regular--rocket-launch]',
				label: 'Launch at Startup',
				desc: 'Automatically start service when you log in',
				bindValue: systemStore.setting.runOnStartup,
				onToggle: () => (systemStore.setting.runOnStartup = !systemStore.setting.runOnStartup)
			})}

			<div class="h-px w-full bg-muted/10"></div>

			{@render toggleItem({
				icon: 'icon-[regular--command]',
				label: 'Run in Background',
				desc: 'Keep application running when closed',
				bindValue: systemStore.setting.runBackground,
				onToggle: () => (systemStore.setting.runBackground = !systemStore.setting.runBackground)
			})}
		</div>
	</section>

	<section class="flex flex-col gap-3">
		<h2 class="pl-1 text-xs font-bold tracking-wider text-muted/60 uppercase">Notifications</h2>
		<div class="overflow-hidden rounded-xl border border-muted/10 bg-surface/50">
			{@render toggleItem({
				icon: 'icon-[regular--bell]',
				label: 'System Notifications',
				desc: 'Get notified when a service status changes',
				bindValue: systemStore.setting.notifications,
				onToggle: () => (systemStore.setting.notifications = !systemStore.setting.notifications)
			})}
		</div>
	</section>
</div>

{#snippet toggleItem({
	icon,
	label,
	desc,
	bindValue,
	onToggle
}: {
	icon: string;
	label: string;
	desc: string;
	bindValue: boolean;
	onToggle: () => void;
})}
	<button
		class="group flex w-full items-center justify-between p-4 text-left transition-colors outline-none hover:bg-muted/4"
		onclick={onToggle}
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
