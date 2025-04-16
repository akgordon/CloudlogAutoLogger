package agg_udp

// Ref: https://forum.golangbridge.org/t/help-with-udp-broadcast/22036
// On Unix remember to open the UDP firewall.  https://help.ubuntu.com/community/UFW
//    Example: sudo ufw allow 4992/udp

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

func init() {
	agg_logger.Get().Log("Networking initialized", "")
}

func WaitOnUDP(port int, timeoutMilliseconds int, doLog bool) (bool, []byte) {

	hostName := "0.0.0.0"
	portNum := strconv.Itoa(port)
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return false, nil
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)
	ln.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutMilliseconds)))

	if err != nil {
		agg_logger.Get().Log(err.Error(), "")
		return false, nil
	}

	if doLog {
		agg_logger.Get().Log("UDP server up and listening on "+service, "")
	}

	defer ln.Close()

	// Now block until get UDP packet
	buffer := make([]byte, 1024)

	_, addr, err := ln.ReadFromUDP(buffer)

	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return false, nil
		}
		agg_logger.Get().Log(err.Error(), "")
	}

	if doLog {
		agg_logger.Get().Log("Received from UDP client : ", addr.String())
	}
	//agg_logger.Get().Log("Received from UDP client :  ", string(buffer[:n]))
	return true, buffer
}
