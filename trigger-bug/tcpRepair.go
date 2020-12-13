package main

// #include <linux/tcp.h>
// #include <sys/types.h>
// #include <sys/socket.h>
// #include <stdint.h>
import "C"

import (
	"fmt"
	"net"
	"unsafe"
)

// via https://github.com/hrntknr/term-lb/blob/4ababc2d6d29a5d5c3de7f6c75ae81078ab428ee/lb/tcpRepair.go

type TCPRepair struct {
	Saddr  net.IP          `json:"saddr"`
	Sport  uint16          `json:"sport"`
	Daddr  net.IP          `json:"daddr"`
	Dport  uint16          `json:"dport"`
	Window TCPRepairWindow `json:"window"`
	SndSeq int             `json:"snd_seq"`
	RcvSeq int             `json:"rcv_seq"`
	Mss    int             `json:"mss"`
}

type TCPRepairWindow struct {
	SndWl1    uint32 `json:"send_wl1"`
	SndWnd    uint32 `json:"snd_wnd"`
	MaxWindow uint32 `json:"max_window"`
	RcvWnd    uint32 `json:"rcv_wnd"`
	RcvWup    uint32 `json:"rcv_wup"`
}

var TCP_SEND_QUEUE = C.TCP_SEND_QUEUE
var TCP_RECV_QUEUE = C.TCP_RECV_QUEUE

func GetsockoptTcpRepairWindow(fd int, level int, opt int) (TCPRepairWindow, error) {
	val := C.struct_tcp_repair_window{}
	len := C.uint(C.sizeof_struct_tcp_repair_window)
	result, err := C.getsockopt(C.int(fd), C.int(level), C.int(opt), unsafe.Pointer(&val), &len)
	if result < 0 {
		return TCPRepairWindow{}, fmt.Errorf("getsockopt() failed. %s", err)
	}
	return TCPRepairWindow{
		SndWl1:    uint32(val.snd_wl1),
		SndWnd:    uint32(val.snd_wnd),
		MaxWindow: uint32(val.max_window),
		RcvWnd:    uint32(val.rcv_wnd),
		RcvWup:    uint32(val.rcv_wup),
	}, nil
}

func SetsockoptTcpRepairWindow(fd int, level int, opt int, window TCPRepairWindow) error {
	val := C.struct_tcp_repair_window{
		snd_wl1:    C.uint32_t(window.SndWl1),
		snd_wnd:    C.uint32_t(window.SndWnd),
		max_window: C.uint32_t(window.MaxWindow),
		rcv_wnd:    C.uint32_t(window.RcvWnd),
		rcv_wup:    C.uint32_t(window.RcvWup),
	}
	len := C.uint(C.sizeof_struct_tcp_repair_window)
	result, err := C.setsockopt(C.int(fd), C.int(level), C.int(opt), unsafe.Pointer(&val), len)
	if result < 0 {
		return fmt.Errorf("setsockopt() failed. %s", err)
	}
	return nil
}
