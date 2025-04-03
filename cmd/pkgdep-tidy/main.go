package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/cloverrose/pkgdep/pkg/log"
	"github.com/cloverrose/pkgdep/pkg/tidy"
)

type config struct {
	tidy tidy.Config
	log  log.Config
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		slog.Error("failed to run", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	closer, err := log.SetDefault(cfg.log)
	if err != nil {
		return err
	}
	defer func() {
		if err := closer(); err != nil {
			fmt.Println(err)
		}
	}()

	tidyCommand := tidy.New(cfg.tidy)
	if err := tidyCommand.Tidy(); err != nil {
		return fmt.Errorf("failed to tidy: %w", err)
	}

	return nil
}

func parseConfig(args []string) (config, error) {
	cfg := config{
		tidy: tidy.Config{
			ConfigFile:  ".pkgdep.yaml",
			InspectFile: "used_rules.csv",
			OutputFile:  ".pkgdep.tidy.yaml",
		},
		log: log.Config{
			Level:  "info",
			Format: "json",
		},
	}

	fs := flag.NewFlagSet("tidy", flag.ContinueOnError)

	fs.StringVar(&cfg.tidy.ConfigFile, "config", cfg.tidy.ConfigFile, "config file path")
	fs.StringVar(&cfg.tidy.InspectFile, "inspector.file", cfg.tidy.InspectFile, "file containing used rules")
	fs.StringVar(&cfg.tidy.OutputFile, "output", cfg.tidy.OutputFile, "output file path")
	fs.StringVar(&cfg.log.Level, "log.level", cfg.log.Level, "log level (debug, info, warn, error)")
	fs.StringVar(&cfg.log.Format, "log.format", cfg.log.Format, "log format (json, text)")

	if err := fs.Parse(args); err != nil {
		return config{}, fmt.Errorf("failed to parse flags: %w", err)
	}

	return cfg, nil
}
