<script>
	import '../app.css';
	import { addMessages, init } from 'svelte-i18n';
	import en from '../lib/en.json';
	import zh from '../lib/zh.json';

	addMessages('en', en);
	addMessages('zh', zh);
	init({ fallbackLocale: 'en' });
</script>

<slot></slot>

<style lang="postcss">
	:global(body) {
		@apply min-h-screen bg-background text-primary;
	}
</style>
