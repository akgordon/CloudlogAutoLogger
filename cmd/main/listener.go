package main

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"CloudlogAutoLogger/internal/agg_udp"
	"strconv"
	"time"
)

type Listeners struct {
	Cloudlog_api_key   string
	Station_profile_id string
	Port               int

	client_name string // Name of broadcasting app (e.g. WSJTX, JS8CALL, VARAC)

	// thread control
	verbose    bool
	endFlag    bool
	threadFlag bool
}

func (cd *Listeners) Start() {
	var portstg = strconv.Itoa(cd.Port)
	agg_logger.Get().Log("***** Begin listener on port:", portstg)
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

func (cd *Listeners) Stop() {
	cd.endFlag = true
	sleepTime := 2 * time.Second

	for {
		if !cd.threadFlag {
			return
		}

		time.Sleep(sleepTime)
	}
}
