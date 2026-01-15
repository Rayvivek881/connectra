package main

import (
	"runtime"
	"runtime/debug"
	"time"
	"vivek-ray/cmd"
	"vivek-ray/conf"
	"vivek-ray/connections"

	"github.com/rs/zerolog/log"
)

const (
	MB = 1024 * 1024
)

func MemoryStats() {
	ticker := time.NewTicker(time.Duration(conf.AppConfig.MemoryLogInterval) * time.Second)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.Info().Msgf("Memory: Alloc=%dMB | HeapInUse=%dMB | Sys=%dMB | NumGC=%d",
			m.Alloc/MB,
			m.HeapInuse/MB,
			m.Sys/MB,
			m.NumGC,
		)
	}
}

func main() {
	// Aggressive GC: minimize memory usage (trades CPU for lower memory)
	debug.SetMemoryLimit(128 * MB) // Soft limit: GC becomes very aggressive near this
	debug.SetGCPercent(50)         // GC triggers at 50% heap growth (4x more frequent than default)

	v := conf.Viper{}
	v.Init()

	connections.InitDB()
	connections.InitSearchEngine()
	connections.InitS3()

	defer func() {
		connections.CloseDB()
		connections.CloseSearchEngine()
	}()
	go MemoryStats()
	cmd.Execute()
}
