package tree

// ValueType represents the underlying data type of a configuration node, 
// ensuring type safety during structural parsing and schema inference.
type ValueType int

const (
	TypeString ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeArray
	TypeObject
	TypeNull
)

// SourceType identifies the origin of a configuration value. 
// It is used by the merging engine to resolve conflicts between multiple config sources.
type SourceType int

const (
	// SourceDefault represents hardcoded fallback values within the application.
	SourceDefault       SourceType = iota
	// SourceRemote identifies values from external configuration providers (e.g., Etcd, Consul).
	SourceRemote
	// SourceFile identifies values loaded from local persistent storage (YAML, JSON, TOML).
	SourceFile
	// SourceRuntimeOverride represents dynamic changes applied via runtime APIs.
	SourceRuntimeOverride
	// SourceSessionEnv represents transient environment variables specific to the current process.
	SourceSessionEnv
	// SourceSystemEnv represents persistent, system-wide environment variables.
	SourceSystemEnv
	// SourceCodeOverride represents explicit programmatic assignments (the highest priority).
	SourceCodeOverride
)

// SourcePriority defines the precedence hierarchy for configuration resolution. 
// A higher integer value indicates a stronger precedence.
//
// Architectural Principles:
// 1. Persistence over Transience: System-wide env vars override process-specific session vars.
// 2. 12-Factor Compliance: Environment variables override local files for seamless container orchestration.
// 3. Local-First Development: Local files override remote stores to simplify debugging and local testing.
// 4. Intentionality: Explicit programmatic overrides (CodeOverride) represent the developer's 
//    absolute intent and thus possess the highest precedence.
var SourcePriority = map[SourceType]int{
	SourceDefault:         0,   // Baseline fallback.
	SourceRemote:          10,  // Dynamic remote configuration.
	SourceRuntimeOverride: 15,  // Mid-tier runtime updates via API.
	SourceFile:            20,  // Static local configuration files.
	SourceSessionEnv:      30,  // Ephemeral environment variables.
	SourceSystemEnv:       40,  // Persistent global environment variables.
	SourceCodeOverride:    100, // Absolute precedence for explicit code assignments.
}

// String returns the canonical string representation of the ValueType.
func (t ValueType) String() string {
	names := []string{
		"String", "Int", "Float", "Bool", "Array", "Object", "Null",
	}
	if t >= 0 && int(t) < len(names) {
		return names[t]
	}
	return "Unknown"
}

// String returns the canonical string representation of the SourceType.
func (s SourceType) String() string {
	names := []string{
		"Default",
		"Remote",
		"File",
		"RuntimeOverride",
		"SessionEnv",
		"SystemEnv",
		"CodeOverride",
	}
	if s >= 0 && int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}