package generate

// Options 定义代码生成的配置选项
type Options struct {
	// InputPath 配置文件路径 (YAML/JSON)
	InputPath string

	// OutputPath 生成的 Go 文件路径
	OutputPath string

	// PackageName 生成代码的包名
	PackageName string

	// WatchMode 是否启用监听模式（文件变化时自动重新生成）
	WatchMode bool

	// DryRun 只打印生成的代码，不写入文件
	DryRun bool

	// StructName 生成的结构体名称，默认为 "Config"
	StructName string

	// Tags 要生成的标签类型，如 ["json", "yaml", "mapstructure"]
	Tags []string
}

// Validate 验证选项的有效性
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