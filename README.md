# bups

一个Go语言写的用于备份Wordpress/Typecho网站数据至云端的小工具

`build for linux`
```shell script
make mod
make create
make build-linux-main
```

`build for darwin`
```shell script
make mod
make create
make build-darwin
```

`run`

```shell script
./bups # simple
./bups -s start # daemon
```

#### 使用`supervisor`管理进程

- 请参阅官方文档的配置: `http://supervisord.org/configuration.html`

#### 使用`systemd`管理进程

- 将项目中的`bupsd.service`拷贝至`service`文件目录

