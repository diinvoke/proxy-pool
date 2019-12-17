# 基于 Golang 的免费代理池

本项目基于 https://github.com/diinvoke/proxy-pool ，由于做了比较大的更改所以建立了个项目单独维护。

<!-- TOC -->

- [更新日志](#更新日志)
- [安装和部署](#安装和部署)
  - [代码层面，可以调用](#代码层面可以调用)
  - [Docker 方式部署](#docker-方式部署)
- [使用](#使用)
  - [配置](#配置)
  - [HTTP 接口](#http-接口)
- [扩展代理来源](#扩展代理来源)
- [@TODO](#todo)

<!-- /TOC -->

## 更新日志

* `20191217` 初始化版本
* `20191215` 扩展原有的项目，并重写部分代码，感谢 @diinvoke 的原项目

## 安装和部署

``` shell
go get github.com/mingcheng/proxypool
```

### 代码层面，可以调用

``` golang
	config := proxypool.Config{
		FetchInterval:   15 * time.Minute,
		CheckInterval:   2 * time.Minute,
		CheckConcurrent: 10,
	}

	go proxypool.Start(config)
	defer proxypool.Stop()
```

### Docker 方式部署

// ...

## 使用

// ...

### 配置

```golang
	config := proxypool.Config{
		FetchInterval:   15 * time.Minute,
		CheckInterval:   2 * time.Minute,
		CheckConcurrent: 10,
	}
```

### HTTP 接口

* `/all` 所有可用的代理列表
* `/random` 随机获取一个可用的代理

## 扩展代理来源

// ...

## @TODO

* 更多的免费代理来源
* 对接 Prometheus ，获取状态信息

`- eof -`