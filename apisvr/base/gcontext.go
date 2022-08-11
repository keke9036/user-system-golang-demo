// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/12

package base

import (
	"net/http"
	"sync"
)

type GContext struct {
	Request *http.Request
	Writer  http.ResponseWriter
	URL     string
	Keys    map[string]any

	mu sync.RWMutex
}

func (c *GContext) Set(key string, value any) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *GContext) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}
