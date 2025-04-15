package AGG_udp

// Ref: https://forum.golangbridge.org/t/help-with-udp-broadcast/22036
// On Unix remember to open the UDP firewall.  https://help.ubuntu.com/community/UFW
//    Example: sudo ufw allow 4992/udp

import (
	"CloudlogAutoLogger/internal/AGG_logger"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func init() {
	AGG_logger.Get().Log("Networking initialized", "")
}

func WaitOnUDP(port int, timeoutMilliseconds int, doLog bool) (bool, []byte) {

	hostName := "0.0.0.0"
	portNum := strconv.Itoa(port)
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		AGG_logger.Get().Log(err.Error(), "")
		return false, nil
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)
	ln.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutMilliseconds)))

	if err != nil {
		AGG_logger.Get().Log(err.Error(), "")
		return false, nil
	}

	if doLog {
		AGG_logger.Get().Log("UDP server up and listening on "+service, "")
	}

	defer ln.Close()

	// Now block until get UDP packet
	buffer := make([]byte, 1024)

	_, addr, err := ln.ReadFromUDP(buffer)

	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return false, nil
		}
		AGG_logger.Get().Log(err.Error(), "")
	}

	if doLog {
		AGG_logger.Get().Log("Received from UDP client : ", addr.String())
	}
	//AGG_logger.Get().Log("Received from UDP client :  ", string(buffer[:n]))
	return true, buffer
}

func WaitOnVita49(port int, timeoutMilliseconds int, doLog bool) (bool, []byte, int) {

	// Create a raw socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	if err != nil {
		msg := fmt.Sprintf("Failed to create raw socket: %v", err)
		AGG_logger.Get().Log("ERROR:", msg)
		return false, nil, 0
	}
	defer syscall.Close(fd)

	// Set the socket to promiscuous mode
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		msg := fmt.Sprintf("Failed to set socket options: %v", err)
		AGG_logger.Get().Log("ERROR:", msg)
		return false, nil, 0
	}

	// Create a buffer to hold incoming packets
	buf := make([]byte, 4096)

	for {
		// Read packets from the socket
		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			msg := fmt.Sprintf("Failed to receive packet: %v", err)
			AGG_logger.Get().Log("ERROR:", msg)
			return false, nil, 0
		}

		// Parse the source address
		srcAddr := from.(*syscall.SockaddrInet4)
		srcIP := net.IPv4(srcAddr.Addr[0], srcAddr.Addr[1], srcAddr.Addr[2], srcAddr.Addr[3])
		srcPort := srcAddr.Port

		// Print the packet details
		msg := fmt.Sprintf("Received packet from %s:%d\n", srcIP, srcPort)
		AGG_logger.Get().Log(msg, "")
		//fmt.Printf("Data: %x\n", buf[:n])

		return true, buf, n
	}
}
