# geo — usage

## Points & distance
```go
p := geo.Point{Lat: 40.7128, Lng: -74.0060}
geo.Haversine(a, b)      // meters
geo.DistanceKm(a, b); geo.DistanceMiles(a, b)
geo.Bearing(a, b)        // initial bearing 0–360°
geo.Parse("40.71,-74.0") // "lat,lng" → Point
```

## Bounding box
```go
box := geo.BoundsAround(center, 1000) // 1 km box
box.Contains(p)
```
Use a box as a cheap pre-filter, then confirm with `Haversine`.

## Geofencing
```go
geo.Circle{Center: c, RadiusM: 500}.Contains(p)        // circular fence
geo.Polygon{{lat,lng}, ...}.Contains(p)                // point-in-polygon (ray casting)
```

## Proximity
```go
geo.Nearest(target, points, k)          // k closest []Point, nearest first
geo.Within(target, points, radiusMeters)// []Point within radius
```

## Place index
```go
ix := geo.NewIndex()
ix.Add(geo.Place{ID:"1", Point:p, Meta: map[string]any{...}})
ix.Near(target, radiusMeters)  // []PlaceHit (with DistanceM), sorted
ix.Nearest(target, k)
ix.Remove(id); ix.Len()
```
