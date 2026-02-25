package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role constants
const (
	RoleDeptAdmin = "DEPT_ADMIN"
	RoleDeptHead  = "DEPT_HEAD"
	RoleCproAdmin = "CPRO_ADMIN"
)

// Gift Type constants
const (
	GiftTypeNormal = "NORMAL"
	GiftTypeVIP    = "VIP"
)

// Application Status constants
const (
	StatusPendingHead = "PENDING_HEAD"
	StatusApprovedHead = "APPROVED_HEAD"
	StatusPendingCpro  = "PENDING_CPRO"
	StatusApprovedCpro = "APPROVED_CPRO"
	StatusRejected     = "REJECTED"
	StatusCompleted    = "COMPLETED"
)

// User represents a system user (Mock)
type User struct {
	ID         string `json:"id" bson:"_id"`
	Username   string `json:"username" bson:"username"`
	Role       string `json:"role" bson:"role"`
	Department string `json:"department" bson:"department"`
}

// GiftItem represents a gift in inventory
type GiftItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl"`
	Type        string             `json:"type" bson:"type"` // NORMAL or VIP
	Quantity    int                `json:"quantity" bson:"quantity"`
}

// Application represents a gift request
type Application struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ApplicantID string             `json:"applicantId" bson:"applicantId"` // User ID
	GiftID      primitive.ObjectID `json:"giftId" bson:"giftId"`
	GiftName    string             `json:"giftName" bson:"giftName"` // Denormalized for display
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	History     []StatusLog        `json:"history" bson:"history"`
}

// StatusLog tracks status changes
type StatusLog struct {
	Status    string    `json:"status" bson:"status"`
	ChangedBy string    `json:"changedBy" bson:"changedBy"` // User ID
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Comment   string    `json:"comment" bson:"comment"`
}
