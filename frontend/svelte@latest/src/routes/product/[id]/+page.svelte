<script>
  export let params;

  const product = {
    id: params.id,
    name: `Product ${params.id}`,
    price: "$49.99",
    description: "This is a detailed description of the product...",
    image: "https://via.placeholder.com/300x300.png?text=Product+Image"
  };
</script>

<div class="modal-backdrop">
  <div class="modal-content">
    <img class="product-image" src={product.image} alt={product.name} />
    <div class="product-details">
      <h1>{product.name}</h1>
      <p class="price">{product.price}</p>
      <p>{product.description}</p>
      <a href="/" class="close">Close</a>
    </div>
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