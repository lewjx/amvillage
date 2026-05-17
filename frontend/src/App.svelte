<script lang="ts">
	import { onMount } from "svelte"
	import { _, addMessages, init as initI18n, locale } from "svelte-i18n"
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

<div class="popup" class:show={$state.popup && $state.notice === ""}>
	<div>
		<div class="popup-title">📣 {$_("admin.label.notice").split(" ")[0]}</div>
		<div class="popup-content">{@html $state.popup.replace(/\n/g, "<br>")}</div>
	</div>
</div>
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
		@apply fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm opacity-0 pointer-events-none transition-opacity duration-300;
	}
	.popup.show {
		@apply opacity-100 pointer-events-auto;
	}
	.popup > div {
		@apply bg-white rounded-3xl shadow-2xl shadow-slate-900/20 p-8 w-full max-w-sm text-center transform scale-95 transition-transform duration-300 border border-slate-100;
	}
	.popup.show > div {
		@apply scale-100;
	}
	.popup-title {
		@apply text-2xl font-bold text-violet-700 mb-4;
	}
	.popup-content {
		@apply text-lg text-slate-700 font-medium leading-relaxed;
	}
	.main-container {
		@apply relative block flex-grow;
	}
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
