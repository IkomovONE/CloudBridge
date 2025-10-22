package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// new type that matches the fields used in the hardcoded products
type FullProduct struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	Store       string `json:"store"`
	Image       string `json:"image"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

var products = []Product{
	{ID: "1", Name: "Laptop X", Price: 999.99},
	{ID: "2", Name: "Smartphone Y", Price: 499.50},
	{ID: "3", Name: "Headphones Z", Price: 89.90},
}

func main() {
	r := gin.Default()

	// CORS middleware - allow your frontend origin(s)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // add origins you use
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Ping endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Fake products endpoint (use FullProduct and valid composite literals)
	r.GET("/products", func(c *gin.Context) {
		products := []FullProduct{
			{
				ID:          "s25",
				Title:       "Samsung Galaxy S25",
				Price:       "799€",
				Store:       "Verkkokauppa",
				Image:       "/bg.svg",
				Category:    "Phones",
				Description: "The Samsung Galaxy S24 features a stunning 6.2-inch Dynamic AMOLED display, triple-lens camera system, and the latest Snapdragon processor for blazing-fast performance. Enjoy all-day battery life and 5G connectivity in a sleek, modern design.",
			},
			{
				ID:          "mba-m2",
				Title:       "Apple MacBook Air M2",
				Price:       "1199€",
				Store:       "Gigantti",
				Image:       "/bg.svg",
				Category:    "Laptops",
				Description: "Apple’s MacBook Air M2 delivers incredible speed and efficiency with the new Apple M2 chip. The ultra-thin, lightweight design features a 13.6-inch Liquid Retina display, Magic Keyboard, and up to 18 hours of battery life—perfect for work or play on the go.",
			},
			{
				ID:          "lg-ug27",
				Title:       `LG UltraGear 27"`,
				Price:       "299€",
				Store:       "Power",
				Image:       "/bg.svg",
				Category:    "Monitors",
				Description: `The LG UltraGear 27" gaming monitor offers a 144Hz refresh rate, 1ms response time, and vibrant IPS panel for smooth, immersive gameplay. G-SYNC compatibility ensures tear-free visuals, while the ergonomic stand provides optimal viewing comfort.`,
			},
			{
				ID:          "kingston-1tb",
				Title:       "Kingston 1TB SSD",
				Price:       "69€",
				Store:       "Jimms",
				Image:       "/bg.svg",
				Category:    "Storage",
				Description: "Upgrade your storage with the Kingston 1TB SSD. Enjoy lightning-fast read and write speeds, enhanced reliability, and silent operation. Ideal for laptops and desktops, this SSD ensures quick boot times and rapid file transfers.",
			},
			{
				ID:          "apple-watch",
				Title:       "Apple Watch Series 9",
				Price:       "399€",
				Store:       "Apple Store",
				Image:       "/bg.svg",
				Category:    "Smartwatches",
				Description: "Stay connected and track your health with the Apple Watch Series 9. Featuring advanced fitness tracking, ECG, blood oxygen monitoring, and seamless integration with your iPhone, it’s the ultimate companion for an active lifestyle.",
			},
			{
				ID:          "ipad-pro",
				Title:       `iPad Pro 11"`,
				Price:       "999€",
				Store:       "Verkkokauppa",
				Image:       "/bg.svg",
				Category:    "Tablets",
				Description: `The iPad Pro 11" combines the power of the M2 chip with a stunning Liquid Retina display. Perfect for creative professionals and multitaskers, it supports the Apple Pencil and Magic Keyboard for a versatile, laptop-like experience.`,
			},
			{
				ID:          "sony-wh1000xm5",
				Title:       "Sony WH-1000XM5",
				Price:       "349€",
				Store:       "Power",
				Image:       "/bg.svg",
				Category:    "Headphones",
				Description: "Experience industry-leading noise cancellation and superior sound quality with the Sony WH-1000XM5 headphones. Enjoy up to 30 hours of battery life, touch controls, and a comfortable, lightweight design for all-day listening.",
			},
			{
				ID:          "logitech-g915",
				Title:       "Logitech G915 Keyboard",
				Price:       "199€",
				Store:       "Gigantti",
				Image:       "/bg.svg",
				Category:    "Keyboards",
				Description: "The Logitech G915 is a premium wireless mechanical gaming keyboard featuring low-profile GL switches, customizable RGB lighting, and ultra-fast LIGHTSPEED wireless technology. Its sleek, durable aluminum design is perfect for serious gamers.",
			},
		}
		c.JSON(200, products)
	})

	r.GET("/product", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.GET("/favourites", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.POST("/register", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.POST("/login", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.PUT("/addfavourite", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.PUT("/removefavourite", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop X", Price: 999.99},
			{ID: "2", Name: "Smartphone Y", Price: 499.50},
			{ID: "3", Name: "Headphones Z", Price: 89.90},
		}
		c.JSON(200, products)
	})

	r.Run(":8080") // listen on http://localhost:8080
}
