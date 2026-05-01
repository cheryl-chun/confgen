package runtime

import (
	"os"
	"strings"

	"github.com/cheryl-chun/confgen/internal/tree"
)

// EnvSource implements the Source interface for loading configuration parameters 
// from the operating system's environment variables.
type EnvSource struct {
	// Prefix defines an optional namespace (e.g., "APP_") to filter 
	// and isolate relevant environment variables.
	Prefix string 
}

// Load iterates through the process environment block and merges matching 
// variables into the provided ConfigTree. It handles prefix stripping and 
// path normalization automatically.
func (s *EnvSource) Load(configTree *tree.ConfigTree) error {
	// Retrieve the full environment block as a slice of "KEY=VALUE" strings.
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Filter variables that do not match the specified namespace prefix.
		if s.Prefix != "" && !strings.HasPrefix(key, s.Prefix) {
			continue
		}

		// Strip the namespace prefix to obtain the raw configuration key.
		if s.Prefix != "" {
			key = strings.TrimPrefix(key, s.Prefix)
		}

		// Transform the uppercase underscore-delimited environment key 
		// into a lowercase dot-delimited configuration path.
		// e.g., "SERVER_HOST" -> "server.host"
		path := s.envKeyToPath(key)

		// Persist the value to the tree with System Environment precedence.
		// Note: All environment variables are treated as strings during initial ingestion.
		configTree.Set(path, value, tree.SourceSystemEnv, tree.TypeString)
	}

	return nil
}

// Priority returns the precedence level assigned to environment-sourced configuration.
func (s *EnvSource) Priority() tree.SourceType {
	return tree.SourceSystemEnv
}

// envKeyToPath converts an environment variable identifier into a hierarchical 
// configuration path using lowercase normalization and delimiter replacement.
// Transformation Example: "SERVER_HOST" -> "server.host"
func (s *EnvSource) envKeyToPath(key string) string {
	// Normalize to lowercase and replace underscores with dot separators.
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", ".")
	return key
}