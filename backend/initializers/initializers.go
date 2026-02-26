package initializers

import (
	"context"
	"fmt"
	"giftflow/config"
	"giftflow/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SeedMockData populates the database with initial data
func SeedMockData() {
	// Seed Users (Optional if we use middleware map, but good for completeness if we switch to DB)
	// Seed Gifts
	SeedData()
}

// SeedData seeds initial data
func SeedData() {
	count, _ := config.GetCollection("giftItems").CountDocuments(context.Background(), bson.M{})
	if count == 0 {
		gifts := []interface{}{
			models.GiftItem{
				ID:          primitive.NewObjectID(),
				Name:        "Notebook computer(VIP)",
				Description: "A high-performance computer",
				Type:        models.GiftTypeVIP,
				IsPublished: "true",
				Quantity:    100,
				ImageURL:    "https://cdn.shopify.com/s/files/1/0743/8857/0275/files/0b6fc75f918d97fd97dffcf79c009c36.png?v=1772093746",
			},
			models.GiftItem{
				ID:          primitive.NewObjectID(),
				Name:        "Christmas gift",
				Description: "Santa Claus parent-child delivery",
				Type:        models.GiftTypeNormal,
				IsPublished: "true",
				Quantity:    39,
				ImageURL:    "https://cdn.shopify.com/s/files/1/0743/8857/0275/files/e3b31d34ec823b3956af70d271290d5c.webp?v=1772093679",
			},
			models.GiftItem{
				ID:          primitive.NewObjectID(),
				Name:        "Desktop robot assistant(VIP)",
				Description: "A highly intelligent silicon-based being",
				Type:        models.GiftTypeVIP,
				IsPublished: "true",
				Quantity:    10,
				ImageURL:    "https://cdn.shopify.com/s/files/1/0743/8857/0275/files/266309dca9faaace1f1bc6a9833bb077.png?v=1772093736",
			},
		}
		config.GetCollection("giftItems").InsertMany(context.Background(), gifts)
		fmt.Println("Seeded gifts")
	}
}
