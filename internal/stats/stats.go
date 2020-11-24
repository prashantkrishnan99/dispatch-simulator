package stats

import "sync"

//Stats :
type Stats struct {
	TotalOrdersProcessed int
	AverageTime          int
	TotalTime            int
	mu                   sync.RWMutex
}

//NewStats :
func NewStats() *Stats {
	return &Stats{
		TotalOrdersProcessed: 0,
		TotalTime:            0,
		AverageTime:          0,
	}
}

//IncrOrdersProcessed :
func (stats *Stats) IncrOrdersProcessed() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.TotalOrdersProcessed++
}

//IncrTotalTime :
func (stats *Stats) IncrTotalTime(time int) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.TotalTime = stats.TotalTime + time
}

//CalculateAverage :
func (stats *Stats) CalculateAverage() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.AverageTime = stats.TotalTime / stats.TotalOrdersProcessed
}

//GetTotalOrdersProcessed :
func (stats *Stats) GetTotalOrdersProcessed() int {
	stats.mu.RLock()
	defer stats.mu.RUnlock()
	return stats.TotalOrdersProcessed
}

//GetTotalTime :
func (stats *Stats) GetTotalTime() int {
	stats.mu.RLock()
	defer stats.mu.RUnlock()
	return stats.TotalTime
}

//GetAVerageTime :
func (stats *Stats) GetAVerageTime() int {
	stats.mu.RLock()
	defer stats.mu.RUnlock()
	return stats.AverageTime
}
