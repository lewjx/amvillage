<script lang="ts">
	import { _ } from 'svelte-i18n';
	import PauseOverlay from './PauseOverlay.svelte';
	import WarningIcon from './icons/WarningIcon.svelte';
</script>

<PauseOverlay text={$_('error.connection')}>
	<WarningIcon class="size-16 text-secondary" slot="icon" />
</PauseOverlay>
