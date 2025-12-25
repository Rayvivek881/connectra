package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vivek-ray/conf"
	"vivek-ray/constants"
	"vivek-ray/jobs"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var s3JobCmd = &cobra.Command{
	Use:   "s3-job",
	Short: "Start the S3 file insert job",
	Long:  "Start the S3 file insert job",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup
		if conf.JobConfig.JobType == constants.NormalJobType {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Info().Msg("Starting S3 file insert job...")
				jobs.InsertFileJob(ctx)
			}()
		}
		if conf.JobConfig.JobType == constants.RetryJobType {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Info().Msg("Starting S3 file retry insert job...")
				jobs.RetryInsertFileJob(ctx)
			}()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info().Msg("Shutting down S3 file insert job...")
		cancel()
		wg.Wait()
		log.Info().Msg("S3 file insert job stopped")
	},
}

func init() {
	rootCmd.AddCommand(s3JobCmd)
}
