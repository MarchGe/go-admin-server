SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./.build/go-admin-agent_linux-amd64 -v ./main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
go build -o ./.build/go-admin-agent_linux-arm64 -v ./main.go