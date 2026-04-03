package main

//go:generate confg --path=config.yaml --out=config_gen.go --package=main

func main() {
	// 这个文件用于演示 go:generate 的使用方式
	// 运行: go generate
	// 将会生成 config_gen.go 文件
}