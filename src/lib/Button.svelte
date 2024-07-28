<script lang="ts">
	export let type: 'solid' | 'outline';
	export let accent: 'primary' | 'secondary' | 'highlight' = 'primary';
	export let size: 'sm' | 'md' | 'lg' = 'lg';
	export let disabled: boolean = false;

	$: colorClasses =
		type === 'solid'
			? `text-background bg-${accent} border-${accent}`
			: `text-${accent} bg-background border-${accent}`;
	const tailwindClasses =
		'text-primary text-secondary text-highlight bg-primary bg-secondary bg-highlight border-primary border-secondary border-highlight';
</script>

<button class={`${colorClasses} button-${size}`} class:disabled {disabled} on:click>
	<slot />
</button>

<style lang="postcss">
	button {
		@apply min-w-40 rounded-lg border-2 px-4 py-2 text-lg transition-all;
		@apply flex items-center justify-center gap-2;
	}
	.button-sm {
		@apply min-w-20 border px-2 py-1 text-sm;
	}
	.button-md {
		@apply min-w-28 border-2 px-2 py-1 text-base;
	}
	.button-outline {
		@apply bg-background text-primary;
	}
	.button-solid {
		@apply bg-primary text-background;
	}
	.disabled {
		@apply border-gray-400 bg-gray-400;
	}
	button:hover {
		@apply shadow-lg;
	}
</style>
