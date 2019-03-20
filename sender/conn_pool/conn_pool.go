package conn_pool

import (
	"fmt"
	"io"
	"sync"
	"time"
)

//TODO: 保存所有的连接, 而不是只保存连接计数

var ErrMaxConn = fmt.Errorf("maximum connections reached")

// 连接接口
type NConn interface {
	io.Closer
	Name() string
	Closed() bool
}

type ConnPool struct {
	sync.RWMutex

	Name     string
	Address  string
	MaxConns int
	MaxIdle  int
	Cnt      int64 // 创建连接计数器
	New      func(name string) (NConn, error)

	active int              // 活动连接数
	free   []NConn          // 保存空闲连接
	all    map[string]NConn // 保存所有连接
}

// 创建ConnPool实例
func NewConnPool(name string, address string, maxConns int, maxIdle int) *ConnPool {
	return &ConnPool{Name: name, Address: address, MaxConns: maxConns, MaxIdle: maxIdle, Cnt: 0, all: make(map[string]NConn)}
}

// 输出自身信息
func (this *ConnPool) Proc() string {
	this.RLock()
	defer this.RUnlock()

	return fmt.Sprintf("Name:%s,Cnt:%d,active:%d,all:%d,free:%d",
		this.Name, this.Cnt, this.active, len(this.all), len(this.free))
}

// 获取一个连接
func (this *ConnPool) Fetch() (NConn, error) {
	this.Lock()
	defer this.Unlock()

	// get from free
	conn := this.fetchFree()
	if conn != nil {
		return conn, nil
	}

	if this.overMax() {
		return nil, ErrMaxConn
	}

	// create new conn
	conn, err := this.newConn()
	if err != nil {
		return nil, err
	}

	this.increActive()
	return conn, nil
}

// 归还连接,归还到free或删除连接
func (this *ConnPool) Release(conn NConn) {
	this.Lock()
	defer this.Unlock()

	if this.overMaxIdle() {
		this.deleteConn(conn)
		this.decreActive()
	} else {
		this.addFree(conn)
	}
}

// 强制关闭active的连接
func (this *ConnPool) ForceClose(conn NConn) {
	this.Lock()
	defer this.Unlock()

	this.deleteConn(conn)
	this.decreActive()
}

// 释放所有连接
func (this *ConnPool) Destroy() {
	this.Lock()
	defer this.Unlock()

	for _, conn := range this.free {
		if conn != nil && !conn.Closed() {
			conn.Close()
		}
	}

	for _, conn := range this.all {
		if conn != nil && !conn.Closed() {
			conn.Close()
		}
	}

	this.active = 0
	this.free = []NConn{}
	this.all = map[string]NConn{}
}

// 创建新rpcConn，名称：address_创建计数_创建时间
// internal, concurrently unsafe
func (this *ConnPool) newConn() (NConn, error) {
	name := fmt.Sprintf("%s_%d_%d", this.Name, this.Cnt, time.Now().Unix())
	conn, err := this.New(name)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return nil, err
	}

	this.Cnt++
	this.all[conn.Name()] = conn
	return conn, nil
}

// 删除连接
func (this *ConnPool) deleteConn(conn NConn) {
	if conn != nil {
		conn.Close()
	}
	delete(this.all, conn.Name())
}

// 增加空闲连接
func (this *ConnPool) addFree(conn NConn) {
	this.free = append(this.free, conn)
}

// 获取空闲连接
func (this *ConnPool) fetchFree() NConn {
	if len(this.free) == 0 {
		return nil
	}

	conn := this.free[0]
	this.free = this.free[1:]
	return conn
}

// 增加active
func (this *ConnPool) increActive() {
	this.active += 1
}

// 减少active
func (this *ConnPool) decreActive() {
	this.active -= 1
}

// 是否超过最大活动连接数
func (this *ConnPool) overMax() bool {
	return this.active >= this.MaxConns
}

// 是否超过最大空闲数
func (this *ConnPool) overMaxIdle() bool {
	return len(this.free) >= this.MaxIdle
}
