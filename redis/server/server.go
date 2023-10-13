package server

/*
	redis 服务器 1.0
*/

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"redis-based-on-go/config"
	"redis-based-on-go/database"
	"redis-based-on-go/log"
	"redis-based-on-go/redis/parser"
	"redis-based-on-go/redis/protocol"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// StartServer 启动服务器
func StartServer(cfg *config.ServerConfig) error {
	return bindAndServeWithSignal(cfg, newServer())
}

// server redis 服务器
type server struct {
	activeConn    sync.Map     // 保存所有可用连接
	clientCounter atomic.Int32 // 所有可用连接数量
	isClosed      atomic.Bool  // 表示服务器是否关闭 关闭时拒接新连接
	db            *database.DB // 内存数据库
}

func newServer() *server {
	return &server{db: database.NewDB()}
}

// closeClient 关闭指定连接
func (s *server) closeClient(conn net.Conn) {
	_ = conn.Close()
	s.activeConn.Delete(conn)
}

// handle 处理指定连接上的请求
func (s *server) handle(ctx context.Context, conn net.Conn) {
	if s.isClosed.Load() {
		_ = conn.Close()
		return
	}
	s.activeConn.Store(conn, struct{}{})
	// 解析请求并存入 channel 后续取出处理
	ch := parser.Parse(conn)
	for payLoad := range ch {
		if payLoad.Err != nil {
			// 客户端连接已关闭
			if payLoad.Err == io.EOF ||
				payLoad.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payLoad.Err.Error(), "use of closed network connection") {
				s.closeClient(conn)
				mylog.Logger.Info("关闭连接:" + conn.RemoteAddr().String())
				return
			}
			// 协议解析错误 回复给客户端
			errReply := &protocol.Error{Data: payLoad.Err.Error()}
			_, err := conn.Write(errReply.ToBytes())
			if err != nil {
				s.closeClient(conn)
				mylog.Logger.Info("关闭连接:" + conn.RemoteAddr().String())
				return
			}
			continue
		}
		if payLoad.Data == nil {
			mylog.Logger.Info("消息内容为空")
			continue
		}
		mylog.Logger.Info(fmt.Sprintf("接收到数据: %v", payLoad.Data.ToBytes()))
		// TODO 执行命令
		array := payLoad.Data.(*protocol.Array)
		result := s.db.Execute(array.ToArgs())
		conn.Write(result.ToBytes())
	}
	mylog.Logger.Info("接收完成")
}

// close 关闭服务器并释放客户端连接
func (s *server) close() error {
	s.isClosed.Store(true)
	s.activeConn.Range(func(key, value interface{}) bool {
		clientConn := key.(net.Conn)
		_ = clientConn.Close()
		return true
	})
	return nil
}

// bindAndServeWithSignal 绑定端口并等待处理请求 会阻塞至收到指定信号
func bindAndServeWithSignal(cfg *config.ServerConfig, server *server) error {
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal)
	// 指定程序捕获和处理来自操作系统的信号 用于关闭退出
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// 等待退出信号
	go func() {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	address := fmt.Sprintf("%s:%v", cfg.Bind, cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		mylog.Logger.Error(fmt.Sprintf("监听端口错误:%s", err.Error()))
		return err
	}
	mylog.Logger.Info(fmt.Sprintf("绑定成功: %s, 开始监听...", address))
	listenAndServe(listener, server, closeChan)
	return nil
}

// listenAndServe 等待并处理请求
func listenAndServe(listener net.Listener, server *server, closeChan <-chan struct{}) {
	errChan := make(chan error, 1)
	defer close(errChan)
	// 控制程序退出
	go func() {
		select {
		case <-closeChan:
			mylog.Logger.Info("捕获到退出信号")
		case er := <-errChan:
			mylog.Logger.Info(fmt.Sprintf("监听出现错误: %s", er.Error()))
		}
		mylog.Logger.Info("关闭服务器...")
		_ = listener.Close()
		_ = server.close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			// 超时错误 listener 未关闭 重试
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				mylog.Logger.Infof("连接超时: %v, 5秒后重试", err)
				time.Sleep(5 * time.Millisecond)
				continue
			}
			// listener.Close() 后 listener.Accept() 返回错误
			errChan <- err
			break
		}
		// 连接成功 处理该连接上的请求
		mylog.Logger.Info(fmt.Sprintf("%s 连接成功", conn.RemoteAddr().String()))
		server.clientCounter.Add(1)
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
				server.clientCounter.Add(-1)
			}()
			mylog.Logger.Info("处理请求")
			server.handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
