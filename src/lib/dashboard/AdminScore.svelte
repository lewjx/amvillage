<script lang="ts">
	import type { Session } from '$lib/login/login';
	import { finalScore, scoreBreakdown, type Config } from '$lib/score';
	import { maxHeight } from '$lib/transition';
	import DetailedScore from '../../routes/DetailedScore.svelte';

	export let cfg: Config;
	export let session: Session;
	export let openTransfer: (teamID: number) => void;

	let expand = cfg.teams.map((_) => false);
	$: scores = session.score
		.map((x, i) => {
			const breakdown = scoreBreakdown(x);
			return {
				...breakdown,
				totalScore: breakdown.minResource * breakdown.gemType,
				teamID: i
			};
		})
		.filter((_, i) => !cfg.teams[i].admin)
		.sort((x, y) => y.totalScore - x.totalScore);

	const toggle = (teamID: number) => (expand[teamID] = !expand[teamID]);
</script>

<div class="scores">
	{#each scores as score, i}
		<p class="ml-2">{i + 1}.</p>
		<p>{cfg.teams[score.teamID].name}</p>
		<p class="score">{score.minResource} x {score.gemType} = {score.totalScore}</p>
		<div class="flex">
			<button on:click={() => openTransfer(score.teamID)}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10"
					/>
				</svg>
			</button>
			<button on:click={() => toggle(score.teamID)}>
				{#if expand[score.teamID]}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="size-6"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5" />
					</svg>
				{:else}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="size-6"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
					</svg>
				{/if}
			</button>
		</div>
		{#if expand[score.teamID]}
			<div class="col-span-4 overflow-hidden" transition:maxHeight>
				<DetailedScore {cfg} score={session.score[score.teamID]} />
			</div>
		{/if}
	{/each}
</div>

<style lang="postcss">
	.scores {
		@apply mx-2 my-4 grid gap-2;
		grid-template-columns: auto auto 1fr auto;
	}
	.score {
		@apply text-right;
	}
</style>
