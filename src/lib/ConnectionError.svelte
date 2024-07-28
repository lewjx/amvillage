<script lang="ts">
	import { _ } from 'svelte-i18n';
	import PauseOverlay from './PauseOverlay.svelte';
</script>

<PauseOverlay icon="⚠️" text={$_('error.connection')} />
