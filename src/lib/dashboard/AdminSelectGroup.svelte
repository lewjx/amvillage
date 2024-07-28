<script lang="ts">
	import { _ } from 'svelte-i18n';
	import Button from '$lib/Button.svelte';
	import type { Session } from '$lib/login/login';
	import type { Config } from '$lib/score';

	export let cfg: Config;
	export let session: Session;
	export let selectTeam: (teamID: number) => void;
</script>

<p class="team-buttons-label">{$_('admin.label.chooseBroadcastTarget')}</p>
<div class="team-buttons">
	{#each cfg.teams as team, i}
		{#if i !== session.team_id}
			<Button type="outline" on:click={() => selectTeam(i)}>{team.name}</Button>
		{/if}
	{/each}
</div>

<style lang="postcss">
	.team-buttons-label {
		@apply text-center;
	}
	.team-buttons {
		@apply flex flex-wrap justify-center gap-2;
	}
</style>
