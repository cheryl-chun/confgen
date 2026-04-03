package main

import (
	"fmt"
	"os"

	"github.com/cheryl-chun/confgen/internal/generate"
	"github.com/spf13/cobra"
)

var (
	// 命令行参数
	pathFlag    string
	outFlag     string
	packageFlag string
	watchFlag   bool
	dryRunFlag  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "confg",
		Short: "Generate Go config structs from YAML/JSON files",
		Long: `confg is a code generator that creates type-safe Go configuration structs
from YAML or JSON configuration files, with runtime support for multiple config sources.

Features:
  • Auto-generate strongly-typed config structs
  • Support YAML and JSON formats
  • Trie-based runtime for efficient config queries
  • Multi-source configuration (file, env, etcd, zookeeper)
  • Hot reload support

Examples:
  # Generate config from YAML file
  confg --path=config.yaml --out=config.go

  # Specify package name
  confg --path=config.json --out=internal/config/config.go --package=config

  # Watch mode for development
  confg --path=config.yaml --out=config.go --watch

  # Dry run (print to stdout)
  confg --path=config.yaml --dry-run`,
		RunE: runGenerate,
	}

	// 定义命令行参数
	rootCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to config file (YAML/JSON) [required]")
	rootCmd.Flags().StringVarP(&outFlag, "out", "o", "config_gen.go", "Output file path")
	rootCmd.Flags().StringVar(&packageFlag, "package", "main", "Package name for generated code")
	rootCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch config file and regenerate on changes")
	rootCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Print generated code to stdout without writing file")

	// 标记必需参数
	rootCmd.MarkFlagRequired("path")

	// 添加版本命令
	rootCmd.AddCommand(versionCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// 构建生成选项
	opts := generate.Options{
		InputPath:   pathFlag,
		OutputPath:  outFlag,
		PackageName: packageFlag,
		WatchMode:   watchFlag,
		DryRun:      dryRunFlag,
	}

	// 执行生成
	return generate.Run(opts)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("confg v0.1.0")
		fmt.Println("A config struct generator for Go")
	},
}