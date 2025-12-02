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
            CloudBridge is a modern deal-aggregator platform designed to help users discover the best 
            electronics offers across Finland — from well-known retailers such as Elisa, DNA, Telia, 
            Jimm’s, Multitronic, Verkkokauppa and more.  
            The goal is simple: make it easier, faster and smarter for people to compare prices and 
            find the right tech at the right moment.
        </p>

        <p class="about-text">
            CloudBridge is built as a full end-to-end cloud application using a clean, scalable, 
            and production-oriented architecture:
        </p>

        <ul class="about-text" style="list-style: disc; padding-left: 20px;">
            <li>Frontend — <strong>SvelteKit + TypeScript</strong> for a fast, reactive user experience</li>
            <li>Backend — <strong>Golang + Gin</strong> API running on AWS EC2</li>
            <li>Database — <strong>AWS DynamoDB</strong> for high-performance, serverless storage</li>
            <li>Authentication — <strong>AWS Cognito</strong> with secure user management</li>
            <li>File hosting — <strong>AWS S3</strong> for lightweight image and asset delivery</li>
            <li>API gateway layer — <strong>AWS API Gateway</strong> providing HTTPS access to backend services</li>
        </ul>

        <p class="about-text">
            While the current version focuses on fast browsing and deal discovery, 
            CloudBridge is designed with extensibility in mind.  
            A Python-based web-scraper (BeautifulSoup) is planned as a future enhancement to 
            automate data collection from multiple retailer sources, further improving freshness and coverage.
        </p>

        <p class="about-text">
            CloudBridge was created as part of a university Full-Stack Development project, 
            demonstrating real-world skills in cloud infrastructure, distributed system design, 
            secure APIs, and modern UI/UX development.  
            The platform integrates multiple AWS services, CI/CD workflows, containerized components, 
            and a fully deployed production pipeline.
        </p>

        <div class="about-highlight">
            Built by Daniil Komov & Umandi Sachini Weerasekara
        </div>
    </div>
</div>