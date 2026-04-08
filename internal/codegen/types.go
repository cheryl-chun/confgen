package codegen

// Options 代码生成选项
type Options struct {
	PackageName string // 包名
	AddComments bool   // 是否添加字段注释
}

// DefaultOptions 返回默认选项
func DefaultOptions() Options {
	return Options{
		PackageName: "main",
		AddComments: false,
	}
}