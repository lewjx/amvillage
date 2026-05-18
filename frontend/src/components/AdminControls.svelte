<script lang="ts">
	import { _ } from "svelte-i18n"
	import Button from "../components/Button.svelte"
	import { state, ws } from "../lib/amvillage"

	let notice = ""
	let popup = ""
	let chargeAmount = 0
	let chargeResource = 0
	const sendNotice = () => {
		$ws.send("notice " + notice)
	}
	const sendPopup = () => {
		$ws.send("popup " + popup)
	}
	const hidePopup = () => {
		$ws.send("popup")
	}
	const chargeAll = () => {
		$ws.send(`charge_all ${chargeResource} ${chargeAmount}`)
		
		let resName = chargeResource < $state.config.currencies.length 
			? $state.config.currencies[chargeResource] 
			: $state.config.gems[chargeResource - $state.config.currencies.length]
		
		let autoNotice = $_("admin.text.chargePopup", { values: { amount: chargeAmount, resource: resName } })
		$ws.send("popup " + autoNotice)
		
		alert($_("admin.text.chargeSuccess"))
		chargeAmount = 0
	}
</script>

<div class="controls">
	<span>{$_("admin.label.notice")}</span>
	<input bind:value={notice} />
	<Button on:click={sendNotice}>{$_("admin.button.sendNotice")}</Button>
	<span>{$_("admin.label.banner")}</span>
	<input bind:value={popup} />
	<div class="flex flex-wrap justify-center gap-1">
		<Button on:click={sendPopup}>{$_("admin.button.sendBanner")}</Button>
		<Button on:click={hidePopup}>{$_("admin.button.hideBanner")}</Button>
	</div>
	<span>{$_("admin.label.chargeAll")}</span>
	<div class="flex flex-wrap items-center gap-2">
		<input type="number" bind:value={chargeAmount} class="w-16" min="0" />
		<select bind:value={chargeResource} class="border border-black p-1">
			{#each $state.config.currencies as cur, i}
				<option value={i}>{cur}</option>
			{/each}
			{#each $state.config.gems as gem, i}
				<option value={$state.config.currencies.length + i}>{gem}</option>
			{/each}
		</select>
	</div>
	<Button on:click={chargeAll} disabled={chargeAmount < 0}>{$_("admin.button.charge")}</Button>
</div>

<style lang="postcss">
	.controls {
		@apply grid justify-center items-center p-2 gap-2;
		grid-template-columns: auto 10rem auto;
	}
	input {
		@apply border border-black;
	}
</style>
