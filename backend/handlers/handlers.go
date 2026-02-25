package handlers

import (
	"context"
	"giftFlow/config"
	"giftFlow/models"
	"math"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Seed Mock Data
func SeedData() {
	coll := config.GetCollection("giftItems")
	count, _ := coll.CountDocuments(context.Background(), bson.M{})
	if count == 0 {
		gifts := []interface{}{
			models.GiftItem{
				Name:        "Notebook Set",
				Description: "Standard university notebook set.",
				ImageURL:    "https://via.placeholder.com/150?text=Notebook",
				Type:        models.GiftTypeNormal,
				Quantity:    100,
			},
			models.GiftItem{
				Name:        "T-Shirt",
				Description: "Cotton T-Shirt with Logo.",
				ImageURL:    "https://via.placeholder.com/150?text=T-Shirt",
				Type:        models.GiftTypeNormal,
				Quantity:    50,
			},
			models.GiftItem{
				Name:        "VIP Hamper",
				Description: "Luxury gift hamper for special guests.",
				ImageURL:    "https://via.placeholder.com/150?text=VIP+Hamper",
				Type:        models.GiftTypeVIP,
				Quantity:    10,
			},
			models.GiftItem{
				Name:        "Crystal Trophy",
				Description: "High-end crystal trophy.",
				ImageURL:    "https://via.placeholder.com/150?text=Trophy",
				Type:        models.GiftTypeVIP,
				Quantity:    5,
			},
		}
		coll.InsertMany(context.Background(), gifts)
	}
}

// GetMe returns current user info
func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}

// GetGifts returns all gift items with pagination
func GetGifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	collection := config.GetCollection("giftItems")

	// Count total documents
	total, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var gifts []models.GiftItem
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &gifts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       gifts,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

// CreateApplication submits a new gift request
func CreateApplication(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	if user.Role != models.RoleDeptAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Dept Admin can apply"})
		return
	}

	var req struct {
		GiftID string `json:"giftId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, _ := primitive.ObjectIDFromHex(req.GiftID)

	// Check and Decrease Stock atomically
	// Only proceed if quantity > 0
	updateResult, err := config.GetCollection("giftItems").UpdateOne(
		context.Background(),
		bson.M{"_id": objID, "quantity": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"quantity": -1}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during stock update"})
		return
	}
	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Out of stock or Gift not found"})
		return
	}

	// Fetch gift details for the application record (snapshot)
	var gift models.GiftItem
	err = config.GetCollection("giftItems").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&gift)
	if err != nil {
		// This should theoretically not happen if UpdateOne succeeded, but handle it just in case
		// Rollback stock if needed, or just log error. For simplicity, we return error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching gift details"})
		return
	}

	app := models.Application{
		ID:          primitive.NewObjectID(),
		ApplicantID: user.ID,
		GiftID:      gift.ID,
		GiftName:    gift.Name,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// GetApplications returns applications relevant to the user with pagination
func GetApplications(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	// Sort by created date desc
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var filter bson.M

	// Role-based filtering
	switch user.Role {
	case models.RoleDeptAdmin:
		// See own applications
		filter = bson.M{"applicantId": user.ID}
	case models.RoleDeptHead:
		// See all applications from their department (mock logic: see all PENDING_HEAD or generic filter)
		// For simplicity in mock, Head sees all PENDING_HEAD or APPROVED_HEAD (waiting for CPRO)
		// In real app, filter by Department. Mock User has "Computer Science".
		// But Applications don't store Department yet. Let's assume Head sees ALL PENDING_HEAD for now for simplicity,
		// or filter by ApplicantID if we had a user store.
		// Let's just return all for Head/CPRO to keep it simple, or filter by status.
		filter = bson.M{}
	case models.RoleCproAdmin:
		// See VIP items pending CPRO approval
		filter = bson.M{}
	}

	collection := config.GetCollection("applications")
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var apps []models.Application
	if err = cursor.All(context.Background(), &apps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       apps,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

// UpdateApplicationStatus handles approvals/rejections
func UpdateApplicationStatus(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	appID := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(appID)

	var req struct {
		Action  string `json:"action"` // APPROVE, REJECT
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var app models.Application
	err := config.GetCollection("applications").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&app)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Fetch gift to check type
	var gift models.GiftItem
	config.GetCollection("giftItems").FindOne(context.Background(), bson.M{"_id": app.GiftID}).Decode(&gift)

	newStatus := app.Status

	if req.Action == "REJECT" {
		// Only restore stock if we are rejecting a non-rejected application
		if app.Status != models.StatusRejected && app.Status != models.StatusCompleted {
			_, err := config.GetCollection("giftItems").UpdateOne(
				context.Background(),
				bson.M{"_id": app.GiftID},
				bson.M{"$inc": bson.M{"quantity": 1}},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore stock"})
				return
			}
		}
		newStatus = models.StatusRejected
	} else if req.Action == "APPROVE" {
		if app.Status == models.StatusPendingHead && user.Role == models.RoleDeptHead {
			if gift.Type == models.GiftTypeVIP {
				newStatus = models.StatusPendingCpro
			} else {
				newStatus = models.StatusCompleted
			}
		} else if app.Status == models.StatusPendingCpro && user.Role == models.RoleCproAdmin {
			newStatus = models.StatusCompleted
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid state transition for this user"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	update := bson.M{
		"$set": bson.M{"status": newStatus},
		"$push": bson.M{
			"history": models.StatusLog{
				Status:    newStatus,
				ChangedBy: user.Username, // Use username for display
				Timestamp: time.Now(),
				Comment:   req.Comment,
			},
		},
	}

	_, err = config.GetCollection("applications").UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": newStatus})
}
