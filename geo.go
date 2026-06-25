// Package geo provides geospatial helpers for togo — great-circle distances,
// bearings, bounding boxes, geofencing (polygon + circle), and nearest/within
// queries. Pure-Go (stdlib math), no external dependencies (GeoDjango-lite).
//
//	d := geo.DistanceKm(geo.Point{40.7128, -74.0060}, geo.Point{34.0522, -118.2437}) // NYC→LA ≈ 3935 km
//	near := geo.Within(me, places, 1000) // places within 1 km
package geo

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// EarthRadiusM is the mean Earth radius in meters.
const EarthRadiusM = 6371000.0

// Point is a WGS-84 latitude/longitude in decimal degrees.
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// String renders the point as "lat,lng".
func (p Point) String() string { return fmt.Sprintf("%g,%g", p.Lat, p.Lng) }

// Parse reads a "lat,lng" string.
func Parse(s string) (Point, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return Point{}, fmt.Errorf("geo: expected \"lat,lng\", got %q", s)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Point{}, fmt.Errorf("geo: bad lat: %w", err)
	}
	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Point{}, fmt.Errorf("geo: bad lng: %w", err)
	}
	return Point{lat, lng}, nil
}

func rad(d float64) float64 { return d * math.Pi / 180 }
func deg(r float64) float64 { return r * 180 / math.Pi }

// Haversine returns the great-circle distance between two points in meters.
func Haversine(a, b Point) float64 {
	la1, la2 := rad(a.Lat), rad(b.Lat)
	dLat := rad(b.Lat - a.Lat)
	dLng := rad(b.Lng - a.Lng)
	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(la1)*math.Cos(la2)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * EarthRadiusM * math.Asin(math.Min(1, math.Sqrt(h)))
}

// DistanceKm / DistanceMiles are convenience wrappers around Haversine.
func DistanceKm(a, b Point) float64    { return Haversine(a, b) / 1000 }
func DistanceMiles(a, b Point) float64 { return Haversine(a, b) / 1609.344 }

// Bearing returns the initial compass bearing from a to b in degrees [0,360).
func Bearing(a, b Point) float64 {
	la1, la2 := rad(a.Lat), rad(b.Lat)
	dLng := rad(b.Lng - a.Lng)
	y := math.Sin(dLng) * math.Cos(la2)
	x := math.Cos(la1)*math.Sin(la2) - math.Sin(la1)*math.Cos(la2)*math.Cos(dLng)
	return math.Mod(deg(math.Atan2(y, x))+360, 360)
}

// Bounds is an axis-aligned lat/lng rectangle.
type Bounds struct {
	MinLat, MinLng, MaxLat, MaxLng float64
}

// Contains reports whether p falls inside the bounds.
func (bx Bounds) Contains(p Point) bool {
	return p.Lat >= bx.MinLat && p.Lat <= bx.MaxLat &&
		p.Lng >= bx.MinLng && p.Lng <= bx.MaxLng
}

// BoundsAround returns a bounding box centered on p that covers radius meters
// (useful as a cheap pre-filter before an exact Haversine check).
func BoundsAround(p Point, meters float64) Bounds {
	latDelta := deg(meters / EarthRadiusM)
	cosLat := math.Cos(rad(p.Lat))
	if cosLat < 1e-12 {
		cosLat = 1e-12
	}
	lngDelta := deg(meters / (EarthRadiusM * cosLat))
	return Bounds{p.Lat - latDelta, p.Lng - lngDelta, p.Lat + latDelta, p.Lng + lngDelta}
}

// Polygon is a closed ring of points (the ring auto-closes).
type Polygon []Point

// Contains reports whether p is inside the polygon using ray casting.
func (poly Polygon) Contains(p Point) bool {
	n := len(poly)
	if n < 3 {
		return false
	}
	in := false
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		yi, yj := poly[i].Lat, poly[j].Lat
		xi, xj := poly[i].Lng, poly[j].Lng
		if (yi > p.Lat) != (yj > p.Lat) &&
			p.Lng < (xj-xi)*(p.Lat-yi)/(yj-yi)+xi {
			in = !in
		}
	}
	return in
}

// Circle is a circular geofence.
type Circle struct {
	Center  Point
	RadiusM float64
}

// Contains reports whether p is within the circle.
func (c Circle) Contains(p Point) bool { return Haversine(c.Center, p) <= c.RadiusM }

// Nearest returns the k points closest to target, nearest first.
func Nearest(target Point, points []Point, k int) []Point {
	type d struct {
		p Point
		m float64
	}
	ds := make([]d, len(points))
	for i, p := range points {
		ds[i] = d{p, Haversine(target, p)}
	}
	sort.Slice(ds, func(i, j int) bool { return ds[i].m < ds[j].m })
	if k > len(ds) {
		k = len(ds)
	}
	if k < 0 {
		k = 0
	}
	out := make([]Point, k)
	for i := 0; i < k; i++ {
		out[i] = ds[i].p
	}
	return out
}

// Within returns the points within radiusMeters of target (unordered).
func Within(target Point, points []Point, radiusMeters float64) []Point {
	var out []Point
	for _, p := range points {
		if Haversine(target, p) <= radiusMeters {
			out = append(out, p)
		}
	}
	return out
}
