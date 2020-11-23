package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatch-simulator/internal/runner"
	"github.com/spf13/cobra"
	"go.melnyk.org/mlog"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()

		if cfg, err := getConfig(); err == nil {
			defer cfg.cleanup()

			log := cfg.logbook.Joiner().Join("cmd")

			// Initialize runner
			runner := runner.NewRunner(cfg.Runner, cfg.logbook.Joiner())

			if err = runner.Run(); err != nil {
				log.Event(mlog.Error, func(ev mlog.Event) {
					ev.String("msg", "Cannot start service process")
					ev.String("err", err.Error())
				})
				return
			}

			sig := make(chan os.Signal)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		forever:
			for {
				select {
				case <-runner.Stopped():
					log.Event(mlog.Error, func(ev mlog.Event) {
						ev.String("msg", "Service internally stopped")
					})
					break forever
				case sgn := <-sig:
					log.Event(mlog.Warning, func(ev mlog.Event) {
						ev.String("msg", "Signal received, process is stopping")
						ev.String("signal", sgn.String())
					})
					runner.Stop()
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
