package main

import (
	"sort"
	"sync"
)

type Top struct {
	mu           sync.Mutex
	capacity     int
	positions    map[string]float64
	threshold    float64
	thresholdKey string
}

func NewTop(capacity int) Top {
	return Top{
		capacity:  capacity,
		positions: make(map[string]float64, capacity),
	}
}

func (t *Top) Add(key string, divergence float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// When structure is not filled, add all values.
	if len(t.positions) < t.capacity {
		t.positions[key] = divergence
		if divergence > t.threshold {
			t.threshold = divergence
			t.thresholdKey = key
		}
		return
	}

	// When structure is filled, do not add values outside of the threshold.
	if divergence > t.threshold {
		return
	}

	// To add new value, remove threshold value.
	delete(t.positions, t.thresholdKey)
	t.positions[key] = divergence

	// Find new threshold value in updated map.
	var (
		newThreshold    float64
		newThresholdVal string
	)
	for key, divergence := range t.positions {
		if divergence > newThreshold {
			newThreshold = divergence
			newThresholdVal = key
		}
	}
	t.threshold = newThreshold
	t.thresholdKey = newThresholdVal
}

func (t *Top) List() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	res := make([]string, 0, len(t.positions))
	for key := range t.positions {
		res = append(res, key)
	}

	sort.Slice(res, func(i, j int) bool {
		return t.positions[res[i]] < t.positions[res[j]]
	})

	return res
}

func (t *Top) Value(key string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.positions[key]
}
