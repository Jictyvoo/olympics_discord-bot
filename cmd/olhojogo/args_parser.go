package main

import "flag"

type cliOptions struct {
	Subcommand string
	ConfigPath string
}

func parseArgs(args []string) cliOptions {
	opts := cliOptions{Subcommand: "serve"}
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		opts.Subcommand = args[0]
		args = args[1:]
	}

	fs := flag.NewFlagSet(opts.Subcommand, flag.ExitOnError)
	fs.StringVar(
		&opts.ConfigPath,
		"config",
		"",
		"Path to config file (default: search standard locations)",
	)
	_ = fs.Parse(args)
	return opts
}
