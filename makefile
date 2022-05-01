build_cli:
	@echo "Building CLI"
	rm -rf bin/cli && mkdir -p bin/cli
	go build -o bin/cli/gsync cmd/cli/main.go
	@echo "Building CLI done"
