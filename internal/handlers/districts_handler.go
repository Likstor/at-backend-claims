package handlers

import (
	"at-backend-claims/internal/handlers/middleware"
	"at-backend-claims/internal/pkg/districts"
	"at-backend-claims/internal/pkg/responses"
	"at-backend-claims/internal/pkg/roles"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

type districtsHandler struct {
	service *districts.DistrictsService
}

func NewDistrictsHandler(service *districts.DistrictsService) *districtsHandler {
	return &districtsHandler{
		service: service,
	}
}

func (d districtsHandler) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("GET /", d.getDistricts)
	muxWithAuth.HandleFunc("GET /check-point", d.isPointInPolygon)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.User, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}

func (d districtsHandler) isPointInPolygon(w http.ResponseWriter, r *http.Request) {
	latString := r.URL.Query().Get("lat")
	longString := r.URL.Query().Get("long")

	if latString == "" || longString == "" {
		slog.WarnContext(r.Context(), "empty query parameters")

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	lat, err := strconv.ParseFloat(latString, 64)
	if err != nil || !(-90 <= lat || lat <= 90) {
		slog.WarnContext(r.Context(), "invalid query parameters")

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	long, err := strconv.ParseFloat(longString, 64)
	if err != nil || !(-180 <= long || long <= 180) {
		slog.WarnContext(r.Context(), "invalid query parameters")

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	resp := map[string]any{
		"inside": d.service.IsPointInPolygon(long, lat),
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

var preparedDistricts map[string]any

var onceGetDistricts sync.Once

func (d districtsHandler) getDistricts(w http.ResponseWriter, r *http.Request) {
	onceGetDistricts.Do(func() {
		districts := d.service.Get()

		resp := make(map[string]any)

		for key, polygon := range districts {
			polygonResp := make([]map[string]any, 0, len(polygon))

			for _, point := range polygon {
				polygonResp = append(polygonResp, map[string]any{
					"latitude": point.Y,
					"longitude": point.X,
				})
			}

			resp[key] = polygonResp
		}

		preparedDistricts = resp
	})

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		preparedDistricts,
	)
}
