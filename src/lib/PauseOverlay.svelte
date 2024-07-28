<script lang="ts">
	export let text: string;
</script>

<div>
	<h1><slot name="icon" /></h1>
	<p>{text}</p>
</div>

<style lang="postcss">
	div {
		@apply h-full w-full;
		@apply flex flex-col items-center justify-center gap-2;
		@apply text-center;
	}
	h1 {
		@apply text-4xl;
	}
	p {
		@apply text-lg;
	}
</style>
