package generate

import "errors"

var (
	ErrInputPathRequired = errors.New("input path is required")
	ErrOutputPathRequired = errors.New("output path is required when not in dry-run mode")
	ErrUnsupportedFileType = errors.New("unsupported file type, only YAML and JSON are supported")
	ErrFileNotFound = errors.New("config file not found")
	ErrInvalidConfig = errors.New("invalid config file format")
	ErrCodeGeneration = errors.New("code generation failed")
)