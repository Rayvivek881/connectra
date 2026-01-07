package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vivek-ray/jobs"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var JobCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Start the jobs",
	Long:  "Start the jobs",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Msg("Starting jobs...")
			jobs.RunJobs(ctx, args)
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info().Msg("Shutting down jobs...")
		cancel()
		wg.Wait()
		log.Info().Msg("Jobs stopped")
	},
}

func init() {
	rootCmd.AddCommand(JobCmd)
}
