# bups

![Go Report Card](https://goreportcard.com/badge/github.com/abingzo/bups)![GitHub](https://img.shields.io/github/license/abingzo/bups)

**一个Go语言写的用于备份Wordpress/Typecho网站数据至云端的小工具，支持自定义插件，目前只支持`linux`&`darwin`**

---

#### 构建

```shell
make build-linux
```

#### 基本配置

---

- [Config](./CONFIG.md)

#### 启动

---

```shell
./bups
```

使用自带的守护进程插件

```shell
./bups --plugin daemon --args '<-s start>'
```

#### 编写自己的插件

- [插件doc](./PLUGIN_DEV.md)

