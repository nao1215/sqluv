// Package main is the entry point of the sqluv command.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/di"
)

// main is the entry point of the sqluv command.
func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args))
}

// run executes the sqluv command.
func run(stdout, stderr io.Writer, args []string) int {
	arg, err := config.NewArgument(args)
	if err != nil {
		fmt.Fprintf(stderr, "failed to parse arguments: %v\n", err)
		return 1
	}
	if arg.CanUsage() {
		fmt.Fprintf(stdout, "%s", arg.Usage())
		return 0
	}
	if arg.CanVersion() {
		fmt.Fprintf(stdout, "%s", arg.Version())
		return 0
	}

	sqluv, err := di.NewSqluv(arg)
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize TUI: %v\n", err)
		return 1
	}
	if err := sqluv.Run(); err != nil {
		fmt.Fprintf(stderr, "failed to run TUI: %v\n", err)
		return 1
	}
	return 0
}
