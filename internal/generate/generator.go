package generate

import (
	"fmt"
	"os"

	"github.com/cheryl-chun/confgen/internal/analyzer"
	"github.com/cheryl-chun/confgen/internal/codegen"
	"github.com/cheryl-chun/confgen/internal/parser"
)

// Run 执行代码生成
func Run(opts Options) error {
	// 1. 验证选项
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	// 2. 检查输入文件是否存在
	if _, err := os.Stat(opts.InputPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, opts.InputPath)
	}

	// 3. 使用工厂模式获取解析器并解析文件
	fmt.Printf("[1/4] Parsing config file: %s\n", opts.InputPath)

	result, err := parser.ParseFile(opts.InputPath)
	if err != nil {
		if err == parser.ErrParserNotFound || err == parser.ErrUnsupportedFormat {
			return fmt.Errorf("%w: %s (supported: %v)",
				ErrUnsupportedFileType, opts.InputPath, parser.SupportedFormats())
		}
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("     ✓ Parsed successfully (root has %d fields)\n", len(result.Root.Children))

	// 4. 类型推断和分析
	fmt.Printf("[2/4] Analyzing types...\n")
	analyzeResult, err := analyzer.Analyze(result.Root)
	if err != nil {
		return fmt.Errorf("failed to analyze types: %w", err)
	}
	fmt.Printf("     ✓ Generated %d struct definitions\n", len(analyzeResult.SubStructs)+1)

	// 5. 生成代码
	fmt.Printf("[3/4] Generating code...\n")
	codegenOpts := codegen.Options{
		PackageName: opts.PackageName,
		AddComments: false,
	}
	code, err := codegen.Generate(analyzeResult, codegenOpts)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}
	fmt.Printf("     ✓ Generated %d lines of code\n", countLines(code))

	// 6. 输出结果
	if opts.DryRun {
		// 打印到 stdout
		fmt.Printf("[4/4] Output mode: stdout\n\n")
		fmt.Println(code)
	} else {
		// 写入文件
		fmt.Printf("[4/4] Writing to: %s\n", opts.OutputPath)
		if err := os.WriteFile(opts.OutputPath, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("     ✓ Successfully written to %s\n", opts.OutputPath)
	}

	fmt.Println("\n✓ Generation completed successfully!")
	return nil
}

// countLines 统计代码行数
func countLines(s string) int {
	count := 0
	for _, ch := range s {
		if ch == '\n' {
			count++
		}
	}
	return count
}