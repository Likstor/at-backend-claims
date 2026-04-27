package handlers

import (
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"net/http"
)

func (h categoriesHandler) getAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.usecase.GetAll(r.Context())
	if err != nil {
		logs.Error(r.Context(), err)

		responses.InternalServerError(r.Context(), w)
		return
	}

	categoriesResp := make([]map[string]any, 0, len(categories))

	for _, cat := range categories {
		subcategories := make([]map[string]any, 0, len(cat.Subcategories))

		for _, subcat := range cat.Subcategories {
			subcategories = append(subcategories, map[string]any{
				"id":   subcat.ID,
				"name": subcat.Name,
			})
		}

		categoriesResp = append(categoriesResp, map[string]any{
			"id":            cat.ID,
			"name":          cat.Name,
			"subcategories": subcategories,
		})
	}

	resp := map[string]any{
		"categories": categoriesResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
