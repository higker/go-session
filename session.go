// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:10 PM - UTC/GMT+08:00

package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Session Unite struct
type Session struct {
	ID     string
	Expire time.Time
}

// Builder build  session store
func Builder(store StoreType, conf *Config) error {
	//if conf.MaxAge < DefaultMaxAge {
	//	return errors.New("session maxAge no less than 30min")
	//}
	_Cfg = conf
	switch store {
	default:
		return errors.New("build session error, not implement type store")
	case Memory:
		_Store = newMemoryStore()
		_Cfg._st = Memory
		return nil
	case Redis:
		redisStore, err := newRedisStore()
		if err != nil {
			return err
		}
		_Store = redisStore
		_Cfg._st = Redis
		return nil
	}
}

// Ctx return request session object
func Ctx(writer http.ResponseWriter, request *http.Request) (*Session, error) {
	// 检测是否有这个session数据
	cookie, err := request.Cookie(_Cfg.CookieName)
	// 如果没有session数据就重新创建一个
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		// 重新生成一个cookie 和唯一 sessionID
		nc, sid, err := generate(writer)
		if err != nil {
			return nil, err
		}
		fmt.Println(nc.Expires.UnixNano())
		return &Session{
			ID:     sid,
			Expire: nc.Expires,
		}, nil
	} else if cookie.Expires.UnixNano() < time.Now().UnixNano() {
		if checkID(cookie.Value) {
			fmt.Println("SID 有效")
			// 如果这个id在存储器里面存在就直接返回
			return &Session{ID: id, Expire: cookie.Expires}, nil
		}
	}

	// 防止浏览器关闭重新打开抛异常
	id, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}

	c, s := generate(writer)
	return &Session{
		ID:     s,
		Expire: c.Expires,
	}, nil
}

func generate(writer http.ResponseWriter) (*http.Cookie, string) {
	nc := newCookie(writer)
	return nc, nc.Value
}

// Get get session data by key
func (s *Session) Get(key string) ([]byte, error) {
	if key == "" || len(key) <= 0 {
		return nil, ErrorKeyNotExist
	}
	//var result Value
	//result.Key = key

	b, err := _Store.Reader(s.ID, key)
	if err != nil {
		return nil, err
	}
	//result.Value = b
	return b, nil
}

// Set set session data by key
func (s *Session) Set(key string, data interface{}) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	// 把id和到期时间传过去方便后面使用
	cv := map[string]interface{}{contextValueID: s.ID, contextValueExpire: &s.Expire}
	return _Store.Writer(context.WithValue(context.TODO(), contextValue, cv), key, data)
}

// Del delete session data by key
func (s *Session) Del(key string) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}

	_Store.Remove(s.ID, key)
	return nil
}

// Clean clean session data
func (s *Session) Clean(w http.ResponseWriter) {
	_Store.Clean(s.ID)
	cookie := &http.Cookie{
		Name:     _Cfg.CookieName,
		Value:    "",
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		Secure:   _Cfg.Secure,
		MaxAge:   int(_Cfg.MaxAge),
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge) * time.Second),
		HttpOnly: _Cfg.HttpOnly,
	}
	http.SetCookie(w, cookie)
}

func newCookie(w http.ResponseWriter) *http.Cookie {
	// 创建一个cookie
	s := Random(32, RuleKindAll)
	cookie := &http.Cookie{
		Name: _Cfg.CookieName,
		//这里是并发不安全的，但是这个方法已上锁
		Value:    string(s), //转义特殊符号@#￥%+*-等
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		HttpOnly: _Cfg.HttpOnly,
		Secure:   _Cfg.Secure,
		MaxAge:   int(_Cfg.MaxAge),                                         // 这个是按秒算的生命周期
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge) * time.Second), // 这个是具体的过期时间
	}
	http.SetCookie(w, cookie) //设置到响应中
	return cookie
}

// 检测sessionID是否有效
func checkID(id string) bool {
	return _Store.(*MemoryStore).values[id] == nil
}
