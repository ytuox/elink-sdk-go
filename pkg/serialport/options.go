package serialport

import (
	"time"
)

type Options struct {
	// Device path (/dev/ttyS0)
	address string
	// Baud rate (default 19200)
	baudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	dataBits int
	// Stop bits: 1 or 2 (default 1)
	stopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	parity string
	// Read (Write) timeout.
	timeout time.Duration
	// Configuration related to RS485
	// Enable RS485 support
	enabled bool
	// Delay RTS prior to send
	delayRtsBeforeSend time.Duration
	// Delay RTS after send
	delayRtsAfterSend time.Duration
	// Set RTS high during send
	rtsHighDuringSend bool
	// Set RTS high after send
	rtsHighAfterSend bool
	// Rx during Tx
	rxDuringTx   bool
	retries      int
	poll         bool
	pollInterval time.Duration
}

func NewOptions() *Options {
	o := &Options{
		timeout:      130 * time.Millisecond,
		retries:      3,
		pollInterval: 5000 * time.Millisecond,
	}
	time, err := calculateTimeout(o.baudRate, o.dataBits, o.stopBits, o.parity, 15)
	if err == nil {
		o.timeout = time
	}
	return o
}

func (o *Options) SetAddress(address string) {
	o.address = address
}

func (o *Options) SetBaudRate(baudRate int) {
	o.baudRate = baudRate
}

func (o *Options) SetDataBits(dataBits int) {
	o.dataBits = dataBits
}

func (o *Options) SetStopBits(stopBits int) {
	o.stopBits = stopBits
}

func (o *Options) SetParity(parity string) {
	o.parity = parity
}

func (o *Options) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

func (o *Options) SetRetries(retries int) {
	o.retries = retries
}

func (o *Options) SetPoll(poll bool) {
	o.poll = poll
}

func (o *Options) SetPollInterval(interval time.Duration) {
	o.pollInterval = interval
}
