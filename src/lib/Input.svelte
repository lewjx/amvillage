<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	// NOTE: min is assumed to be 0 or negative. max is assumed to be positive.

	export let problematic: boolean;
	export let value: number;
	export let min: number;
	export let max: number;

	const dispatch = createEventDispatcher();
	const onChange = () => {
		dispatch('change', {});
	};

	const capValue = (e: Event & { currentTarget: HTMLInputElement | null }) => {
		if (e.currentTarget !== null) {
			const inputMax = sign === 1 ? max : -min;
			inputValue = Math.max(Math.min(+e.currentTarget.value, inputMax), 0);
			// Force value to be a valid number.
			e.currentTarget.value = +inputValue + '';
			onChange();
		}
	};
	const clearValue = (e: Event & { currentTarget: HTMLInputElement | null }) => {
		if (e.currentTarget !== null) {
			e.currentTarget.value = '';
		}
	};

	let inputValue = Math.abs(value);
	let sign = value >= 0 ? 1 : -1;
	$: value = inputValue * sign;

	const toggleNegative = () => (sign *= -1);
	let hold = 0;
	let speed = 1;
	const setHold = (newHold: number) => {
		hold = newHold;
		speed = 1;
	};

	$: {
		if (sign === -1 && min >= 0) sign = 1;
	}

	const update = (hold: number) => {
		if (hold === 0) {
			speed = 1;
			return;
		}
		speed *= 1.1;
		inputValue += Math.floor(speed) * hold;
		const inputMax = sign === 1 ? max : -min;
		inputValue = Math.max(Math.min(inputValue, inputMax), 0);
		onChange();
	};
	onMount(() => {
		let interval = setInterval(() => update(hold), 100);
		return () => clearInterval(interval);
	});
</script>

<div>
	<button class:problematic on:click={toggleNegative}>
		{#if sign === 1}x{:else}x -{/if}
	</button>
	<input
		class:problematic
		class:gray={value === 0}
		type="number"
		bind:value={inputValue}
		min="0"
		{max}
		on:input={capValue}
		on:focus={clearValue}
		on:blur={capValue}
	/>
	<p class="total" class:problematic>&nbsp;/ {sign === 1 ? max : -min}</p>
	<div class="modifier">
		<button
			on:click={() => update(-1)}
			on:keypress={() => update(-1)}
			on:pointerdown={() => setHold(-1)}
			on:pointerup={() => setHold(0)}
		>
			-
		</button>
		<button
			on:click={() => update(1)}
			on:keypress={() => update(1)}
			on:pointerdown={() => setHold(1)}
			on:pointerup={() => setHold(0)}
		>
			+
		</button>
	</div>
</div>

<style lang="postcss">
	div {
		@apply grid items-center;
		grid-template-columns: auto 1fr auto auto;
	}
	button {
		@apply grid min-w-8 items-center justify-center text-xl;
	}

	input {
		@apply min-w-0 bg-background text-right text-xl;
	}
	.total {
		@apply min-w-16 self-end text-left text-sm;
	}
	input.gray {
		@apply text-gray-400;
	}
	.problematic {
		@apply text-highlight;
	}

	.modifier {
		@apply ml-2 flex gap-1 border border-primary;
	}

	input[type='number']::-webkit-outer-spin-button,
	input[type='number']::-webkit-inner-spin-button,
	input[type='number'] {
		-webkit-appearance: none;
		margin: 0;
		-moz-appearance: textfield !important;
	}
</style>
