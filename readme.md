# bililive-danmaku

## Summary

This project is a bilibili live danmaku(bullet screen) backend

The current version of bililive-danmaku can run in the following mode:

+ Single mode: you can run this program directly on your computer, and it will display the danmaku of the specified room
+ Backend mode: you can run this program on your server, and it will provide http services

The backend mode provides the following services:

+ Connect to specified room of live.bilibili.com
+ Send danmaku in the specified room to specified frontend
+ Send gift-message in specified room to specified frontend

This program can also provide services for multiple rooms or front-end

## Features

+ Support single mode and backend mode, can be used with frontend/vscode-plugin/idea-plugin
+ High performance, run under multiple-goroutine
+ High robust, implements automatically restart when unknown error occurs
+ High compatibility, support short room id and numerous types of message
+ High availability, don't miss any message when client serving normally
+ Strong expansion, reserved framework for custom handlers

## How to use

### Single mode

**From source**

If you want to run this program from source, please confirm your go version is at least go1.18.6

Run the following scripts:

```shell
git clone https://github.com/sunist-c/bililive-danmaku.git danmaku
cd danmaku

# if you are a Mainland Chinese user, please use proxy
# export GOPROXY=https://goproxy.cn

go mod tidy
go run main.go -b=false
```

**With Docker**

Run the following scripts:

```shell
git clone https://github.com/sunist-c/bililive-danmaku.git danmaku
cd danmaku

docker build -t bililive-danmaku:latest .
docker run --rm --name=danmaku bililive-danmaku:latest
```

## Customize

Turn to [customize.md](./customize.md) 

## License

MIT License

## Contributing

1. fork the project
2. coding
3. send a pull request

## Thanks

Parts(websocket) of the project is inspired by [sh1luo/BiliDanmu_go](https://github.com/sh1luo/BiliDanmu_go)

## More

If the project helps you, please give me a star, thanks :)

## Related Resources

### Bililive-Danmaku Idea Plugin

**From Plugin Store**

<iframe width="245px" height="48px" src="https://plugins.jetbrains.com/embeddable/install/19967"/>

**From Source**

[sunist-c/bililive-danmaku](https://github.com/sunist-c/bililive-danmaku): a danmaku plugin in idea ides
