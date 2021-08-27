package services

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/constant"
	"github.com/zaaksam/dproxy/go/logger"
	"github.com/zaaksam/dproxy/go/model"
)

// Proxy 代理服务对象
var Proxy proxyService

func init() {
	Proxy.proxys = make(map[int64]*proxy)

	// 重置所有状态为：已停止
	err := PortMap.resetState(constant.PORTMAP_STATE_STOP)
	if err != nil {
		logger.Warning(err)
	}

	// 监控待停止、待运行的数据，并进行停止、运行
	go Proxy.Watch()
}

type proxyService struct {
	proxys map[int64]*proxy
	mx     sync.Mutex
}

func (s *proxyService) addProxy(id int64, p *proxy) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.proxys[id] = p

	// 更新状态
	err := PortMap.updateState(constant.PORTMAP_STATE_START, id)
	if err != nil {
		idStr := strconv.FormatInt(id, 10)
		logger.Warning("对 " + idStr + " 进行 " + err.Error())
	}
}

func (s *proxyService) delProxy(id int64) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.proxys, id)

	// 更新状态
	err := PortMap.updateState(constant.PORTMAP_STATE_STOP, id)
	if err != nil {
		idStr := strconv.FormatInt(id, 10)
		logger.Warning("对 " + idStr + " 进行 " + err.Error())
	}
}

// StartAll 启动所有端口映射
func (s *proxyService) StartAll() error {
	// 将区域内有效数据更新为：待运行
	err := PortMap.resetState(constant.PORTMAP_STATE_START_WAIT)
	if err != nil {
		return err
	}

	// list, err := PortMap.Find(1, 500, config.AppConf.Region, "", "", "", "", "", []constant.PortmapStateType{constant.PORTMAP_STATE_STOP})
	// if err != nil {
	// 	return err
	// }

	// if list.Total <= 0 {
	// 	return nil
	// }

	// var mds []*model.PortMapModel
	// if items, ok := list.Items.(*[]*model.PortMapModel); ok {
	// 	mds = *items
	// } else {
	// 	return nil
	// }

	// for i, l := 0, len(mds); i < l; i++ {
	// 	if !mds[i].IsStart {
	// 		err = s.start(mds[i])
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

// StopAll 停止所有运行
func (s *proxyService) StopAll() {
	// list, err := PortMap.Find(1, 1, config.AppConf.Region, "", "", "", "", "", []constant.PortmapStateType{constant.PORTMAP_STATE_START, constant.PORTMAP_STATE_START_WAIT})
	// if err != nil {
	// 	return
	// }

	// if list.Total <= 0 {
	// 	return
	// }

	for partMapID := range Proxy.proxys {
		s.stop(partMapID)
	}
}

// Start 启动代理请求服务
func (s *proxyService) Start(portMapID int64) error {
	// 获取端口映射信息
	md, err := PortMap.Get(portMapID)
	if err != nil {
		return err
	}

	if md.Region == config.AppConf.Region {
		// 检查是否已启动监听
		if _, ok := s.proxys[portMapID]; ok {
			return nil
		}
	}

	// 检查区域实例是否正常
	isTimeout, err := Region.CheckTimeout(md.Region)
	if err != nil {
		return err
	}
	if isTimeout {
		err = errors.New(md.Region + " 区域不存在运行实例")
		return err
	}

	// 更新状态为：待运行
	err = PortMap.updateState(constant.PORTMAP_STATE_START_WAIT, portMapID)
	if err != nil {
		return err
	}

	// 4 秒后退出
	i := 0
	for {
		if i > 2 {
			err = errors.New("启动超时，详情请查看日志")
			break
		}

		if md.Region == config.AppConf.Region {
			// 检查是否已启动监听
			if _, ok := s.proxys[portMapID]; ok {
				break
			}
		} else {
			md, err = PortMap.Get(portMapID)
			if err != nil {
				return err
			}
			if md.State == constant.PORTMAP_STATE_START {
				break
			}
		}

		i++
		time.Sleep(2 * time.Second)
	}

	return err

	//检查是否存在该配置
	// md, err := PortMap.Get(portMapID)
	// if err != nil {
	// 	return err
	// }

	// // 更新 portmap 状态为 待启动
	// return s.start(md)
	// return nil
}

// start 启动代理请求服务
func (s *proxyService) start(md *model.PortMapModel) error {

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
	s.addProxy(md.ID, p)

	go func() {
		for {
			//监听请求
			clientConn, err := p.listener.Accept()
			if err != nil {
				logger.Error("客户端请求响应失败：", err)
				break
			}

			go s.Serve(p, clientConn, sourceAddr, targetAddr)
		}
	}()

	return nil
}

func (s *proxyService) Serve(p *proxy, clientConn net.Conn, sourceAddr, targetAddr string) {
	defer clientConn.Close()

	//白名单检查
	clientAddr := strings.ToLower(clientConn.RemoteAddr().String())
	ip := strings.Split(clientAddr, ":")[0]
	if !WhiteList.Check(ip) {
		logger.Warning("非法请求：", clientAddr, "->", sourceAddr, "(", targetAddr, ")")
		return
	}

	//创建 server
	serverConn, err := net.DialTimeout("tcp", targetAddr, 30*time.Second)
	if err != nil {
		logger.Error("服务端请求响应失败：", err)
		return
	}
	defer serverConn.Close()

	//注册客户端
	invalidChan := p.addClient(ip)

	//代理停止命令监测
	endChan := make(chan bool)
	go func() {
		select {
		case <-p.stopChan:
			//代理停止命令
			serverConn.Close()
		case <-invalidChan:
			//白名单失效通知，可能删除，可能过期
			serverConn.Close()
		case <-endChan:
			//单次请求连接结束通知
		}
	}()
	defer close(endChan)

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader, info string) {
		_, err := io.Copy(dst, src)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				err = nil
			} else {
				err = errors.New(info + "失败：" + err.Error())
			}
		}
		errc <- err
	}

	go cp(serverConn, clientConn, "代理请求发送")
	go cp(clientConn, serverConn, "服务响应接收")
	err = <-errc
	if err != nil {
		logger.Error(err)
	}
}

// Stop 停止
func (s *proxyService) Stop(portMapID int64) error {
	// 获取端口映射信息
	md, err := PortMap.Get(portMapID)
	if err != nil {
		return err
	}

	if md.Region == config.AppConf.Region {
		// 检查是否已启动监听
		if _, ok := s.proxys[portMapID]; !ok {
			return nil
		}
	}

	// 检查区域实例是否正常
	isTimeout, err := Region.CheckTimeout(md.Region)
	if err != nil {
		return err
	}
	if isTimeout {
		err = errors.New(md.Region + " 区域不存在运行实例")
		return err
	}

	// 更新状态为：待停止
	err = PortMap.updateState(constant.PORTMAP_STATE_STOP_WAIT, portMapID)
	if err != nil {
		return err
	}

	// 4 秒后退出
	i := 0
	for {
		if i > 2 {
			err = errors.New("停止超时，详情请查看日志")
			break
		}

		if md.Region == config.AppConf.Region {
			// 检查是否已停止监听
			if _, ok := s.proxys[portMapID]; !ok {
				break
			}
		} else {
			md, err = PortMap.Get(portMapID)
			if err != nil {
				return err
			}

			if md.State == constant.PORTMAP_STATE_STOP {
				break
			}
		}

		i++
		time.Sleep(2 * time.Second)
	}

	return err

	// 更新 portmap 状态为 待停止
	// s.stop(portMapID)
}

func (s *proxyService) stop(portMapID int64) {
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

// Watch 监控需要运行的代理，进行启动、停止，2 秒循环
func (s *proxyService) Watch() {
	for {
		// list, err := PortMap.Find(1, 500, config.AppConf.Region, "", "", "", "", "", []constant.PortmapStateType{constant.PORTMAP_STATE_START_WAIT, constant.PORTMAP_STATE_STOP_WAIT})
		list, err := PortMap.findWait(1, 500)
		if err != nil {
			logger.Error("获取端口映射列表失败：", err)
		} else {
			if items, ok := list.Items.(*[]*model.PortMapModel); ok {
				items2 := *items

				for i, l := 0, len(items2); i < l; i++ {
					if items2[i].State == constant.PORTMAP_STATE_START_WAIT {
						// 启动
						err = s.start(items2[i])
						if err != nil {
							logger.Warning(err)
						}
					} else if items2[i].State == constant.PORTMAP_STATE_STOP_WAIT {
						// 停止
						s.stop(items2[i].ID)
					}
				}
			}
		}

		time.Sleep(2 * time.Second)
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
