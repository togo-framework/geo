package geo

import (
	"sort"
	"sync"
)

// Place is an identified point with optional metadata, for the in-memory index.
type Place struct {
	ID    string         `json:"id"`
	Point Point          `json:"point"`
	Meta  map[string]any `json:"meta,omitempty"`
}

// PlaceHit is a Place plus its distance from the query target.
type PlaceHit struct {
	Place
	DistanceM float64 `json:"distance_m"`
}

// Index is a concurrency-safe in-memory place index with proximity queries.
type Index struct {
	mu     sync.RWMutex
	places map[string]Place
}

// NewIndex returns an empty place index.
func NewIndex() *Index { return &Index{places: map[string]Place{}} }

// Add inserts or replaces a place by ID.
func (ix *Index) Add(p Place) {
	ix.mu.Lock()
	ix.places[p.ID] = p
	ix.mu.Unlock()
}

// Remove deletes a place by ID.
func (ix *Index) Remove(id string) {
	ix.mu.Lock()
	delete(ix.places, id)
	ix.mu.Unlock()
}

// Len reports how many places are indexed.
func (ix *Index) Len() int {
	ix.mu.RLock()
	defer ix.mu.RUnlock()
	return len(ix.places)
}

// Near returns places within radiusMeters of target, nearest first.
func (ix *Index) Near(target Point, radiusMeters float64) []PlaceHit {
	ix.mu.RLock()
	hits := make([]PlaceHit, 0)
	for _, p := range ix.places {
		if d := Haversine(target, p.Point); d <= radiusMeters {
			hits = append(hits, PlaceHit{Place: p, DistanceM: d})
		}
	}
	ix.mu.RUnlock()
	sort.Slice(hits, func(i, j int) bool { return hits[i].DistanceM < hits[j].DistanceM })
	return hits
}

// Nearest returns the k closest places to target, nearest first.
func (ix *Index) Nearest(target Point, k int) []PlaceHit {
	ix.mu.RLock()
	hits := make([]PlaceHit, 0, len(ix.places))
	for _, p := range ix.places {
		hits = append(hits, PlaceHit{Place: p, DistanceM: Haversine(target, p.Point)})
	}
	ix.mu.RUnlock()
	sort.Slice(hits, func(i, j int) bool { return hits[i].DistanceM < hits[j].DistanceM })
	if k >= 0 && k < len(hits) {
		hits = hits[:k]
	}
	return hits
}
