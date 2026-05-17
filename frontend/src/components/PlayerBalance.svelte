<script lang="ts">
	import { _ } from "svelte-i18n"
	import Bar from "../components/Bar.svelte"
	import I18n from "../components/I18n.svelte"
	import { state } from "../lib/amvillage"

	$: teamName = $state.config.teams[$state.team]?.name ?? ""
	$: currencyCount = $state.config.currencies.length
	$: balance = $state.balances[$state.team].slice(0, currencyCount) ?? Array(currencyCount)
	$: gemBalance = $state.balances[$state.team].slice(currencyCount)
	$: sum = Math.min(...balance)
	$: max = Math.max(...balance)
</script>

<p class="label"><I18n id="mainmenu.label.balance" values={{ teamName }} /></p>
<span class="summary">{$_("mainmenu.label.combined", { values: { total: sum } })}</span>
<hr />
<p class="label">{$_("mainmenu.label.combinedExplanation")}</p>
<table>
	{#each $state.config.currencies as currency, i}
		<tr class:isMin={balance[i] === sum}>
			<td><Bar value={balance[i]} {max} isMin={balance[i] === sum} /></td>
			<td>{balance[i]} {currency}</td>
		</tr>
	{/each}
</table>
<hr />
<div class="gem">
	{#each $state.config.gems as gem, i}
		<span class:noGem={gemBalance[i] === 0}>{gemBalance[i]} {gem}</span>
	{/each}
</div>

<style lang="postcss">
	* {
		@apply text-slate-800;
	}
	.summary {
		@apply text-6xl font-extrabold text-violet-600 py-2;
	}
	hr {
		@apply border-slate-200 w-full;
	}
	table {
		@apply mx-auto border-separate border-spacing-x-4 border-spacing-y-2;
	}
	tr {
		@apply text-left text-xl font-medium;
	}
	.isMin td {
		@apply text-red-500;
	}
	.gem {
		@apply flex flex-wrap gap-x-6 gap-y-2 items-center justify-center text-lg font-semibold bg-slate-50 p-4 rounded-2xl w-full border border-slate-100;
	}
	.noGem {
		@apply text-red-500;
	}
</style>
