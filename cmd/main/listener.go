package main

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"CloudlogAutoLogger/internal/agg_udp"
	"strconv"
	"time"
)

func (cd *listeners) start() {
	var portstg = strconv.Itoa(cd.Port)
	agg_logger.Get().Log("***** Begin WSJT-X listener on port:", portstg)
	for !cd.endFlag {
		stat, udp_pkt := agg_udp.WaitOnUDP(cd.Port, 250, false)
		if stat {
			agg_logger.Get().Log("***** Begin packet ***** Port:", portstg)

			pkt_stg := string(udp_pkt)
			agg_logger.Get().Log(pkt_stg, "")

			//flex_util.BinaryDump(udp_pkt)

			//pkt_stg := string(udp_pkt[:])
			// flds := strings.FieldsFunc(pkt_stg, split)
			// for _, v := range flds {
			// 	agg_logger.Get().Log(portStg, " "+v)
			// }
			agg_logger.Get().Log("***** End packet ***** Port:", portstg)
		}
		//time.Sleep(time.Millisecond)
	}
	cd.threadFlag = false
}

func (cd *listeners) stop() {
	cd.endFlag = true
	sleepTime := 2 * time.Second

	for {
		if !cd.threadFlag {
			return
		}

		time.Sleep(sleepTime)
	}
}
