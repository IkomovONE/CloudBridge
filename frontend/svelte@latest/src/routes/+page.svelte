<script lang="ts">
    import { scale, fade, slide, fly } from 'svelte/transition';
    import { onMount, onDestroy } from 'svelte';
    import { writable } from 'svelte/store';
    import { addToast } from '$lib/toastStore';
	import { refreshAll } from '$app/navigation';
	import { resolveRoute } from '$app/paths';
    import RangeSlider from 'svelte-range-slider-pips';
	import { spring } from 'svelte/motion';
	
    
    

    function disableScroll() {
        if (typeof document !== "undefined") {
            document.body.style.overflow = "hidden";
        }
    }

    function enableScroll() {
        if (typeof document !== "undefined") {
            document.body.style.overflow = "";
        }
    }


    $: {
        if (selectedProduct) {
            disableScroll();   // modal opened
        } else {
            enableScroll();    // modal closed
        }
    }

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
    let search = '';
    let suggestions = [];

    $: {
        const q = search.trim().toLowerCase();
        if (!q) {
            suggestions = [];
        } else {
            suggestions = deals
                .filter((d) =>
                    d.title.toLowerCase().includes(q) ||
                    d.store.toLowerCase().includes(q)
                )
                .slice(0, 5);
        }
    }

    $: filteredDeals = deals.filter((deal) => {
        const q = search.trim().toLowerCase();
        if (!q) return true;

        return (
            deal.title.toLowerCase().includes(q) ||
            deal.store.toLowerCase().includes(q) ||
            deal.category?.toLowerCase().includes(q) ||
            (deal.description && deal.description.toLowerCase().includes(q))
        );
    });

    let priceRange = [0, 2000]; // default
    let priceMin = 0;
    let priceMax = 0;

    $: maxPossiblePrice = Math.max(...deals.map(d => Number(d.price) || 0));

    $: if (deals.length > 0 && priceMax === 0) {
    priceMax = maxPossiblePrice;
}

    $: finalDeals = filteredDeals
    // PRICE RANGE
    .filter(d => {
        const p = Number(d.price);
        if (priceMin && p < priceMin) return false;
        if (priceMax && p > priceMax) return false;
        return true;
    })
    // SORTING
    .sort((a, b) => {
        if (sortType === "priceAsc") return Number(a.price) - Number(b.price);
        if (sortType === "priceDesc") return Number(b.price) - Number(a.price);
        if (sortType === "titleAsc") return a.title.localeCompare(b.title);
        if (sortType === "titleDesc") return b.title.localeCompare(a.title);
        if (sortType === "ratingDesc") return (b.rating || 0) - (a.rating || 0);
        return 0;
    });

    onMount(() => {
        if (selectedProduct) {
            document.body.style.overflow = "hidden"; // disable background scrolling
        }

        return () => {
            document.body.style.overflow = "auto"; // re-enable when unmounted
        };
    });

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
                description: item.description ?? '',
                rating: item.rating ?? '',
                special: item.special ?? '',
                color: item.color ?? '',
                link: item.link ?? ''
            }));
        } catch (err: any) {
            loadError = err?.message ?? 'Failed to load products';
            console.error('loadProducts error', err);
        } finally {
            loading = false;
            deals = deals.filter(d => d.id && d.id !== "");
        }
    }

    let gridRef: HTMLDivElement | null = null;
    

    type Deal = {
        id: string;
        title: string;
        price: string;
        description?: string;
        store: string;
        image: string;
        category?: string;
        rating?: number;
        link?: string;
        special?: string;
        color?: string;
    };

    let selectedProduct: Deal | null = null;
    let accountCardSelected: boolean | null = null;
    let profileCardSelected: boolean | null = null;
    let selectedCategory: string | null = null;
    let showSubHeader = false;
    let showSortMenu = false;

    let accountMode: 'login' | 'register' = 'login';
    let email = '';
    let password = '';
    let passwordRepeat = '';
    let nickname = ''; // ADD THIS

    let showPassword = false;

    let favouritesSelected= false;

    let favouriteDealsMemory = [];

    const user = writable<{ email: string; nickname: string; idToken: string } | null>(null);

    let user_id= '';

    let favouriteDeals = [];
    


    let pendingEmail = '';
    let confirmationCode = '';
    let oldPassword = '';
    let newPassword = '';
    let expectingConfirmation = false;
    let sortType = "Default";
   

    // add this state var
    let registerStep: 'form' | 'verify' = 'form'; 
    let profileStep: 'none' | 'password' | 'favs' = 'none';// track which step we're on

    function scrollToGrid() {
        gridRef?.scrollIntoView({ behavior: 'smooth' });
    }

    function openProduct(product: Deal) {
        selectedProduct = product;
    }

    

    function handleClickOutside(e) {
        if (!e.target.closest(".search-box")) {
            suggestions = [];
        }
    }

    onMount(() => {
        document.addEventListener("click", handleClickOutside);
        return () => document.removeEventListener("click", handleClickOutside);
    });

    
    function closeModal() {
        selectedProduct = null;
        accountCardSelected = false;
        profileCardSelected = false;
    }

    let last = 0;
    function throttleScroll() {
        const now = performance.now();
        if (now - last < 100) return;   // 100ms throttle = smooth
        last = now;

        handleScroll();
    }

    function handleScroll() {
        if (!gridRef) return;
        const rect = gridRef?.getBoundingClientRect();
        // Show sub-header if grid is at or above the top of the viewport
        showSubHeader = rect && rect.top <= 200;
        showSortMenu = rect && rect.top <= 200;

    }

    onMount(() => {
        // fetch products and attach scroll listener in the browser
        loadProducts();
        if (typeof window !== 'undefined') {
            window.addEventListener('scroll', throttleScroll);
            handleScroll();
        }
    });
    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('scroll', handleScroll);
        }
    });

    const S3_BASE = 'https://cloudbridge-pictures.s3.amazonaws.com/';

    const S3_BASE_LOGOS = 'https://cloudbridge-store-logos.s3.eu-north-1.amazonaws.com/';

    // build exactly 3 candidate URLs from backend base id like "iphone16"
    function buildImageSlidesFromId(id?: string): string[] {
        if (!id) return ['/bg.svg'];
        const base = String(id).replace(/\.png$/i, '').replace(/_(\d+)$/i, '');
        return [1, 2, 3].map(i => `${S3_BASE}${base}_${i}.png`);
    }

    function buildLogoFromId(store?: string): string {
    if (!store) return "/bg.svg";

    const storeName = store.trim().toLowerCase().replace(/\.png$/i, "");

    return S3_BASE_LOGOS + storeName + ".png";
}




    // reactive slides derived from selectedProduct.image (which is a base id)
    $: imageSlides = selectedProduct ? buildImageSlidesFromId(selectedProduct.image) : [];
    $: storeLogo = selectedProduct ? buildLogoFromId(selectedProduct.store) : "";
    $: if (selectedProduct) currentImageIndex = 0;

    let currentImageIndex = 0;

    function prevImage() {
        if (!imageSlides.length) return;
        currentImageIndex = (currentImageIndex + imageSlides.length - 1) % imageSlides.length;
    }
    function nextImage() {
        if (!imageSlides.length) return;
        currentImageIndex = (currentImageIndex + 1) % imageSlides.length;
    }

    function closeSuggestions() {
        suggestions = [];
    }

    // helper: build thumbnail URL from backend image field (backend returns base id like "iphone16")
    function getThumbUrl(imageField?: string) {
        if (!imageField) return '/bg.svg';

        // if it's a full URL
        if (/^https?:\/\//i.test(imageField)) {
            // already numbered (iphone16_1.png) -> use as-is
            if (/_\d+\.png$/i.test(imageField)) return imageField;
            // ends with .png but not numbered -> strip extension and append _1
            const m = imageField.match(/(.+?)\.png$/i);
            if (m) return `${m[1]}_1.png`;
            return imageField;
        }

        // raw id like "iphone16" or "iphone16.png" or "iphone16_2" -> normalize base and append _1
        const base = String(imageField).replace(/\.png$/i, '').replace(/_(\d+)$/i, '');
        return `${S3_BASE}${base}_1.png`;
    }

    async function handleRegister() {
        if (password !== passwordRepeat) {
            addToast("Passwords do not match", "error");
            return;
        }

        try {
            const res = await fetch('http://localhost:8080/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password, nickname })
            });

            const data = await res.json();
            if (!res.ok) {
                const error = data.error;
                const friendly = error.split(":").pop().trim();
                addToast(friendly || res.statusText, "error");
                return;
            }

            // move to verify step
            pendingEmail = email;
            registerStep = 'verify';
            confirmationCode = '';
        } catch (err: any) {
            
            
            addToast('Register error: ' + err.message, "error");
        }
    }

    async function handleConfirm() {
        try {
            const res = await fetch('http://localhost:8080/confirm', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: pendingEmail, code: confirmationCode })
            });
            const data = await res.json();
            if (!res.ok) {
                const error = data.error;
                const friendly = error.split(":").pop().trim();
                addToast('Confirm failed: '+ friendly || res.statusText, "error");
                return;
            }
           
            addToast('Verification successful — you can now log in.', "success");
            
            // reset and go back to login
            registerStep = 'form';
            accountMode = 'login';
            email = '';
            password = '';
            passwordRepeat = '';
            nickname = '';
            confirmationCode = '';
            pendingEmail = '';
        } catch (err: any) {
            addToast(`Confirm error: ${err.message}`, "error");
        }
    }

    async function handleChangePassword() {
    try {
        const token = localStorage.getItem('accessToken'); // or idToken depending on Cognito setup
        if (!token) {
            addToast("No user token found. Please log in.", "error");
            return;
        }

        const res = await fetch('http://localhost:8080/change-password', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}` // send token in header
            },
            body: JSON.stringify({ 
                old_password: oldPassword, 
                new_password: newPassword
            })
        });

        const data = await res.json();

        if (!res.ok) {
            const error = data.error;
            const friendly = error?.split(":").pop().trim();
            addToast('Change password failed: ' + (friendly || res.statusText), "error");
            return;
        }

        addToast('Password changed successfully — you can now log in.', "success");

        // reset form
        oldPassword = '';
        newPassword = '';

    } catch (err: any) {
        addToast(`Change password error: ${err.message}`, "error");
    }
    }

    async function handleResend() {
        try {
            const res = await fetch('http://localhost:8080/resend-confirm', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: pendingEmail })
            });
            const data = await res.json();
            if (!res.ok) {
                const error = data.error;
                const friendly = error.split(":").pop().trim();
                addToast(friendly || res.statusText, "error");
                return;
            }
            
            addToast('Code resent — check your email.', "info");
        } catch (err: any) {
            addToast(`Resend error: ${err.message}`, "error");
        }
    }

    

    

    async function handleLogin() {
        

        try {
            const res = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            const data = await res.json();
            
            if (!res.ok) {
                const error = data.error;
                const friendly = error.split(":").pop().trim();
                addToast(friendly || res.statusText, "error");
                
                return;
            }

            // store tokens
            localStorage.setItem('idToken', data.id_token);
            localStorage.setItem('accessToken', data.access_token);
            localStorage.setItem('refreshToken', data.refresh_token);

            // decode id_token to get user info
            const decoded = decodeToken(data.id_token);

            
            user.set({
                email: decoded.email,
                nickname: decoded.nickname || decoded.preferred_username,
                idToken: data.id_token
            });

            handleFavourites();

            loadFavouritesMemory();

            



            // reset form and close modal
            accountCardSelected = false;
            email = '';
            password = '';

            return true



            
        } catch (err: any) {
            addToast(`Login error: ${err.message}`, "error");
            return false
        }
    }


    async function handleFavourites() {
        try {
            let decodedID = decodeToken($user.idToken);
            let userId = decodedID.sub;

            const res = await fetch("http://localhost:8080/favourites", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ userId })  // <-- FIXED
            });

            const data = await res.json();
            return data.favouriteProducts || [];

        } catch (err: any) {
            addToast(`Favourites error: ${err.message}`, "error");
            return [];
        }
    }

    async function loadFavourites() {
        const favIds = await handleFavourites();
        favouriteDeals = deals.filter(p => favIds.includes(p.id));
        favouriteDealsMemory = favIds;

        favouriteDeals = favouriteDeals.filter(d => d && d.id && d.id !== "");
        
        
    }

    async function loadFavouritesMemory() {
        const favIds = await handleFavourites();
        favouriteDealsMemory = favIds;
        
    }


    async function addFavourite(deal) {
        try {

            let decodedID = decodeToken($user.idToken);
            let userId = decodedID.sub;

            

            const res = await fetch("http://localhost:8080/addfavourite", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    userId: String(userId),
                    dealId: String(deal.id) // <-- force string
                })
            });

            

            const data = await res.json();

            if (res.ok && (data.status === "added" || data.status === "already_in_favourites")) {
                if (!favouriteDeals.includes(String(deal.id))) { // keep consistent
                    favouriteDeals.push(String(deal.id));
                }
            }
            favouriteDealsMemory = [...favouriteDealsMemory, String(deal.id)];
            addToast("Added to favourites", "success");


        } catch (err) {
            
            console.error("Failed to add favourite:", err);
            addToast("Failed to add favourite: " + err, "error");
        }
    }



    async function removeFavourite(deal) {
        try {

            let decodedID = decodeToken($user.idToken);
            let userId = decodedID.sub;

            

            

            const res = await fetch("http://localhost:8080/removefavourite", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    userId: String(userId),
                    dealId: String(deal.id) // <-- force string
                })
            });

            

            const data = await res.json();

            if (res.ok && (data.status === "removed" || data.status === "not_in_favourites")) {
                if (!favouriteDeals.includes(String(deal.id))) { // keep consistent
                    favouriteDeals.push(String(deal.id));
                }
            }
            
            
            addToast("Removed from favourites", "success");

            removeFromLocal(deal.id);

            favouriteDeals = favouriteDeals
                .filter(d => d && d.id && String(d.id) !== String(deal.id));

            favouriteDealsMemory = favouriteDealsMemory
                .filter(id => id && String(id) !== String(deal.id));


        } catch (err) {
            
            console.error("Failed to remove favourite:", err);
            addToast("Failed to remove favourite: " + err, "error");
        }
    }

    function removeFromLocal(id) {
        favouriteDeals = [...favouriteDeals].filter(d => String(d.id) !== String(id));
        favouriteDealsMemory = [...favouriteDealsMemory].filter(x => String(x) !== String(id));
    }

    // helper: decode JWT payload (no verification needed client-side for display)
    function decodeToken(token: string) {
        try {
            const parts = token.split('.');
            if (parts.length !== 3) throw new Error('Invalid token');
            const payload = JSON.parse(atob(parts[1]));
            return payload;
        } catch (err) {
            console.error('Failed to decode token:', err);
            return {};
        }
    }

    // on mount: check if user is already logged in
    onMount(() => {
        const idToken = localStorage.getItem('idToken');
        if (idToken) {
            const decoded = decodeToken(idToken);
            user.set({
                email: decoded.email,
                nickname: decoded.nickname || decoded.preferred_username,
                idToken
            });
            loadFavourites();
            loadFavouritesMemory();
        }
    });

    // logout handler
    function logout() {
        localStorage.removeItem('idToken');
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        user.set(null);
        user_id= '';
        favouriteDealsMemory = [];
    }

    // ...existing code...
</script>

<style>
    @keyframes bounce {
        0%, 100% { transform: translateY(0);}
        50% { transform: translateY(10px);}
    }
    img.hero-svg {
      width: 40vw;
      height: 50vh;
      z-index: 10;
      object-fit: contain;
      display: block;
      position: relative;
    }
    

    .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    backdrop-filter: blur(2px);
    z-index: 50;
    
}

    .grid-area {
        flex-direction: row;
        align-items: flex-start;
        gap: 32px; 
        
    }

    .grid-shift {
        padding-left: 110px; /* or whatever looks good with your filter width */
        transition: padding 0.25s ease;
    }

    .sorting {
        
        position: sticky;
        top: 200px;
        z-index: 60;
        padding: 16px;                  /* p-4 */
        border-radius: 5px;            /* rounded-xl */
        border: 1px solid #e5e7eb;      /* border-gray-200 */
        background: white;
        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        max-width: 200px;            /* w-56 */
        display: flex;                  /* flex */
        flex-direction: column;         /* flex-col */
        gap: 16px;                      /* gap-4 */
    }


  .modal-content {
    max-height: 90vh;      /* allow it to fit screen */
    overflow: hidden;      /* enable scrolling */
    background: white;
    border-radius: 12px;
    padding-top: 50px;
    padding-left: 1px;
    padding-right: 1px;
    
    display: flex;
    position: relative;
    flex-direction: column;
    gap: 40px;
    width: 90%;
    box-shadow: 0 12px 40px rgba(2,6,23,0.12);
}

.modal-backdrop-profile {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 50;
    
}

.modal-content-profile {
    background: white;
    padding: 1.5rem;
    border-radius: 1rem;
    max-width: 920px;
    width: calc(100% - 2rem);
    display: flex;
    flex-direction: row;
    gap: 1.5rem;
    position: relative;
    box-shadow: 0 12px 40px rgba(2,6,23,0.12);
}


.modal-scroll {
    overflow-y: auto;       
    padding: 24px;
    width: 100%;
    flex: 1;
    border-radius: 20px;
    display: flex;
    gap: 40px;
    flex-direction: column;
    
}




.modal-basic-info {
    overflow-y: auto;       
    padding: 24px;
    width: 100%;
    flex: 1;
    border-radius: 20px;
    display: flex;
    flex-direction: column;
     
}

    .scroll-spacer {
        height: 40px;   /* or 64px, whatever feels good */
        width: 100px;
        
    }
    



  /* larger, fixed-ish image area that preserves aspect ratio */
  .modal-image {
    width: 360px;
    height: 360px;
    object-fit: contain; /* use cover if you prefer cropping */
    border-radius: 0.75rem;
    flex-shrink: 0;
    background: #ffffff;
  }

  /* details column stretches and wraps nicely */
  .modal-details {
    flex-direction: row;
    justify-content: center;
    display: flex;
    flex: 1;
    min-width: 0;
  }

  

  .store-logo {
    width: 100px;
    height: auto;
    object-fit: contain;
    opacity: 0.9;
    margin-top: 6px;    
}

    .store-link {
        color: #3d4658;
        font-weight: 600;
        font-size: 0.9rem;
        text-decoration: none;
    }

    .store-link:hover {
        color: #5376be;
        
    }

  .modal-details-profile {
    display: flex;
    flex-direction: column;
    justify-content: center;
    flex: 1;
    min-width: 0;
  }

  /* smaller screens: stack and make image fill width but limit height */
  @media (max-width: 768px) {
    .modal-content {
      flex-direction: column;
      gap: 1rem;
      max-width: 92%;
      padding: 1rem;
    }
    .modal-image {
      width: 100%;
      height: auto;
      max-height: 100vh;
      object-fit: contain;
    }
  }

    .modal-details h1 {
        margin: 0 0 0.5rem 0;
        font-size: 2rem;
        font-weight: bold;
    }
    .modal-scroll h1 {
        margin: 0 0 0.5rem 0;
        font-size: 2.5rem;
        font-weight: bold;
    }

    .price {
        color: #555;
    }

    .star-rating {
        display: flex;
        align-items: center;
        gap: 4px;
        margin: 8px 0 12px 0;
    }

    .star {
        width: 22px;
        height: 22px;
        fill: #e5e7eb; /* default gray */
    }

    .star.full {
        fill: #facc15; /* yellow */
    }

    .star.half {
        fill: url(#half);
    }

    .star.empty {
        fill: #d1d5db; /* lighter gray */
    }

    .rating-number {
        margin-left: 6px;
        font-size: 0.9rem;
        color: #555;
    }

    :global(.description h2) {
        font-size: 1.5rem;      /* bigger */
        font-weight: 700;
        margin-top: 1rem;
        margin-bottom: 0.25rem;
    }

    :global(.description p) {
        font-size: 0.95rem;
        line-height: 1.5;
        margin-bottom: 0.5rem;
    }

    
    :global(.description ul) {
        margin: 1rem 0;
        padding-left: 1.5rem;
    }

    :global(.description li) {
        margin-bottom: 0.5rem;
    }

    :global(.description strong) {
        font-weight: 600;
        color: #1e293b;
    }

    .subheader-wrapper {
        position: fixed;
        top: 64px;
        width: 100%;
        z-index: 40;
        pointer-events: none;
    }

    .subheader-wrapper > div {
        pointer-events: auto;
    }

    .sorting-wrapper {
        position: fixed;
        top: 200px;
        left: 30px;
        z-index: 60;
    }


    .modal-details .price {
        color: #2563eb;
        font-size: 1.5rem;
        font-weight: 600;
        margin-bottom: 0.5rem;
    }

    .modal-details .color {
        color: #515869;
        font-size: 1.2rem;
        font-weight: 500;
        margin-bottom: 0.5rem;
    }
    .modal-details .store {
        color: #888;
        font-size: 1.2rem;
        margin-bottom: 1rem;
    }
    .close-btn {
        position: absolute;
        top: 8px;
        right: 20px;
        background: transparent;
        border: none;
        font-size: 2.3rem;
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

    /* carousel container */
  .carousel {
    width: 360px;
    height: 360px;
    position: relative;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f7fafc;
    border-radius: 0.75rem;
  }

  .carousel-view {
    width: 100%;
    height: 100%;
    position: relative;
    overflow: hidden;
  }

  .carousel-view img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: contain; /* or cover if you want cropping */
    object-position: center;
    opacity: 0;
    transform: scale(0.98);
    transition: opacity 220ms ease, transform 220ms ease;
    border-radius: 0.5rem;
    background: #fff;
  }

  .carousel-view img.selected {
    opacity: 1;
    transform: scale(1);
    z-index: 2;
  }

  .carousel-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    width: 36px;
    height: 36px;
    border-radius: 999px;
    background: rgba(255,255,255,0.85);
    border: 1px solid rgba(0,0,0,0.06);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    cursor: pointer;
    z-index: 5;
  }
  .carousel-btn.left { left: 8px; }
  .carousel-btn.right { right: 8px; }

  .carousel-btn:hover { background: rgba(255,255,255,0.98); }

  .carousel-indicators {
    position: absolute;
    bottom: 8px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 6px;
    z-index: 6;
  }

  .indicator {
    width: 8px;
    height: 8px;
    border-radius: 999px;
    background: rgba(255,255,255,0.7);
    border: 1px solid rgba(0,0,0,0.06);
    padding: 0;
    cursor: pointer;
  }
  .indicator.active {
    background: #2563eb;
    width: 14px;
    border-radius: 8px;
  }

  /* responsive: stack on small screens */
  @media (max-width: 768px) {
    .modal-content { flex-direction: column; gap: 1rem; max-width: 92%; }
    .carousel { width: 100%; height: auto; max-height: 50vh; }
    .carousel-view img { position: relative; height: 100%; }
    .carousel-btn { display: none; } /* use indicators / swipe on touch */
  }

  

</style>

<!-- Gradient background for the whole page -->
<div class="min-h-screen w-full" style="background: linear-gradient(90deg, #e0e7ff 0%, #f0fdfa 100%);">

    <!-- Main Header -->
    <header class="fixed top-0 left-0 z-50 px-6 py-4 bg-white shadow w-full flex items-center">
        <div class="text-2xl font-extrabold text-blue-700 tracking-tight select-none">CloudBridge</div>

        
        <button
            class="px-3 py-1 fixed left-50 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
        
            on:click={() => {
                

                window.onbeforeunload = () => window.scrollTo(0, 0);
                location.reload();
                
                
            }}
        >
            Home
        </button>

        <button
            class="px-3 py-1 fixed left-70 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
        
            on:click={() => {
                selectedCategory = null;
                scrollToGrid();
                favouritesSelected= false;
                
            }}
        >
            About
        </button>
        

        {#if $user}
            <button
                class="px-3 py-1 fixed right-35 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
            
                on:click={() => {
                    loadFavourites();
                    selectedCategory = null;
                    scrollToGrid();
                    favouritesSelected= true; 
                }}
            >
                ☆Favourites 
            </button>
        {/if}

       

        <p class="px-3 py-1 fixed right-24" style="font-size: large;">|</p>

        <button
            class="px-3 py-1 fixed right-5 text-sm rounded transition border border-blue-200 bg-blue-100 text-blue-700 hover:bg-blue-200 hover:text-blue-900"
            on:click={() => {
                if ($user) {
                    profileCardSelected = !profileCardSelected;  // toggle profile modal
                } else {
                    accountCardSelected = !accountCardSelected;  // toggle login/account modal
                }
            }}
        >
            {$user?.nickname || 'Login / Sign Up'}
        </button>
    </header>

    <!-- Sub-header with categories -->
    

    
    {#if showSubHeader}

        <div class="subheader-wrapper">

            <div transition:fade|local={{ duration: 200 }}
                class:visible={showSubHeader}
                class:hidden={!showSubHeader}>

                {#if favouritesSelected}

                    <div in:fade={{ duration: 180 }} out:fade={{ duration: 140 }} class="fixed left-0 w-full bg-blue-50 border-b border-blue-200 flex items-center justify-center gap-2 px-6 py-2 z-40" style="top:64px; min-height:80px;">
                        <nav class="flex gap-2 flex-wrap justify-center w-full">
                            <span class="text-blue-700 font-semibold">
                                Viewing {$user?.nickname}'s Favourite Deals
                            </span>

                            <button
                                class="ml-4 px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition"
                                on:click={() => {
                                    selectedCategory = null;
                                    scrollToGrid();
                                    favouritesSelected = false;
                                }}
                            >
                                Back to latest deals
                            </button>
                        </nav>
                    </div>
                {:else}
                    <div in:fade={{ duration: 180 }} out:fade={{ duration: 140 }} class="fixed left-0 w-full bg-blue-50 border-b border-blue-200 flex items-center justify-center gap-2 px-6 py-2 z-40" style="top:64px; min-height:32px;">
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
                                    favouritesSelected= false;
                                    search = '';
                                }}
                            >
                                Reset filter
                            </button>
                        </nav>
                    </div>
                {/if}



            </div>

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
                <div class="search-box relative">

                    <div class="mb-2 mt-8">
                        <p class="text-left text-lg font-medium mb-2">Search whatever you wish</p>
                        <div>
                            <input
                                type="text"
                                bind:value={search}
                                on:keydown={async (e) => {
                                    if (e.key === 'Enter' || e.key === ' ') {
                                    suggestions = [];
                                    favouritesSelected = false;
                                    scrollToGrid();}
                                }}
                                placeholder="Search for products, brands, stores..."
                                class="w-full max-w-3xl text-xl px-6 py-4 rounded-lg border-2 border-blue-400 focus:outline-none focus:border-blue-600 shadow"
                            />
                        </div>
                    </div>



                </div>
                
                {#if suggestions.length > 0}
                    <div class="absolute mt-0 w-full max-w-2xl bg-white shadow rounded-lg border">
                        {#each suggestions as s}
                            <div
                                class="px-4 py-2 hover:bg-blue-100 cursor-pointer flex gap-3 items-center"
                                on:click={() => {
                                    search = s.title;
                                    suggestions = [];
                                    openProduct(s);
                                }}
                            >
                                <img src={getThumbUrl(s.image)} class="w-10 h-10 rounded object-contain" />
                                <div>
                                    <p class="font-medium">{s.title}</p>
                                    <p class="text-sm text-gray-500">{s.store}</p>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
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
                                favouritesSelected= false;
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
                on:click={() => {scrollToGrid(); favouritesSelected= false;}}
                aria-label="Browse for latest deals"
            >
                <span class="text-lg font-semibold animate-pulse">Browse for latest deals</span>
                <svg class="w-8 h-8 text-blue-600" style="animation: bounce 1.2s infinite;" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
                </svg>
            </button>
        </div>
    </div>

    <div class="grid-area">

        {#if showSortMenu && !selectedProduct}

            <div class="sorting-wrapper">
                <div transition:fade|local={{ duration: 200 }}
                    class:visible={showSortMenu}
                    class:hidden={!showSortMenu}>


                    <div
                        class="sorting"
                        
                        style="height: auto;"
                    >
                        <h3 class="text-lg font-bold text-gray-800 mb-2">Filters</h3>

                        <!-- SORT -->
                        <div class="flex flex-col gap-1">
                            <label class="text-sm font-semibold text-gray-700">Sort by</label>
                            <select
                                bind:value={sortType}
                                class="border rounded p-2"
                                
                            >
                                <option value="Default" >Default</option>
                                <option value="priceAsc">Price: Low → High</option>
                                <option value="priceDesc">Price: High → Low</option>
                                <option value="titleAsc">Title A → Z</option>
                                <option value="titleDesc">Title Z → A</option>
                                <option value="ratingDesc">Rating: High → Low</option>
                            </select>
                        </div>

                        <!-- PRICE RANGE -->
                        <div class="flex flex-col gap-1">
                            <label class="text-sm font-semibold text-gray-700">Price range</label>
                            <RangeSlider
                                min={0}
                                max={2000}
                                step={1}
                                values={priceRange}
                                on:change={(e) => {
                                    priceRange = e.detail.values;
                                    priceMin = priceRange[0];
                                    priceMax = priceRange[1];
                                }}
                                spring={false}

                               
                                
                                

                                  


                            />

                            <div class="flex justify-between text-sm text-gray-600 mt-0">
                                <span>{priceMin} €</span>
                                <span>{priceMax} €</span>
                            </div>
                        </div>

                        <button
                            class="mt-2 bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
                            on:click={() => {
                                priceMin = 0;
                                priceMax = 0;
                                sortType = "Default";
                            }}
                        >
                            Reset filters
                        </button>
                    </div>
                </div>
            </div>

        {/if}

        

        <div bind:this={gridRef} 
        class="max-w-6xl mx-auto px-0 py-8 min-h-[80vh] grid-shift" 
        style="scroll-margin-top: 200px;"> 

        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">

            

            {#if favouritesSelected}

                

                

        

                {#if favouriteDeals.length === 0}
                    <p class="text-center col-span-full text-gray-500">
                        You have no favourite deals yet.
                    </p>
                {/if}

                {#each favouriteDeals as deal}

                    <div transition:fade class="relative bg-white rounded shadow p-4 flex flex-col items-center hover:shadow-lg transition cursor-pointer">

                        <button
                        class="absolute top-2 right-2 text-black hover:text-yellow-500 focus:outline-none"
                        on:click|stopPropagation={() => {removeFavourite(deal); removeFromLocal(deal.id);}}
                        aria-label="Add to favourites"
                        >
                        
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                                <path d="M12 .587l3.668 7.431L24 9.753l-6 5.853 1.42 8.294L12 18.896 4.58 23.9 6 15.606 0 9.753l8.332-1.735z"/>
                            </svg>
                        </button>
                    
                        <a href={`/product/${deal.id}`} 
                        class="flex flex-col items-center"
                        on:click|preventDefault={() => openProduct(deal)}>

                            <img
                                src={getThumbUrl(deal.image)}
                                alt={deal.title}
                                class="w-24 h-24 object-contain mb-3 rounded"
                                on:error={(e) => (e.currentTarget.src = '/bg.svg')}
                            />
                            <h2 class="text-lg font-semibold mb-1 text-center">{deal.title}</h2>
                            <p class="text-blue-600 font-bold mb-1">{deal.price}€</p>
                            <p class="text-gray-500 text-sm">{deal.store}</p>
                        </a>

                    </div>

                    
                {/each}

            {:else}

               

                {#each finalDeals.filter(deal => !selectedCategory || deal.category === selectedCategory) as deal}
                    <div class="relative bg-white rounded shadow p-4 flex flex-col items-center hover:shadow-lg transition cursor-pointer">

                        <!-- Star button -->
                        {#if $user}
                            
                            {#if favouriteDealsMemory.includes(deal.id)}

                                <button
                                class="absolute top-2 right-2 text-black hover:text-yellow-500 focus:outline-none"
                                on:click|stopPropagation={() => {removeFavourite(deal); favouriteDealsMemory = favouriteDealsMemory.filter(id => id !== String(deal.id));}}
                                aria-label="Add to favourites"
                                >
                                <!-- Filled star -->
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M12 .587l3.668 7.431L24 9.753l-6 5.853 1.42 8.294L12 18.896 4.58 23.9 6 15.606 0 9.753l8.332-1.735z"/>
                                </svg>
                                </button>
                            {:else}

                                <button
                                class="absolute top-2 right-2 text-black hover:text-yellow-500 focus:outline-none"
                                on:click|stopPropagation={() => {addFavourite(deal)}}
                                aria-label="Add to favourites"
                                >
                                
                                <!-- Empty star -->
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 .587l3.668 7.431L24 9.753l-6 5.853 1.42 8.294L12 18.896 4.58 23.9 6 15.606 0 9.753l8.332-1.735z"/>
                                </svg>
                                </button>
                            {/if}
                           
                        {/if}
                        {#if !$user}
                            <button
                                class="absolute top-2 right-2 text-black hover:text-yellow-500 focus:outline-none"
                                on:click|stopPropagation={() => {
                                    accountCardSelected = true;
                                }}
                                aria-label="Login to add favourites"
                            >
                                <!-- Empty star -->
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 .587l3.668 7.431L24 9.753l-6 5.853 1.42 8.294L12 18.896 4.58 23.9 6 15.606 0 9.753l8.332-1.735z"/>
                                </svg>
                            </button>
                        {/if}

                        <a href={`/product/${deal.id}`} 
                            class="flex flex-col items-center"
                            on:click|preventDefault={() => openProduct(deal)}
                        >
                            <img
                                src={getThumbUrl(deal.image)}
                                alt={deal.title}
                                class="w-24 h-24 object-contain mb-3 rounded"
                                on:error={(e) => (e.currentTarget.src = '/bg.svg')}
                            />
                            <h2 class="text-lg font-semibold mb-1 text-center">{deal.title}</h2>
                            <p class="text-blue-600 font-bold mb-1">{deal.price}€</p>
                            <p class="text-gray-500 text-sm">{deal.store}</p>
                        </a>
                    </div>
                {/each}
            {/if}

        </div>
    </div>



    </div>

    <!-- Deals grid -->
        

    {#if selectedProduct}
        <div class="modal-backdrop" on:click={closeModal}>
            <div class="modal-content" on:click|stopPropagation transition:scale={{ duration: 250, start: 0.8 }}>
                
                <!-- CLOSE BUTTON -->
                <button class="close-btn" on:click={closeModal}>×</button>

                <!-- SCROLLABLE CONTAINER -->
                <div class="modal-scroll">

                   
                    

                    <!-- DETAILS -->
                    <div class="modal-details">

                        <div class="carousel">
                            <button class="carousel-btn left" on:click={prevImage}>‹</button>

                            <div class="carousel-view">
                                {#each imageSlides as src, idx}
                                    <img 
                                        src={src}
                                        class:selected={idx === currentImageIndex}
                                        on:error={(e) => (e.currentTarget.src = '/bg.svg')}
                                    />
                                {/each}
                            </div>

                            <button class="carousel-btn right" on:click={nextImage}>›</button>

                            <div class="carousel-indicators">
                                {#each imageSlides as _, idx}
                                    <button 
                                    class="indicator {idx === currentImageIndex ? 'active' : ''}"
                                    on:click={() => currentImageIndex = idx}
                                    />
                                {/each}
                            </div>
                        </div>

                        <div class="scroll-spacer"></div>

                        <div class="modal-basic-info">

                            <h1>{selectedProduct.title}</h1>

                            <p class="color">{selectedProduct.color} - {selectedProduct.special}</p>
                            
                            <div class="star-rating">
                                {#each Array(5) as _, i}
                                    {#if selectedProduct.rating >= i + 1}
                                        <!-- full star -->
                                        <svg class="star full" viewBox="0 0 24 24">
                                            <path d="M12 .587l3.668 7.568 8.332 1.151-6.064 5.79 1.52 8.258L12 18.897l-7.456 4.457 1.52-8.257L0 9.306l8.332-1.151L12 .587z"/>
                                        </svg>
                                    {:else if selectedProduct.rating > i && selectedProduct.rating < i + 1}
                                        <!-- half star -->
                                        <svg class="star half" viewBox="0 0 24 24">
                                            <defs>
                                                <linearGradient id="half">
                                                    <stop offset="50%" stop-color="#facc15"/>
                                                    <stop offset="50%" stop-color="#e5e7eb"/>
                                                </linearGradient>
                                            </defs>
                                            <path fill="url(#half)" d="M12 .587l3.668 7.568 8.332 1.151-6.064 5.79 1.52 8.258L12 18.897l-7.456 4.457 1.52-8.257L0 9.306l8.332-1.151L12 .587z"/>
                                        </svg>
                                    {:else}
                                        <!-- empty star -->
                                        <svg class="star empty" viewBox="0 0 24 24">
                                            <path d="M12 .587l3.668 7.568 8.332 1.151-6.064 5.79 1.52 8.258L12 18.897l-7.456 4.457 1.52-8.257L0 9.306l8.332-1.151L12 .587z"/>
                                        </svg>
                                    {/if}
                                {/each}

                                <span class="rating-number">{selectedProduct.rating}</span>
                            </div>
                            <p class="price">{selectedProduct.price}€</p>
                            <img class="store-logo" src={storeLogo}/>
                            <a href={selectedProduct.link} class="store-link" target="_blank" rel="noopener noreferrer">Go to store page</a>



                        </div>

                        
                        
                        
                    </div>

                    <p></p>

                    <h1>About the product</h1>

                    <div class="description mt-2"> {@html selectedProduct.description}</div>

                    
                    

                </div>
                
            </div>
        </div>
    {/if}

    

    {#if !$user}
        {#if accountCardSelected}
            <!-- show login/register modal only if NOT logged in -->
            <div
                class="modal-backdrop-profile"
                role="button"
                tabindex="0"
                aria-label="Close modal"
                on:click={() => { accountCardSelected = false; registerStep = 'form';}}
                on:keydown={async (e) => {
                    // only trigger on Enter or Space
                    if (e.key === 'Enter' || e.key === ' ') {
                        if (accountMode === 'login') {
                            const success = await handleLogin();
                            if (success) {
                                closeModal();
                                addToast('Welcome, ' + $user.nickname + '!', 'success');
                            }
                        } else if (accountMode === 'register' && registerStep !== 'verify') {
                            handleRegister(); // no modal closing here
                        }
                    }
                }}
                
            >
                <div
                    class="modal-content-profile"
                    role="dialog"
                    aria-modal="true"
                    tabindex="0"
                    on:click|stopPropagation
                    transition:scale={{ duration: 250, start: 0.8 }}
                    style="max-width: 400px; flex-direction: column; align-items: stretch;"
                >
                    <!-- FORM STEP (login or register form) -->
                    {#if registerStep === 'form'}
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
                        <div class="modal-details-profile" style="align-items: stretch;">
                            <h1 class="mb-4 text-xl font-bold text-center">{accountMode === 'login' ? 'Login' : 'Register'}</h1>
                            <label class="mb-2 text-sm font-medium">E-mail</label>
                            <input
                                class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                type="email"
                                bind:value={email}
                                placeholder="Enter your e-mail"
                            />

                            {#if accountMode === 'register'}
                                <label class="mb-2 text-sm font-medium">Nickname</label>
                                <input
                                    class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                    type="text"
                                    bind:value={nickname}
                                    placeholder="Choose a nickname"
                                />
                            {/if}

                            <label class="mb-2 text-sm font-medium">Password</label>
                            <input
                                class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                type={showPassword ? "text" : "password"}
                                bind:value={password}
                                placeholder="Enter your password"
                            />
                            {#if accountMode === 'register'}
                                <label class="mb-2 text-sm font-medium">Repeat Password</label>
                                <input
                                    class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                    type={showPassword ? "text" : "password"}
                                    bind:value={passwordRepeat}
                                    placeholder="Repeat your password"
                                />
                            {/if}
                            <button type="button" on:click={() => showPassword = !showPassword}>
                                {showPassword ? "Hide password" : "Show password"}
                            </button>
                            <button
                                class="w-full mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition font-semibold"
                                type="button"
                                on:click={() => { accountMode === 'login' ? handleLogin() : handleRegister(); }}
                            >
                                {accountMode === 'login' ? 'Login' : 'Register'}
                            </button>
                        </div>
                    {/if}

                    <!-- VERIFY STEP (code entry) -->
                    {#if registerStep === 'verify'}
                        <div class="modal-details-profile" style="align-items: stretch;">
                            <h1 class="mb-4 text-xl font-bold text-center">Verify Email</h1>
                            <p class="text-sm text-gray-600 mb-4">We sent a code to <strong>{pendingEmail}</strong></p>
                            <label class="mb-2 text-sm font-medium">Confirmation Code</label>
                            <input
                                class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                type="text"
                                bind:value={confirmationCode}
                                placeholder="Enter the code from your email"
                            />
                            <button
                                class="w-full mb-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition font-semibold"
                                on:click={handleConfirm}
                            >
                                Confirm
                            </button>
                            <button
                                class="w-full mb-2 px-4 py-2 border border-blue-600 text-blue-600 rounded hover:bg-blue-50 transition"
                                on:click={handleResend}
                            >
                                Resend Code
                            </button>
                            <button
                                class="w-full px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50 transition"
                                on:click={() => { registerStep = 'form'; confirmationCode = ''; }}
                                on:keydown={async (e) => {
                                    if (e.key === 'Enter' || e.key === ' ' && registerStep === 'verify') {
                                    handleConfirm();}
                                }}
                            >
                                Back
                            </button>
                        </div>
                    {/if}

                    <button class="close-btn" on:click={() => { accountCardSelected = false; registerStep = 'form'; }}>×</button>
                </div>
            </div>
        {/if}
    {:else}
        {#if profileCardSelected}
        <!-- show profile modal only if logged in -->
            {#if profileStep === 'none'}
                <div
                    class="modal-backdrop-profile fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center"
                    role="button"
                    tabindex="0"
                    aria-label="Close modal"
                    on:click={() => { profileCardSelected = false; }}
                    on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { profileCardSelected = false; } }}
                >
                    <div
                        class="modal-content-profile bg-white rounded-lg p-6 flex flex-col gap-4 shadow-lg"
                        role="dialog"
                        aria-modal="true"
                        tabindex="0"
                        on:click|stopPropagation
                        transition:scale={{ duration: 250, start: 0.8 }}
                        style="max-width: 400px; flex-direction: column; align-items: stretch;"
                    >
                        <!-- Profile Info -->
                        <div class="text-center mb-4">
                            <p class="text-xl font-bold mb-1">{$user.nickname}</p>
                            <p class="text-gray-600 mb-2">{$user.email}</p>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex flex-col gap-2">
                            <button
                                class="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition font-semibold"
                                type="button"
                                on:click={() => {profileStep = 'password'; showPassword = false; oldPassword = ''; newPassword = '';}}
                            >
                                Change Password
                            </button>

                            <button
                                class="w-full px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition font-semibold"
                                type="button"
                                on:click={() => {profileStep = 'none'; handleFavourites(); loadFavourites(); favouritesSelected= true; closeModal(); scrollToGrid();} }
                            >
                                View Favorites
                            </button>

                            
                        </div>

                        <!-- Close Button -->
                        <button
                                class="w-full px-4 py-2 bg-gray-600 text-white rounded hover:bg-red-700 transition font-semibold"
                                type="button"
                                on:click={() => { logout(); addToast('Logged out successfully', 'info'); location.reload();}}
                            >
                                Log out
                        </button>
                    </div>
                </div>
                {/if}
            {#if profileStep === 'password'}
                <div
                    class="modal-backdrop-profile fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center"
                    role="button"
                    tabindex="0"
                    aria-label="Close modal"
                    on:click={() => { profileCardSelected = false; profileStep = 'none'; }}
                    on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { profileCardSelected = false; } }}
                >
                    <div
                        class="modal-content-profile bg-white rounded-lg p-6 flex flex-col gap-4 shadow-lg"
                        role="dialog"
                        aria-modal="true"
                        tabindex="0"
                        on:click|stopPropagation
                        transition:scale={{ duration: 250, start: 0.8 }}
                        style="max-width: 400px; flex-direction: column; align-items: stretch;"
                    >
                       
                        

                        <!-- Action Buttons -->
                        <div class="flex flex-col gap-2">
                            <h1 class="mb-4 text-xl font-bold text-center">Change password for {$user.nickname}</h1>
                            <input
                                class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                type={showPassword ? "text" : "password"}
                                bind:value={oldPassword}
                                placeholder="Enter old password"
                            />
                            
                            <input
                                class="mb-4 px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
                                type={showPassword ? "text" : "password"}
                                bind:value={newPassword}
                                placeholder="Enter new password"
                            />

                            <button type="button" on:click={() => showPassword = !showPassword}>
                                {showPassword ? "Hide passwords" : "Show passwords"}
                            </button>

                            
                            <button
                                class="w-full mb-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition font-semibold"
                                on:click={handleChangePassword}
                            >
                                Confirm
                            </button>

                            

                            
                        </div>

                        <!-- Close Button -->
                        <button
                                class="w-full px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition font-semibold"
                                type="button"
                                on:click={() => { profileStep = 'none'; oldPassword = ''; newPassword = '';}}
                            >
                                Go back
                        </button>
                    </div>
                </div>
                {/if}
            
        {/if} <!-- profileCardSelected -->

        
   
    

    <div class="account-info text-center mt-4">
        <p>Logged in as: <strong>{$user.email}</strong></p>
        <p>Nickname: <strong>{$user.nickname}</strong></p>
    </div>
{/if} <!-- outer if -->
</div>






