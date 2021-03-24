# 设置环境变量
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOPROXY="https://goproxy.cn,direct"
export GIN_MODE=release

# 创建文件夹
mkdir -p ./_bulid/data
mkdir -p ./_bulid/cache/backup
mkdir -p ./_bulid/cache/rsa
mkdir -p ./_bulid/conf/dev
mkdir -p ./_bulid/log

# 创建空文件
touch ./_bulid/log/app.log
touch ./_bulid/log/gin.log

# 赋值文件至编译缓存区
cp ./conf/dev/app.conf.toml ./_bulid/conf/dev/app.conf.toml


# 编译文件
go build -o bups

# 移动编译的文件至指定位置
mv bups ./_bulid