# 对配置文件的释义

> `./`指的是程序的安装目录
>
> 配置文件默认被打包在`./config/config.toml`

> 现在来释义一些配置项的作用

---

以下的配置项`lopp_time`表示循环调用插件的时间，以小时计算`n*hour`，不包括`Init`插件，`install`则如注释所说，自带的插件一般被打包在`./plugins`

```toml
[main]
# 安装的插件
install = ["backup.so","upload.so","web_config.so","daemon.so"]
# 循环的时间，即备份开始的时间
lopp_time = 7200
```

#### 自带的插件定义的一些配置项

---

`root`和`static`实际上是等价的，取不同的名字主要是为了做区分，`database`就如字段名一样`数据库的类型·主机·端口·用户·密码·多数据库`,`upload`插件目前只支持将备份的文件上传到`Cos`

```toml
# 插件的配置由plugin.name.case组成,比如plugin.backup.file_path
[plugin.backup.file_path]
# 项目根路径
root = "/User/harder/html"
# 静态资源的路径
static = "/User/harder/static"
[plugin.backup.database]
driver = "mysql"
host = "localhost"
port = "3306"
user = "harder"
password = "83nnfd.."
databases = ["youyu"]
[plugin.upload.cos]
sId = "1"
sKey = "1"
bucketUrl = "1"
serviceUrl = "1"
```

