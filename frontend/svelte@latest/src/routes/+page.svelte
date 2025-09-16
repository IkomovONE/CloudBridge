<script lang="ts">
    const categories = ['Phones', 'Laptops', 'Monitors', 'Components', 'Smartwatches'];
    const deals = [
        {
            title: 'Samsung Galaxy S24',
            price: '799€',
            store: 'Verkkokauppa',
            image: 'https://via.placeholder.com/150'
        },
        {
            title: 'Apple MacBook Air M2',
            price: '1199€',
            store: 'Gigantti',
            image: 'https://via.placeholder.com/150'
        },
        {
            title: 'LG UltraGear 27"',
            price: '299€',
            store: 'Power',
            image: 'https://via.placeholder.com/150'
        },
        {
            title: 'Kingston 1TB SSD',
            price: '69€',
            store: 'Jimms',
            image: 'https://via.placeholder.com/150'
        }
    ];

    let gridRef: HTMLDivElement | null = null;
    let search = '';

    function scrollToGrid() {
        gridRef?.scrollIntoView({ behavior: 'smooth' });
    }
</script>

<style>
    @keyframes bounce {
        0%, 100% { transform: translateY(0);}
        50% { transform: translateY(10px);}
    }
    /* Remove .svg-container and .big-svg styles, use Tailwind for sizing */
    img.hero-svg {
      width: 40vw;
      height: 40vh;
      object-fit: contain;
      display: block;
    }
</style>

<!-- Gradient background for the whole page -->
<div class="min-h-screen w-full" style="background: linear-gradient(90deg, #e0e7ff 0%, #f0fdfa 100%);">

    <!-- Header with logo, px-6 for padding -->
    <header class="fixed top-0 left-0 z-50 px-6 py-4 bg-white shadow w-full">
        <div class="text-2xl font-extrabold text-blue-700 tracking-tight select-none">CloudBridge</div>
    </header>

    <!-- Main content centered, px-6 for padding, pt-20 for header space -->
    <div class="max-w-8xl mx-auto px-6 py-8 min-h-screen flex flex-col justify-between pt-20">
        <div class="flex flex-col md:flex-row items-start md:items-center">
            <!-- Left: Text content -->
            <div class="flex-1 px-6 py-30">
                <h1 class="text-6xl font-bold text-left mb-2">Finland's Best Tech Deals</h1>
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
                            on:click={scrollToGrid}
                        >
                            {category}
                        </button>
                    {/each}
                </div>
            </div>

            <!-- Right: SVG image, only on md+ screens -->
            <div class="hidden md:flex flex-1 justify-center items-center">
                <img src="/bg.svg" alt="CloudBridge Illustration" class="hero-svg" />
            </div>
        </div>

        <!-- Animated "browse for latest deals" text with arrow" at the bottom and centered -->
        <div class="flex-1 flex flex-col justify-end">
            <div class="flex flex-col items-center mb-8 cursor-pointer select-none" on:click={scrollToGrid}>
                <span class="text-lg font-semibold animate-pulse">Browse for latest deals</span>
                <svg class="w-8 h-8 mt-2 text-blue-600" style="animation: bounce 1.2s infinite;" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
                </svg>
            </div>
        </div>
    </div>

    <!-- Deals grid, initially below the fold -->
    <div bind:this={gridRef} class="max-w-5xl mx-auto px-6 py-30 min-h-screen">
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {#each deals as deal}
                <div class="bg-white rounded shadow p-4 flex flex-col items-center">
                    <img src={deal.image} alt={deal.title} class="w-24 h-24 object-cover mb-3 rounded" />
                    <h2 class="text-lg font-semibold mb-1 text-center">{deal.title}</h2>
                    <p class="text-blue-600 font-bold mb-1">{deal.price}</p>
                    <p class="text-gray-500 text-sm">{deal.store}</p>
                </div>
            {/each}
        </div>
    </div>

   
</div>





