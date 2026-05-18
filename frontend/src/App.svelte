<script lang="ts">
	import { onMount } from "svelte"
	import { _, addMessages, init as initI18n, locale } from "svelte-i18n"
	import { connect, connected, state, ws, type State } from "./lib/amvillage"
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

	let localPopupHiddenFor = ""
	$: if ($state.popup !== localPopupHiddenFor) {
		localPopupHiddenFor = ""
	}
	
	$: isAdmin = $state?.team >= 0 && $state.config?.teams?.[$state.team]?.is_admin
	
	let showNoticePopout = false
	let lastNotice = ""
	let noticeTimeout: any
	
	$: if ($state.notice !== "" && $state.notice !== lastNotice) {
		lastNotice = $state.notice
		if (!isAdmin) {
			showNoticePopout = true
			clearTimeout(noticeTimeout)
			noticeTimeout = setTimeout(() => {
				showNoticePopout = false
			}, 6000)
		}
	} else if ($state.notice === "") {
		lastNotice = ""
		showNoticePopout = false
	}
</script>

<svelte:window on:click={handleGlobalClick} on:touchstart={handleTouchStart} />

{#if $state.notice !== ""}
	<div class="notice-banner">
		<span class="flex-grow font-bold text-center">{$state.notice}</span>
		{#if isAdmin}
			<button class="stop-btn" on:click={() => $ws.send("notice ")}>{$_("notice.button.stop")}</button>
		{/if}
	</div>
{/if}

<div class="popup" class:show={$state.popup && localPopupHiddenFor !== $state.popup && !isAdmin}>
	<div>
		<div class="popup-title">📣 {$_("admin.label.banner").split(" ")[0]}</div>
		<div class="popup-content">{@html $state.popup.replace(/\n/g, "<br>")}</div>
		<button class="close-btn" on:click={() => localPopupHiddenFor = $state.popup}>✕</button>
	</div>
</div>

<div class="popup" class:show={showNoticePopout}>
	<div>
		<div class="popup-title">📣 {$_("admin.label.notice").split(" ")[0]}</div>
		<div class="popup-content">{@html $state.notice.replace(/\n/g, "<br>")}</div>
		<button class="close-btn" on:click={() => showNoticePopout = false}>✕</button>
	</div>
</div>

<div class="main-container">
	<svelte:component this={views[$status.status]} />
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
		@apply relative bg-white rounded-3xl shadow-2xl shadow-slate-900/20 p-8 w-full max-w-sm text-center transform scale-95 transition-transform duration-300 border border-slate-100;
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
	.close-btn {
		@apply absolute top-3 right-3 text-slate-400 transition-colors duration-200 text-xl w-8 h-8 flex items-center justify-center rounded-full;
	}
	.close-btn:hover {
		@apply text-slate-700 bg-slate-100;
	}
	.main-container {
		@apply relative block flex-grow;
	}
	.notice-banner {
		@apply bg-rose-600 text-white flex items-center justify-between px-4 py-3 shadow-md flex-shrink-0 w-full z-10;
	}
	.stop-btn {
		@apply bg-white/20 text-white px-3 py-1 ml-4 rounded-full text-sm font-bold transition-colors whitespace-nowrap;
	}
	.stop-btn:hover {
		@apply bg-white/30;
	}
	.warning {
		@apply absolute bottom-2 right-2 text-5xl animate-bounce;
	}
</style>
