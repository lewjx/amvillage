<script lang="ts">
	import { _ } from "svelte-i18n"
	import { state } from "../lib/amvillage"
	import I18n from "./I18n.svelte"

	const currencyCount = $state.config.currencies.length
</script>

<div><I18n id="admin.text.explanation" /></div>
<div class="table-container">
	<table>
		<tr>
			<th class="first" />
			{#each $state.config.teams as team}
				{#if team.is_admin}
					<td>{$_("admin.label.given", { values: { team: team.name } })}</td>
				{:else}
					<td>{team.name}</td>
				{/if}
			{/each}
		</tr>
		{#each $state.config.currencies as currency, i}
			<tr>
				<td class="first">{currency}</td>
				{#each $state.config.teams as team, j}
					{#if team.is_admin}
						<td>{-$state.balances[j][i]}</td>
					{:else}
						<td class:zero={$state.balances[j][i] <= 0}>{$state.balances[j][i]}</td>
					{/if}
				{/each}
			</tr>
		{/each}
		{#each $state.config.gems as gem, i}
			<tr>
				<td class="first">{gem}</td>
				{#each $state.config.teams as team, j}
					{#if team.is_admin}
						<td>{-$state.balances[j][currencyCount + i]}</td>
					{:else}
						<td class:zero={$state.balances[j][currencyCount + i] <= 0}>{$state.balances[j][currencyCount + i]}</td>
					{/if}
				{/each}
			</tr>
		{/each}
	</table>
</div>

<style lang="postcss">
	.table-container {
		@apply relative max-w-full overflow-x-auto border border-slate-200 rounded-2xl bg-white shadow-sm;
	}
	table {
		@apply border-separate border-spacing-0 w-full text-sm text-slate-800;
	}
	th,
	td {
		@apply px-4 py-3 border-b border-r border-slate-200 min-w-[3rem] text-center font-medium;
	}
	.first {
		@apply sticky left-0 bg-slate-50 min-w-[6rem] text-left border-r-slate-300 shadow-[2px_0_5px_rgba(0,0,0,0.05)];
	}
	.zero {
		@apply text-red-500;
	}
</style>
