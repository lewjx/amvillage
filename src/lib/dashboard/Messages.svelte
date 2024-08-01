<script lang="ts">
	import { _ } from 'svelte-i18n';
	import type { Config, Notice } from '$lib/score';
	import type { Session } from '$lib/login/login';
	import { maxHeight } from '$lib/transition';
	import Message from './Message.svelte';

	export let cfg: Config;
	export let session: Session;
	export let messages: Notice[];

	const priority = {
		pause: 3,
		highlight: 2,
		warning: 1,
		message: 0
	};

	$: isAdmin = cfg.teams[session.team_id].admin;
	$: preprocessedMessage = messages
		.map((x) => {
			if (isAdmin && typeof x.team_id === 'number') {
				return { ...x, targetTeam: cfg.teams[x.team_id].name };
			}
			return { ...x, targetTeam: null };
		})
		.map((x) => {
			if ('message' in x) return x;
			return { ...x, message: $_(x.translation_key, { values: x.translation_value }) };
		})
		.sort((x, y) => {
			const xPriority = priority[x.level];
			const yPriority = priority[y.level];
			if (xPriority != yPriority) return yPriority - xPriority;
			return y.timestamp - x.timestamp;
		});
	$: archived = preprocessedMessage.filter((x) => x.dismissed);
	$: unarchived = preprocessedMessage.filter((x) => !x.dismissed);

	const adminOnly = ['pause', 'highlight'];
	const typeToColor: { [typ: string]: 'highlight' | 'secondary' | 'primary' } = {
		pause: 'highlight',
		highlight: 'highlight',
		warning: 'secondary',
		chat: 'primary'
	};

	const archiveCallback = (id: number) => {
		const i = messages.findIndex((x) => x.id === id);
		if (i === -1) {
			console.warn('expected unarchived message to exist while creating callback:', id, messages);
			return undefined;
		}
		const message = messages[i];
		if (adminOnly.includes(message.level)) {
			// Users can only archive if they are admin.
			return isAdmin ? () => adminArchive(id) : undefined;
		}
		return () => archive(id);
	};
	const adminArchive = (id: number) => {
		session.ws?.send(
			JSON.stringify({
				type: 'notice_status_update',
				id: id,
				dismissed: true
			})
		);
	};
	const archive = (id: number) => {
		const i = messages.findIndex((x) => x.id === id);
		if (i === -1) {
			console.warn('expected unarchived message to exist while archiving locally:', id, messages);
			return undefined;
		}
		messages[i].dismissed = true;
	};
</script>

<div class="messages">
	{#each unarchived as message}
		<div transition:maxHeight>
			<Message
				color={typeToColor[message.level]}
				icon={message.level}
				onArchive={archiveCallback(message.id)}
				targetTeam={message.targetTeam}
			>
				{message.message}
			</Message>
		</div>
	{:else}
		<p class="no-message">{$_('dashboard.label.noMessage')}</p>
	{/each}

	<h2>{$_('player.label.pastMessages')}</h2>

	{#each archived as message}
		<div transition:maxHeight>
			<Message icon={message.level} targetTeam={message.targetTeam}>{message.message}</Message>
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
		@apply mx-4 mb-4;
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
