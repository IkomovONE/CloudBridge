<script lang="ts">
  import { user } from '$lib/stores';
  let email = '', password = '';

  async function handleLogin() {
    await user.login(email, password);
  }
</script>

<h2>Login</h2>
<input bind:value={email} placeholder="Email" />
<input type="password" bind:value={password} placeholder="Password" />
<button on:click={handleLogin}>Login</button>
