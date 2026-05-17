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
</script>

<div>
	<button
		on:mousedown={hold(-1)}
		on:mouseup={release}
		on:touchstart|preventDefault={hold(-1)}
		on:touchend|preventDefault={release}
		disabled={typeof min === "number" && (value ?? 0) <= min}>−</button
	>
	<input type="number" inputmode="numeric" {min} {max} bind:value on:change={update} />
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
		@apply inline-flex w-12 h-12 bg-black text-white items-center justify-center text-2xl font-bold rounded-md;
	}
	button:disabled {
		@apply bg-gray-300 text-gray-700;
	}
	input {
		@apply border border-black mx-1 p-1 w-20 h-12 text-center text-xl rounded-md;
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
