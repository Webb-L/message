# 消息服务

简单又好用的消息服务。快来给你"项目"添加消息功能。

## 启动

> 如果您需要测试该项目。我推荐您使用DockerCompose方法使用。

### 使用编译好的文件

> 打开[releases](https://github.com/Webb-L/message/releases)页面。
>
> 1.下载`other.zip`文件（启动消息服务需要的文件）
>
> 解压`other.zip`文件夹。修改`config/config.yaml`中`mysql`配置。
>
> 2.下载并启动消息服务（使用`linux`演示）
>
> 启动`./message-linux-amd64`。
>
> 测试是否启动成功`curl localhost:1204`。如果返回OK表示启动消息服务成功。

### Docker

您可以使用以下命令来拉取Docker镜像：

```bash
docker pull webbl1/message
```

然后，您可以修改`config.yaml`文件中的MySQL配置。

```yaml
app:
  # 网站标题
  title: Message
  # 是否为调试模式
  debug: true
  # 网站语言
  language: zh
  # 文件存储设置，设置上传文件的存储路径以及路由前缀
  store:
    path: ./uploads
    prefix: uploads

  # 开发环境：本地 EnvLocal / 测试 EnvTest / 生产 EnvProd
  env: local

  # 日志本地存储路径
  log:
    info: logs/info.log
    error: logs/error.log
    access: logs/access.log

database:
  host: 127.0.0.1
  port: 3306
  user: message
  pwd: message
  name: message
  max_idle_con: 5
  max_open_con: 10
  # params为驱动需要的额外的传参
  params:
    character: utf8mb4

api:
  # 是否开启SwaggerApi
  test: true
  # 返回最多数量
  maxLimit: 15
```

接下来，您可以使用以下命令在容器中运行该镜像，并将容器的端口1204映射到主机的端口1204上，并将容器命名为message：

```bash
docker run -p 1204:1204 -d --name message webbl1/message
```

然后，您可以使用以下命令将修改后的`config.yaml`文件复制到正在运行的容器中的`/app/config/config.yaml`位置：

```bash
docker cp config.yaml message:/app/config/config.yaml
```

最后，您可以使用curl命令来测试应用程序是否正常运行，例如检查/ping端点：

```bash
curl localhost:1204/ping
```

### DockerCompose

编辑`config.yaml`配置文件，把`config.yaml`文件中的`  host: 127.0.0.1`替换成`  host: mysql`

```yaml
app:
  # 网站标题
  title: Message
  # 是否为调试模式
  debug: true
  # 网站语言
  language: zh
  # 文件存储设置，设置上传文件的存储路径以及路由前缀
  store:
    path: ./uploads
    prefix: uploads

  # 开发环境：本地 EnvLocal / 测试 EnvTest / 生产 EnvProd
  env: local

  # 日志本地存储路径
  log:
    info: logs/info.log
    error: logs/error.log
    access: logs/access.log

database:
  host: 192.168.200.136
  port: 3306
  user: message
  pwd: message
  name: message
  max_idle_con: 5
  max_open_con: 10
  # params为驱动需要的额外的传参
  params:
    character: utf8mb4

api:
  # 是否开启SwaggerApi
  test: true
  # 返回最多数量
  maxLimit: 15
```

新建一个`docker-compose.yml`文件把下方的内容写入到文件中。

```yaml
version: '3'
services:
  app:
    image: webbl1/message:latest
    networks:
      - default
    links:
      - mysql
    ports:
      - "1204:1204"
    depends_on:
      - mysql
  mysql:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    networks:
      - default
    environment:
      MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
      MYSQL_DATABASE: 'message'
      MYSQL_USER: 'message'
      MYSQL_PASSWORD: 'message'
    volumes:
      - mysql-data:/var/lib/mysql
    ports:
      - "3306:3306"
volumes:
  mysql-data:
networks:
  default:
```

启动容器。

```bash
docker-compose up -d
```

把我们的文件修改好的`config.yaml`文件复制到`message_app_1`容器中。

```bash
docker cp config.yaml message_app_1:/app/config/config.yaml
```

重启消息服务

```bash
docker restart message_app_1
```

最后，您可以使用curl命令来测试应用程序是否正常运行，例如检查/ping端点：

```bash
curl localhost:1204/ping
```

## 开发

访问`SwaggerApi文档`。[http://localhost:1204/swagger/index.html](http://localhost:1204/swagger/index.html)

### 过滤语法

格式：

```text
列 比较 值
```

例子：

```text
title = 标题

title = 标题,content = 简易内容

status = 0|1|2

big_content = %内容%
```

列：

|                | 介绍     |
|----------------|--------|
| created_at     | 创建时间   | 
| updated_at     | 修改时间   | 
| sender_ids     | 发送者的id |
| title          | 标题     |
| content        | 简短内容   |
| category       | 类别     |
| big_content    | 长内容    |
| introducer_ids | 接受者的id |
| status         | 状态     |


比较：

|      | 介绍    |
|------|-------|
| \>   | 大于    | 
| =    | 等于    | 
| <    | 小于    |
| \>=  | 大于等于  |
| <=   | 小于等于  |
| !=   | 不等于   |
| like | 模糊比较  |
| in   | 多个值比较 |