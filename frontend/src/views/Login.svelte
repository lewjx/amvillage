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
		@apply bg-white/10 backdrop-blur-xl border border-white/20 p-8 rounded-3xl shadow-2xl flex flex-col items-center gap-6 w-full max-w-md;
	}
	h1 {
		@apply font-bold text-4xl text-center text-white mb-2;
	}
	.input-group {
		@apply flex flex-col gap-2 w-full;
	}
	label {
		@apply text-white/80 font-medium ml-1;
	}
	input {
		@apply bg-black/20 border border-white/20 p-4 text-lg w-full text-white placeholder-white/40 rounded-2xl outline-none transition-all;
	}
	input:focus {
		@apply ring-2 ring-indigo-400/50 border-indigo-400;
	}
	.error {
		@apply text-red-400 font-semibold mt-[-10px];
	}
</style>
