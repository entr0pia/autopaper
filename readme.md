## golang实现的跨平台自动切换壁纸程序

每15分钟切换一次壁纸, 来源: [爱壁纸UWP](https://www.microsoft.com/zh-mo/p/%E7%88%B1%E5%A3%81%E7%BA%B8uwp/9nblggh5kccf)

### 安装依赖
**推荐使用go mod**
```shell
go env -w GO111MODULE=on
go mod download
```

~~go get github.com/reujab/wallpaper~~

### 交叉编译

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui
```