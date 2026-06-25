---
name: geo
description: Add geospatial features to a togo app with the geo plugin — distances, bearings, bounding boxes, polygon/circle geofencing, and nearest/within proximity queries (pure-Go, no PostGIS needed).
---

# togo geo

Use the `geo` plugin for location math.

## Distance & direction
```go
geo.DistanceKm(a, b)   // a,b are geo.Point{Lat,Lng}
geo.Haversine(a, b)    // meters
geo.Bearing(a, b)      // 0–360°
```

## Geofencing
```go
geo.Circle{Center:c, RadiusM:500}.Contains(p)
geo.Polygon{...}.Contains(p)   // point-in-polygon
```

## Find nearby
```go
geo.Nearest(target, points, k)        // k closest
geo.Within(target, points, meters)    // within radius
ix := geo.NewIndex(); ix.Add(geo.Place{ID, Point}); ix.Near(target, meters)
```

## Tips
- For DB-scale "near me" queries, pre-filter with `BoundsAround` (an indexable lat/lng rectangle) then refine with `Haversine`; reach for PostGIS / pgvector only when you outgrow this.
- Store coordinates as `{lat,lng}` (Point marshals that way); validate ranges (lat −90..90, lng −180..180).
