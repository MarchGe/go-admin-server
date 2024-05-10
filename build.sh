# ==============================
# 该脚本用于构建不同平台的二进制执行文件
# ==============================
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -o ./.build/go-admin_linux-amd64 -v ./main.go

CGO_ENABLED=0
GOOS=linux
GOARCH=arm64
go build -o ./.build/go-admin_linux-arm64 -v ./main.go