package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
	"unsafe"

	"github.com/benjojo/fgbgp/messages"
	fgbgp "github.com/benjojo/fgbgp/server"
	"golang.org/x/sys/unix"
)

// Collector a
type Collector struct {
}

// Notification a
func (col *Collector) Notification(msg *messages.BGPMessageNotification, n *fgbgp.Neighbor) bool {
	log.Print("Notification")
	return true
}

var (
	RouteSendCount        = flag.Int("flood.routes", 99, "How many routes to send on the flood")
	TrickleRouteSendCount = flag.Int("trickle.routes", 1, "How many routes to send per -trickle.wait")
	TrickleWait           = flag.Duration("trickle.wait", time.Second*10, "How long to wait to send -trickle.routes many routes")
)

// ProcessReceived a
func (col *Collector) ProcessReceived(v interface{}, n *fgbgp.Neighbor) (bool, error) {
	log.Print("ProcessReceived")

	switch v.(type) {
	case *messages.BGPMessageKeepAlive:
		if n.HasFlooded {
			break
		}
		n.HasFlooded = true

		log.Print("Sending update flood")

		for i := uint64(0); i < uint64(*RouteSendCount); i++ {

			afi := messages.AfiSafi{
				Afi:  1,
				Safi: 1,
			}

			blahByte := make([]byte, 8)
			binary.PutUvarint(blahByte, i)

			rpfx := make([]byte, 3)
			rand.Read(rpfx)
			_, a, _ := net.ParseCIDR(fmt.Sprintf("%d.%d.%d.0/24", blahByte[2]+1, blahByte[1], blahByte[0]))
			pfx := messages.NLRI_IPPrefix{
				Prefix: *a,
				PathId: 1,
			}
			n.SendRoute(
				afi,
				[]messages.NLRI{pfx},
				nil,
				net.IP{192, 168, 2, 50},
				[]uint32{},
				[]uint32{65001},
				1,
				1)
		}

		// Set the window to 0
		sockFile, sockErr := n.Tcpconn.File()
		if sockErr == nil {

			// got socket file handle. Getting descriptor.
			fd := int(sockFile.Fd())

			err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.SO_RCVBUF, 257)
			if err != nil {
				log.Fatal("on SO_RCVBUFFORCE -- ", err.Error())
			}

			// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.TCP_REPAIR, 1)
			// if err != nil {
			// 	log.Fatal("on switching on repair -- ", err.Error())
			// }

			// // _, _, err := unix.Syscall(unix.SYS_GETSOCKOPT, uintptr(fd), unix.TCP_REPAIR, uintptr(unsafe.Pointer(&twin)))
			// // twin := tcpRepairWindowStruct{}

			// // Window Set
			// /*
			// 	//// SetsockoptByte(fd, syscall.IPPROTO_TCP, unix.TCP_REPAIR_WINDOW, 3)
			// 	func SetsockoptByte(fd, level, opt int, value byte) (err error) {
			// 		return setsockopt(fd, level, opt, unsafe.Pointer(&value), 1)
			// 	}
			// 	func SetsockoptByte(fd, level, opt int, value byte) (err error) {
			// 		return setsockopt(fd, level, opt, unsafe.Pointer(&value), 1)
			// 	}
			// */
			// // _, _, err2 := unix.Syscall6(unix.SYS_GETSOCKOPT,
			// // 	uintptr(fd),
			// // 	uintptr(syscall.IPPROTO_TCP),
			// // 	uintptr(unix.TCP_REPAIR_WINDOW),
			// // 	uintptr(unsafe.Pointer(&twin)),
			// // 	uintptr(20),
			// // 	uintptr(0))

			// window, err := GetsockoptTcpRepairWindow(fd, unix.IPPROTO_TCP, unix.TCP_REPAIR_WINDOW)
			// log.Printf("aaa %#v", window)

			// window.RcvWnd = 1
			// err = SetsockoptTcpRepairWindow(fd, unix.IPPROTO_TCP, unix.TCP_REPAIR_WINDOW, window)
			// // err2 := getsockopt(fd, syscall.IPPROTO_TCP, unix.TCP_REPAIR_WINDOW, unsafe.Pointer(&twin), 20)
			// // // err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.TCP_REPAIR_WINDOW, 3)
			// // log.Printf("DeBUG : %#v", twin)
			// if err != nil {
			// 	log.Fatal("on set -- ", err.Error())
			// }

			// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.TCP_REPAIR, 0)
			// if err != nil {
			// 	log.Fatal("on switching off repair -- ", err.Error())
			// }
			// // don't forget to close the file. No worries, it will *not* cause the connection to close.
			sockFile.Close()
		} else {
			log.Fatal("on setting socket keepalive", sockErr.Error())
		}
		// unix.SetsockoptByte(0, 0, 0, 0)
		ka := messages.BGPMessageKeepAlive{}
		n.OutQueue <- ka
		n.StopRecv = true
		globalsend <- v
		go trickleRoutes(n)
	default:
		break
	}
	return true, nil
}

func trickleRoutes(n *fgbgp.Neighbor) {
	i := uint64(*RouteSendCount) + 1
	for {
		time.Sleep(*TrickleWait)
		if n.Connected {
			for j := 0; j < *TrickleRouteSendCount; j++ {
				afi := messages.AfiSafi{
					Afi:  1,
					Safi: 1,
				}

				blahByte := make([]byte, 8)
				binary.PutUvarint(blahByte, i)

				rpfx := make([]byte, 3)
				rand.Read(rpfx)
				_, a, _ := net.ParseCIDR(fmt.Sprintf("%d.%d.%d.0/24", blahByte[2]+1, blahByte[1], blahByte[0]))
				pfx := messages.NLRI_IPPrefix{
					Prefix: *a,
					PathId: 1,
				}
				n.SendRoute(
					afi,
					[]messages.NLRI{pfx},
					nil,
					net.IP{192, 168, 2, 50},
					[]uint32{},
					[]uint32{65001},
					1,
					1)

				i++
			}
		} else {
			return
		}
	}
}

// func setReadBuffer(fd *netFD, bytes int) error {
// 	err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
// 	runtime.KeepAlive(fd)
// 	return wrapSyscallError("setsockopt", err)
// }

func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err syscall.Errno) {
	_, _, e1 := unix.Syscall6(unix.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
	return e1
}

func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err syscall.Errno) {
	_, _, e1 := unix.Syscall6(unix.SYS_SETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(vallen), 0)
	return e1
}

type tcpRepairWindowStruct struct {
	/*
		struct tcp_repair_window {
			__u32	snd_wl1;
			__u32	snd_wnd;
			__u32	max_window;

			__u32	rcv_wnd;
			__u32	rcv_wup;
		};
	*/
	sndwl1    uint32
	sndwnd    uint32
	maxwindow uint32
	rcvwnd    uint32
	rcvwup    uint32
}

func setWindow() {

}

// ProcessSend a
func (col *Collector) ProcessSend(v interface{}, n *fgbgp.Neighbor) (bool, error) {
	log.Print("ProcessSend")

	globalsend <- v
	return true, nil
}

// ProcessUpdateEvent a
func (col *Collector) ProcessUpdateEvent(e *messages.BGPMessageUpdate, n *fgbgp.Neighbor) (add bool) {
	log.Print("ProcessUpdateEvent")

	globalsend <- *e
	return false
}

// NewNeighbor a
func (col *Collector) NewNeighbor(msg *messages.BGPMessageOpen, n *fgbgp.Neighbor) bool {
	log.Print("NewNeighbor")

	log.Printf("New connection: %#v", n)

	sockFile, sockErr := n.Tcpconn.File()
	if sockErr == nil {
		// got socket file handle. Getting descriptor.
		fd := int(sockFile.Fd())

		err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.SO_RCVBUF, 257)
		if err != nil {
			log.Fatal("on SO_RCVBUFFORCE -- ", err.Error())
		}
	}

	return true // ANYONE!
}

// DisconnectedNeighbor a
func (col *Collector) DisconnectedNeighbor(n *fgbgp.Neighbor) {
	log.Printf("Disconnected")
	return
}

// OpenSend a
func (col *Collector) OpenSend(msg *messages.BGPMessageOpen, n *fgbgp.Neighbor) bool {
	log.Print("OpenSend")

	return true
}

var globalsend chan interface{}

func main() {
	BgpAddr := flag.String("address", "0.0.0.0:1179", "where to listen on")
	flag.Parse()

	globalsend = make(chan interface{}, 1)

	m := fgbgp.NewManager(65001, net.ParseIP("1.3.3.7"), false, false)
	m.UseDefaultUpdateHandler(10)
	col := &Collector{}
	m.SetEventHandler(col)
	m.SetUpdateEventHandler(col)
	err := m.NewServer(*BgpAddr)
	if err != nil {
		log.Fatal(err)
	}

	go m.Start()

	for msg := range globalsend {
		log.Printf("aaa -%T- %v", msg, msg)
	}

}
