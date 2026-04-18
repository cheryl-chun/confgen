.PHONY: help build install test clean run example benchmark

help:
	@echo "Available commands:"
	@echo "  make build      - Build the confg binary"
	@echo "  make install    - Install confg to GOPATH/bin"
	@echo "  make test       - Run tests"
	@echo "  make benchmark  - Run all benchmarks and save results"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Run example generation"
	@echo "  make example    - Generate config from example YAML"

build:
	@echo "Building confg..."
	go build -o bin/confg ./cmd/confg

install:
	@echo "Installing confg..."
	go install ./cmd/confg

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f examples/basic/config_gen.go

run: build
	@echo "Running confg..."
	./bin/confg --help

example: build
	@echo "Generating config from examples/basic/config.yaml..."
	./bin/confg --path=examples/basic/config.yaml --out=examples/basic/config_gen.go --package=main

benchmark:
	@echo "Running complete benchmark suite with correctness verification..."
	@mkdir -p test/results
	@echo "Benchmark Results - $(shell date)" > test/results/benchmark_results.txt
	@echo "" >> test/results/benchmark_results.txt
	@echo "=== CORRECTNESS VERIFICATION ===" >> test/results/benchmark_results.txt
	@echo "" >> test/results/benchmark_results.txt
	@go test ./test -run=TestCorrectness_ConfgenVsViper -v 2>&1 | tee -a test/results/benchmark_results.txt
	@echo "" >> test/results/benchmark_results.txt
	@echo "=== PERFORMANCE BENCHMARKS ===" >> test/results/benchmark_results.txt
	@echo "" >> test/results/benchmark_results.txt
	@for bench in BenchmarkGenerated_Confgen_Load BenchmarkGenerated_Viper_Load BenchmarkGenerated_Confgen_FieldAccess BenchmarkGenerated_Viper_FieldAccess BenchmarkGenerated_Confgen_BatchAccess BenchmarkGenerated_Viper_BatchAccess BenchmarkGenerated_Confgen_AllFields BenchmarkGenerated_Viper_AllFields; do \
		echo "=== $$bench ===" >> test/results/benchmark_results.txt; \
		go test ./test -bench="^$$bench\$$" -benchmem -run='none' 2>&1 | grep -E "Benchmark|ns/op|allocs" >> test/results/benchmark_results.txt; \
		echo "" >> test/results/benchmark_results.txt; \
	done
	@echo "✓ Full benchmark results saved to test/results/benchmark_results.txt"
	@echo ""
	@echo "===== SUMMARY ====="
	@grep -E "✅|✓|ns/op|allocs/op" test/results/benchmark_results.txt | head -30

.DEFAULT_GOAL := help