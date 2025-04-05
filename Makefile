swag-gen:
	swag init -g $(shell find internal/http -name "*.go" | head -n 1) -o docs --parseDependency --parseInternal
