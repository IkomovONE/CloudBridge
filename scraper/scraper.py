import boto3
import requests
import uuid


# ---------------------------
#  DynamoDB Setup
# ---------------------------

dynamodb = boto3.resource("dynamodb", region_name="eu-north-1")
table = dynamodb.Table("Products")


# ---------------------------
#  ELISA DEVELOPER API SCRAPER
# ---------------------------

def scrape_elisa_api():
    """
    Scrapes Elisa products using their official developer API.
    """

    API_URL = "https://api.elisa.fi/v1/products?category=phones"

    try:
        print("Fetching Elisa API data...")
        response = requests.get(API_URL)
        response.raise_for_status()
        data = response.json()

        products = []

        for item in data:
            product = {
                "id": str(uuid.uuid4()),                # Generate UUID
                "title": item.get("name", "Unknown"),   # Elisa API field
                "price": item.get("price", "N/A"),
                "store": "Elisa",
                "description": item.get("description", "No description available."),
                "image": item["images"][0] if item.get("images") else "",
                "carousel": item.get("images", []),
                "category": item.get("category", "Phones")
            }

            products.append(product)

        print(f"Elisa API returned {len(products)} products.")
        return products

    except Exception as e:
        print("Error scraping Elisa API:", e)
        return []


# ---------------------------
#  DNA SCRAPER (placeholder)
# ---------------------------

def scrape_dna():
    """Return empty for now."""
    print("DNA scraper not implemented yet.")
    return []


# ---------------------------
#  TELIA SCRAPER (placeholder)
# ---------------------------

def scrape_telia():
    """Return empty for now."""
    print("Telia scraper not implemented yet.")
    return []


# ---------------------------
#  MERGE ALL SCRAPER RESULTS
# ---------------------------

def run_all_scrapers():
    all_products = []

    for scraper in [scrape_elisa_api, scrape_dna, scrape_telia]:
        try:
            products = scraper()
            all_products.extend(products)
        except Exception as e:
            print(f"Error in scraper {scraper.__name__}: {e}")

    print(f"Total scraped products: {len(all_products)}")
    return all_products


# ---------------------------
#  SAVE TO DYNAMODB
# ---------------------------

def save_to_dynamo(products):
    print("Uploading products to DynamoDB...")

    for p in products:
        try:
            table.put_item(Item=p)
            print("Saved:", p["title"])
        except Exception as e:
            print("Error saving:", p.get("title", "Unknown"), "|", e)

    print("Upload complete.")


# ---------------------------
#  MAIN ENTRYPOINT
# ---------------------------

def main():
    print("Starting scraper...")

    products = run_all_scrapers()

    print("Saving to DynamoDB...")
    save_to_dynamo(products)

    print("Scraping finished!")


if __name__ == "__main__":
    main()