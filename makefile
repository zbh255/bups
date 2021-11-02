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

source:
	mkdir -p $(build_path)/cache $(build_path)/plugins $(build_path)/log $(build_path)/config
	touch $(build_path)/log/bups.log
	touch $(build_path)/config/config.toml
	cp -r ./config/pro/config.toml $(build_path)/config/config.toml

build-plugins:
	@for str in $(wildcard ./plugins/*);\
	do \
	  	echo $$str | awk '{len=split($$0,a,"/");print a[len] > "./tmp.txt"}' ;\
        tmp=$$(cat ./tmp.txt) ;\
        echo '插件:'$$tmp':编译中' ;\
	  	$(GOBUILD) -buildmode=plugin -o $$tmp.so $$str/$$tmp.go;\
	  	echo '插件:'$$tmp':编译完成' ;\
		mv ./$$tmp.so $(build_path)/plugins;\
	done
	@rm -rf ./tmp.txt

build-darwin:export GOOS=darwin
build-darwin:export GOARCH=amd64
build-darwin:source
build-darwin:build-plugins
build-darwin:
	@echo 'GOOS='$(GOOS)
	@echo 'GOARCH='$(GOARCH)
	@$(GOBUILD) -o $(BINARY_MAIN_NAME) ./
	# 移动编译之后的文件
	@mv ./$(BINARY_MAIN_NAME) $(build_path)

# 编译linux版本的局部变量
build-linux:export GOOS=linux
build-linux:export GOARCH=amd64
build-linux:source
build-linux:build-plugins
build-linux:
	@echo 'GOOS='$(GOOS)
	@echo 'GOARCH='$(GOARCH)
	@$(GOBUILD) -o $(BINARY_MAIN_NAME) ./
	# 移动编译之后的文件
	@mv ./$(BINARY_MAIN_NAME) $(build_path)