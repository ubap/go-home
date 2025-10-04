# Cross compile

On windows
```
$env:GOOS = ""; $env:GOARCH = "";
go build .
```


On Mac

-s -w linker flags strip debug information, making the binary smaller

```
GOOS=linux GOARCH=arm64 go build -o goHome-rpi -ldflags="-s -w" .
```

To copy to RPI, use for example mc (midnight commander)
See the script deploy.sh
