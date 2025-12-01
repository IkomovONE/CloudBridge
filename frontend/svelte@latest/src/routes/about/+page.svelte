<script lang="ts">
    import { goto } from "$app/navigation";
</script>

<style>
    img.hero-svg {
        width: 40vw;
        height: 50vh;
        object-fit: contain;
        display: block;
        position: relative;
    }

    header {
        height: 64px;
    }

    .about-section {
        max-width: 900px;
        margin: 0 auto;
        padding: 60px 20px;
        text-align: left;
    }

    .about-title {
        font-size: 2.5rem;
        font-weight: 700;
        margin-bottom: 20px;
        color: #1e3a8a;
    }

    .about-text {
        font-size: 1.1rem;
        line-height: 1.7;
        color: #374151;
        margin-bottom: 20px;
    }

    .about-highlight {
        background: #e0e7ff;
        padding: 10px 16px;
        border-radius: 8px;
        display: inline-block;
        margin-top: 20px;
        font-weight: 600;
        color: #1d4ed8;
    }
</style>

<div class="min-h-screen w-full" style="background: linear-gradient(90deg, #e0e7ff 0%, #f0fdfa 100%);">

    <!-- HEADER -->
    <header class="fixed top-0 left-0 z-50 px-6 py-4 bg-white shadow w-full flex items-center">
        <div class="text-2xl font-extrabold text-blue-700 tracking-tight select-none">CloudBridge</div>

        <button
            class="px-3 py-1 fixed left-50 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700"
            on:click={() => goto('/')}
        >
            Home
        </button>

        <button
            class="px-3 py-1 fixed left-70 text-sm rounded transition border border-blue-200 bg-blue-200 text-blue-900"
        >
            About
        </button>
    </header>

    <!-- PAGE CONTENT -->
    <div class="about-section pt-24">
        <h1 class="about-title">About CloudBridge</h1>

        <p class="about-text">
            CloudBridge is a modern deal-aggregator platform that gathers the best electronics
            offers from Finland’s biggest retailers — Elisa, DNA, Telia, Jimm’s, Multitronic,
            Verkkokauppa and many more.
        </p>

        <p class="about-text">
            The platform is built using cutting-edge technologies:
        </p>

        <ul class="about-text" style="list-style: disc; padding-left: 20px;">
            <li>Frontend — <strong>SvelteKit + TypeScript</strong></li>
            <li>Backend — <strong>Golang + Gin API</strong></li>
            <li>Database — <strong>AWS DynamoDB</strong></li>
            <li>Authentication — <strong>AWS Cognito</strong></li>
            <li>Scraper service — <strong>Python + BeautifulSoup</strong></li>
            <li>Storage — <strong>AWS S3 (product images, store logos)</strong></li>
        </ul>

        <p class="about-text">
            CloudBridge was created as part of a Software Systems Architecture university project,
            showcasing cloud infrastructure design, distributed components, secure APIs and modern
            UI/UX principles.
        </p>

        <div class="about-highlight">
            Built by Daniil Komov & Umandi Sachini Weerasekara
        </div>
    </div>
</div>