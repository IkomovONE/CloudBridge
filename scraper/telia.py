import requests
from bs4 import BeautifulSoup
import uuid

def scrape_telia():
    url = "https://www.telia.fi/kauppa/puhelimet"
    headers = { "User-Agent": "Mozilla/5.0" }

    html = requests.get(url, headers=headers).text
    soup = BeautifulSoup(html, "html.parser")

    items = soup.select(".product-card")
    products = []

    for item in items:
        title = item.select_one(".product-title")
        price = item.select_one(".product-price")
        img = item.select_one("img")

        if not title or not price:
            continue

        products.append({
            "id": str(uuid.uuid4()),
            "title": title.text.strip(),
            "price": price.text.strip(),
            "image": img["src"] if img else "",
            "description": "Product from Telia scraped automatically.",
            "store": "Telia",
            "category": "Phones"
        })

    return products