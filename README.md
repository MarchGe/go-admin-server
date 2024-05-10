### 执行方式

1. 先执行build.bat（或build.sh）编译为二进制文件，编译后生成的二进制文件在.build/目录下。
2. 如果采用Docker部署，在项目根目录下执行 `docker build` ，如：

    ```bash
    docker build -t go-admin-server:1.0.0 .
    ```

   默认Dockerfile中会自动寻找编译生成的 `go-admin-server_linux-amd64` 二进制文件，如果在非linux或非amd64处理器架构下运行服务，需要手动修改Dockerfile文件中对应的二进制文件。