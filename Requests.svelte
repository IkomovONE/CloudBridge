<script lang="ts">
  import { apiFetch } from '$lib/api';
  let query = '';
  let response = '';

  async function sendRequest() {
    const res = await apiFetch('/api/request', { query });
    response = JSON.stringify(res, null, 2);
  }
</script>

<h2>Request Data</h2>
<input bind:value={query} placeholder="Enter keyword" />
<button on:click={sendRequest}>Fetch</button>

<pre>{response}</pre>
