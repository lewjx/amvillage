<script lang="ts">
	import { _ } from 'svelte-i18n';
	import type { Message as MessageType } from '$lib/score';
	import { maxHeight } from '$lib/transition';
	import Message from './Message.svelte';

	export let messages: MessageType[];
	const archived = messages.filter((x) => x.archived);
	const unarchived = messages.filter((x) => !x.archived);

	const typeToColor: { [typ: string]: 'highlight' | 'secondary' | 'primary' } = {
		highlight: 'highlight',
		warning: 'secondary',
		chat: 'primary'
	};
	const archive = (i: number) => () => {
		messages[i].archived = true;
	};
</script>

<div class="messages">
	{#each unarchived as message, i}
		<div transition:maxHeight>
			<Message
				color={typeToColor[message.type]}
				icon={message.type}
				onArchive={message.type !== 'highlight' ? archive(i) : undefined}
			>
				{message.content}
			</Message>
		</div>
	{:else}
		<p class="no-message">{$_('dashboard.label.noMessage')}</p>
	{/each}

	<h2>{$_('player.label.pastMessages')}</h2>

	{#each archived as message}
		<div transition:maxHeight>
			<Message icon={message.type}>{message.content}</Message>
		</div>
	{:else}
		{#if unarchived.length === 0}
			<p class="no-message">{$_('dashboard.label.noPastOrPresentMessage')}</p>
		{:else}
			<p class="no-message">{$_('dashboard.label.noPastMessage')}</p>
		{/if}
	{/each}
</div>

<style lang="postcss">
	.messages {
		@apply mx-4;
	}
	div {
		@apply overflow-hidden;
	}
	h2 {
		@apply my-6 text-center text-2xl;
	}
	.no-message {
		@apply text-center text-secondary;
	}
</style>
