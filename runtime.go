package gmx

// pkg/runtime instrumentation

import (
	"runtime"
)


func init() {
	NewGaugeWithCallback("runtime.gomaxprocs", func() int64 {
		return (int64)(runtime.GOMAXPROCS(0))
	})

	NewGaugeWithCallback("runtime.numcgocall", func()int64{
		return (int64)(runtime.NumCgoCall())
	})

	NewGaugeWithCallback("runtime.numcpu", func()int64{
		return (int64)(runtime.NumCPU())
	})

	NewGaugeWithCallback("runtime.numgoroutine", func()int64{
		return(int64)(runtime.NumGoroutine())
	})

	Publish("runtime.version", (MetricFunc)(func() interface{}{
		return runtime.Version()
	}))

	//Publish("runtime.memstats", runtimeMemStats)
	NewGaugeWithCallback("runtime.memstats.alloc", func() int64 {
		return (int64)(runtimeMemStats().Alloc)
	})

	NewGaugeWithCallback("runtime.memstats.buckhashsyc", func() int64 {
		return (int64)(runtimeMemStats().BuckHashSys)
	})

	NewGaugeWithCallback("runtime.memstats.frees", func() int64 {
		return (int64)(runtimeMemStats().Frees)
	})

	NewGaugeWithCallback("runtime.memstats.heapalloc", func() int64 {
		return (int64)(runtimeMemStats().HeapAlloc)
	})

	NewGaugeWithCallback("runtime.memstats.heapidle", func() int64 {
		return (int64)(runtimeMemStats().HeapIdle)
	})

	NewGaugeWithCallback("runtime.memstats.heapinuse", func() int64 {
		return (int64)(runtimeMemStats().HeapInuse)
	})

	NewGaugeWithCallback("runtime.memstats.heapobjects", func() int64 {
		return (int64)(runtimeMemStats().HeapObjects)
	})

	NewGaugeWithCallback("runtime.memstats.heapreleased", func() int64 {
		return (int64)(runtimeMemStats().HeapReleased)
	})

	NewGaugeWithCallback("runtime.memstats.heapsys", func() int64 {
		return (int64)(runtimeMemStats().HeapSys)
	})

	NewGaugeWithCallback("runtime.memstats.lastgc", func() int64 {
		return (int64)(runtimeMemStats().LastGC)
	})

	NewGaugeWithCallback("runtime.memstats.lookups", func() int64 {
		return (int64)(runtimeMemStats().Lookups)
	})

	NewGaugeWithCallback("runtime.memstats.mallocs", func() int64 {
		return (int64)(runtimeMemStats().Mallocs)
	})

	NewGaugeWithCallback("runtime.memstats.mcacheinuse", func() int64 {
		return (int64)(runtimeMemStats().MCacheInuse)
	})

	NewGaugeWithCallback("runtime.memstats.mcachesys", func() int64 {
		return (int64)(runtimeMemStats().MCacheSys)
	})

	NewGaugeWithCallback("runtime.memstats.mspaninuse", func() int64 {
		return (int64)(runtimeMemStats().MSpanInuse)
	})

	NewGaugeWithCallback("runtime.memstats.mspansys", func() int64 {
		return (int64)(runtimeMemStats().MSpanSys)
	})

	NewGaugeWithCallback("runtime.memstats.nextgc", func() int64 {
		return (int64)(runtimeMemStats().NextGC)
	})

	NewGaugeWithCallback("runtime.memstats.numgc", func() int64 {
		return (int64)(runtimeMemStats().NumGC)
	})

	NewGaugeWithCallback("runtime.memstats.gccpufraction", func() int64 {
		return (int64)(runtimeMemStats().GCCPUFraction)
	})

	NewGaugeWithCallback("runtime.memstats.stackinuse", func() int64 {
		return (int64)(runtimeMemStats().StackInuse)
	})


	NewGaugeWithCallback("runtime.memstats.stacksys", func() int64 {
		return (int64)(runtimeMemStats().StackSys)
	})

	NewGaugeWithCallback("runtime.memstats.sys", func() int64 {
		return (int64)(runtimeMemStats().Sys)
	})

	NewGaugeWithCallback("runtime.memstats.totalalloc", func() int64 {
		return (int64)(runtimeMemStats().TotalAlloc)
	})

	NewGaugeWithCallback("runtime.memstats.stackinuse", func() int64 {
		return (int64)(runtimeMemStats().StackInuse)
	})
}


func runtimeMemStats() *runtime.MemStats{
	var memstats runtime.MemStats

	runtime.ReadMemStats(&memstats)

	return &memstats
}

