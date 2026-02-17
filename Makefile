Name = fkmcps
Version = 0.0.1
BuildTime = $(shell date +'%Y-%m-%d %H:%M:%S')

CURRENT_OS = $(shell go env GOOS)
CURRENT_ARCH = $(shell go env GOARCH)

LDFlags = -ldflags "-s -w -X '${Name}/version.version=$(Version)' -X '${Name}/version.buildTime=${BuildTime}'"

targets ?= darwin:arm64 windows:amd64 linux:amd64

.DEFAULT_GOAL := native

native:
	@$(MAKE) build t="$(CURRENT_OS):$(CURRENT_ARCH)"

all:
	@$(MAKE) build t="$(targets)"

build:
	@if [ -z "$(t)" ]; then \
		echo "error: please specify a target, e.g., make build t=linux:amd64"; \
		exit 1; \
	fi
	@$(foreach n, $(t),\
		os=$$(echo "$(n)" | cut -d : -f 1);\
		arch=$$(echo "$(n)" | cut -d : -f 2);\
		suffix=""; \
		if [ "$${os}" = "windows" ]; then suffix=".exe"; fi; \
		output_name="./release/${Name}_$${os}_$${arch}$${suffix}"; \
		echo "Compiling: $${os}/$${arch}..."; \
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} go build -trimpath $(LDFlags) -o $${output_name} ./main.go;\
		echo "Compilation finished: $${output_name}";\
	)

clean:
	rm -rf ./release

.PHONY: native all build clean