PROJECT=$(shell grep module go.mod | rev | cut -d'/' -f1-2 | rev)
NAMESPACE=$(shell echo $(PROJECT) | cut -d'/' -f1)
REPO=$(shell echo $(PROJECT) | cut -d'/' -f2)
GO_BUILD_FLAGS=-ldflags="-s -w"

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
