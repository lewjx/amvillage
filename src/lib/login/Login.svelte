<script lang="ts">
	import { onMount } from 'svelte';
	import { type Readable } from 'svelte/store';
	import { _ } from 'svelte-i18n';

	import Button from '$lib/Button.svelte';
	import { createWebsocket, type Session } from './login';

	export let onComplete: (ws: Readable<Session>) => void;

	let loginName = '';
	let secret = '';
	onMount(() => {
		if (window.location.hash) {
			secret = window.location.hash.substring(1);
		}
	});

	const tryLogin = async () => {
		try {
			// TODO: Disable login button.
			onComplete(await createWebsocket(loginName, secret));
		} catch (err: any) {
			// TODO: Handle error correctly.
			console.error(err);
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
	<Button type="solid" on:click={tryLogin} disabled={loginName.trim() === ''}>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="size-6"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z"
			/>
		</svg>
		{$_('login.button.play')}
	</Button>
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
