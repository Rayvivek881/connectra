package cmd

import (
	"vivek-ray/jobs"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var s3JobCmd = &cobra.Command{
	Use:   "s3-job",
	Short: "Start the S3 file insert job",
	Long:  "Start the S3 file insert job",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting S3 file insert job...")
		jobs.InsertFileJob()
	},
}

func init() {
	rootCmd.AddCommand(s3JobCmd)
}
