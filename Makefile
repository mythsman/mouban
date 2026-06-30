SWAG := $(shell go env GOPATH)/bin/swag

.PHONY: install-tools build-frontend build-swagger build-backend build

install-tools:
	@command -v $(SWAG) >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@v1.16.6

build-frontend:
	@if [ -d "web" ]; then \
		cd web; \
		if [ -f "package-lock.json" ]; then npm ci --no-audit --no-fund; else npm install --no-audit --no-fund; fi; \
		npm run build; \
	fi

build-swagger: install-tools
	@$(SWAG) init --outputTypes json,yaml -g main.go -o docs
	@mkdir -p build/swagger
	@cp docs/swagger.json build/swagger/swagger.json
	@cp tools/swagger/index.html build/swagger/index.html

build-backend:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main .

build: build-frontend build-swagger build-backend
