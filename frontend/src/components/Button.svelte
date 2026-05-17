<script lang="ts">
	export let classes = ""
	export let disabled = false
</script>

<button on:click class={classes} {disabled}><slot /></button>

<style lang="postcss">
	button {
		@apply border border-black px-6 py-3 rounded-xl cursor-pointer transition-all text-lg font-medium shadow-sm;
		@apply hover:bg-black hover:text-white active:bg-gray-800;
		@apply disabled:bg-gray-300 disabled:text-gray-700 disabled:shadow-none;
	}
</style>
