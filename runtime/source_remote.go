package runtime

import (
	"strings"

	"github.com/cheryl-chun/confgen/internal/tree"
)

// RemoteConfigSource serves as the foundational abstraction for remote configuration 
// providers (e.g., Etcd, Consul, Zookeeper). It encapsulates shared logic for 
// namespace management and bidirectional key-to-path transformation.
type RemoteConfigSource struct {
	// Prefix defines the root namespace for the configuration store.
	Prefix string
}

// normalizePrefix performs a sanitization pass on the configured prefix, 
// stripping leading and trailing delimiters to ensure consistent matching logic.
func (s *RemoteConfigSource) normalizePrefix() string {
	return strings.Trim(s.Prefix, "/.")
}

// KeyToPath transforms a raw remote provider key into a canonical configuration 
// tree path. It handles namespace stripping, delimiter translation (e.g., '/' to '.'), 
// and case normalization to ensure cross-source compatibility.
func (s *RemoteConfigSource) KeyToPath(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}

	prefix := s.normalizePrefix()
	if prefix != "" {
		// Normalize both the key and the prefix for robust comparison.
		trimmedKey := strings.TrimLeft(key, "/.")
		trimmedPrefix := strings.TrimLeft(prefix, "/.")
		
		// If the key is within the defined namespace, strip the prefix.
		if strings.HasPrefix(trimmedKey, trimmedPrefix) {
			key = strings.TrimPrefix(trimmedKey, trimmedPrefix)
		} else {
			key = trimmedKey
		}
	}

	// Clean up residual delimiters from both ends.
	key = strings.TrimLeft(key, "/.")
	key = strings.TrimRight(key, "/.")
	if key == "" {
		return ""
	}

	// Translate Unix-style or environment-style separators to dot-notation.
	key = strings.ReplaceAll(key, "/", ".")
	key = strings.ReplaceAll(key, "_", ".")
	key = strings.Trim(key, ".")
	
	// Ensure the internal path is strictly lowercase for case-insensitive resolution.
	return strings.ToLower(key)
}

// PathToKey performs the inverse mapping, converting a hierarchical tree path 
// back into a delimited string suitable for remote key-value stores.
func (s *RemoteConfigSource) PathToKey(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return s.normalizePrefix()
	}

	// Translate internal dot-notation back to standard Unix-style hierarchical keys.
	key := strings.ReplaceAll(path, ".", "/")
	key = strings.Trim(key, "/")
	
	prefix := s.normalizePrefix()
	if prefix == "" {
		return key
	}
	
	// Inject the namespace prefix into the resulting key.
	return prefix + "/" + key
}

// ConfigurePrefix updates the namespace prefix for the remote source.
func (s *RemoteConfigSource) ConfigurePrefix(prefix string) {
	s.Prefix = prefix
}

// Priority returns the default precedence level assigned to remote configuration 
// sources during the merge process.
func (s *RemoteConfigSource) Priority() tree.SourceType {
	return tree.SourceRemote
}