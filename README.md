<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/geo</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/geo"><img src="https://pkg.go.dev/badge/github.com/togo-framework/geo.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Geospatial helpers for <a href="https://to-go.dev">togo</a> — distances, bounding boxes, geofencing & nearest.</strong></p>
</div>

## Install

```bash
togo install togo-framework/geo
```

Pure-Go geospatial helpers (GeoDjango-lite): great-circle distances, bearings, bounding boxes, polygon/circle geofencing, and nearest/within proximity queries — no external dependencies.

## Usage

```go
nyc := geo.Point{Lat: 40.7128, Lng: -74.0060}
la  := geo.Point{Lat: 34.0522, Lng: -118.2437}

geo.DistanceKm(nyc, la)        // ≈ 3935 km
geo.Haversine(nyc, la)         // meters
geo.Bearing(nyc, la)           // initial compass bearing (0–360°)

// Geofencing
fence := geo.Circle{Center: nyc, RadiusM: 5000}
fence.Contains(p)              // within 5 km of NYC?
poly := geo.Polygon{{0,0},{0,1},{1,1},{1,0}}
poly.Contains(geo.Point{0.5, 0.5}) // point-in-polygon (ray casting)

// Proximity
geo.Nearest(nyc, points, 5)            // 5 closest points, nearest first
geo.Within(nyc, points, 1000)          // points within 1 km

// Bounding box (cheap pre-filter before exact distance)
box := geo.BoundsAround(nyc, 1000)
box.Contains(p)
```

### Place index

```go
ix := geo.NewIndex()
ix.Add(geo.Place{ID: "store-1", Point: nyc, Meta: map[string]any{"name": "HQ"}})
hits := ix.Near(me, 2000)     // []PlaceHit sorted by distance, within 2 km
top  := ix.Nearest(me, 10)    // 10 closest
```

`Point` marshals as `{"lat":…,"lng":…}`; `geo.Parse("40.7128,-74.0060")` reads `"lat,lng"`.

---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
