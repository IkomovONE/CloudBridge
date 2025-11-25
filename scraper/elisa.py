

import requests
from bs4 import BeautifulSoup
from models import Product

def scrape_elisa():
    url = "https://elisa.fi/kauppa/puhelimet"
    res = requests.get(url)
    soup = BeautifulSoup(res.text, "html.parser")

    products = []


    

    url2 = "https://elisa.fi/kauppa/rest/products/catalog"
    data = requests.get(url2).json()

    

    # EXAMPLE selectors — we’ll adjust once you inspect HTML
    for item in soup.select(".product-card"):
        title = item.select_one(".product-title").get_text(strip=True)
        price = item.select_one(".price").get_text(strip=True)
        description = item.select_one(".product-description").get_text(strip=True)
        img = item.select_one("img")["src"]
        pid = item["data-product-id"]

        products.append(
            Product(
                id=pid,
                title=title,
                price=price,
                store="Elisa",
                description=description,
                image_urls=[img],
                category="Phones",
            )
        )

    return products