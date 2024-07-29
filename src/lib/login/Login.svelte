<script lang="ts">
	import { onMount } from 'svelte';
	import { type Readable } from 'svelte/store';
	import { _ } from 'svelte-i18n';

	import Button from '$lib/Button.svelte';
	import { createWebsocket, type Session } from './login';
	import PlayIcon from '$lib/icons/PlayIcon.svelte';

	export let onComplete: (ws: Readable<Session>) => void;

	let loginName = '';
	let secret = '';
	onMount(() => {
		if (window.location.hash) {
			secret = window.location.hash.substring(1);
		}
	});

	let loggingIn = false;
	let currentError = '';
	const tryLogin = async () => {
		try {
			loggingIn = true;
			currentError = '';
			onComplete(await createWebsocket(loginName, secret));
		} catch (err: any) {
			currentError = err.message;
		} finally {
			loggingIn = false;
		}
	};
</script>

<div class="login">
	<p>{$_('login.label.prompt')}</p>
	<div class="inputs">
		<div class="group">
			<label for="name">{$_('login.label.name')}</label>
			<input id="name" bind:value={loginName} />
		</div>
		<div class="group">
			<label for="secret">{$_('login.label.secret')}</label>
			<input id="secret" type="password" bind:value={secret} />
		</div>
	</div>
	<Button type="solid" on:click={tryLogin} disabled={loginName.trim() === '' || loggingIn}>
		<PlayIcon />
		{$_('login.button.play')}
	</Button>
	<div class="text-highlight">{currentError}</div>
</div>

<style lang="postcss">
	.login {
		@apply flex flex-col items-center gap-8;
	}
	.inputs {
		@apply flex flex-col gap-2;
	}
	.group {
		@apply flex flex-col;
	}
	input {
		@apply border border-primary px-2 py-1;
	}
</style>
