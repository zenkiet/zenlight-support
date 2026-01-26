<script lang="ts">
	import Footer from '$lib/components/Footer.svelte';
	import UpdateModal from '$lib/components/modals/UpdateModal.svelte';
	import { serviceStore } from '$lib/stores/services.svelte';
	import { systemStore } from '$lib/stores/system.svelte';
	import './layout.css';
	let { children } = $props();

	serviceStore.init();
	systemStore.init();
</script>

<div class="flex h-screen w-screen flex-col overflow-hidden selection:bg-blue-500/30 pt-5">
	<main class="relative flex-1 overflow-hidden">
		{#if serviceStore.initialized}
			{@render children()}
		{:else}
			<div class="flex items-center justify-center h-full">
				<div class="text-center">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
					<p class="text-gray-600">Loading services...</p>
				</div>
			</div>
		{/if}
	</main>

	<Footer />

	<UpdateModal />
</div>