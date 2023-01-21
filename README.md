# CHECK24 in Go

## Build & Deployment
```shell
$env:GOARCH = "arm"
$env:GOARM = "6"
$env:GOOS = "linux"
go env

go build -o check24-in-go


rsync -az -e ssh check24-in-go [user]@[ip]:[directory]
```
