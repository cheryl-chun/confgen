package codegen

import (
	"bytes"
	"text/template"

	"github.com/cheryl-chun/confgen/internal/analyzer"
)

// singletonTemplate defines the boilerplate for thread-safe singleton management.
// It implements the 'sync.Once' pattern to ensure idempotent initialization 
// and provides a functional options API for flexible runtime configuration.
const singletonTemplate = `
import (
	"sync"

	"github.com/cheryl-chun/confgen/runtime"
)

var (
	instance *{{.RootStructName}}
	once     sync.Once
)

// Load initializes the global configuration singleton. It is thread-safe and 
// guaranteed to execute only once. Subsequent calls will not re-trigger loading.
//
// Example:
//    config.Load("config.yaml")
//    config.Load("config.yaml", config.WithEnv("APP_"))
func Load(path string, opts ...Option) error {
	var err error
	once.Do(func() {
		loader := runtime.NewLoader()
		loader.AddFile(path)

		// Apply functional options to customize the loader's behavior.
		for _, opt := range opts {
			opt.Apply(loader)
		}

		instance = &{{.RootStructName}}{}
		// Populate the struct fields using the runtime loader.
		err = loader.Fill(instance)
	})
	return err
}

// Get retrieves the initialized configuration singleton. 
// It triggers a panic if invoked before a successful call to Load().
func Get() *{{.RootStructName}} {
	if instance == nil {
		panic("config not loaded: ensure Load() is called during application bootstrap")
	}
	return instance
}

// MustLoad performs a mandatory configuration load and panics upon failure.
// Ideal for strict application startup sequences.
func MustLoad(path string, opts ...Option) {
	if err := Load(path, opts...); err != nil {
		panic(err)
	}
}

// Option defines the interface for functional configuration of the runtime loader.
type Option interface {
	Apply(*runtime.Loader)
}

type optionFunc func(*runtime.Loader)

func (f optionFunc) Apply(l *runtime.Loader) {
	f(l)
}

// WithEnv returns an Option that configures the loader to bind environment variables.
// Example: WithEnv("APP_") maps "APP_SERVER_PORT" to "server.port".
func WithEnv(prefix string) Option {
	return optionFunc(func(l *runtime.Loader) {
		l.AddEnv(prefix)
	})
}
`

// dynamicMethodsTemplate provides type-safe accessors and dynamic schema traversal 
// capabilities by delegating calls to the underlying runtime tree.
const dynamicMethodsTemplate = `
// Set updates the value at the specified hierarchical path with automated 
// type inference and source tracking.
func (c *{{.RootStructName}}) Set(path string, value interface{}, source runtime.SourceType) error {
	if c.ConfigTree == nil {
		return nil
	}
	return c.ConfigTree.Set(path, value, source)
}

// Watch registers a reactive callback to monitor value transitions at the given path.
// It returns a cancel function to stop the observer.
func (c *{{.RootStructName}}) Watch(path string, callback runtime.WatchCallback) func() {
	if c.ConfigTree == nil {
		return func() {}
	}
	return c.ConfigTree.Watch(path, callback)
}

// GetString performs a dynamic query for a string-typed value at the given path.
func (c *{{.RootStructName}}) GetString(path string) string {
	if c.ConfigTree == nil {
		return ""
	}
	return c.ConfigTree.GetString(path)
}

// GetInt performs a dynamic query for an integer-typed value at the given path.
func (c *{{.RootStructName}}) GetInt(path string) int {
	if c.ConfigTree == nil {
		return 0
	}
	return c.ConfigTree.GetInt(path)
}

// GetBool performs a dynamic query for a boolean-typed value at the given path.
func (c *{{.RootStructName}}) GetBool(path string) bool {
	if c.ConfigTree == nil {
		return false
	}
	return c.ConfigTree.GetBool(path)
}

// GetFloat performs a dynamic query for a floating-point value at the given path.
func (c *{{.RootStructName}}) GetFloat(path string) float64 {
	if c.ConfigTree == nil {
		return 0
	}
	return c.ConfigTree.GetFloat(path)
}

// Get retrieves the raw interface representation of the value at the specified path.
func (c *{{.RootStructName}}) Get(path string) interface{} {
	if c.ConfigTree == nil {
		return nil
	}
	return c.ConfigTree.Get(path)
}
`

// generateSingletonCodeWithTemplate renders the thread-safe initialization 
// logic by executing the singleton template with the inferred root metadata.
func (g *Generator) generateSingletonCodeWithTemplate(result *analyzer.AnalyzeResult) (string, error) {
	tmpl, err := template.New("singleton").Parse(singletonTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		RootStructName string
	}{
		RootStructName: result.RootStruct.Name,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateDynamicMethodsWithTemplate renders type-safe dynamic accessors 
// by executing the method template with the inferred root metadata.
func (g *Generator) generateDynamicMethodsWithTemplate(result *analyzer.AnalyzeResult) (string, error) {
	tmpl, err := template.New("methods").Parse(dynamicMethodsTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		RootStructName string
	}{
		RootStructName: result.RootStruct.Name,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}