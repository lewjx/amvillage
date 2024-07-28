<script>
	import '../app.css';
	import { addMessages, getLocaleFromNavigator, init } from 'svelte-i18n';
	import en from '../lib/en.json';

	addMessages('en', en);
	init({
		initialLocale: getLocaleFromNavigator(),
		fallbackLocale: 'en'
	});
</script>

<slot></slot>

<style lang="postcss">
	:global(body) {
		@apply min-h-screen bg-background text-primary;
	}
</style>
