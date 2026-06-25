package geo

import (
	"math"
	"testing"
)

var (
	nyc = Point{40.7128, -74.0060}
	la  = Point{34.0522, -118.2437}
	sf  = Point{37.7749, -122.4194}
)

func TestHaversineKnownDistance(t *testing.T) {
	km := DistanceKm(nyc, la)
	if math.Abs(km-3935) > 30 { // NYC↔LA ≈ 3935 km
		t.Fatalf("NYC→LA = %.0f km, want ~3935", km)
	}
	if mi := DistanceMiles(nyc, la); math.Abs(mi-2445) > 25 {
		t.Fatalf("NYC→LA = %.0f mi, want ~2445", mi)
	}
	if Haversine(nyc, nyc) != 0 {
		t.Fatal("distance to self must be 0")
	}
}

func TestBearing(t *testing.T) {
	// Due east along the equator → ~90°.
	b := Bearing(Point{0, 0}, Point{0, 1})
	if math.Abs(b-90) > 0.5 {
		t.Fatalf("east bearing = %.1f, want ~90", b)
	}
	// Due north → ~0°.
	if n := Bearing(Point{0, 0}, Point{1, 0}); math.Abs(n) > 0.5 && math.Abs(n-360) > 0.5 {
		t.Fatalf("north bearing = %.1f, want ~0", n)
	}
}

func TestBoundsAround(t *testing.T) {
	b := BoundsAround(nyc, 1000) // 1 km box
	if !b.Contains(nyc) {
		t.Fatal("box must contain its center")
	}
	// A point ~10 km away must be outside a 1 km box.
	far := Point{nyc.Lat + 0.1, nyc.Lng}
	if b.Contains(far) {
		t.Fatal("far point should be outside the 1 km box")
	}
}

func TestPolygonContains(t *testing.T) {
	// A unit square around the origin.
	sq := Polygon{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	if !sq.Contains(Point{0.5, 0.5}) {
		t.Fatal("center should be inside")
	}
	if sq.Contains(Point{2, 2}) {
		t.Fatal("far point should be outside")
	}
	if (Polygon{{0, 0}, {1, 1}}).Contains(Point{0.5, 0.5}) {
		t.Fatal("degenerate polygon (<3 pts) should contain nothing")
	}
}

func TestCircleFence(t *testing.T) {
	c := Circle{Center: nyc, RadiusM: 5000}
	if !c.Contains(nyc) {
		t.Fatal("center inside circle")
	}
	if c.Contains(la) {
		t.Fatal("LA should be outside a 5 km NYC circle")
	}
}

func TestNearestAndWithin(t *testing.T) {
	pts := []Point{la, sf, {40.730, -73.935}} // last is ~near NYC
	near := Nearest(nyc, pts, 2)
	if len(near) != 2 {
		t.Fatalf("Nearest k=2 → %d", len(near))
	}
	// Closest should be the near-NYC point, then SF, then LA — so [0] is near-NYC.
	if Haversine(nyc, near[0]) > Haversine(nyc, near[1]) {
		t.Fatal("Nearest not sorted ascending")
	}
	if near[0] != pts[2] {
		t.Fatalf("closest = %v, want near-NYC point", near[0])
	}
	within := Within(nyc, pts, 50000) // 50 km
	if len(within) != 1 || within[0] != pts[2] {
		t.Fatalf("Within 50km = %v, want just the near-NYC point", within)
	}
}

func TestParseAndString(t *testing.T) {
	p, err := Parse("40.7128, -74.0060")
	if err != nil || p.Lat != 40.7128 || p.Lng != -74.0060 {
		t.Fatalf("Parse = %+v, %v", p, err)
	}
	if _, err := Parse("nope"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestIndexNear(t *testing.T) {
	ix := NewIndex()
	ix.Add(Place{ID: "la", Point: la})
	ix.Add(Place{ID: "sf", Point: sf})
	ix.Add(Place{ID: "near", Point: Point{40.730, -73.935}})
	if ix.Len() != 3 {
		t.Fatalf("Len = %d", ix.Len())
	}
	hits := ix.Near(nyc, 50000)
	if len(hits) != 1 || hits[0].ID != "near" {
		t.Fatalf("Near = %+v, want just 'near'", hits)
	}
	all := ix.Nearest(nyc, 3)
	if len(all) != 3 || all[0].ID != "near" || all[0].DistanceM > all[1].DistanceM {
		t.Fatalf("Nearest order wrong: %+v", all)
	}
	ix.Remove("near")
	if ix.Len() != 2 {
		t.Fatal("Remove failed")
	}
}
