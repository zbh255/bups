package app

import (
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// CFG 解决配置文件缓冲读的问题
type CFG struct {
	mu sync.Mutex
	// 是否以读取
	isRead bool
	// 上层缓存区
	topBuf []byte
	// 底层缓冲区
	buf []byte
	_fd *os.File
}

func NewCFGBuffer(fd *os.File) io.ReadWriteCloser {
	return &CFG{
		_fd: fd,
	}
}

func (c *CFG) Open(fd *os.File) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._fd = fd
}

func (c *CFG) Update(fd *os.File) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._fd = fd
	c.isRead = false
	c.buf = nil
}

// 兼容io.Reader
func (c *CFG) Read(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isRead {
		bytes, err := ioutil.ReadAll(c._fd)
		if err != nil {
			return 0, err
		}
		defer func() {
			c.isRead = true
		}()
		c.buf = bytes
	}
	if c.topBuf == nil || len(c.topBuf) == 0 {
		c.topBuf = append(c.buf)
	}
	lens := 0
	if len(p) > len(c.topBuf) {
		lens = len(c.topBuf)
	} else {
		lens = len(p)
	}
	for i := 0; i < lens; i++ {
		p[i] = c.topBuf[i]
	}
	c.topBuf = c.topBuf[lens:]
	if len(c.topBuf) == 0 {
		return lens, io.EOF
	}
	return lens, nil
}

func (c *CFG) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._fd.Write(p)
}

func (c *CFG) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._fd.Close()
}
