GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
NAME=$$(grep TerraformProviderName version.go | grep -o -P 'terraform-provider-[a-z]+')
TARGETS=darwin linux windows
VERSION=$$(grep TerraformProviderVersion version.go | grep -o -P '\d\.\d\.\d')

default: build

build:
	go build -o "bin/$(NAME)_v$(VERSION)-custom_x4"

example:
	go build -o "example/$(NAME)_v$(VERSION)-custom_x4"
	@read -p "Enter IronAuth token: " token && \
		cd ./example && \
		terraform init && \
		terraform apply -auto-approve -var "token=$$token" && \
		terraform apply -auto-approve -var "token=$$token" && \
		terraform destroy -auto-approve -var "token=$$token"

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test -v

init:
	go get ./...

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build \
		-o "dist/$@/$(NAME)_v$(VERSION)-custom_x4" \
		-a -ldflags '-extldflags "-static"'
	zip \
		-j "dist/$(NAME)_v$(VERSION)-custom_$@_amd64.zip" \
		"dist/$@/$(NAME)_v$(VERSION)-custom_x4"

.PHONY: build example fmt test init targets $(TARGETS)
