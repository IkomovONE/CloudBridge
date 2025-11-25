from dataclasses import dataclass

@dataclass
class Product:
    id: str
    title: str
    price: str
    store: str
    description: str
    image_urls: list
    category: str