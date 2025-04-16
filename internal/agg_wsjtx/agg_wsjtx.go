package agg_wsjtx

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"CloudlogAutoLogger/internal/agg_udp"
	"strconv"
	"time"
)

type Wsjtx struct {
	Cloudlog_api_key   string
	Station_profile_id string
	Port               int

	// thread control
	verbose    bool
	endFlag    bool
	threadFlag bool
}

// Factory
var theWSJTX *Wsjtx = nil

func Get() *Wsjtx {
	if theWSJTX == nil {
		theWSJTX = &Wsjtx{}
		theWSJTX.init()
	}
	return theWSJTX
}

func (cd *Wsjtx) init() {
	cd.Cloudlog_api_key = ""
	cd.Station_profile_id = ""
	cd.Port = 0
	cd.verbose = false
	cd.endFlag = false
	cd.threadFlag = false
}

func (cd *Wsjtx) start() {
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

func (cd *Wsjtx) stop() {
	cd.endFlag = true
	sleepTime := 2 * time.Second

	for {
		if !cd.threadFlag {
			return
		}

		time.Sleep(sleepTime)
	}
}
