<script lang="ts">
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';

	import ConnectionError from '$lib/ConnectionError.svelte';
	import PauseOverlay from '$lib/PauseOverlay.svelte';
	import { type Session } from '$lib/login/login';
	import type { Config } from '$lib/score';
	import { estimatedQueueingTime, isHoldingLock, isQueueing } from './lock';
	import TransferResource from './TransferResource.svelte';
	import ClockIcon from '$lib/icons/ClockIcon.svelte';

	export let cfg: Config;
	export let session: Session;
	export let close: () => void;
	export let teamSelected = -1;
	let locked = false;
	$: lockHolders = session.lock_holder[session.team_id].players;
	$: timeRemaining = session.lock_holder[session.team_id].seconds_remaining;
	$: holdingLock = isHoldingLock(cfg, session);
	$: queueing = isQueueing(cfg, session);
	$: waitTime = estimatedQueueingTime(cfg, session);
	$: {
		if (holdingLock) locked = true;
		if (locked) break $;
		if (!queueing) {
			// Try to acquire lock.
			session.ws?.send('{"type": "lock"}');
		}
	}

	onMount(() => () => {
		// Send unlock on unmount.
		session.ws?.send('{"type": "unlock"}');
	});
</script>

{#if !session.ws}
	<ConnectionError />
{:else if locked && !holdingLock}
	<PauseOverlay text={$_('error.tradingExpired')}>
		<ClockIcon slot="icon" class="size-16" />
	</PauseOverlay>
{:else if !holdingLock && lockHolders.length === 0}
	<PauseOverlay text={$_('error.tradingWait')}>
		<ClockIcon slot="icon" class="size-16" />
	</PauseOverlay>
{:else if !holdingLock && lockHolders.length > 0}
	<PauseOverlay
		text={$_('error.trading', {
			values: { name: lockHolders[0].nickname, seconds: waitTime }
		})}
	>
		<ClockIcon slot="icon" class="size-16" />
	</PauseOverlay>
{:else if holdingLock}
	{#if !cfg.teams[session.team_id].admin}
		<p class="timeNotice">
			{$_('trade.label.timeRemaining', { values: { seconds: timeRemaining } })}
		</p>
		<hr />
	{/if}
	<TransferResource {cfg} {session} {close} {teamSelected} />
{/if}

<style lang="postcss">
	.timeNotice {
		@apply mr-6 mt-2 text-center;
	}
	hr {
		@apply mx-4 my-4;
		@apply border-b border-primary;
	}
</style>
