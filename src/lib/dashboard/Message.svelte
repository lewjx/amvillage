<script lang="ts">
	import { _ } from 'svelte-i18n';

	import ArchiveIcon from '$lib/icons/ArchiveIcon.svelte';
	import GroupIcon from '$lib/icons/GroupIcon.svelte';
	import HighlightIcon from '$lib/icons/HighlightIcon.svelte';
	import MessageIcon from '$lib/icons/MessageIcon.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';
	import WarningIcon from '$lib/icons/WarningIcon.svelte';
	import type { NoticeType } from '$lib/score';

	export let color: 'highlight' | 'secondary' | 'primary' = 'primary';
	export let icon: NoticeType;
	export let targetTeam: string | null = null;

	export let onArchive: (() => void) | undefined = undefined;

	const tailwindClasses = 'text-highlight text-primary text-secondary';
</script>

<div class={`message text-${color}`}>
	<div class="icon">
		{#if icon === 'pause'}
			<PauseIcon class="size-8" />
		{:else if icon === 'highlight'}
			<HighlightIcon class="size-8" />
		{:else if icon === 'warning'}
			<WarningIcon class="size-8" />
		{:else if icon === 'message'}
			<MessageIcon class="size-8" />
		{/if}
	</div>
	<div class="message-content">
		{#if targetTeam}
			<GroupIcon />
			{$_('dashboard.label.toGroup', { values: { groupName: targetTeam } })}
		{/if}
		<slot />
	</div>
	<div class="archive">
		{#if onArchive}
			<button on:click={onArchive}><ArchiveIcon /></button>
		{/if}
	</div>
</div>

<style lang="postcss">
	.message {
		@apply my-2 grid gap-4;
		grid-template-columns: auto 1fr auto;
	}
	.icon {
		@apply self-start;
	}
	.message-content {
		@apply flex items-center gap-2 self-center text-justify leading-tight;
	}
	.archive {
		@apply self-end;
	}
</style>
