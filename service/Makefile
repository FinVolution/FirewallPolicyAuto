BIN_NAME := service_api

PHONY: dev d
dev d:
	@ENV=dev go run .

PHONY: dist
dist:
	@if [ -f "bin/${BIN_NAME}" ];then rm -rf bin/${BIN_NAME}; else echo "OK!"; fi \
	&& go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct \
	&& go env -w CGO_ENABLED=0 && go env && go mod tidy \
	&& go build -trimpath -gcflags "-N -l" -ldflags "-s -w" -o bin/${BIN_NAME}

PHONY: clean c
clean c:
	@rm -rf .DS_Store ${BIN_NAME} bin/ logs/
	@find . -type f -name '*.md' -or -name '*.log' -print0 | xargs -0 $(RM) -v
	@find ./vendor -type f \( -name '.travis.yml' -o -name 'AUTHORS' -o -name 'LICENSE' -o -name '.gitignore' -o -name '*.md' \) -print0 | xargs -0 $(RM) -v

PHONY: clean-win
clean-win:
	fd -p 'vendor/*' -e md -e yml -x rm
