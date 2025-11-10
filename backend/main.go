package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// new type that matches the fields used in frontend
type FullProduct struct {
	ID          string `json:"id" dynamodbav:"id"`
	Title       string `json:"title" dynamodbav:"title"`
	Price       string `json:"price" dynamodbav:"price"`
	Store       string `json:"store" dynamodbav:"store"`
	Image       string `json:"image" dynamodbav:"image"`
	Category    string `json:"category" dynamodbav:"category"`
	Description string `json:"description" dynamodbav:"description"`
}

// legacy sample products (kept for other handlers if needed)
var products = []Product{
	{ID: "1", Name: "Laptop X", Price: 999.99},
	{ID: "2", Name: "Smartphone Y", Price: 499.50},
	{ID: "3", Name: "Headphones Z", Price: 89.90},
}

// cachedProducts will be populated once at startup from DynamoDB
var cachedProducts []FullProduct

func main() {
	r := gin.Default()

	// load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	tableName := "Products"

	// fetch items once at startup and cache them
	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalf("Failed to scan table: %v", err)
	}

	fmt.Println("Items found in DynamoDB:", len(out.Items))
	for _, item := range out.Items {
		var fp FullProduct
		if err := attributevalue.UnmarshalMap(item, &fp); err != nil {
			// log and skip malformed item
			log.Printf("failed to unmarshal item: %v", err)
			continue
		}
		// fallback: if Title is empty but Name exists, copy it
		if fp.Title == "" {
			if v, ok := item["title"]; ok {
				_ = v // attributevalue.UnmarshalMap should already handle this; kept for clarity
			}
		}
		cachedProducts = append(cachedProducts, fp)
		fmt.Printf("Loaded product: %s %s\n", fp.ID, fp.Title)
	}

	// CORS middleware - allow your frontend origin(s)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
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

	// Serve cached DynamoDB results here. This endpoint does NOT hit DynamoDB per request.
	r.GET("/products", func(c *gin.Context) {
		c.JSON(200, cachedProducts)
	})

	// other endpoints remain unchanged (they can be adapted later)
	r.GET("/product", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.GET("/favourites", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.POST("/register", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.POST("/login", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.PUT("/addfavourite", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.PUT("/removefavourite", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.Run(":8080") // listen on http://localhost:8080
}
