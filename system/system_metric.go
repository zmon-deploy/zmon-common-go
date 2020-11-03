package system

import (
	"expvar"
	"runtime"
)

type SystemMetrics struct {
	NumGoroutines int     `json:"numGoroutines"`
	Alloc         uint64  `json:"alloc"`
	TotalAlloc    uint64  `json:"totalAlloc"`
	Mallocs       uint64  `json:"mallocs"`
	Frees         uint64  `json:"frees"`
	HeapAlloc     uint64  `json:"heapAlloc"`
	HeapSys       uint64  `json:"heapSys"`
	HeapIdle      uint64  `json:"heapIdle"`
	HeapInuse     uint64  `json:"heapInuse"`
	HeapReleased  uint64  `json:"heapReleased"`
	HeapObjects   uint64  `json:"heapObjects"`
	StackInuse    uint64  `json:"stackInuse"`
	StackSys      uint64  `json:"stackSys"`
	PauseTotalNs  uint64  `json:"pauseTotalNs"`
	NumGC         uint32  `json:"numGC"`
	NumForcedGC   uint32  `json:"numForcedGC"`
	GcCpuFraction float64 `json:"gcCpuFraction"`
}

func CurrentSystemMetrics() *SystemMetrics {
	memstatsFn := expvar.Get("memstats").(expvar.Func)
	memstats := memstatsFn().(runtime.MemStats)

	return &SystemMetrics{
		NumGoroutines: runtime.NumGoroutine(),
		Alloc:         memstats.Alloc,
		TotalAlloc:    memstats.TotalAlloc,
		Mallocs:       memstats.Mallocs,
		Frees:         memstats.Frees,
		HeapAlloc:     memstats.HeapAlloc,
		HeapSys:       memstats.HeapSys,
		HeapIdle:      memstats.HeapIdle,
		HeapInuse:     memstats.HeapInuse,
		HeapReleased:  memstats.HeapReleased,
		HeapObjects:   memstats.HeapObjects,
		StackInuse:    memstats.StackInuse,
		StackSys:      memstats.StackSys,
		PauseTotalNs:  memstats.PauseTotalNs,
		NumGC:         memstats.NumGC,
		NumForcedGC:   memstats.NumForcedGC,
		GcCpuFraction: memstats.GCCPUFraction,
	}
}

