IMG = inf-repo-registry.cn-wulanchabu.cr.aliyuncs.com/infly-dev/i-device-plugin:v0.1.0

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/i-device-plugin cmd/main.go

.PHONY:build-image
build-image:
	docker build -t ${IMG} .