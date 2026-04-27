// Package raycast предоставляет функции для работы с полигонами и проверки,
// находится ли точка внутри полигона.
package districts

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	EPS = 1e-9
)

var (
	RootDir string
	ErrLess3Points = errors.New("less than 3 points")
	ErrInvalidPath = errors.New("invalid path")
)

// Point представляет собой точку в двумерном пространстве с координатами X и Y.
type Point struct {
	X, Y float64
}

// Polygon представляет собой многоугольник, состоящий из сегментов.
type Polygon []Point

func init() {
	var err error

	RootDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

// NewPolygon создает новый полигон из заданного массива точек.
// Возвращает ошибку ErrLess3Points, если передано менее 3 точек.
func NewPolygon(points []Point) (Polygon, error) {
	if len(points) < 3 {
		return nil, ErrLess3Points
	}

	p := Polygon(points)
	p = append(p, points[0])

	return p, nil
}

// IsPointInPolygon проверяет, находится ли заданная точка внутри полигона.
// Возвращает true, если точка находится внутри или на границе полигона, иначе false.
func IsPointInPolygon(point Point, polygon Polygon) bool {
	intersection := 0

	for i := 0; i < len(polygon) - 1; i++ {
		if !(point.X <= max(polygon[i].X, polygon[i + 1].X) &&
			point.Y >= min(polygon[i].Y, polygon[i + 1].Y) &&
			point.Y <= max(polygon[i].Y, polygon[i + 1].Y)) {
			continue
		}

		if point.Y == polygon[i].Y && point.Y == polygon[i + 1].Y {
			continue
		}

		xIntr := (polygon[i + 1].X-polygon[i].X)*
			(point.Y-polygon[i].Y)/
			(polygon[i + 1].Y-polygon[i].Y) +
			polygon[i].X

		if xIntr-EPS < point.X && point.X < xIntr+EPS {
			return true
		}

		if point.X < xIntr+EPS {
			intersection++
		}
	}

	return intersection%2 == 1
}

type geometryJSON struct {
	Coordinates [][2]float64 `json:"coordinates"`
}

// LoadPolygonFromFile загружает полигон из файла, где каждая строка содержит
// координаты точки в формате "широта, долгота".
func LoadPolygonFromFile(path string) (Polygon, error) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	
	geo := geometryJSON{}

	if err := json.NewDecoder(file).Decode(&geo); err != nil {
		return nil, err
	}

	points := make([]Point, 0, len(geo.Coordinates) + 1)
	for _, point := range geo.Coordinates {
		points = append(points, Point{
			X: point[0],
			Y: point[1],
		})
	}

	return NewPolygon(points)
}