# Bups

![Go Report Card](https://goreportcard.com/badge/github.com/abingzo/bups) ![GitHub](https://img.shields.io/github/license/abingzo/bups) ![GitHub](https://badgen.net/github/release/zbh255/bups) ![GitHub](https://github.com/zbh255/bups/actions/workflows/go.yml/badge.svg)

## About Bups

`bups`是一个Go实现的用于自动化备份配置文件数据和博客数据库的工具，最初是用于备份`WordPress/Typecho`这类博客应用

## Feature(功能)

- [ ] 自定义`Hook`的支持
- [x] 跨平台支持
- [x] 备份多数据库
- [x] 备份Typecho博客
- [x] 备份其他的文件
- [x] 容易拓展的`plugins`组件
- [x] 文件归档的支持
- [ ] 归档文件的加密
- [x] 归档文件上传至单一云端的支持

## Implement(实现)

`bups`的核心是一个组件资源的管理器与分配器,它负责管理和分配组件使用到的`Source`

`bups`内部定义了4种状态，到了设定的时间时，内部就会在除`Init`状态之外的状态直接切换，处于不同的状态则调用不同类型的组件。

```go
const (
	Init      Type = 0 // 初始化时调用的插件,程序运行周期只会调用一次
	BCollect  Type = 1 // 需要搜集备份数据时调用的插件
	BHandle   Type = 2 // 需要处理备份的数据时调用的插件
	BCallBack Type = 3 // 处理完备份数据时调用的插件
)
```

- Context
    - State Swtich
    - Handle Args
    - Set Plugin Source
- Plugins
    - backup
        - `Collect`
    - daemon
        - `Init`
    - encrypt
        - `Handle`
    - upload
        - `CallBack`
    - web_config
        - `Init`

## Usage(使用)

### Download(下载)

- [Github Release](https://github.com/zbh255/bups/releases)

### Config(配置)

> `bups`使用的是Toml格式的配置文件，如果不熟悉该格式的话，可以看这里

```toml
[project]
	# 安装的组件
	install = ["backup","upload","web_config","daemon","encrypt"]
	# 循环的时间，即备份开始的时间，以分钟来计算
	lopp_time = 14400

[project.log]
	access_log = "./access.log"
	error_log = "./error.log"

# 插件的配置由plugin.name.case组成,比如plugin.backup.file_path
[plugin.backup.file_path]
	# 要备份的文件路径
	root = "/User/harder/html"
	# 静态资源的路径
	static = "/User/harder/static"
	# 下面的名字可以随便写，zip文件以key命名
	nginx_conf = /etc/local/nginx/nginx.conf
	apache_conf = /etc/local/apache/apache.conf
[plugin.backup.database]
	# 要备份的数据库类型
	driver = "mysql"
	# 数据库主机
	host = "localhost"
	# 数据库端口
	port = "3306"
	# 指定的数据库用户名
	user = "harder"
	# 用户的密码
	password = "83nnfd.."
	# 要备份的库，可以备份多个，但要使用同一个用户
	databases = ["youyu"]
[plugin.upload.cos]
	# Tencent Cos相关，具体含义请查看腾讯云SDK文档
	sId = "1"
	sKey = "1"
	bucketUrl = "1"
	serviceUrl = "1"
```

### Start(启动)

默认的启动方式

```shell
./bups
```

使用自带的守护进程插件

```shell
./bups --plugin daemon --args '<-s start>'
```

`Supervisor`启动，将`/Users/xx`替换为程序所在的绝对路径，[原模板文件](./bupsd.ini)

```ini
[program:bupsd]
# run folder
directory=/Users/xx
# run command
command=/Users/xx/bups

autostart=true
autorestart=false
startsecs=3

user=root
stdout_logfile=/Users/xx/bups/logs/stdout.log
redirect_stderr=true
// stdout log size
stdout_logfile_maxbytes=30MB
```

`Systemctl`启动，修改[service模板文件](./bupsd.service)，将模板中的目录改为自己的安装目录即可，**注意：**由于`service`文件的编写时间比较久远可能有问题，这一部分抽时间再改

### Expand(拓展)

> `bups`的结构设计使你很容易就能编写自己的的组件，不过在此之前，你需要准备以下环境和知识
>
> - Go 编译器
> - Git的使用
> - Go基本知识
> - Go Mod的使用

查看程序的版本并拉取相应标签的代码到本地

```sh
./bups --option version
```

在`plugins`包下建立自己的包，比如`custom`，建立完成之后它们的结构看起来是这样的：`bups/plugins/custom/custom.go`，示例代码可以在这里找到:link:

```go
package example

import (
	"github.com/abingzo/bups/common/logger"
	"github.com/abingzo/bups/common/plugin"
	"time"
)

// Simple 实现plugin.Plugin接口
type Simple struct {
	name string
	sup []uint32
	_type plugin.Type
	stdLog logger.Logger
}

// New plugin包定义的必须实现的函数，用于注册组件使用
func New() plugin.Plugin {
	return &Simple{
		name:   "simple",
		sup:    []uint32{plugin.SUPPORT_STDLOG},
		_type:  plugin.Init,
		stdLog: nil,
	}
}

// Start 启动时调用的方法
func (s *Simple) Start(args []string) {
	go func() {
		time.Sleep(time.Second)
		s.stdLog.Info("my is " + s.name)
	}()
}

// Caller 接收可接受的Posix信号时调用的方法
func (s *Simple) Caller(single plugin.Single) {
	return
}

// GetName 返回组件的名字
func (s *Simple) GetName() string {
	return s.name
}

// GetType 返回组件的类型
func (s *Simple) GetType() plugin.Type {
	return s._type
}

// GetSupport 返回组件的需要的资源类型
func (s *Simple) GetSupport() []uint32 {
	return s.sup
}

// SetSource 设置从Context准备好的资源
// context会根据Support Slice声明的类型设置，没有声明的为nil
func (s *Simple) SetSource(source *plugin.Source) {
	s.stdLog = source.StdLog
}
```

在代码和配置文件中加载编写的组件并重新编译主程序

在plugin_load.go中的LoaderPlugin函数中有与下面的代码相似的语句的地方加上下面的代码即可加载

```go
iocc.RegisterPlugin(custom.New)
```

并在配置文件(config.toml)的`Project.Install`中添加`"custom"`

接下来就是重新构建程序了

```go
go build -o bups main.go
```

## Lisense

The Bups Use Mit licensed. More is See [Lisence](https://github.com/zbh255/bups/blob/dev/LICENSE)

