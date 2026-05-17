<script lang="ts">
	import { onMount } from "svelte"
	import { addMessages, init as initI18n, locale } from "svelte-i18n"
	import { connect, connected, state, type State } from "./lib/amvillage"
	import { status } from "./lib/state"
	import localeEn from "./locale/en.json"
	import localeZh from "./locale/zh.json"
	import Login from "./views/Login.svelte"
	import MainMenu from "./views/MainMenu.svelte"
	import Notice from "./views/Notice.svelte"
	import Trade from "./views/Trade.svelte"

	const views = {
		login: Login,
		mainMenu: MainMenu,
		trade: Trade,
	}

	state.subscribe((state: State) => locale.set(state.config.lang))
	addMessages("en", localeEn)
	addMessages("zh", localeZh)
	initI18n({
		initialLocale: "en",
		fallbackLocale: "en",
		ignoreTag: false,
	})

	onMount(() => connect())

	const handleTouchStart = (e: TouchEvent) => {
		const target = e.target as HTMLElement
		if (target.tagName !== "INPUT" && target.tagName !== "BUTTON") {
			if (document.activeElement?.tagName === "INPUT") {
				(document.activeElement as HTMLElement).blur()
			}
		}
	}

	const handleGlobalClick = () => {
		// Only attempt auto-fullscreen on mobile devices
		if (window.innerWidth <= 1024 && !document.fullscreenElement) {
			document.documentElement.requestFullscreen().catch(() => {})
		}
	}
</script>

<svelte:window on:click={handleGlobalClick} on:touchstart={handleTouchStart} />

<div class="popup" class:show={$state.popup && $state.notice === ""}><div>{$state.popup}</div></div>
<div class="main-container">
	{#if $state.notice === ""}
		<svelte:component this={views[$status.status]} />
	{:else}
		<Notice />
	{/if}
</div>

{#if !$connected && $status.status !== "login"}
	<div class="warning">⚠️</div>
{/if}

<style lang="postcss">
	.popup {
		@apply w-full bg-red-300 max-h-0 overflow-hidden transition-all;
	}
	.popup.show {
		@apply max-h-full;
	}
	.popup > div {
		@apply p-2;
	}
	.main-container {
		@apply relative block flex-grow;
	}
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
