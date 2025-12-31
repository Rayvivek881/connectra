package main

import (
	"runtime"
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
