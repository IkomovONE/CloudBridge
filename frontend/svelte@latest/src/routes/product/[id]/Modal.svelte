<script>
  export let children;
</script>

<div class="modal-backdrop">
  <div class="modal-content">
    <slot />
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 50;
  }
  .modal-content {
    background: white;
    padding: 2rem;
    border-radius: 1rem;
    max-width: 600px;
    width: 90%;
    display: flex;
    flex-direction: row;
    gap: 2rem;
  }
</style>