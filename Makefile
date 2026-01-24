PROJECT=$(shell grep module go.mod | rev | cut -d'/' -f1-2 | rev)
NAMESPACE=$(shell echo $(PROJECT) | cut -d'/' -f1)
REPO=$(shell echo $(PROJECT) | cut -d'/' -f2)

# Module path from go.mod used for ldflags -X assignments
MODULE=$(shell go list -m)

# If running in CI (GitHub Actions), use GITHUB_REF_NAME which contains the tag or branch name
TAG=$(GITHUB_REF_NAME)

# Default version, commit and date (can be overridden via env or in build targets)
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

ifneq ($(TAG),)
VERSION=$(TAG)
endif

# GO_BUILD_FLAGS injects build-time version info into the binary.
GO_BUILD_FLAGS=-ldflags="-s -w -X '$(MODULE)/cmd.version=$(VERSION)' -X '$(MODULE)/cmd.commit=$(COMMIT)' -X '$(MODULE)/cmd.date=$(DATE)'"

build:
	@go build $(GO_BUILD_FLAGS) .

docker:
	@docker build -t $(PROJECT) .

docker-compose:
	@docker compose -f ./deploy/docker-compose/docker-compose.yaml up

docs:
	@docker run --rm \
	  -v $$PWD:/tmp/project \
	  -w /tmp/project \
	  -e NODE_NO_WARNINGS=1 \
	  -e npm_config_loglevel=silent \
	  node:23-alpine3.20 \
	  npx @redocly/cli build-docs api/openapi.yaml --output internal/server/api.html

k8s:
	@tilt up --stream=true

lint:
	@docker run --rm \
	  -v $$PWD:/tmp/project \
	  -w /tmp/project \
	  golangci/golangci-lint:v1.62.0-alpine \
	  golangci-lint run -v

run:
	@go run . server

sdk:
	@go generate ./...

sure-docs-are-updated:
	@bash -c 'if [[ $$(git status --porcelain) ]]; then \
		echo "API docs not updated. Run make docs"; \
		exit 1; \
	fi'

template:
	@docker run --rm \
	  -v $$PWD:/tmp/project \
	  -w /tmp/project \
	  --entrypoint '' \
	  busybox \
	  /bin/sh -c 'grep -rl "circa10a/go-rest-template" . | xargs sed -i "s/circa10a\/go-rest-template/$(NAMESPACE)\/$(REPO)/g"; grep -rl "go-rest-template" . | xargs sed -i "s/go-rest-template/$(REPO)/g" '

test:
	@go test -v ./...
