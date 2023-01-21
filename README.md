# CHECK24 in Go

## Build & Deployment
```shell
$env:GOARCH = "arm"
$env:GOARM = "6"
$env:GOOS = "linux"
# $env:GOARCH = "amd64"
# $env:GOOS = "windows"
go env

go build -o check24-in-go


rsync -az -e ssh check24-in-go [user]@[ip]:[directory]
rsync -az -e ssh public [user]@[ip]:[directory]/public

rsync -az -e ssh check24-in-go test@172.30.124.39:/home/test
rsync -az -e ssh public test@172.30.124.39:/home/test
```
