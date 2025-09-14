package main

import "github.com/gin-gonic/gin"

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func main() {
    r := gin.Default()

    // Ping endpoint
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    // Fake products endpoint
    r.GET("/products", func(c *gin.Context) {
        products := []Product{
            {ID: "1", Name: "Laptop X", Price: 999.99},
            {ID: "2", Name: "Smartphone Y", Price: 499.50},
            {ID: "3", Name: "Headphones Z", Price: 89.90},
        }
        c.JSON(200, products)
    })

    r.Run(":8080") // listen on http://localhost:8080
}