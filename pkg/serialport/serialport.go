package serialport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/antlabs/timer"
	"github.com/goburrow/serial"
)

type Client interface {
	Start() error
	Opened() bool
	Send(data []byte) error
	AddPollCmd(pollCmd []byte)
	SetReceiveCallback(fn func(data []byte))
	StopPolling()
	Close() error
}

type client struct {
	ctx         context.Context
	cancel      context.CancelFunc
	conn        serial.Port
	options     Options
	readBuf     []byte
	receiveFunc func(data []byte)
	writeChl    chan []byte
	sending     bool
	timer       timer.Timer
	pollCmds    [][]byte
}

func NewClient(o *Options) Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &client{
		ctx:      ctx,
		cancel:   cancel,
		options:  *o,
		writeChl: make(chan []byte),
		readBuf:  make([]byte, 15),
		pollCmds: make([][]byte, 0),
	}
}

func (c *client) Start() error {
	port, err := serial.Open(&serial.Config{
		Address:  c.options.address,
		BaudRate: c.options.baudRate,
		DataBits: c.options.dataBits,
		StopBits: c.options.stopBits,
		Parity:   c.options.parity,
		Timeout:  c.options.timeout, // 读写超时
	})
	if err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}
	c.conn = port

	go c.reader()
	go c.write()

	if c.options.poll {
		go c.startPolling()
	}
	return nil
}

func (c *client) AddPollCmd(pollCmd []byte) {
	c.pollCmds = append(c.pollCmds, pollCmd)
}

func (c *client) SetReceiveCallback(fn func(data []byte)) {
	c.receiveFunc = fn
}

func (c *client) Send(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("serial port closed")
	}

	// 设置发送标识
	c.sending = true
	c.writeChl <- data
	return nil
}

func (c *client) startPolling() {
	c.timer = timer.NewTimer()
	c.timer.ScheduleFunc(c.options.pollInterval, c.polling)
	c.timer.Run()
}

func (c *client) polling() {
	// 本次循环完成，恢复默认轮询间隔
	startTs := time.Now()

	for _, poll := range c.pollCmds {
		select {
		case <-c.ctx.Done():
			return
		default:
			if c.sending {
				continue
			}
			c.writeChl <- poll
		}
	}
	// 等待最后一包完成
	c.addCmdInterval()

	endTs := time.Now()
	fmt.Printf("[%s] Polled Time: %dms\n", c.options.address, endTs.Sub(startTs).Milliseconds())
}

func (c *client) write() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.writeChl:
			err := c.writeWithRetries(data, c.options.retries)
			if err != nil {
				fmt.Printf("Failed to send data: %v\n", err)
			}
		}
	}
}

func (c *client) writeWithRetries(data []byte, retries int) error {
	for i := 0; i < retries; i++ {
		select {
		case <-c.ctx.Done():
			return errors.New("context canceled")
		default:
			_, err := c.conn.Write(data)
			if err != nil {
				// 连接错误直接返回
				_, ok := err.(*net.OpError)
				if ok || err.Error() == "no such file or directory" || err.Error() == "device not configured" {
					c.conn.Close()
					return err
				}

				// 业务数据错误则进行重试
				continue
			}

			fmt.Printf("Send: % 0X\n", data)
			return nil
		}
	}
	return fmt.Errorf("failed to send data after %d retries", retries)
}

func (c *client) reader() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			n, err := c.conn.Read(c.readBuf)
			if err != nil {
				// 连接错误直接返回
				_, ok := err.(*net.OpError)
				if ok || err.Error() == "no such file or directory" || err.Error() == "device not configured" {
					c.conn.Close()
					return
				}

				// 业务数据错误则进行重试
				continue
			}

			if c.receiveFunc != nil {
				c.receiveFunc(c.readBuf[:n])
				c.sending = false
			}
		}
	}
}

func (c *client) Opened() bool {
	return c.conn != nil
}

func (c *client) StopPolling() {
	// c.pollTimeNoder.Stop()
	c.timer.Stop()
}

func (c *client) Close() error {
	if c.timer != nil {
		c.StopPolling()
	}

	c.cancel()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *client) addCmdInterval() {
	time.Sleep(c.options.timeout)
}
