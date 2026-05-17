<script lang="ts">
	export let min: number | undefined
	export let max: number | undefined
	export let value = min

	let timeout: number | null = null
	const cap = (num: number): number => {
		if (typeof min === "number") num = Math.max(min, num)
		if (typeof max === "number") num = Math.min(max, num)
		return num
	}
	const hold = (num: number) => () => {
		release()
		let speed = 0.1
		let internalValue = (value ?? 0) - 0
		const update = () => {
			const prev = internalValue
			internalValue += speed * num
			internalValue = cap(internalValue)
			if (internalValue === prev) {
				if (timeout !== null) clearInterval(timeout)
				return
			}
			if (num > 0) {
				value = Math.ceil(internalValue)
			} else {
				value = Math.floor(internalValue)
			}
			speed *= 1.1
		}
		update()
		timeout = setInterval(update, 50)
	}
	const release = () => {
		if (timeout === null) return
		clearInterval(timeout)
		timeout = null
	}

	const update = () => {
		value = cap(value ?? 0)
	}

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			(e.target as HTMLElement).blur()
		}
	}
</script>

<div>
	<button
		on:mousedown={hold(-1)}
		on:mouseup={release}
		on:touchstart|preventDefault={hold(-1)}
		on:touchend|preventDefault={release}
		disabled={typeof min === "number" && (value ?? 0) <= min}>−</button
	>
	<input type="number" inputmode="numeric" {min} {max} bind:value on:change={update} on:keydown={handleKeydown} />
	<button
		on:mousedown={hold(1)}
		on:mouseup={release}
		on:touchstart|preventDefault={hold(1)}
		on:touchend|preventDefault={release}
		disabled={typeof max === "number" && (value ?? 0) >= max}>+</button
	>
</div>

<style lang="postcss">
	div {
		@apply flex;
	}
	button {
		@apply inline-flex w-8 h-8 bg-slate-200 text-slate-700 items-center justify-center text-lg font-bold rounded-lg transition-all duration-200 shadow-sm;
	}
	button:hover {
		@apply bg-slate-300 scale-[1.05];
	}
	button:active {
		@apply bg-slate-400 scale-95;
	}
	button:disabled {
		@apply bg-slate-100 text-slate-300 scale-100 cursor-not-allowed;
	}
	input {
		@apply bg-white border border-slate-300 text-slate-900 mx-1 p-1 w-14 h-8 text-center text-lg font-semibold rounded-lg outline-none transition-all shadow-sm;
	}
	input:focus {
		@apply ring-2 ring-violet-500/20 border-violet-500;
	}

	/* Reset. */
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		@apply m-0;
		-webkit-appearance: none;
	}

	input[type="number"] {
		-moz-appearance: textfield;
	}
</style>
