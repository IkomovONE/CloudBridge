<script lang="ts">
    import { scale, fade, slide } from 'svelte/transition';
    import { onMount, onDestroy } from 'svelte';

    const categories = [
        'Phones',
        'Laptops',
        'Monitors',
        'Components',
        'Smartwatches',
        'Tablets',
        'Headphones',
        'Keyboards',
        'Mice',
        'Printers',
        'Networking',
        'Storage',
        'Cameras',
        'Drones',
        'Wearables'
    ];

    // products will be fetched from backend
    let deals: Deal[] = [];
    let loading = false;
    let loadError: string | null = null;

    async function loadProducts() {
        loading = true;
        loadError = null;
        try {
            const res = await fetch('http://localhost:8080/products');
            if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
            const data = await res.json();
            // map backend objects to Deal shape defensively
            deals = (Array.isArray(data) ? data : []).map((item: any) => ({
                id: item.id ?? item.ID ?? '',
                title: item.title ?? item.name ?? '',
                price: typeof item.price === 'number' ? String(item.price) : (item.price ?? ''),
                store: item.store ?? '',
                image: item.image ?? '/bg.svg',
                category: item.category ?? '',
                description: item.description ?? ''
            }));
        } catch (err: any) {
            loadError = err?.message ?? 'Failed to load products';
            console.error('loadProducts error', err);
        } finally {
            loading = false;
        }
    }

    let gridRef: HTMLDivElement | null = null;
    let search = '';

    type Deal = {
        id: string;
        title: string;
        price: string;
        description?: string;
        store: string;
        image: string;
    };

    let selectedProduct: Deal | null = null;
    let accountCardSelected: boolean | null = null;
    let selectedCategory: string | null = null;
    let showSubHeader = false;

    let accountMode: 'login' | 'register' = 'login';
    let email = '';
    let password = '';
    let passwordRepeat = '';

    function scrollToGrid() {
        gridRef?.scrollIntoView({ behavior: 'smooth' });
    }

    function openProduct(product: Deal) {
        selectedProduct = product;
    }

    
    function closeModal() {
        selectedProduct = null;
    }

    function handleScroll() {
        if (!gridRef) return;
        const rect = gridRef.getBoundingClientRect();
        // Show sub-header if grid is at or above the top of the viewport
        showSubHeader = rect.top <= 200; // 64px = header height
    }

    onMount(() => {
        // fetch products and attach scroll listener in the browser
        loadProducts();
        if (typeof window !== 'undefined') {
            window.addEventListener('scroll', handleScroll);
            handleScroll();
        }
    });
    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('scroll', handleScroll);
        }
    });
</script>

<style>
    @keyframes bounce {
        0%, 100% { transform: translateY(0);}
        50% { transform: translateY(10px);}
    }
    img.hero-svg {
      width: 40vw;
      height: 40vh;
      object-fit: contain;
      display: block;
    }
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
        position: relative;
        box-shadow: 0 8px 32px rgba(0,0,0,0.18), 0 1.5px 6px rgba(0,0,0,0.12);
    }
    .modal-image {
        width: 180px;
        height: 180px;
        object-fit: cover;
        border-radius: 1rem;
        flex-shrink: 0;
        background: #f3f3f3;
    }
    .modal-details {
        display: flex;
        flex-direction: column;
        justify-content: center;
        flex: 1;
        min-width: 0;
    }
    .modal-details h1 {
        margin: 0 0 0.5rem 0;
        font-size: 1.5rem;
        font-weight: bold;
    }
    .modal-details .price {
        color: #2563eb;
        font-size: 1.2rem;
        font-weight: 600;
        margin-bottom: 0.5rem;
    }
    .modal-details .store {
        color: #888;
        font-size: 1rem;
        margin-bottom: 1rem;
    }
    .close-btn {
        position: absolute;
        top: 1rem;
        right: 1rem;
        background: transparent;
        border: none;
        font-size: 2rem;
        color: #888;
        cursor: pointer;
        transition: color 0.2s;
    }
    .close-btn:hover {
        color: #e53e3e;
    }

    /* New styles for the header and sub-header */
    header {
        height: 64px;
    }
    .sub-header {
        height: 32px;
    }

    button.opacity-60 {
        opacity: 0.6;
        pointer-events: none;
    }

    .account-switch button {
        transition: background 0.2s, color 0.2s;
    }
    .account-switch button.bg-blue-600 {
        background: #2563eb !important;
        color: #fff !important;
    }
    .account-switch button.bg-blue-100 {
        background: #e0e7ff !important;
        color: #2563eb !important;
    }
</style>

<!-- Gradient background for the whole page -->
<div class="min-h-screen w-full" style="background: linear-gradient(90deg, #e0e7ff 0%, #f0fdfa 100%);">

    <!-- Main Header -->
    <header class="fixed top-0 left-0 z-50 px-6 py-4 bg-white shadow w-full flex items-center">
        <div class="text-2xl font-extrabold text-blue-700 tracking-tight select-none">CloudBridge</div>

        <button
            class="px-3 py-1 fixed right-30 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
            
            on:click={() => {
                
                scrollToGrid();
            }}
        >
        ☆Favourites 
        </button>

        <p class="px-3 py-1 fixed right-24" style="font-size: large;">|</p>

        <button
            class="px-3 py-1 fixed right-5 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
            on:click={() => {
                
                accountCardSelected = true;
            }}

            
        >
        Account 
        </button>
    </header>

    <!-- Sub-header with categories -->
    

    <!-- Sub-header visible only when scrolling down -->
    {#if showSubHeader}
        <div transition:slide class="fixed left-0 w-full bg-blue-50 border-b border-blue-200 flex items-center justify-center gap-2 px-6 py-2 z-40" style="top:64px; min-height:32px;">
            <nav class="flex gap-2 flex-wrap justify-center w-full">
                {#each categories as category}
                    <button
                        class="px-3 py-1 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
                        class:opacity-60={selectedCategory === category}
                        class:font-semibold={selectedCategory === category}
                        on:click={() => {
                            selectedCategory = category;
                            scrollToGrid();
                        }}
                    >
                        {category}
                    </button>
                {/each}
                <button
                    class="ml-4 px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition"
                    on:click={() => {
                        selectedCategory = null;
                        scrollToGrid();
                    }}
                >
                    Reset filter
                </button>
            </nav>
        </div>
    {/if}

    <!-- Main content centered, px-6 for padding, pt-20 for header space -->
    <div class="max-w-8xl mx-auto px-6 py-8 min-h-screen flex flex-col justify-between pt-20">
        <div class="flex flex-col md:flex-row items-start md:items-center">
            <!-- Left: Text content -->
            <div class="flex-1 px-6 py-20">
                <h1 class="text-5xl font-bold text-left mb-2">Finland's Best Tech Deals</h1>
                <p class="text-left text-gray-600 mb-20">Find the hottest electronics deals – all in one place!</p>

                <!-- Search section -->
                <div class="mb-2 mt-8">
                    <p class="text-left text-lg font-medium mb-2">Search whatever you wish</p>
                    <div>
                        <input
                            type="text"
                            bind:value={search}
                            placeholder="Search for products, brands, stores..."
                            class="w-full max-w-3xl text-xl px-6 py-4 rounded-lg border-2 border-blue-400 focus:outline-none focus:border-blue-600 shadow"
                        />
                    </div>
                </div>
                <p class="text-left text-gray-500 mt-3 mb-4">or pick the category</p>

                <div class="flex flex-wrap gap-2 mb-2">
                    {#each categories as category}
                        <button
                            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                            class:selected={selectedCategory === category}
                            on:click={() => {
                                selectedCategory = category;
                                scrollToGrid();
                            }}
                        >
                            {category}
                        </button>
                    {/each}
                    {#if selectedCategory}
                        <button
                            class="px-4 py-2 bg-gray-300 text-gray-800 rounded hover:bg-gray-400 transition"
                            on:click={() => {
                                selectedCategory = null;
                                scrollToGrid();
                            }}
                        >
                            Show All
                        </button>
                    {/if}
                </div>
            </div>

            <!-- Right: SVG image, only on md+ screens -->
            <div class="hidden md:flex flex-1 justify-center items-center">
                <img src="/bg.svg" alt="CloudBridge Illustration" class="hero-svg" />
            </div>
        </div>

        <!-- Animated "browse for latest deals" text with arrow" at the bottom and centered -->
        <div class="flex-1 flex flex-col justify-end">
            <button
                type="button"
                class="flex flex-col items-center mb-8 cursor-pointer select-none bg-transparent border-none focus:outline-none"
                on:click={scrollToGrid}
                aria-label="Browse for latest deals"
            >
                <span class="text-lg font-semibold animate-pulse">Browse for latest deals</span>
                <svg class="w-8 h-8 text-blue-600" style="animation: bounce 1.2s infinite;" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
                </svg>
            </button>
        </div>
    </div>

    <!-- Deals grid, initially below the fold -->
    <div bind:this={gridRef} class="max-w-5xl mx-auto px-6 py-8 min-h-[80vh]" style="scroll-margin-top: 200px;"> 
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {#each deals.filter(deal => !selectedCategory || deal.category === selectedCategory) as deal}
                <a href={`/product/${deal.id}`} class="bg-white rounded shadow p-4 flex flex-col items-center hover:shadow-lg transition cursor-pointer" on:click|preventDefault={() => openProduct(deal)}>
                    <img src={deal.image} alt={deal.title} class="w-24 h-24 object-cover mb-3 rounded" />
                    <h2 class="text-lg font-semibold mb-1 text-center">{deal.title}</h2>
                    <p class="text-blue-600 font-bold mb-1">{deal.price}</p>
                    <p class="text-gray-500 text-sm">{deal.store}</p>
                </a>
            {/each}
        </div>
    </div>

    {#if selectedProduct}
        <div
            class="modal-backdrop"
            role="button"
            tabindex="0"
            aria-label="Close modal"
            on:click={closeModal}
            on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeModal(); }}
        >
            <div
                class="modal-content"
                role="dialog"
                aria-modal="true"
                tabindex="0"
                on:click|stopPropagation
                on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeModal(); }}
                transition:scale={{ duration: 250, start: 0.8 }}
            >
                <img class="modal-image" src={selectedProduct.image} alt={selectedProduct.title} />
                <div class="modal-details">
                    <h1>{selectedProduct.title}</h1>
                    <p class="price">{selectedProduct.price}</p>
                    {#if selectedProduct.description}
                        <p>{selectedProduct.description}</p>
                    {/if}
                    <p class="store">{selectedProduct.store}</p>
                    <button class="close-btn" on:click={closeModal}>×</button>
                </div>
            </div>
        </div>
    {/if}

    {#if accountCardSelected}
        <div
            class="modal-backdrop"
            role="button"
            tabindex="0"
            aria-label="Close modal"
            on:click={closeModal}
            on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeModal(); }}
        >
            <div
                class="modal-content"
                role="dialog"
                aria-modal="true"
                tabindex="0"
                on:click|stopPropagation
                on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeModal(); }}
                transition:scale={{ duration: 250, start: 0.8 }}
                style="max-width: 400px; flex-direction: column; align-items: stretch;"
            >
                <div class="account-switch flex justify-center mb-4 gap-2">
                    <button
                        class="px-3 py-1 rounded-t bg-blue-600 text-white font-semibold"
                        class:bg-blue-600={accountMode === 'login'}
                        class:bg-blue-100={accountMode !== 'login'}
                        class:text-blue-700={accountMode !== 'login'}
                        on:click={() => accountMode = 'login'}
                    >
                        Login
                    </button>
                    <button
                        class="px-3 py-1 rounded-t font-semibold"
                        class:bg-blue-600={accountMode === 'register'}
                        class:bg-blue-100={accountMode !== 'register'}
                        class:text-blue-700={accountMode !== 'register'}
                        on:click={() => accountMode = 'register'}
                    >
                        Register
                    </button>
                </div>
                <div class="modal-details" style="align-items: stretch;">
                    <h1 class="mb-4 text-xl font-bold text-center">{accountMode === 'login' ? 'Login' : 'Register'}</h1>
                    <label class="mb-2 text-sm font-medium">E-mail</label>
                    <input
                        class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                        type="email"
                        bind:value={email}
                        placeholder="Enter your e-mail"
                    />
                    <label class="mb-2 text-sm font-medium">Password</label>
                    <input
                        class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                        type="password"
                        bind:value={password}
                        placeholder="Enter your password"
                    />
                    {#if accountMode === 'register'}
                        <label class="mb-2 text-sm font-medium">Repeat Password</label>
                        <input
                            class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                            type="password"
                            bind:value={passwordRepeat}
                            placeholder="Repeat your password"
                        />
                    {/if}
                    <button
                        class="w-full mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition font-semibold"
                    >
                        {accountMode === 'login' ? 'Login' : 'Register'}
                    </button>
                </div>
                <button class="close-btn" on:click={() => {accountCardSelected = false}}>×</button>
            </div>
        </div>
    {/if}
</div>





