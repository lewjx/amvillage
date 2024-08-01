<script lang="ts">
	import { _ } from 'svelte-i18n';

	const tutorialMessage = $_('player.label.tutorialContent')
		.split('\n')
		.map((x) => {
			let index = x.indexOf(':');
			if (index === -1) index = x.indexOf('：');
			if (index === -1) return [x];
			return [x.slice(0, index + 1), x.slice(index + 1)];
		});
</script>

<div class="tutorial">
	<h1>{$_('player.label.tutorial')}</h1>
	<div class="content">
		{#each tutorialMessage as msg}
			{#if msg.length > 1}
				<p><b>{msg[0]}</b>{msg[1]}</p>
			{:else}
				<p>{msg[0]}</p>
			{/if}
		{/each}
	</div>
</div>

<style lang="postcss">
	.tutorial {
		@apply flex flex-col items-center gap-4;
	}
	h1 {
		@apply text-center text-xl font-semibold;
	}
	b {
		@apply font-semibold;
	}
	.content {
		@apply flex max-w-md flex-col gap-2;
	}
</style>
