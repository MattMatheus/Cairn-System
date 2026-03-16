package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"toolcli/internal/telemetry"
)

func main() {
	shutdown, initErr := telemetry.InitOTel(context.Background())
	if initErr != nil {
		fmt.Fprintf(os.Stderr, "warning: otel initialization failed: %v\n", initErr)
		shutdown = func(context.Context) error { return nil }
	}
	runShutdown := func() {
		if err := shutdown(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: otel shutdown failed: %v\n", err)
		}
	}
	defer runShutdown()

	if len(os.Args) < 2 {
		exitErr(errors.New("usage: codegraph-cli <analyze|status|context|impact> [flags]"))
	}

	var err error
	switch os.Args[1] {
	case "analyze":
		err = runAnalyze(os.Args[2:])
	case "status":
		err = runStatus(os.Args[2:])
	case "context":
		err = runContext(os.Args[2:])
	case "impact":
		err = runImpact(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}
	if err != nil {
		runShutdown()
		exitErr(err)
	}
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
