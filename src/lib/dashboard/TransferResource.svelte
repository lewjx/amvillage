<script lang="ts">
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';

	import Button from '$lib/Button.svelte';
	import type { Session } from '$lib/login/login';
	import type { Config } from '$lib/score';
	import { isHoldingLock } from './lock';
	import { isValid } from './transfer';
	import Input from '$lib/Input.svelte';

	export let cfg: Config;
	export let session: Session;
	export let close: () => void;
	onMount(() => {
		/* This component assumes that the player has acquired the trading lock. */
		if (!isHoldingLock(cfg, session)) {
			throw new Error(
				'Invariance check failed! Expected the active player to hold the trading lock.'
			);
		}
	});

	$: isAdmin = cfg.teams[session.team_id].admin;
	let state: 'resources' | 'gems' = 'resources';
	let teamSelected = -1;

	let resourceCount = cfg.resource_names.map(() => 0);
	let gemCount = cfg.gem_names.map(() => 0);
	$: resourceNotable = resourceCount.some((x) => x !== 0);
	$: gemNotable = gemCount.some((x) => x !== 0);
	$: okChecks = isValid(cfg, session, teamSelected, resourceCount, gemCount);

	// Used for generalizing resource/gem tabs.
	$: labels = state === 'resources' ? cfg.resource_names : cfg.gem_names;
	$: count = state === 'resources' ? resourceCount : gemCount;
	$: teamCount =
		state === 'resources'
			? session.score[session.team_id].resources
			: session.score[session.team_id].gems;
	$: targetTeamCount =
		teamSelected === -1
			? []
			: state === 'resources'
				? session.score[teamSelected].resources
				: session.score[teamSelected].gems;
	$: ok = state === 'resources' ? okChecks.resource : okChecks.gem;

	$: allGemOK = okChecks.gem.every((x) => x);
	$: allResourceOK = okChecks.resource.every((x) => x);
	$: canTransfer = allGemOK && allResourceOK && (resourceNotable || gemNotable);
	const transfer = () => {
		const msg = JSON.stringify({
			type: 'transfer',
			from: session.team_id,
			to: teamSelected,
			gem_amount: gemCount,
			resource_amount: resourceCount
		});
		if (!session.ws) {
			// We lost connection. Keep this modal open for now.
			console.warn('Refusing to transfer without connection');
			return;
		}
		session.ws.send(msg);
		close();
	};

	const selectTeam = (select: number) => {
		teamSelected = select;
		state = 'resources';
	};
	const onChange = () => {
		resourceCount = resourceCount;
		gemCount = gemCount;
	};
</script>

{#if teamSelected === -1}
	<p class="team-buttons-label">{$_('trade.label.chooseTarget')}</p>
	<div class="team-buttons">
		{#each cfg.teams as team, i}
			{#if i !== session.team_id}
				<Button type="outline" on:click={() => selectTeam(i)}>{team.name}</Button>
			{/if}
		{/each}
	</div>
{:else}
	<div class="flex flex-col items-center">
		<div class="trade-target">
			<p>{$_('trade.label.target', { values: { teamName: cfg.teams[teamSelected].name } })}</p>
			<Button type="outline" size="sm" on:click={() => (teamSelected = -1)}>
				{$_('trade.button.changeTarget')}
			</Button>
		</div>
		<div class="resource-tabs">
			<Button
				size="md"
				type={state === 'resources' ? 'solid' : 'outline'}
				on:click={() => (state = 'resources')}
			>
				{$_('trade.button.selectResource')}
				{#if resourceNotable && allResourceOK}
					✅
				{:else if resourceNotable}
					❌
				{/if}
			</Button>
			<Button
				size="md"
				type={state === 'gems' ? 'solid' : 'outline'}
				on:click={() => (state = 'gems')}
			>
				{$_('trade.button.selectGem')}
				{#if gemNotable && allGemOK}
					✅
				{:else if gemNotable}
					❌
				{/if}
			</Button>
		</div>
		<div class="resources">
			{#key state}
				{#each count as _, i}
					<p class="problem-label">
						{#if !ok[i]}❌{/if}
					</p>
					<p class:problematic={!ok[i]}>{labels[i]}</p>
					<Input
						problematic={!ok[i]}
						bind:value={count[i]}
						on:change={onChange}
						min={isAdmin ? -targetTeamCount[i] : 0}
						max={isAdmin ? 999999 : teamCount[i]}
					/>
				{/each}
			{/key}
		</div>
		<Button type="solid" on:click={transfer} disabled={!canTransfer}>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="size-5"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M12 6v12m-3-2.818.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
				/>
			</svg>
			{$_('trade.button.transfer')}
		</Button>
	</div>
{/if}

<style lang="postcss">
	.team-buttons-label {
		@apply mb-2 text-center;
	}
	.team-buttons {
		@apply flex flex-wrap justify-center gap-2;
	}
	.trade-target {
		@apply flex items-center justify-center gap-2;
		@apply mb-4;
	}
	.trade-target p {
		@apply text-center text-lg;
	}
	.resource-tabs {
		@apply flex justify-center gap-1;
	}
	.resources {
		@apply my-4 grid items-center gap-x-1 gap-y-1;
		width: min(100%, 30rem);
		grid-template-columns: 2rem auto minmax(auto, 1fr);
	}
	.resources .problem-label {
		@apply text-xl;
	}
</style>
