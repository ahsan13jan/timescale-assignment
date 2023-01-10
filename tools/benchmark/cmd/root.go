/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"timescale/internal/config"
	"timescale/internal/run"
)

var (
	fileName             string
	maxConcurrentWorkers int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark timescale instance",
	Long:  `benchmark timescale instance`,
	Run: func(cmd *cobra.Command, args []string) {
		var c config.Config
		envconfig.MustProcess("", &c)
		fmt.Println(c.String())
		metrics := run.Execute(fileName, maxConcurrentWorkers, c)
		metrics.Print()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&fileName, "file", "f", "", "Queries in csv format")
	rootCmd.PersistentFlags().IntVarP(&maxConcurrentWorkers, "maxWorkers", "m", -1, "Maximum concurrent workers")

}
