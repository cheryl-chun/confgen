package generate

import (
	"fmt"
	"os"

	"github.com/cheryl-chun/confgen/internal/analyzer"
	"github.com/cheryl-chun/confgen/internal/codegen"
	"github.com/cheryl-chun/confgen/internal/parser"
)

// Run orchestrates the entire code generation pipeline. 
// It follows a multi-stage process: validation, parsing, static analysis, 
// and code rendering. It returns a wrapped error if any stage of the pipeline fails.
func Run(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	if _, err := os.Stat(opts.InputPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, opts.InputPath)
	}

	fmt.Printf("[1/4] Parsing config file: %s\n", opts.InputPath)

	// Parse the source file into an Intermediate Representation (IR).
	// This stage abstracts away the specific format (YAML/JSON).
	result, err := parser.ParseFile(opts.InputPath)
	if err != nil {
		if err == parser.ErrParserNotFound || err == parser.ErrUnsupportedFormat {
			return fmt.Errorf("%w: %s (supported: %v)",
				ErrUnsupportedFileType, opts.InputPath, parser.SupportedFormats())
		}
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("     ✓ Parsed successfully (root has %d fields)\n", len(result.Root.Children))

	fmt.Printf("[2/4] Analyzing types...\n")
	analyzeResult, err := analyzer.Analyze(result.Root)
	if err != nil {
		return fmt.Errorf("failed to analyze types: %w", err)
	}
	fmt.Printf("     ✓ Generated %d struct definitions\n", len(analyzeResult.SubStructs)+1)

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

	if opts.DryRun {
		fmt.Printf("[4/4] Output mode: stdout\n\n")
		fmt.Println(code)
	} else {
		fmt.Printf("[4/4] Writing to: %s\n", opts.OutputPath)
		if err := os.WriteFile(opts.OutputPath, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("     ✓ Successfully written to %s\n", opts.OutputPath)
	}

	fmt.Println("\n✓ Generation completed successfully!")
	return nil
}

func countLines(s string) int {
	count := 0
	for _, ch := range s {
		if ch == '\n' {
			count++
		}
	}
	return count
}