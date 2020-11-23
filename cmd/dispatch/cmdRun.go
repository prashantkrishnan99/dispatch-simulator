package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatch-simulator/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run service",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg, err := getConfig(); err == nil {
			runner := runner.NewRunner(cfg.Runner)
			if err := runner.Run(); err != nil {
				fmt.Println("Cannot Start Runner Process")
				return
			}

			sig := make(chan os.Signal)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		forever:
			for {
				select {
				case sgn := <-sig:
					fmt.Println("Signal received, process is stopping" + sgn.String())
					break forever
				}
			}
		} else {
			fmt.Println("Config file validation check failed" + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
