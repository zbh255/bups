# Go
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_MAIN_NAME=bups


# Source
build_path = ./build_release
plugin_path = ./plugins

# build app info
# 需要写入的版本信息
gitTag=$(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
buildDate=$(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit=$(shell git rev-parse --short HEAD)
gitTreeState=$(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-s -w -X main.buildDate=${buildDate} -X main.gitCommit=${gitCommit} -X main.gitTreeState=${gitTreeState} -X main.version=${gitTag} -X main.gitBranch=${gitBranch}"

clean:
	@rm -rf $(build_path)

source:
	mkdir -p $(build_path)/cache
	touch $(build_path)/config.toml
	cp -r ./config/pro/config.toml $(build_path)/config.toml

build-plugins:
	@for str in $(wildcard ./plugins/*);\
	do \
	  	echo $$str | awk '{len=split($$0,a,"/");print a[len] > "./tmp.txt"}' ;\
        tmp=$$(cat ./tmp.txt) ;\
        echo '插件:'$$tmp':编译中' ;\
	  	$(GOBUILD) -buildmode=plugin -gcflags="all=-N -l" -o $$tmp.so $$str/$$tmp.go;\
	  	echo '插件:'$$tmp':编译完成' ;\
		mv ./$$tmp.so $(build_path)/plugins;\
	done
	@rm -rf ./tmp.txt

shell-build:
	@echo 'GOOS='$(GOOS)
	@echo 'GOARCH='$(GOARCH)
	# generate
	$(GOCMD) generate ./
	$(GOBUILD) -gcflags="all=-N -l" -o $(BINARY_MAIN_NAME) -ldflags ${ldflags} ./
	# 移动编译之后的文件
	mv ./$(BINARY_MAIN_NAME) $(build_path)

.PHONY:test
test:
	make -I ./test delete_file -C ./test
	make -I ./test create_file -C ./test
	$(GOCMD) test -race -v ./...
	make -I ./test delete_file -C ./test

test-env:
	make -v
	mysql -v
	mysqldump -v

test-clean:
	make -I ./test delete_file -C ./test

build-darwin:export GOOS=darwin
build-darwin:export GOARCH=amd64
build-darwin:source
build-darwin:shell-build

# 编译linux版本的局部变量
build-linux:export GOOS=linux
build-linux:export GOARCH=amd64
build-linux:source
build-linux:shell-build

build-windows:export GOOS=windows
build-windows:export GOARCH=amd64
build-windows:source
build-windows:shell-build
