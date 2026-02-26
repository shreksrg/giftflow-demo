package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"giftflow/config"
	"giftflow/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMe returns current user info
func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}

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

// CreateApplication creates a new application
func CreateApplication(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req struct {
		GiftID   string `json:"giftId"`
		Reason   string `json:"reason"`
		Quantity int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	giftID, err := primitive.ObjectIDFromHex(req.GiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Gift ID"})
		return
	}

	// Find and decrement stock atomically
	var gift models.GiftItem
	err = config.GetCollection("giftItems").FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": giftID, "quantity": bson.M{"$gte": req.Quantity}, "isPublished": true},
		bson.M{"$inc": bson.M{"quantity": -req.Quantity}},
	).Decode(&gift)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock, gift not found, or not published"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock: " + err.Error()})
		}
		return
	}

	app := models.Application{
		ID:          primitive.NewObjectID(),
		ApplicantID: user.ID,
		GiftID:      giftID,
		GiftName:    gift.Name,
		ImageURL:    gift.ImageURL,
		Quantity:    req.Quantity,
		Status:      models.StatusPendingHead,
		CreatedAt:   time.Now(),
		History: []models.StatusLog{
			{
				Status:    models.StatusPendingHead,
				ChangedBy: user.ID,
				Timestamp: time.Now(),
				Comment:   "--",
			},
		},
	}

	_, err = config.GetCollection("applications").InsertOne(context.Background(), app)
	if err != nil {
		// Rollback stock decrement
		config.GetCollection("giftItems").UpdateOne(
			context.Background(),
			bson.M{"_id": giftID},
			bson.M{"$inc": bson.M{"quantity": req.Quantity}},
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// GetApplications returns applications based on role
func GetApplications(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	view := c.DefaultQuery("view", "pending") // pending or history

	filter := bson.M{}
	if user.Role == models.RoleDeptAdmin {
		filter["applicantId"] = user.ID

		// Search filters for DEPT_ADMIN
		if status := c.Query("status"); status != "" {
			filter["status"] = status
		}
		if giftName := c.Query("giftName"); giftName != "" {
			filter["giftName"] = bson.M{"$regex": primitive.Regex{Pattern: giftName, Options: "i"}}
		}

		startDateStr := c.Query("startDate")
		endDateStr := c.Query("endDate")
		if startDateStr != "" || endDateStr != "" {
			dateFilter := bson.M{}
			if startDateStr != "" {
				if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
					dateFilter["$gte"] = t
				}
			}
			if endDateStr != "" {
				if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
					// End of the day
					dateFilter["$lte"] = t.Add(24 * time.Hour)
				}
			}
			if len(dateFilter) > 0 {
				filter["createdAt"] = dateFilter
			}
		}

	} else if user.Role == models.RoleDeptHead {
		if view == "history" {
			filter["status"] = bson.M{"$ne": models.StatusPendingHead}
		} else {
			filter["status"] = models.StatusPendingHead
		}
	} else if user.Role == models.RoleCproAdmin {
		if view == "rejected" {
			filter["status"] = models.StatusRejectedCpro
		} else {
			filter["status"] = bson.M{"$in": []string{models.StatusApprovedHead, models.StatusPendingCpro, models.StatusApprovedCpro}}
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	skip := (page - 1) * limit

	collection := config.GetCollection("applications")
	total, _ := collection.CountDocuments(context.Background(), filter)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.M{"createdAt": -1})

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var apps []models.Application
	if err = cursor.All(context.Background(), &apps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       apps,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})
}

// UpdateApplicationStatus updates status
func UpdateApplicationStatus(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var req struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch Application
	var app models.Application
	err := config.GetCollection("applications").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&app)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// 2. Fetch Gift to check Type
	var gift models.GiftItem
	err = config.GetCollection("giftItems").FindOne(context.Background(), bson.M{"_id": app.GiftID}).Decode(&gift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gift details"})
		return
	}

	// 3. Determine New Status
	var newStatus string

	// Handle Rejection
	if req.Status == "REJECTED" || req.Status == models.StatusRejectedHead || req.Status == models.StatusRejectedCpro {
		// Allow rejection if user has authority over current stage
		if user.Role == models.RoleDeptHead && app.Status == models.StatusPendingHead {
			newStatus = models.StatusRejectedHead
		} else if user.Role == models.RoleCproAdmin && (app.Status == models.StatusPendingCpro || app.Status == models.StatusApprovedHead) {
			newStatus = models.StatusRejectedCpro
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized rejection"})
			return
		}
	} else {
		// Handle Approval
		if user.Role == models.RoleDeptHead {
			if app.Status != models.StatusPendingHead {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Application is not pending Department Head approval"})
				return
			}
			// Logic: Normal -> ApprovedCpro (Final); VIP -> PendingCpro
			if gift.Type == models.GiftTypeVIP {
				newStatus = models.StatusPendingCpro
			} else {
				// Normal item: Skip CPRO approval, go straight to Final (ApprovedCpro)
				newStatus = models.StatusApprovedCpro
			}
		} else if user.Role == models.RoleCproAdmin {
			// CPRO can only approve VIP items that are pending CPRO
			if app.Status != models.StatusPendingCpro && app.Status != models.StatusApprovedHead {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Application is not pending CPRO approval"})
				return
			}
			newStatus = models.StatusApprovedCpro
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized approval"})
			return
		}
	}

	update := bson.M{
		"$set": bson.M{"status": newStatus},
		"$push": bson.M{"history": models.StatusLog{
			Status:    newStatus,
			ChangedBy: user.ID,
			Timestamp: time.Now(),
			Comment:   req.Comment,
		}},
	}

	_, err = config.GetCollection("applications").UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Restore stock if application is rejected
	if newStatus == models.StatusRejectedHead || newStatus == models.StatusRejectedCpro {
		restoreQty := app.Quantity
		if restoreQty <= 0 {
			restoreQty = 1 // Default to 1 for older records
		}
		_, err = config.GetCollection("giftItems").UpdateOne(
			context.Background(),
			bson.M{"_id": app.GiftID},
			bson.M{"$inc": bson.M{"quantity": restoreQty}},
		)
		if err != nil {
			// Log error but don't fail the request since status is updated
			// Ideally should use transaction or robust error handling
			c.JSON(http.StatusOK, gin.H{"status": "updated", "newStatus": newStatus, "warning": "Failed to restore stock: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "newStatus": newStatus})
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

// SeedData seeds initial data
func SeedData() {
	count, _ := config.GetCollection("giftItems").CountDocuments(context.Background(), bson.M{})
	if count == 0 {
		gifts := []interface{}{
			models.GiftItem{
				ID:          primitive.NewObjectID(),
				Name:        "Coffee Mug",
				Description: "A nice coffee mug",
				Type:        models.GiftTypeNormal,
				Quantity:    100,
				ImageURL:    "https://via.placeholder.com/150",
			},
			models.GiftItem{
				ID:          primitive.NewObjectID(),
				Name:        "Luxury Pen",
				Description: "Gold plated pen",
				Type:        models.GiftTypeVIP,
				Quantity:    10,
				ImageURL:    "https://via.placeholder.com/150",
			},
		}
		config.GetCollection("giftItems").InsertMany(context.Background(), gifts)
		fmt.Println("Seeded gifts")
	}
}
