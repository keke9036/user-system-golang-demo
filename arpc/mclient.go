// @Description 支持和多个server通信
// @Author weitao.yin@shopee.com
// @Since 2022/7/21

package arpc

import (
	"context"
	"entry-task/util"
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect SelectMode = iota // random
)

type MClient struct {
	r             *rand.Rand
	originalAddrs []string
	addrs         []string
	mode          SelectMode
	timeout       time.Duration
	mu            sync.Mutex
	clients       map[string]*Client
}

var _ io.Closer = (*MClient)(nil)

func NewMClientAndInit(addrs []string, mode SelectMode, timeout time.Duration) *MClient {
	mclient := MClient{
		r:             rand.New(rand.NewSource(time.Now().UnixNano())),
		originalAddrs: addrs,
		addrs:         addrs,
		mode:          mode,
		timeout:       timeout,
		clients:       make(map[string]*Client),
	}

	for _, addr := range mclient.addrs {
		client, err := mclient.dial(addr)
		if err != nil || client == nil {
			util.Logger.Errorf("Init client error, %s", addr)
			return nil
		}
	}

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			<-ticker.C
			mclient.detectServers()
		}
	}()

	return &mclient
}

func (mc *MClient) detectServers() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for _, addr := range mc.originalAddrs {
		if _, ok := mc.clients[addr]; ok {
			continue
		}
		client, err := Dial("tcp", addr, mc.timeout)
		if err == nil {
			return
		}
		mc.clients[addr] = client
	}
}

func (mc *MClient) Close() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for key, client := range mc.clients {
		_ = client.Close()
		delete(mc.clients, key)
	}
	return nil
}

func (mc *MClient) dial(rpcAddr string) (*Client, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	client, ok := mc.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(mc.clients, rpcAddr)
		client = nil
	}
	if client == nil {
		var err error
		client, err = Dial("tcp", rpcAddr, mc.timeout)
		if err != nil {
			return nil, err
		}
		mc.clients[rpcAddr] = client
	}
	return client, nil
}

//func (mc *MClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
//	client, err := mc.dial(rpcAddr)
//	if err != nil {
//		return err
//	}
//	return client.Call(ctx, serviceMethod, args, reply)
//}

func (mc *MClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// get a client
	client, err := mc.getClient(mc.mode)
	if err != nil {
		util.Logger.Errorf("call getClient error, %v", err)
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

func (mc *MClient) getClient(mode SelectMode) (*Client, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	// remove terminated server
	for key, client := range mc.clients {
		if !client.IsAvailable() {
			delete(mc.clients, key)
			mc.addrs = util.DeleteStringElement(mc.addrs, key)
		}
	}

	n := len(mc.addrs)
	if n == 0 {
		return nil, errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case RandomSelect:
		return mc.clients[mc.addrs[mc.r.Intn(n)]], nil
	default:
		return nil, errors.New("rpc: not supported select mode")
	}
}
