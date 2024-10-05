PROJECT_SRC := "./"									# main.go 所在目录
PROJECT_NAME := "example" 							# 自行改动
PROJECT_REPO := "github.com/qiaogy91/ioc/example" 		# 自行改动

.PHONY: dep build clean gen

all: clean gen dep build


dep:
	@go mod tidy

build: dep
	@GOOS=linux GOARCH=amd64 go build -o build/${PROJECT_NAME}-linux-x86_64 ${PROJECT_SRC}
	@GOOS=darwin GOARCH=amd64 go build -o build/${PROJECT_NAME}-darwin-x86_64 ${PROJECT_SRC}
	@GOOS=windows GOARCH=amd64 go build -o build/${PROJECT_NAME}-windows-x86_64.exe ${PROJECT_SRC}

clean:
	@rm -fr build/*

gen:
	@protoc -I=./ --go_out=. --go-grpc_out=. --go_opt=module=${PROJECT_REPO} --go-grpc_opt=module=${PROJECT_REPO}  apps/*/pb/*.proto
	@protoc-go-inject-tag -input="apps/*/*.pb.go"


run:
	@go run main.go start -t file -f etc/conf.toml || true