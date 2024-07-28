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
	import AdminView from './AdminView.svelte';
	import AdminSelectGroup from '$lib/dashboard/AdminSelectGroup.svelte';

	let sheet: 'login' | 'send' | 'tutorial' | 'error' | 'select_group' | undefined = 'login';
	let teamID: number;
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
	const openSend = (sendTeamID: number = -1) => {
		sheet = 'send';
		teamID = sendTeamID;
	};
	const openTutorial = () => {
		sheet = 'tutorial';
	};
	const openError = () => {
		sheet = 'error';
	};

	let groupCb: (teamID: number) => void;
	const selectGroup = (cb: (teamID: number) => void) => {
		groupCb = cb;
		sheet = 'select_group';
	};
	const groupCallbackWithClose = (teamID: number) => {
		closeSheet();
		groupCb(teamID);
	};
</script>

{#if sheet === 'login'}
	<BottomSheet><Login onComplete={setWebsocket} /></BottomSheet>
{:else}
	{#if cfg.teams[$session.team_id].admin}
		<AdminView {cfg} {session} {openSend} {selectGroup} />
	{:else}
		<PlayerView {cfg} {session} {openSend} {openTutorial} />
	{/if}
	{#if sheet === 'tutorial'}
		<BottomSheet onClose={closeSheet}>Tutorial</BottomSheet>
	{:else if sheet === 'send'}
		<BottomSheet onClose={closeSheet}>
			<Transfer session={$session} {cfg} teamSelected={teamID} close={closeSheet} />
		</BottomSheet>
	{:else if sheet === 'select_group'}
		<BottomSheet onClose={closeSheet}>
			<AdminSelectGroup {cfg} session={$session} selectTeam={groupCallbackWithClose} />
		</BottomSheet>
	{:else if sheet === 'error'}
		<BottomSheet><ConnectionError /></BottomSheet>
	{/if}
{/if}
