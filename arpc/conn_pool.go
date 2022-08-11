package arpc

import (
	"entry-task/conf"
	errorcode "entry-task/error"
	"errors"
	"sync"
	"time"
)

var (
	once sync.Once
	pool *ConnPool
)

type ConnPool struct {
	conns     chan *Client
	connCount int
	config    *conf.RpcClientConf
	lock      sync.Mutex
}

func NewConnPool(conf *conf.RpcClientConf) *ConnPool {
	if pool == nil { //防止加锁
		once.Do(func() {
			pool = &ConnPool{}
			if conf.MaxConnections < 1 {
				conf.MaxConnections = 1
			}
			pool.conns = make(chan *Client, conf.MaxConnections)
			pool.config = conf
			pool.connCount = 0
		})
	}

	return pool
}

func (c *ConnPool) Get() (*Client, error) {
	c.lock.Lock()
	if c.connCount < c.config.MaxConnections && len(c.conns) < 10 {
		c.connCount++
		c.lock.Unlock()
		return c.newClient()
	}
	c.lock.Unlock()
	select {
	case conn := <-c.conns:
		return conn, nil
	case <-time.After(time.Millisecond * c.config.ConnectTimeout):
		return nil, errors.New(errorcode.Timeout.GetMsg())
	}
}

func (c *ConnPool) Put(client *Client) {
	c.conns <- client
}

func (c *ConnPool) newClient() (*Client, error) {
	client, err := Dial("tcp", conf.RpcServerConfig.Addr, c.config.ConnectTimeout)
	return client, err
}
