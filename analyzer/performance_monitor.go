// analyzer/performance_monitor.go
package analyzer

import (
	"runtime"
	"time"
)

// PerformanceMonitor monitorea el rendimiento del analizador
type PerformanceMonitor struct {
	startTime    time.Time
	memoryBefore runtime.MemStats
	memoryAfter  runtime.MemStats
}

// StartMonitoring inicia el monitoreo de rendimiento
func (pm *PerformanceMonitor) StartMonitoring() {
	pm.startTime = time.Now()
	runtime.GC() // Forzar garbage collection para medición precisa
	runtime.ReadMemStats(&pm.memoryBefore)
}

// StopMonitoring detiene el monitoreo y retorna estadísticas
func (pm *PerformanceMonitor) StopMonitoring() PerformanceStats {
	runtime.ReadMemStats(&pm.memoryAfter)
	duration := time.Since(pm.startTime)
	
	return PerformanceStats{
		Duration:     duration,
		MemoryUsed:   pm.memoryAfter.Alloc - pm.memoryBefore.Alloc,
		Allocations:  pm.memoryAfter.Mallocs - pm.memoryBefore.Mallocs,
		GCRuns:       pm.memoryAfter.NumGC - pm.memoryBefore.NumGC,
	}
}

// PerformanceStats estadísticas de rendimiento
type PerformanceStats struct {
	Duration     time.Duration
	MemoryUsed   uint64
	Allocations  uint64
	GCRuns       uint32
}

