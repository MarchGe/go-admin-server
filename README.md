#### 项目简介
go-admin-server是go-admin项目的服务端程序，基于RBAC的权限模型，
集成运维管理等功能，支持Web Shell、一键部署等，基于go + gin + gorm开发。

#### 项目源码
- GitHub: https://github.com/MarchGe/go-admin-server
- Gitee: https://gitee.com/go-admin_1/go-admin-server

#### 在线体验
- https://go-admin.dy-technology.com
  ```text
  账号：root@example.com
  密码：123456
  ```

#### go-admin-server
- 配置文件
  - 支持使用nacos做配置中心或直接指定本地配置文件两种方式，以本地配置文件为例，启动示例：
    ```bash 
    go-admin-server server -c ./config.example.json
    ```
  项目中使用viper解析配置文件，配置文件格式支持json、yaml等，配置项内容：见config/config.go
- 启动说明
  - go-admin-server目前只支持单机部署，部署前，先执行数据库（MySQL）初始化脚本ml_admin.sql，
    服务在初次启动时，会自动创建root超级账号，账号信息：
      ```
      Email:     root@example.com
      Password:  123456
    ```
#### go-admin-agent
部署在各个主机上，用于收集主机性能数据，通过rpc上报给go-admin-server服务，go-admin-server中的
服务器性能监控功能，需要使用go-admin-agent上报的数据。使用示例：
```bash
go-admin-agent -H 192.168.1.10 -p 9080 -f 5
```
命令行参数含义，可通过-h查看帮助。