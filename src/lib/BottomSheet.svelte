<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { disableScroll, enableScroll } from './modal';
	import CloseIcon from './icons/CloseIcon.svelte';

	export let onClose: (() => void) | undefined = undefined;
	$: canClose = !!onClose;

	onMount(() => {
		disableScroll();
		return enableScroll;
	});

	const duration = 250;
	export const close = () => {
		if (!onClose) return;
		onClose();
	};

	export const keypress = (e: KeyboardEvent) => {
		if (e.key === 'Escape') {
			close();
		}
	};
</script>

<svelte:body on:keydown={keypress} />
<div transition:fade={{ duration }} class="overlay" on:click={close} role="presentation" />
<div class="sheet-container" transition:fly={{ y: 500, duration }}>
	{#if canClose}
		<div class="close-row">
			<button on:click={close}><CloseIcon /></button>
		</div>
	{/if}
	<slot />
</div>

<style lang="postcss">
	.overlay {
		@apply fixed bottom-0 left-0 right-0 top-0 bg-black bg-opacity-70;
		@apply pointer-events-auto select-auto;
	}
	.sheet-container {
		@apply fixed bottom-0 left-0 right-0 h-[80vh] overflow-auto bg-background p-4;
		@apply pointer-events-auto select-auto;
	}
	.close-row {
		@apply absolute right-2 top-2 flex justify-end;
	}
</style>
