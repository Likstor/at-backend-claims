package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ClaimTitleCharactersMaxCount       = 100
	ClaimDescriptionCharactersMaxCount = 1000
)

type Claim struct {
	ID                uint64
	CreatedBy         uuid.UUID
	Title             string
	Description       string
	Category          string
	Status            ClaimStatus
	Photos            []string
	Latitude          float64
	Longitude         float64
	Feedback          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	StatusUpdatedAt   time.Time
	FeedbackUpdatedAt time.Time
}

type ClaimStatus string

const (
	ClaimStatusPending   ClaimStatus = "pending"   // в обработке
	ClaimStatusAccepted  ClaimStatus = "accepted"  // принята в работу
	ClaimStatusCompleted ClaimStatus = "completed" // выполнена
	ClaimStatusDeclined  ClaimStatus = "declined"  // отклонена
	ClaimStatusUnknown   ClaimStatus = "unknown"   // для обработки некорректного ввода
)

func (cs ClaimStatus) String() string {
	return string(cs)
}

func (cs *ClaimStatus) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "pending":
		*cs = ClaimStatusPending
	case "accepted":
		*cs = ClaimStatusAccepted
	case "completed":
		*cs = ClaimStatusCompleted
	case "declined":
		*cs = ClaimStatusDeclined
	default:
		*cs = ClaimStatusUnknown
	}

	return nil
}
