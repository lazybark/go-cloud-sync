package configs

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

// Config stores while app config divided into separate structs
type Config struct {
	FS FilesystemConfig
}

// FilesystemConfig holds variables to config filesystem watcher and event processor
type FilesystemConfig struct {
	Root string
}

// EnvVars is a struct to parse ENV variables & flags
type EnvVars struct {
	RootPath string `arg:"-r, env:SYNC_ROOT_PATH" help:"Root directory for sync filesystem in"`
}

// MakeConfig creates Config struct out of ENV variables and/or CLI arguments
func MakeConfig() (Config, error) {
	var c Config
	var e EnvVars
	var fsc FilesystemConfig

	err := arg.Parse(&e)
	if err != nil {
		return c, fmt.Errorf("[MakeConfig] %w", err)
	}

	if e.RootPath == "" {
		e.RootPath = "filesystem_root"
	}
	fsc.Root = e.RootPath

	c = Config{
		FS: fsc,
	}

	return c, nil
}
