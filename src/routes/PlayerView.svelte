<script lang="ts">
	import type { Readable } from 'svelte/store';
	import Button from '$lib/Button.svelte';
	import ConnectionError from '$lib/ConnectionError.svelte';
	import Messages from '$lib/dashboard/Messages.svelte';
	import Score from '$lib/dashboard/Score.svelte';
	import type { Session } from '$lib/login/login';
	import { type Config, finalScore } from '$lib/score';
	import DetailedScore from './DetailedScore.svelte';

	export let cfg: Config;
	export let session: Readable<Session>;
	export let openSend: () => void;
	export let openTutorial: () => void;

	$: score = $session.score[$session.team_id];
</script>

{#if !$session.ws}
	<div class="p-4">
		<ConnectionError />
	</div>
{:else}
	<Score score={finalScore(score)} />
	<DetailedScore {score} {cfg} />
	<div class="buttons">
		<Button type="solid" on:click={openSend}>Send</Button>
		<Button type="outline" on:click={openTutorial}>Tutorial</Button>
	</div>

	<hr />

	<Messages bind:messages={$session.messages} />
{/if}

<style lang="postcss">
	.buttons {
		@apply flex justify-center gap-2;
	}
	hr {
		@apply mx-4 my-8 border-t border-primary;
	}
</style>
