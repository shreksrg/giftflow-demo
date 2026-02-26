package handlers

import (
	"context"
	"fmt"
	"giftflow/config"
	"giftflow/models"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetGifts returns list of gifts
func GetGifts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	skip := (page - 1) * limit

	// If skip is provided, use it (backward compatibility or flexibility)
	if s, ok := c.GetQuery("skip"); ok {
		skip, _ = strconv.Atoi(s)
	}

	collection := config.GetCollection("giftItems")

	filter := bson.M{}
	if publishedStr, ok := c.GetQuery("published"); ok {
		if published, err := strconv.ParseBool(publishedStr); err == nil {
			filter["isPublished"] = published
		}
	}

	total, _ := collection.CountDocuments(context.Background(), filter)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var gifts []models.GiftItem
	if err = cursor.All(context.Background(), &gifts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       gifts,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})
}

// CreateGift adds a new gift item
func CreateGift(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	if user.Role != models.RoleCproAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only CPRO Admin can manage gifts"})
		return
	}

	var gift models.GiftItem
	if err := c.ShouldBindJSON(&gift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gift.ID = primitive.NewObjectID()
	_, err := config.GetCollection("giftItems").InsertOne(context.Background(), gift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gift)
}

// UpdateGift modifies an existing gift
func UpdateGift(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	if user.Role != models.RoleCproAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only CPRO Admin can manage gifts"})
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(req, "id")
	delete(req, "_id")

	_, err = config.GetCollection("giftItems").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": req},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteGift removes a gift
func DeleteGift(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	if user.Role != models.RoleCproAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only CPRO Admin can manage gifts"})
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	_, err = config.GetCollection("giftItems").DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UploadGiftImage handles image upload
func UploadGiftImage(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	if user.Role != models.RoleCproAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only CPRO Admin can upload images"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
		return
	}

	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	url := fmt.Sprintf("http://127.0.0.1:8080/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": url})
}
