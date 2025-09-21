<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	// Removed invalid $props() usage; Svelte uses <slot /> for children and named slots for modal.
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Background: main page -->
<slot />

<!-- Modal overlay (if active) -->
<slot name="modal" />