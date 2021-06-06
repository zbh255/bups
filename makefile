# Go
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_MAIN_NAME=bups
BINARY_RECOVER_NAME=bups_recover
BUILD=linx

# 环境变量
export GIN_MODE=release
export BUPS_MODE=release
export GOPROXY="https://goproxy.cn,direct"
export GO111MODULE=on
export CGO_ENABLED=0

# 文件与目录路径
build_path=./_bulid_release
config_version=pro
ifeq ($(BUPS_MODE),debug)
build_path=./_build_debug
config_version=dev
else
build_path=./_build_release
config_version=pro
endif
# path
Main_APP_CMD=./cmd/backup_application
Recover_APP_CMD=./cmd/recover_application
config_path=$(build_path)/conf/$(config_version)
# 项目配置文件的所在
project_config_path=./conf/$(config_version)
build_path_cache=$(build_path)/cache
build_path_cache_rsa=$(build_path_cache)/rsa
build_path_cache_backup=$(build_path_cache)/backup
build_path_cache_upload=$(build_path_cache)/upload
build_path_dir=$(build_path)/dir
build_path_log=$(build_path)/log

# file
config_path_file=app.conf.toml
build_path_log_app=app.log
build_path_log_gin=gin.log
build_path_log_service=service.log
build_path_dir_info=app_info.json
build_source_file=main.go

mod:
	@$(GOMOD) download -json

.PHONY:create
create:
	# 创建目录
	mkdir -p $(build_path_cache)
	mkdir -p $(build_path_cache_backup)
	mkdir -p $(build_path_cache_rsa)
	mkdir -p $(build_path_cache_upload)
	mkdir -p $(config_path)
	mkdir -p $(build_path_dir)
	mkdir -p $(build_path_log)
	# 创建文件
	touch $(build_path_log)/$(build_path_log_app)
	touch $(build_path_log)/$(build_path_log_gin)
	touch $(build_path_log)/$(build_path_log_service)
	touch $(build_path_dir)/$(build_path_dir_info)
	# 拷贝项目的文件
	cp $(project_config_path)/$(config_path_file) $(config_path)/$(config_path_file)

build-darwin:
	@echo $(BUPS_MODE)
	@echo $(GOOS)
	cd $(Main_APP_CMD) && $(GOBUILD) -o $(BINARY_MAIN_NAME) -ldflags '-s -w'
	# 移动编译之后的文件
	mv $(Main_APP_CMD)/$(BINARY_MAIN_NAME) $(build_path)

build-windows:
	echo $(Main_APP_CMD)

# 编译linux版本的局部变量
build-linux-main:export GOOS=linux
build-linux-main:export GOARCH=amd64
build-linux-main:
	# 编译主程序文件
	# 切换编译目录
	cd $(Main_APP_CMD) && $(GOBUILD) -o $(BINARY_MAIN_NAME) -ldflags '-s -w'
	# 移动编译之后的文件
	mv $(Main_APP_CMD)/$(BINARY_MAIN_NAME) $(build_path)

build-linux-sub:
	@echo $(GOOS)