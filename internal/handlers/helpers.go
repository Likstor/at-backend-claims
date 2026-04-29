package handlers

import (
	"at-backend-claims/internal/domain"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

const defaultPageSize = 10

func getPageSize(v url.Values) uint64 {
	pageSizeString := v.Get("page_size")

	if pageSizeString == "" {
		return defaultPageSize
	}

	if pageSize, err := strconv.ParseUint(pageSizeString, 10, 64); err == nil {
		return pageSize
	}

	return defaultPageSize
}

func getUserID(v url.Values) (uuid.UUID, bool) {
	userIDString := v.Get("user_id")
	if userIDString == "" {
		return uuid.UUID{}, false
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, false
	}

	return userID, true
}

func sliceClaimsToSliceClaimsForPage(page []domain.Claim) []map[string]any {
	pageResp := make([]map[string]any, 0, len(page))
	for _, claim := range page {
		claimResp := map[string]any{
			"id":          claim.ID,
			"title":       claim.Title,
			"description": claim.Description,
			"category":    claim.Category,
			"status":      claim.Status,
			"latitude":    claim.Latitude,
			"longitude":   claim.Longitude,
			"created_at":  claim.CreatedAt,
			"updated_at":  claim.UpdatedAt,
		}

		pageResp = append(pageResp, claimResp)
	}

	return pageResp
}