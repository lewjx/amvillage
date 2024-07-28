<script lang="ts">
	import { _, number } from 'svelte-i18n';

	export let score: number;
</script>

<div class="score">
	<div class="title">{$_('player.label.score')}</div>
	<div class="value">{$number(score)}</div>
</div>

<style lang="postcss">
	.score {
		@apply bg-primary py-6 text-center text-background;
		@apply flex flex-col gap-1;
	}
	.title {
		@apply text-xl;
	}
	.value {
		@apply text-5xl;
	}
</style>
