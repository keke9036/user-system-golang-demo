// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/15

package arpc

import (
	"encoding/json"
	"entry-task/arpc/codec"
	"entry-task/util"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0xaabb

type Option struct {
	MagicNumber int
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
}

type Server struct {
	serviceMap sync.Map
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value // argv and replyv of request
	mtype        *methodType
	svc          *service
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener, timeout time.Duration) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			util.Logger.Println("arpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn, timeout)
	}
}

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("arpc: service already defined: " + s.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("arpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("arpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("arpc server: can't find method " + methodName)
	}
	return
}

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
func (server *Server) ServeConn(conn io.ReadWriteCloser, timeout time.Duration) {

	// server recovery
	defer func() {
		if e := recover(); e != nil {
			recoveryLog := "Rpc server recovery log:  message: %v, stack: %s"
			util.Logger.Errorf(recoveryLog, e, string(debug.Stack()[:]))
		}
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		util.Logger.Println("arpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		util.Logger.Printf("arpc server: invalid magic number %x", opt.MagicNumber)
		return
	}

	gobF := codec.NewCodecFuncMap[codec.GobType]
	server.serveCodec(gobF(conn), timeout)
}

var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc codec.Codec, timeout time.Duration) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg, timeout)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			util.Logger.Println("arpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		util.Logger.Println("arpc server: read body err:", err)
		return req, err
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		util.Logger.Println("arpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()

	called := make(chan struct{})
	sent := make(chan struct{})

	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-sent
		return
	}
	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func Accept(lis net.Listener, timeout time.Duration) { DefaultServer.Accept(lis, timeout) }
