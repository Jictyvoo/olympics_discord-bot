package confloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const envConfFile = "CONF_FILE"

// Load overlays a TOML file then env overrides onto defaults. The filename comes
// from the argument, else CONF_FILE; if both are empty the file step is skipped.
func Load[T any](filename string, defaults T, envBinder func(*T) error) (T, error) {
	cfg := defaults

	if filename == "" {
		filename = os.Getenv(envConfFile)
	}

	if filename != "" {
		if err := loadTOML(filename, &cfg); err != nil {
			return cfg, err
		}
	}

	if envBinder != nil {
		if err := envBinder(&cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

func loadTOML(path string, dst any) error {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("confloader: open %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	raw, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("confloader: read %s: %w", path, err)
	}

	if err = toml.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("confloader: parse %s: %w", path, err)
	}

	return nil
}
