<script lang="ts">
	import { onMount } from 'svelte';
	import { type Readable } from 'svelte/store';

	import BottomSheet from '$lib/BottomSheet.svelte';
	import ConnectionError from '$lib/ConnectionError.svelte';
	import { blankConfig, type Config } from '$lib/score';
	import Transfer from '$lib/dashboard/Transfer.svelte';
	import Login from '$lib/login/Login.svelte';
	import { type Session } from '$lib/login/login';
	import PlayerView from './PlayerView.svelte';

	let sheet: 'login' | 'send' | 'tutorial' | 'error' | undefined = 'login';
	let session: Readable<Session>;
	let cfg = blankConfig;

	onMount(() => {
		fetch('/api/config')
			.then((resp) => resp.json())
			.then((resp) => (cfg = resp as Config))
			.catch((err) => {
				console.error(err);
				openError();
			});
	});

	const setWebsocket = (establishedSession: Readable<Session>) => {
		session = establishedSession;
		closeSheet();
	};
	const closeSheet = () => {
		sheet = undefined;
	};
	const openSend = () => {
		sheet = 'send';
	};
	const openTutorial = () => {
		sheet = 'tutorial';
	};
	const openError = () => {
		sheet = 'error';
	};
</script>

{#if sheet === 'login'}
	<BottomSheet><Login onComplete={setWebsocket} /></BottomSheet>
{:else}
	<PlayerView {cfg} {session} {openSend} {openTutorial} />
	{#if sheet === 'tutorial'}
		<BottomSheet onClose={closeSheet}>Tutorial</BottomSheet>
	{:else if sheet === 'send'}
		<BottomSheet onClose={closeSheet}>
			<Transfer session={$session} {cfg} close={closeSheet} />
		</BottomSheet>
	{:else if sheet === 'error'}
		<BottomSheet><ConnectionError /></BottomSheet>
	{/if}
{/if}
