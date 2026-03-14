package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"athenause/internal/telemetry"
	"gopkg.in/yaml.v3"
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
		exitErr(errors.New("usage: promote-cli <inspect|note> [flags]"))
	}

	var err error
	switch os.Args[1] {
	case "inspect":
		err = runInspect(os.Args[2:])
	case "note":
		err = runNote(os.Args[2:])
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

func writeOutput(payload any, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		data, err := yaml.Marshal(payload)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	case "text":
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func repoRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	dir := wd
	for {
		candidate := filepath.Join(dir, "products", "athena-mind", "cmd", "memory-cli")
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return wd
		}
		dir = parent
	}
}

func defaultMemoryRoot() string {
	if value := strings.TrimSpace(os.Getenv("ATHENA_MEMORY_ROOT")); value != "" {
		return value
	}
	return filepath.Join(repoRoot(), ".athena", "memory", "default")
}

func memoryCLIPath() string {
	return filepath.Join(repoRoot(), "products", "athena-mind", "cmd", "memory-cli")
}

func goCommand() string {
	if value := strings.TrimSpace(os.Getenv("GO")); value != "" {
		return value
	}
	return "go"
}
