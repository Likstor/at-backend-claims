package districts

type DistrictsService struct {
	districts map[string]Polygon
}

func NewDistrictsService(districts map[string]Polygon) *DistrictsService {
	return &DistrictsService{
		districts: districts,
	}
}

func (ms DistrictsService) IsPointInPolygon(long, lat float64) bool {
	point := Point{
		X: long,
		Y: lat,
	}

	for _, district := range ms.districts {
		if IsPointInPolygon(point, district) {
			return true
		}
	}

	return false
}

func (ms DistrictsService) Get() map[string]Polygon {
	return ms.districts
}