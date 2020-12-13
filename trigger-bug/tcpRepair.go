// +build !linux

package main

import "net"

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

func GetsockoptTcpRepairWindow(fd int, level int, opt int) (TCPRepairWindow, error) {
	window := TCPRepairWindow{}
	return window, nil
}

func SetsockoptTcpRepairWindow(fd int, level int, opt int, window TCPRepairWindow) error {
	return nil
}
