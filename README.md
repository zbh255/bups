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

### 使用这个程序

- 首先，您需要修改配置文件，演示为运行的最小配置

    ```toml
    cloud_api = "COS" # 云接口类型
    save_name = "host"
    save_time = 1 # 自动备份的时间,单位为小时
    ```

    > 之后配置您要备份的本地文件夹

    ```toml
    [local]
    web = "" 			#填写本地绝对路径,为空则该选项不备份
    static = ""   #填写本地绝对路径,为空则该选项不备份
    log = ""			#填写本地绝对路径,为空则该选项不备份
    ```

    > 对应云接口的鉴权参数,以腾讯COS为例子

    ```toml
    [bucket]
    bucket_url = ""
    secretid = ""
    secretkey = ""
    token = ""
    ```

    > 数据库备份不填写`db_name`则不备份

    ```toml
    [database]
    ipaddr = "127.0.0.1"
    port = "3306"
    user_name = "dbname"
    user_passwd = "dbpassword"
    db_name = "dbname1"
    db_name2 = ""
    ```

    > 加密选项

    ```toml
    [encryption]
    switch = "on" # on为开启,off为关闭
    encrypt_mode = "rsa" # aes为固定的key,即aes中填写的16位key;rsa为随机的aes_key并带有rsa保护
    aes = "1234567890123456" # 长度固定为16位
    ```

    > Web配置模块(未完成)

    ```toml
    [web_config]
    switch = "off" # on为开启,off为关闭
    ipaddr = "127.0.0.1"
    port = ":8080"
    user_name = "admin" # basicauth鉴权账号
    user_passwd = "admin" # basicauth鉴权密码
    ```

- 使用`systemctl`管理`service`运行: 您需要修改目录下的`bupsd.service`文件

    > 将`/usr/local/bups`替换为程序所在的目录

    ```shell
    [Unit]
    Description=bups backup service
    After=network.target
    
    [Service]
    Type=simple
    User=nobody
    Restart=1 
    RestartSec=10s
    ExecStart=/usr/local/bups/bups -s start
    ExecStop=/usr/local/bups/bups -s stop
    ExecReload=/usr/local/bups/bups -s restart
    
    [Install]
    WantedBy=multi-user.target
    ```

    

