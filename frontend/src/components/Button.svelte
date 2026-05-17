<script lang="ts">
	export let classes = ""
	export let disabled = false
</script>

<button on:click class={classes} {disabled}><slot /></button>

<style lang="postcss">
	button {
		@apply border border-black px-6 py-3 rounded-xl cursor-pointer transition-all text-lg font-medium shadow-sm;
	}
	button:hover {
		@apply bg-black text-white;
	}
	button:active {
		@apply bg-gray-800;
	}
	button:disabled {
		@apply bg-gray-300 text-gray-700 shadow-none;
	}
</style>
