<script lang="ts">
	import { fly } from "svelte/transition"
	import { _ } from "svelte-i18n"
	import Button from "../components/Button.svelte"
	import { creds, error, state } from "../lib/amvillage"
	import { status } from "../lib/state"

	let username = ""
	let password = location.pathname.split("/").at(-1) || ""
	if (password === "") password = "" // Prevent "undefined" from being shown if no path exists
	const logon = () => {
		$error = ""
		$creds = { username, password }
	}

	$: {
		if ($state.username) {
			$status = {
				status: "mainMenu",
			}
		}
	}
</script>

<main transition:fly={{ y: 500 }}>
	<div class="glass-card">
		<h1>{$_("login.text.welcome")}</h1>
		<div class="input-group">
			<label for="username">{$_("login.label.username")}</label>
			<input
				id="username"
				type="text"
				autocomplete="off"
				bind:value={username}
				placeholder={$_("login.placeholder.username")}
			/>
		</div>
		<div class="input-group">
			<label for="password">{$_("login.label.password")}</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				placeholder={$_("login.placeholder.password")}
			/>
		</div>
		<Button on:click={logon} disabled={!username || !password}>{$_("login.button.start")}</Button>
		<div class="error">{$error}&nbsp;</div>
	</div>
</main>

<style lang="postcss">
	main {
		@apply flex flex-col items-center justify-center w-full h-full p-4;
	}
	.glass-card {
		@apply bg-white border border-slate-100 p-8 rounded-3xl shadow-xl shadow-slate-200/50 flex flex-col items-center gap-6 w-full max-w-md;
	}
	h1 {
		@apply font-bold text-4xl text-center text-slate-900 mb-2;
	}
	.input-group {
		@apply flex flex-col gap-2 w-full;
	}
	label {
		@apply text-slate-600 font-medium ml-1;
	}
	input {
		@apply bg-slate-50 border border-slate-200 p-4 text-lg w-full text-slate-900 placeholder-slate-400 rounded-2xl outline-none transition-all;
	}
	input:focus {
		@apply ring-4 ring-violet-500/10 border-violet-500 bg-white;
	}
	.error {
		@apply text-red-500 font-semibold mt-[-10px];
	}
</style>
