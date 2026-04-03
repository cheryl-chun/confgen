package generate

import "errors"

var (
	// ErrInputPathRequired 输入路径为空
	ErrInputPathRequired = errors.New("input path is required")

	// ErrOutputPathRequired 输出路径为空（非 dry-run 模式）
	ErrOutputPathRequired = errors.New("output path is required when not in dry-run mode")

	// ErrUnsupportedFileType 不支持的文件类型
	ErrUnsupportedFileType = errors.New("unsupported file type, only YAML and JSON are supported")

	// ErrFileNotFound 文件不存在
	ErrFileNotFound = errors.New("config file not found")

	// ErrInvalidConfig 配置文件格式错误
	ErrInvalidConfig = errors.New("invalid config file format")

	// ErrCodeGeneration 代码生成失败
	ErrCodeGeneration = errors.New("code generation failed")
)