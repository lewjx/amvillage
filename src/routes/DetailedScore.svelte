<script lang="ts">
	import { number } from 'svelte-i18n';
	import { scoreBreakdown, type Config, type Score } from '../lib/score';
	import DownIcon from '$lib/icons/DownIcon.svelte';

	export let cfg: Config;
	export let score: Score;

	$: processedScore = scoreBreakdown(score);
	$: resourceWarning = score.resources.map((resource) => processedScore.minResource === resource);
	$: gemWarning = score.gems.map((gem) => gem <= 0);
</script>

<div class="detailed">
	<div class="col">
		<div class="calculation">{$number(processedScore.minResource)}</div>
		<div class="resources">
			{#each score.resources as resource, i}
				<div class="warning">
					{#if resourceWarning[i]}<DownIcon />{/if}
				</div>
				<div class="resourceName" class:warning={resourceWarning[i]}>{cfg.resource_names[i]}</div>
				<div class="amount" class:warning={resourceWarning[i]}>{$number(resource)}</div>
			{/each}
		</div>
	</div>
	<div class="calculation">x</div>
	<div class="col">
		<div class="calculation">{$number(processedScore.gemType)}</div>
		<div class="resources">
			{#each score.gems as gem, i}
				<div class="warning">
					{#if gemWarning[i]}<DownIcon />{/if}
				</div>
				<div class="resourceName" class:warning={gemWarning[i]}>{cfg.gem_names[i]}</div>
				<div class="amount" class:warning={gemWarning[i]}>{$number(gem)}</div>
			{/each}
		</div>
	</div>
</div>

<style lang="postcss">
	.detailed {
		@apply grid py-4;
		grid-template-columns: 1fr auto 1fr;
	}
	.calculation {
		@apply mb-2 flex justify-around text-2xl;
	}
	.col {
		@apply flex flex-col items-center;
	}
	.resources {
		@apply grid justify-center gap-x-3;
		grid-template-columns: auto 1fr auto;
	}
	.resources div.warning {
		@apply text-secondary;
	}
</style>
