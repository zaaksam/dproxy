package services

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/zaaksam/dproxy/go/logger"
	"github.com/zaaksam/dproxy/go/model"
)

// Proxy 代理服务对象
var Proxy proxyService

func init() {
	Proxy.proxys = make(map[int64]*proxy)
}

type proxyService struct {
	proxys map[int64]*proxy
	mx     sync.Mutex
}

func (s *proxyService) addProxy(id int64, p *proxy) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.proxys[id] = p
}

func (s *proxyService) delProxy(id int64) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.proxys, id)
}

func (s *proxyService) StartAll() error {
	list, err := PortMap.Find(1, 500, "", "", "")
	if err != nil {
		return err
	}

	if list.Total <= 0 {
		return nil
	}

	var mds []*model.PortMapModel
	if items, ok := list.Items.(*[]*model.PortMapModel); ok {
		mds = *items
	} else {
		return nil
	}

	for i, l := 0, len(mds); i < l; i++ {
		if !mds[i].IsStart {
			err = s.Start(mds[i].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *proxyService) StopAll() {
	list, err := PortMap.Find(1, 500, "", "", "")
	if err != nil {
		return
	}

	if list.Total <= 0 {
		return
	}

	for partMapID := range Proxy.proxys {
		s.Stop(partMapID)
	}
}

// Start 启动代理请求服务
func (s *proxyService) Start(portMapID int64) error {
	if _, ok := s.proxys[portMapID]; !ok {
		//检查是否存在该配置
		md, err := PortMap.Get(portMapID)
		if err != nil {
			return err
		}

		//开启代理请求协程
		sourceAddr := fmt.Sprintf("%s:%d", md.SourceIP, md.SourcePort)
		targetAddr := fmt.Sprintf("%s:%d", md.TargetIP, md.TargetPort)
		listener, err := net.Listen("tcp", sourceAddr)
		if err != nil {
			return errors.New("代理启动失败：" + err.Error())
		}

		p := &proxy{
			listener: listener,
			stopChan: make(chan bool),
			clients:  make(map[string]chan bool),
		}
		s.addProxy(portMapID, p)

		go func() {
			for {
				//创建 client
				clientConn, err := p.listener.Accept()
				if err != nil {
					if !strings.Contains(err.Error(), "use of closed network connection") {
						logger.Error("客户端请求响应失败：", err)
					}
					break
				}

				//白名单检查
				clientAddr := strings.ToLower(clientConn.RemoteAddr().String())
				ip := strings.Split(clientAddr, ":")[0]
				if !WhiteList.Check(ip) {
					clientConn.Close()

					logger.Warning("非法请求：", clientAddr, "->", sourceAddr, "(", targetAddr, ")")
					continue
				}

				//创建 server
				serverConn, err := net.DialTimeout("tcp", targetAddr, 30*time.Second)
				if err != nil {
					clientConn.Close()

					logger.Error("服务端请求响应失败：", err)
					break
				}

				//注册客户端
				invalidChan := p.addClient(ip)
				//代理停止命令监测
				endChan := make(chan bool)
				go func() {
					select {
					case <-p.stopChan:
						//代理停止命令
					case <-invalidChan:
						//白名单失效通知，可能删除，可能过期
					case <-endChan:
						//单次请求连接结束通知
					}

					//关闭服务端连接，忽略各种状态判断，作为补充措施
					serverConn.Close()
				}()

				// server -> client
				go func() {
					_, err := io.Copy(clientConn, serverConn)
					if err == nil {
						//服务端EOF，关闭客户端连接
						clientConn.Close()
					} else {
						if strings.Contains(err.Error(), "use of closed network connection") {
							clientConn.Close()
						} else {
							logger.Error("服务响应接收失败：", err)
						}
					}
				}()

				// client -> server
				go func() {
					_, err := io.Copy(serverConn, clientConn)
					if err == nil {
						//客户端EOF，关闭服务端连接
						serverConn.Close()
					} else {
						if strings.Contains(err.Error(), "use of closed network connection") {
							serverConn.Close()
						} else {
							logger.Error("代理请求发送失败：", err)
						}
					}

					//本次请求结束
					close(endChan)
				}()
			}
		}()
	}

	return nil
}

func (s *proxyService) Stop(portMapID int64) {
	if p, ok := s.proxys[portMapID]; ok {
		p.listener.Close()
		close(p.stopChan)

		s.delProxy(portMapID)
	}
}

func (s *proxyService) setClientInvalid(ip string) {
	for _, p := range s.proxys {
		p.setClientInvalid(ip)
	}
}

// ================ proxy 定义 ========================

type proxy struct {
	listener net.Listener
	stopChan chan bool
	clients  map[string]chan bool // map[ip]invalidChan
	mx       sync.Mutex

	// clients  map[string]map[string]net.Conn // map[ip]map[port]net.Conn
	// ipMutex   sync.RWMutex
	// portMutex sync.RWMutex
}

func (p *proxy) addClient(ip string) <-chan bool {
	p.mx.Lock()
	defer p.mx.Unlock()

	invalidChan, ok := p.clients[ip]
	if !ok {
		invalidChan = make(chan bool)
		p.clients[ip] = invalidChan
	}

	return invalidChan
}

func (p *proxy) removeClient(ip string) {
	p.mx.Lock()
	defer p.mx.Unlock()

	delete(p.clients, ip)
}

func (p *proxy) setClientInvalid(ip string) {
	if invalidChan, ok := p.clients[ip]; ok {
		close(invalidChan)

		p.removeClient(ip)
	}
}
