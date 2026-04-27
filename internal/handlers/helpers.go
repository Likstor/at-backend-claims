package handlers

import (
	"at-backend-claims/internal/domain"
	"net/url"
	"strconv"
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