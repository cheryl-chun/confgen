package runtime

import (
	"fmt"

	"github.com/cheryl-chun/confgen/internal/parser"
	"github.com/cheryl-chun/confgen/internal/tree"
)

// FileSource implements the Source interface for filesystem-backed 
// configuration ingestion. It supports loading data from various 
// persistent formats (YAML, JSON, etc.) defined by the underlying parser.
type FileSource struct {
	// Path specifies the absolute or relative location of the configuration file.
	Path string
}

// Load executes the ingestion lifecycle for the target file. It parses 
// the file into an intermediate representation and integrates it into 
// the primary configuration tree using source-aware replacement logic.
func (s *FileSource) Load(configTree *tree.ConfigTree) error {
	tempTree, err := s.parseTree()
	if err != nil {
		return err
	}
	
	// Atomically integrate the source-specific values into the global registry.
	configTree.ReplaceSource(tempTree, tree.SourceFile)
	return nil
}

// Priority returns the precedence level assigned to static file-based configuration.
func (s *FileSource) Priority() tree.SourceType {
	return tree.SourceFile
}

// parseTree delegates file parsing to the internal parser engine to generate 
// a temporary configuration tree with attributed source metadata.
func (s *FileSource) parseTree() (*tree.ConfigTree, error) {
	tempTree, err := parser.ParseToTree(s.Path, tree.SourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %w", s.Path, err)
	}
	return tempTree, nil
}

// mergeTrees orchestrates the integration of a source tree into the 
// destination hierarchy. This is primarily used for cross-tree aggregation.
func (s *FileSource) mergeTrees(dst, src *tree.ConfigTree) error {
	// Initiate a recursive walk from the source root to the destination tree.
	return s.mergeNode(dst, "", src.Root)
}

// mergeNode performs a depth-first traversal of the configuration hierarchy. 
// It recursively maps source nodes to their corresponding hierarchical paths 
// in the destination tree, ensuring strict source attribution.
func (s *FileSource) mergeNode(dst *tree.ConfigTree, prefix string, node *tree.ConfigNode) error {
	if node == nil {
		return nil
	}

	// Resolve the canonical hierarchical path for the current node.
	path := prefix
	if prefix != "" && node.Key != "root" {
		path = prefix + "." + node.Key
	} else if node.Key != "root" {
		path = node.Key
	}

	// If the current node encapsulates a scalar or terminal value, 
	// persist it to the destination tree.
	if node.HasValue() {
		if path != "" {
			dst.Set(path, node.GetValue(), tree.SourceFile, node.Type)
		}
	}

	// Recursively process complex structural child nodes (Object-style).
	for _, child := range node.Children {
		if err := s.mergeNode(dst, path, child); err != nil {
			return err
		}
	}

	// Recursively process sequential collection elements (Array-style).
	for _, item := range node.Items {
		if err := s.mergeNode(dst, path, item); err != nil {
			return err
		}
	}

	return nil
}