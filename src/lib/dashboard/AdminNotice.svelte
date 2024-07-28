<script lang="ts">
	import { _ } from 'svelte-i18n';

	import Button from '$lib/Button.svelte';
	import HighlightIcon from '$lib/icons/HighlightIcon.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';
	import WarningIcon from '$lib/icons/WarningIcon.svelte';
	import MessageIcon from '$lib/icons/MessageIcon.svelte';
	import GlobalBroadcastIcon from '$lib/icons/GlobalBroadcastIcon.svelte';
	import GroupIcon from '$lib/icons/GroupIcon.svelte';
	import type { Session } from '$lib/login/login';

	export let session: Session;
	export let selectGroup: (cb: (teamID: number) => void) => void;

	let noticeMessage = '';
	let messageType: 'pause' | 'highlight' | 'warning' | 'message' = 'message';
	const messageColor = {
		pause: 'highlight',
		highlight: 'highlight',
		warning: 'secondary',
		message: 'primary'
	};
	let currentMessageColor: 'highlight' | 'secondary' | 'primary';
	$: currentMessageColor = messageColor[messageType] as typeof currentMessageColor;
	$: hasMessage = noticeMessage.trim().length > 0;

	$: conflictingPause =
		messageType === 'pause' &&
		session.notices.some(
			(x) => x.level === 'pause' && !x.dismissed && typeof x.team_id !== 'number'
		);
	const send = (teamID: number | null) => {
		const msg = JSON.stringify({
			type: 'notice',
			level: messageType,
			message: noticeMessage,
			team_id: teamID
		});
		if (!session.ws) {
			console.warn('not sending, no connection');
			return;
		}
		session.ws.send(msg);
		noticeMessage = '';
	};
	const groupBroadcast = () => {
		selectGroup((teamID: number) => send(teamID));
	};
</script>

<div class="mb-4 flex flex-col items-center gap-2">
	<div class="notice">
		<textarea
			class="notice-message"
			bind:value={noticeMessage}
			placeholder={$_('admin.label.noticePlaceholder')}
		/>
		<div class="level">
			<div class="text-highlight">
				<input
					type="radio"
					name="message_type"
					id="message_type_pause"
					on:click={() => (messageType = 'pause')}
					checked={messageType === 'pause'}
				/>
				<PauseIcon />
				<label for="message_type_pause">{$_('admin.label.noticePause')}</label>
			</div>
			<div class="text-highlight">
				<input
					type="radio"
					name="message_type"
					id="message_type_highlight"
					on:click={() => (messageType = 'highlight')}
					checked={messageType === 'highlight'}
				/>
				<HighlightIcon />
				<label for="message_type_highlight">{$_('admin.label.noticeHighlight')}</label>
			</div>
			<div class="text-secondary">
				<input
					type="radio"
					name="message_type"
					id="message_type_warning"
					on:click={() => (messageType = 'warning')}
					checked={messageType === 'warning'}
				/>
				<WarningIcon />
				<label for="message_type_warning">{$_('admin.label.noticeWarning')}</label>
			</div>
			<div>
				<input
					type="radio"
					name="message_type"
					id="message_type_message"
					on:click={() => (messageType = 'message')}
					checked={messageType === 'message'}
				/>
				<MessageIcon />
				<label for="message_type_message">{$_('admin.label.noticeMessage')}</label>
			</div>
		</div>
	</div>
	<div class="mx-2 flex flex-grow flex-wrap justify-center gap-2">
		{#if conflictingPause}
			<div class="text-center text-highlight">{$_('admin.label.conflictingPause')}</div>
		{:else}
			<Button
				type={messageType === 'pause' ? 'solid' : 'outline'}
				size="md"
				accent={currentMessageColor}
				disabled={!hasMessage}
				on:click={() => send(null)}
			>
				<GlobalBroadcastIcon />
				{#if messageType === 'pause'}
					{$_('admin.button.broadcastGlobalPause')}
				{:else}
					{$_('admin.button.broadcastGlobal')}
				{/if}
			</Button>
			<Button
				type={messageType === 'pause' ? 'solid' : 'outline'}
				size="md"
				accent={currentMessageColor}
				disabled={!hasMessage}
				on:click={groupBroadcast}
			>
				<GroupIcon />
				{#if messageType === 'pause'}
					{$_('admin.button.broadcastGroupPause')}
				{:else}
					{$_('admin.button.broadcastGroup')}
				{/if}
			</Button>
		{/if}
	</div>
</div>

<style lang="postcss">
	.notice {
		@apply mx-2 border border-primary px-4 py-2;
	}
	textarea {
		@apply w-full bg-background text-lg text-primary;
	}
	.level {
		@apply flex flex-col flex-wrap gap-1 lg:flex-row;
	}
	.level > div {
		@apply flex gap-1;
	}
</style>
