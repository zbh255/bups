# 设置环境变量
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOPROXY="https://goproxy.cn,direct"
export GIN_MODE=release
export BUPS_MODE='debug'
export APP_NAME='bups'

# 变量
build_path="./_bulid_1"
build_path_config=$build_path'/conf/dev'
project_path_config='./conf/dev'
conf_name='/app.conf.toml' #配置文件名称
if [ $BUPS_MODE == 'debug' ]
then
    build_path_config=$build_path'/conf/dev'
    project_path_config='./conf/dev'
elif [ $BUPS_MODE == 'release' ]
then
    project_path_config='./conf/pro'
    build_path_config=$build_path'/conf/pro'
fi

build_path_data=$build_path'/data'
build_path_cache=$build_path'/cache'
build_path_backup=$build_path_cache'/backup'
build_path_rsa=$build_path_cache'/rsa'
build_path_log=$build_path'/log'

# 打印信息
echo ${project_path_config}
echo ${build_path}
echo ${build_path_data}
echo ${build_path_cache}
echo ${build_path_backup}
echo ${build_path_rsa}
echo ${build_path_config}

# 创建文件夹
mkdir -p $build_path_data
mkdir -p $build_path_backup
mkdir -p $build_path_rsa
mkdir -p $build_path_config
mkdir -p $build_path_log

# 创建空文件
touch $build_path_log'/app.log'
touch $build_path_log'/gin.log'

# 赋值文件至编译缓存区
cp $project_path_config$conf_name $build_path_config$conf_name


# 编译文件
go build -v -work $APP_NAME
go build -o $APP_NAME -ldflags '-s -w'

# 移动编译的文件至指定位置
mv $APP_NAME $build_path