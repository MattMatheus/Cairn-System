package codegraph

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command string

const (
	CommandAnalyze Command = "analyze"
	CommandStatus  Command = "status"
	CommandContext Command = "context"
	CommandImpact  Command = "impact"
)

type Options struct {
	Command   Command
	Repo      string
	Target    string
	File      string
	Direction string
	Force     bool
}

type Install struct {
	Mode       string
	BinaryPath string
	NodePath   string
	EntryPoint string
	RepoRoot   string
}

func Run(ctx context.Context, opts Options) ([]byte, error) {
	install, err := ResolveInstall(os.Getwd, os.LookupEnv, exec.LookPath)
	if err != nil {
		return nil, err
	}
	cmd, err := BuildExecCommand(install, opts)
	if err != nil {
		return nil, err
	}
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Dir = commandDir(opts)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("gitnexus command failed: %w\n%s", err, strings.TrimSpace(string(output)))
	}
	return output, nil
}

func ResolveInstall(getwd func() (string, error), lookupEnv func(string) (string, bool), lookPath func(string) (string, error)) (Install, error) {
	if value, ok := lookupEnv("CAIRN_GITNEXUS_BIN"); ok && strings.TrimSpace(value) != "" {
		return Install{Mode: "binary", BinaryPath: strings.TrimSpace(value)}, nil
	}
	if path, err := lookPath("gitnexus"); err == nil && strings.TrimSpace(path) != "" {
		return Install{Mode: "binary", BinaryPath: path}, nil
	}
	nodePath, nodeErr := lookPath("node")
	if nodeErr != nil {
		nodePath = ""
	}

	root := ""
	if value, ok := lookupEnv("CAIRN_GITNEXUS_ROOT"); ok && strings.TrimSpace(value) != "" {
		root = strings.TrimSpace(value)
	} else {
		cwd, err := getwd()
		if err == nil {
			root = findGitNexusRoot(cwd)
		}
	}
	if root != "" {
		entry := filepath.Join(root, "dist", "cli", "index.js")
		if nodePath != "" {
			if _, err := os.Stat(entry); err == nil {
				return Install{Mode: "node", NodePath: nodePath, EntryPoint: entry, RepoRoot: root}, nil
			}
		}
		return Install{}, fmt.Errorf("gitnexus checkout found at %s but dist/cli/index.js is missing or node is unavailable", root)
	}
	return Install{}, errors.New("gitnexus is not ready; set CAIRN_GITNEXUS_BIN, or provide a built checkout via CAIRN_GITNEXUS_ROOT with node available")
}

func BuildExecCommand(install Install, opts Options) (*exec.Cmd, error) {
	args, err := buildArgs(opts)
	if err != nil {
		return nil, err
	}
	switch install.Mode {
	case "binary":
		return exec.Command(install.BinaryPath, args...), nil
	case "node":
		return exec.Command(install.NodePath, append([]string{install.EntryPoint}, args...)...), nil
	default:
		return nil, fmt.Errorf("unsupported gitnexus install mode: %s", install.Mode)
	}
}

func commandDir(opts Options) string {
	if strings.TrimSpace(opts.Repo) != "" {
		return opts.Repo
	}
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func buildArgs(opts Options) ([]string, error) {
	switch opts.Command {
	case CommandAnalyze:
		args := []string{"analyze"}
		if strings.TrimSpace(opts.Repo) != "" {
			args = append(args, opts.Repo)
		}
		if opts.Force {
			args = append(args, "--force")
		}
		return args, nil
	case CommandStatus:
		return []string{"status"}, nil
	case CommandContext:
		if strings.TrimSpace(opts.Target) == "" {
			return nil, errors.New("context target is required")
		}
		args := []string{"context", opts.Target}
		if strings.TrimSpace(opts.File) != "" {
			args = append(args, opts.File)
		}
		return args, nil
	case CommandImpact:
		if strings.TrimSpace(opts.Target) == "" {
			return nil, errors.New("impact target is required")
		}
		direction := strings.TrimSpace(opts.Direction)
		if direction == "" {
			direction = "upstream"
		}
		if direction != "upstream" && direction != "downstream" {
			return nil, fmt.Errorf("unsupported impact direction: %s", direction)
		}
		return []string{"impact", opts.Target, "--direction", direction}, nil
	default:
		return nil, fmt.Errorf("unsupported codegraph command: %s", opts.Command)
	}
}

func findGitNexusRoot(start string) string {
	dir := start
	for {
		candidate := filepath.Join(dir, "repos", "untrusted", "GitNexus", "gitnexus", "package.json")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Dir(candidate)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
