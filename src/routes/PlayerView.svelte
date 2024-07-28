<script lang="ts">
	import type { Readable } from 'svelte/store';
	import Button from '$lib/Button.svelte';
	import ConnectionError from '$lib/ConnectionError.svelte';
	import Messages from '$lib/dashboard/Messages.svelte';
	import Score from '$lib/dashboard/Score.svelte';
	import type { Session } from '$lib/login/login';
	import { type Config, finalScore } from '$lib/score';
	import DetailedScore from './DetailedScore.svelte';
	import PauseOverlay from '$lib/PauseOverlay.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';

	export let cfg: Config;
	export let session: Readable<Session>;
	export let openSend: () => void;
	export let openTutorial: () => void;

	$: score = $session.score[$session.team_id];
	$: fullScreenNotice = $session.notices.find((x) => x.level === 'pause' && !x.dismissed);
</script>

{#if !$session.ws}
	<div class="h-screen p-4">
		<ConnectionError />
	</div>
{:else if fullScreenNotice}
	<div class="h-screen p-8">
		<PauseOverlay text={fullScreenNotice.message}>
			<PauseIcon class="size-16 text-highlight" slot="icon" />
		</PauseOverlay>
	</div>
{:else}
	<Score score={finalScore(score)} />
	<DetailedScore {score} {cfg} />
	<div class="buttons">
		<Button type="solid" on:click={() => openSend()}>Send</Button>
		<Button type="outline" on:click={openTutorial}>Tutorial</Button>
	</div>

	<hr />

	<Messages {cfg} session={$session} bind:messages={$session.notices} />
{/if}

<style lang="postcss">
	.buttons {
		@apply flex justify-center gap-2;
	}
	hr {
		@apply mx-4 my-8 border-t border-primary;
	}
</style>
