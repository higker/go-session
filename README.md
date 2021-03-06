# sessionx
Session library for Golang.

[![Go](https://github.com/higker/sessionx/actions/workflows/go-test.yml/badge.svg?event=push)](https://github.com/higker/sessionx/actions/workflows/go-test.yml)
[![codecov](https://codecov.io/gh/higker/sessionx/branch/master/graph/badge.svg?token=btbed5BUUZ)](https://codecov.io/gh/higker/sessionx)
[![DeepSource](https://deepsource.io/gh/higker/sessionx.svg/?label=active+issues&show_trend=true)](https://deepsource.io/gh/higker/sessionx/?ref=repository-badge)
[![DeepSource](https://deepsource.io/gh/higker/sessionx.svg/?label=resolved+issues&show_trend=true)](https://deepsource.io/gh/higker/sessionx/?ref=repository-badge)
[![License](https://img.shields.io/badge/license-MIT-db5149.svg)](https://github.com/higker/sessionx/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/higker/sessionx.svg)](https://pkg.go.dev/github.com/higker/sessionx)


# 介 绍
`sessionx`是适用于`go`的`web`编程的`session`中间件的库，你可以轻松得使用这个包来管理你的`session`。


1. 支持内存存储
2. 支持`redis`存储

## 获取安装库

```go
go get -u github.com/higker/sessionx
```

## 使用例子

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	sessionx "github.com/higker/sessionx"
)

var (
	cfg = &sessionx.Configs{
		TimeOut:        time.Minute * 30,
		RedisAddr:      "127.0.0.1:6379",
		RedisDB:        0,
		RedisPassword:  "redis.nosql",
		RedisKeyPrefix: sessionx.SessionKey,
		PoolSize:       100,
		Domain:         "localhost", // set domain by you
		Name:           sessionx.SessionKey,
		Path:           "/",
		Secure:         true,
		HttpOnly:       true,
	}
	// 如果使用内存存储就直接使用 sessionx.DefaultCfg
	// sessionx.New(sessionx.M, sessionx.DefaultCfg)
)

func main() {
	sessionx.New(sessionx.M, cfg)
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		// 存储K的值
		session.Set("K", time.Now().Format("2006 01-02 15:04:05"))
		fmt.Fprintln(writer, "set time value succeed.")
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		// 获取存储K的值
		v, err := session.Get("K")
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return
		}
		fmt.Fprintln(writer, fmt.Sprintf("The stored value is : %s", v))
	})

	http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)

		// MigrateSession 函数会迁移session会话数据，返回新的session
		session, err := session.MigrateSession()
		if err != nil {
			log.Println(err)
		}

		session.Set("person", "Jarvib Ding")
		fmt.Fprintln(writer, session)
	})
	_ = http.ListenAndServe(":8080", nil)
}

```
## 其他帮助

[点击查看: 本库设计和实现文章！](https://mp.weixin.qq.com/s/z_mLGZKXt0hO1l8UWjukUg)
