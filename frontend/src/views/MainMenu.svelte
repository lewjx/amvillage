<script lang="ts">
	import { fly } from "svelte/transition"
	import AdminBalance from "../components/AdminBalance.svelte"
	import AdminControls from "../components/AdminControls.svelte"
	import PlayerBalance from "../components/PlayerBalance.svelte"
	import Button from "../components/Button.svelte"
	import I18n from "../components/I18n.svelte"
	import { state, ws } from "../lib/amvillage"
	import { status } from "../lib/state"

	$: isAdmin = $state.config.teams[$state.team].is_admin
	$: lock = $state.locks[$state.team]

	let loading = false
	const trade = (i: number) => () => {
		if (!isAdmin) $ws.send("lock")
		loading = true
		let interval: number
		interval = setInterval(() => {
			if ($state.locks[$state.team] !== null || isAdmin) {
				$status = {
					status: "trade",
					target: i,
				}
				clearInterval(interval)
				loading = false
			}
		}, 100)
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="welcome">
		<I18n id="mainmenu.text.welcome" values={{ username: $state.username }} />
	</div>
	<div class="balance">
		{#if isAdmin}
			<AdminBalance />
		{:else}
			<PlayerBalance />
		{/if}
	</div>
	<div class="trade">
		{#if $state.config.teams[$state.team].is_admin || lock === null}
			{#each $state.config.teams as team, i}
				{#if i !== $state.team}
					<Button on:click={trade(i)} disabled={loading}>{team.name}</Button>
				{/if}
			{/each}
		{:else}
			<span class="error">
				<I18n
					id="mainmenu.error.othersTrading"
					values={{ player: lock.member, sec: Math.floor(lock.millis_left / 1000) }}
				/>
			</span>
		{/if}
	</div>
	{#if isAdmin}
		<hr />
		<AdminControls />
	{/if}
</main>

<style lang="postcss">
	main {
		@apply flex flex-col gap-4 py-4;
	}
	.welcome {
		@apply px-4;
	}
	.balance {
		@apply flex flex-col items-center bg-white/10 backdrop-blur-xl border border-white/20 text-white text-center gap-6 p-8 rounded-3xl shadow-2xl;
	}
	.trade {
		@apply flex flex-wrap gap-4 justify-center;
	}
	.error {
		@apply text-red-800;
	}
</style>
