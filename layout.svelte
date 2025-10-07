<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '$lib/stores';
</script>

<nav class="p-4 bg-blue-700 text-white flex justify-between">
  <div>🌩️ CloudBridge</div>
  <div>
    {#if $user}
      <button on:click={() => user.logout()}>Logout</button>
    {:else}
      <a href="/login">Login</a>
    {/if}
  </div>
</nav>

<main class="p-6">
  <slot />
</main>
