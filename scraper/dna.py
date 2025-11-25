import requests
from bs4 import BeautifulSoup
import uuid

def scrape_dna():
    url = "https://www.dna.fi/puhelimet"
    headers = { "User-Agent": "Mozilla/5.0" }

    html = requests.get(url, headers=headers).text
    soup = BeautifulSoup(html, "html.parser")

    items = soup.select(".product-card")
    products = []

    for item in items:
        title = item.select_one(".product__title")
        price = item.select_one(".product__price")
        img = item.select_one("img")

        if not title or not price:
            continue

        products.append({
            "id": str(uuid.uuid4()),
            "title": title.text.strip(),
            "price": price.text.strip(),
            "image": img["src"] if img else "",
            "description": "Product from DNA scraped automatically.",
            "store": "DNA",
            "category": "Phones"
        })

    return products