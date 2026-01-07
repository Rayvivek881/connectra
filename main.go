package main

import (
	"os"
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

// isLambdaMode detects if the application is running in AWS Lambda
func isLambdaMode() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || 
		   os.Getenv("LAMBDA_MODE") == "true"
}

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

// serverMain runs the application in server mode (traditional HTTP server)
func serverMain() {
	v := conf.Viper{}
	v.Init()

	// Validate configuration
	if err := conf.ValidateConfig(); err != nil {
		log.Fatal().Err(err).Msg("Configuration validation failed")
	}

	connections.InitDB()
	connections.InitSearchEngine()
	connections.InitS3()

	defer func() {
		connections.CloseDB()
		connections.CloseSearchEngine()
	}()

	// Start memory stats goroutine (only in server mode)
	go MemoryStats()

	// Execute Cobra commands (api-server, jobs, etc.)
	cmd.Execute()
}

// lambdaMain runs the application in Lambda mode
// Note: Lambda uses lambda/main.go as entry point, this is just a placeholder
func lambdaMain() {
	log.Info().Msg("Lambda mode detected - using lambda package entry point")
	log.Info().Msg("Lambda handler should be built from lambda/main.go")
}

func main() {
	if isLambdaMode() {
		lambdaMain()
	} else {
		serverMain()
	}
}
