package handlers

import (
	"context"
	"giftFlow/config"
	"giftFlow/models"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GetGifts returns all gift items
func GetGifts(c *gin.Context) {
	var gifts []models.GiftItem
	cursor, err := config.GetCollection("giftItems").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())
	cursor.All(context.Background(), &gifts)
	c.JSON(http.StatusOK, gifts)
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
	var gift models.GiftItem
	err := config.GetCollection("giftItems").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&gift)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
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
				Comment:   "Application submitted",
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

// GetApplications returns applications relevant to the user
func GetApplications(c *gin.Context) {
	user := c.MustGet("user").(models.User)
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

	cursor, err := config.GetCollection("applications").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var apps []models.Application
	if err = cursor.All(context.Background(), &apps); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
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
