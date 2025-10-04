# Cross compile

On windows
```
$env:GOOS = ""; $env:GOARCH = "";
go build .
```


On Mac

-s -w linker flags strip debug information, making the binary smaller

```bash
GOOS=linux GOARCH=arm64 go build -o goHome-rpi -ldflags="-s -w" .
```

To copy to RPI, use for example mc (midnight commander)
See the script deploy.sh

```bash
./cmd/deploy.sh
```

To download logs

```bash
./cmd/logs.sh > logs.txt
```
