<script lang="ts">
	import { toasts, removeToast } from './toastStore';
</script>

<div class="toast-container">
	{#each $toasts as toast (toast.id)}
		<div class="toast {toast.type}">
			<span>{toast.message}</span>
			<button class="close-btn" on:click={() => removeToast(toast.id)}>×</button>
		</div>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		top: 100px;
		right: 20px;
		display: flex;
		flex-direction: column;
		gap: 10px;
		z-index: 9999;
	}
	.toast {
		padding: 12px 16px;
		border-radius: 8px;
		color: white;
		font-weight: 500;
		display: flex;
		align-items: center;
		gap: 10px;
		box-shadow: 0 4px 8px rgba(0,0,0,0.2);
		min-width: 220px;
	}
	.toast.info { background: #6e95be; }
	.toast.success { background: #2e97cc; }
	.toast.error { background: #d16a5f; }

	.close-btn {
		background: transparent;
		border: none;
		color: white;
		font-size: 20px;
		cursor: pointer;
	}
</style>