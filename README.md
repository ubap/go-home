# Quick start
```bash
make build
```

```bash
make test
```

To add user

```bash
make user
```

```bash
make deploy
```

To download logs

```bash
make logs
```


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