<script lang="ts">
	export let classes = ""
	export let disabled = false
</script>

<button on:click class={classes} {disabled}><slot /></button>

<style lang="postcss">
	button {
		@apply px-6 py-3 rounded-2xl cursor-pointer transition-all duration-300 text-lg font-semibold;
		@apply bg-violet-600 text-white shadow-md shadow-violet-600/30;
	}
	button:hover {
		@apply bg-violet-700 scale-[1.02] shadow-lg shadow-violet-600/40;
	}
	button:active {
		@apply bg-violet-800 scale-95;
	}
	button:disabled {
		@apply bg-slate-200 text-slate-400 shadow-none scale-100 cursor-not-allowed;
	}
</style>
