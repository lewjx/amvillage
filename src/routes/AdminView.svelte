<script lang="ts">
	import type { Readable } from 'svelte/store';
	import { _ } from 'svelte-i18n';

	import type { Session } from '$lib/login/login';
	import { type Config } from '$lib/score';
	import ConnectionError from '$lib/ConnectionError.svelte';
	import AdminScore from '$lib/dashboard/AdminScore.svelte';
	import Messages from '$lib/dashboard/Messages.svelte';
	import AdminNotice from '$lib/dashboard/AdminNotice.svelte';

	export let cfg: Config;
	export let session: Readable<Session>;
	export let openSend: (teamID: number) => void;
	export let selectGroup: (cb: (teamID: number) => void) => void;
</script>

<div class="header">{$_('admin.label.header')}</div>
{#if !$session.ws}
	<div class="p-4">
		<ConnectionError />
	</div>
{:else}
	<AdminScore {cfg} session={$session} openTransfer={openSend} />
	<hr />
	<AdminNotice session={$session} {selectGroup} />
	<Messages {cfg} session={$session} messages={$session.notices} />
{/if}

<style lang="postcss">
	.header {
		@apply bg-highlight py-4 text-center text-background;
	}
	hr {
		@apply mx-4 my-6 border-t border-primary;
	}
</style>
