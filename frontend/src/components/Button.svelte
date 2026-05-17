<script lang="ts">
	export let classes = ""
	export let disabled = false
</script>

<button on:click class={classes} {disabled}><slot /></button>

<style lang="postcss">
	button {
		@apply px-6 py-3 rounded-2xl cursor-pointer transition-all duration-300 text-lg font-semibold;
		@apply bg-white/10 backdrop-blur-md border border-white/20 text-white shadow-lg;
	}
	button:hover {
		@apply bg-white/20 scale-[1.02] border-white/40 shadow-indigo-500/25;
	}
	button:active {
		@apply bg-white/5 scale-95;
	}
	button:disabled {
		@apply bg-black/20 text-white/30 border-white/5 shadow-none scale-100 cursor-not-allowed;
	}
</style>
