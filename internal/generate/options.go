package generate

// Options encapsulates the configuration parameters required by the generator engine.
// It defines the behavior for source discovery, output destination, and metadata injection.
type Options struct {
	// InputPath specifies the source file path (e.g., YAML or JSON) 
	// containing the configuration schema to be parsed.
	InputPath string

	// OutputPath defines the target filesystem destination 
	// where the generated Go source code will be persisted.
	OutputPath string

	// PackageName specifies the Go package identifier for the generated artifacts. 
	// If omitted, it defaults to "main".
	PackageName string

	// WatchMode enables filesystem monitoring (hot-reloading) to trigger 
	// automatic regeneration whenever the source file undergoes changes.
	WatchMode bool

	// DryRun simulates the generation process by streaming the output 
	// to the console (stdout) without modifying any files on disk.
	DryRun bool

	// StructName defines the identifier for the generated configuration struct. 
	// Defaults to "Config" if not explicitly provided.
	StructName string

	// Tags specifies a collection of struct tag types (e.g., "json", "yaml") 
	// to be injected into the generated fields for serialization purposes.
	Tags []string
}

func (o *Options) Validate() error {
	if o.InputPath == "" {
		return ErrInputPathRequired
	}

	if o.OutputPath == "" && !o.DryRun {
		return ErrOutputPathRequired
	}

	if o.PackageName == "" {
		o.PackageName = "main"
	}

	if o.StructName == "" {
		o.StructName = "Config"
	}

	// 默认生成常用的 tags
	if len(o.Tags) == 0 {
		o.Tags = []string{"json", "yaml", "mapstructure"}
	}

	return nil
}